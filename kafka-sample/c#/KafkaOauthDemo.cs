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

            var clientId = "YOUR_CLIENT_ID";                    // use the client you have received
            var clientSecret = "YOUR_CLIENT_SECRET";            // use the secret you have received
            var topic = $"vehiclesignals.{clientId}";           // use topic for the client you have received
            var consumerGroup = $"{clientId}.GROUP_ID_POSTFIX"; // you can change the postfix of your consumer group
            var rootCaFile = @"cluster-ca.crt";                 // file path of your CA certificate (must be a PEM file)

            var bootstrapUrl = "BOOTSTRAP_URL";                 // use the correct broker url for your region
            var oauthTokenApiUrl = "OAUTH_TOKEN_API_URL";       // use the correct token API url for your region

            var securityProtocol = SecurityProtocol.SaslSsl;
            var saslMechanism = SaslMechanism.OAuthBearer;
            var sslEndpointIdentificationAlgorithm = SslEndpointIdentificationAlgorithm.None;
            var autoOffsetReset = AutoOffsetReset.Earliest;

            var config = new ConsumerConfig
            {
                BootstrapServers = bootstrapUrl,
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
