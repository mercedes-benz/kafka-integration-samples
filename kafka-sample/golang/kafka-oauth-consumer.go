//
// Copyright 2022 Mercedes-Benz Connectivity Services GmbH
//
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
)

const (
	ClientId         = "YOUR_CLIENT_ID"       // If you are an MBCon customer, please use the client id you have received
	ClientSecret     = "YOUR_CLIENT_SECRET"   // If you are an MBCon customer, please use the secret you have received
	Scope            = "SCOPE"                // use the correct scope for your region
	TopicName        = "YOUR_DEDICATED_TOPIC" // If you are an MBCon customer, please use topic name as 'vehiclesignals.<client name>'
	GroupId          = "CONSUMER_GROUP"       // If you are an MBCon customer, please use the received client name as the prefix. eg: '<client name>.GROUP_ID_POSTFIX_OF_YOUR_CHOICE'
	BootstrapUrl     = "BOOTSTRAP_URL"        // use the correct broker url for your region
	OauthTokenApiUrl = "OAUTH_TOKEN_API_URL"  // use the correct token API url for your region
)

func forwardLogs(logsChan chan kafka.LogEvent) {
	for logEvent := range logsChan {
		var logLevel log.Level
		switch logEvent.Level {
		case 7:
			logLevel = log.DebugLevel
		case 6:
		case 5:
			logLevel = log.InfoLevel
		case 4:
			logLevel = log.WarnLevel
		case 3:
			logLevel = log.ErrorLevel
		default:
			log.Fatalf("[%v, %v] %v", logEvent.Name, logEvent.Tag, logEvent.Message)
		}
		log.StandardLogger().Logf(logLevel, "[%v, %v] %v", logEvent.Name, logEvent.Tag, logEvent.Message)
	}
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	})
	level := parseArguments()
	log.SetLevel(level)

	log.Debug("Creating kafka consumer")

	kafkaConsumer, err := kafka.NewConsumer(
		// see https://docs.confluent.io/platform/current/clients/librdkafka/html/md_CONFIGURATION.html for full configuration documentation
		&kafka.ConfigMap{
			"bootstrap.servers":                   BootstrapUrl,
			"group.id":                            GroupId,
			"security.protocol":                   "SASL_SSL",
			"sasl.mechanism":                      "OAUTHBEARER",
			"sasl.oauthbearer.method":             "OIDC",
			"sasl.oauthbearer.client.id":          ClientId,
			"sasl.oauthbearer.client.secret":      ClientSecret,
			"sasl.oauthbearer.scope":              Scope,
			"sasl.oauthbearer.token.endpoint.url": OauthTokenApiUrl,
			"go.logs.channel.enable":              true,
			"debug":                               "consumer,security",
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	kafkaLogs := kafkaConsumer.Logs()

	go consumeMessages(kafkaConsumer)
	go forwardLogs(kafkaLogs)

	// set up a channel for interrupting execution
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func parseArguments() log.Level {
	var logLevelArg = flag.String("log", "info", "sets the log level")
	flag.Parse()
	level, err := log.ParseLevel(*logLevelArg)
	if err != nil {
		log.Fatal(err)
	}
	return level
}

func consumeMessages(kafkaConsumer *kafka.Consumer) {
	log.Debugf("Subscribing to topic %v", TopicName)
	err := kafkaConsumer.Subscribe(TopicName, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Starting polling loop")
	for {
		result := kafkaConsumer.Poll(10000)

		switch event := result.(type) {
		case *kafka.Message:
			// actual event processing possible
			log.Debug("New message received")
			print(string(event.Value))
		case kafka.Error:
			log.Errorf("Kafka consumer error: %v", event)
		default:
			log.Debugf("No new message: %v", event)
		}
	}
}
