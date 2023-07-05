//
// Copyright 2021 Mercedes-Benz Connectivity Services GmbH
//
// SPDX-License-Identifier: MIT

using Confluent.Kafka;

namespace KafkaConsumer
{
    class Program
    {

        static void Main()
        {

            var clientId = "YOUR_CLIENT_ID";               // CHANGE HERE
            var clientSecret = "YOUR_CLIENT_SECRET";       // CHANGE HERE
            var topic = "TOPIC_NAME";                      // CHANGE HERE
            var consumerGroup = "CONSUMER_GROUP";          // CHANGE HERE            
            var rootCaFile = @"cert.pem";                  // CHANGE HERE, put here the location of your CA certificate (must be PEM file)

            var brokerUrl = "BROKER_URL";                  // CHANGE HERE
            var oauthTokenApiUrl = "OAUTH_TOKEN_API_URL";  // CHANGE HERE

            var securityProtocol = SecurityProtocol.SaslSsl;
            var saslMechanism = SaslMechanism.OAuthBearer;
            var sslEndpointIdentificationAlgorithm = SslEndpointIdentificationAlgorithm.None;
            var autoOffsetReset = AutoOffsetReset.Earliest;

            var config = new ConsumerConfig
            {
                BootstrapServers = brokerUrl,
                GroupId = consumerGroup,

                SslCaLocation = rootCaFile,
                SslEndpointIdentificationAlgorithm = sslEndpointIdentificationAlgorithm,
                SecurityProtocol = securityProtocol,
                SaslMechanism = saslMechanism,

                SaslOauthbearerMethod = SaslOauthbearerMethod.Oidc,
                SaslOauthbearerClientId = clientId,
                SaslOauthbearerClientSecret = clientSecret,
                SaslOauthbearerTokenEndpointUrl = oauthTokenApiUrl,

                Debug = "consumer,security",
                AutoOffsetReset = autoOffsetReset,
            };


            using var consumer = new ConsumerBuilder<Ignore, string>(config).Build();
            consumer.Subscribe(topic);

            CancellationTokenSource cts = new();
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

                        Console.WriteLine($"Consumed message '{consumerResult.Message.Value}' at: '{consumerResult.TopicPartitionOffset}'.");
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
}
