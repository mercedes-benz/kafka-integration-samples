//
// Copyright 2022 Mercedes-Benz Connectivity Services GmbH
//
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const (
	PushApiUrl   = "PUSH_API_URL"    // use the correct broker url for your region
	GroupId      = "YOUR_GROUP_ID"   // you can change the postfix of your consumer group
	CertLocation = "PATH_TO_CERT"    // file path of your CA certificate (must be a PEM file)
	TopicName    = "YOUR_TOPIC_NAME" // use topic for the client you have received

	AuthUrl      = "OAUTH_TOKEN_API_URL" // use the correct token API url for your region
	ClientId     = "YOUR_CLIENT_ID"      // use the client you have received
	ClientSecret = "YOUR_CLIENT_SECRET"  // use the secret you have received
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
			"bootstrap.servers":                   PushApiUrl,
			"group.id":                            GroupId,
			"security.protocol":                   "SASL_SSL",
			"sasl.mechanism":                      "OAUTHBEARER",
			"sasl.oauthbearer.method":             "OIDC",
			"sasl.oauthbearer.client.id":          ClientId,
			"sasl.oauthbearer.client.secret":      ClientSecret,
			"sasl.oauthbearer.token.endpoint.url": AuthUrl,
			"ssl.ca.location":                     CertLocation,
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
