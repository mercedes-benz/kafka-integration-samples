//
// Copyright 2021 Mercedes-Benz Connectivity Services GmbH
//
// SPDX-License-Identifier: MIT

using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Confluent.Kafka;
using Newtonsoft.Json;

namespace KafkaConsumer
{
    class Program
    {
        static void Main(string[] args)
        {

            var clientId = "YOUR_CLIENT_ID";                        // CHANGE HERE
            var clientSecret = "YOUR_CLIENT_SECRET";                // CHANGE HERE
            var topic = "TOPIC_NAME";                               // CHANGE HERE
            var consumerGroup = "CONSUMER_GROUP";                   // CHANGE HERE            
            var rootCaFile = @"cert.pem";                           // CHANGE HERE, put here the location of your CA certificate (must be PEM file)

            var brokerUrl = "BROKER_URL";                           // CHANGE HERE
            var oauthTokenApiUrl = new Uri("OAUTH_TOKEN_API_URL");  // CHANGE HERE
            var scope = "";

            var securityProtocol = SecurityProtocol.SaslSsl;
            var saslMechanism = SaslMechanism.OAuthBearer;
            var sslEndpointIdentificationAlgorithm = SslEndpointIdentificationAlgorithm.None;
            var autoOffsetReset = AutoOffsetReset.Earliest;

            void Callback(IClient client, string cfg)
            {
                var issuedAt = DateTimeOffset.UtcNow;
                var expiresAt = issuedAt.AddMinutes(5);

                var oAuthToken = RequestTokenToAuthorizationServer(oauthTokenApiUrl, clientId, scope, clientSecret).GetAwaiter().GetResult();        
                var accessToken = JsonConvert.DeserializeObject<Token>(oAuthToken).access_token;
                client.OAuthBearerSetToken(accessToken, expiresAt.ToUnixTimeMilliseconds(), clientId);
            }

            var config = new ConsumerConfig
            {
                BootstrapServers = brokerUrl,
                GroupId = consumerGroup,
                ClientId = clientId,

                SslCaLocation = rootCaFile,
                SslEndpointIdentificationAlgorithm = sslEndpointIdentificationAlgorithm,
                SecurityProtocol = securityProtocol,
                SaslMechanism = saslMechanism,
                
                Debug = "security",
                AutoOffsetReset = autoOffsetReset,
            };


            using (var consumer = new ConsumerBuilder<Ignore, string>(config).SetOAuthBearerTokenRefreshHandler(Callback).Build())
            {
                consumer.Subscribe(topic);

                CancellationTokenSource cts = new CancellationTokenSource();
                Console.CancelKeyPress += (_, e) =>
                {
                    e.Cancel = true; // prevent the process from terminating.
                    cts.Cancel();
                };

                try
                {
                    while (true)
                    {
                        try
                        {
                            var consumerResult = consumer.Consume(cts.Token);

                            Console.WriteLine($"Consumed message '{consumerResult.Value}' at: '{consumerResult.TopicPartitionOffset}'.");
                        }
                        catch (ConsumeException e)
                        {
                            Console.WriteLine($"Error occurred: {e.Error.Reason}");

                        }
                    }
                }
                catch (OperationCanceledException)
                {
                    // Ensure the consumer leaves the group cleanly and final offsets are committed.
                    consumer.Close();
                }
            }
        }

            private static async Task<string> RequestTokenToAuthorizationServer(Uri uriAuthorizationServer, string clientId, string scope, string clientSecret)
            {
                HttpResponseMessage responseMessage;

            var handler = new HttpClientHandler();
            handler.UseProxy = false;

            using (HttpClient client = new HttpClient(handler))
                {

                    HttpRequestMessage tokenRequest = new HttpRequestMessage(HttpMethod.Post, uriAuthorizationServer);
                    HttpContent httpContent = new FormUrlEncodedContent(
                        new[]
                        {
                    new KeyValuePair<string, string>("grant_type", "client_credentials"),
                    new KeyValuePair<string, string>("client_id", clientId),
                    new KeyValuePair<string, string>("scope", scope),
                    new KeyValuePair<string, string>("client_secret", clientSecret)
                        });
                    tokenRequest.Content = httpContent;
                    responseMessage = await client.SendAsync(tokenRequest);
                }
                return await responseMessage.Content.ReadAsStringAsync();
            }
    }

    public class Token
    {
        public string access_token { get; set; }
    }
}
