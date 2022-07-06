//
// Copyright 2022 Mercedes-Benz Connectivity Services GmbH
//
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

type AuthResponse struct {
	AccessToken   string `json:"access_token" default:""`
	AccessExpires int64  `json:"expires_in" default:"0"`
	TokenType     string `json:"token_type" default:""`
	Scope         string `json:"scope" default:""`
}

func main() {
	kafkaConsumer, err := kafka.NewConsumer(
		// see https://docs.confluent.io/platform/current/clients/librdkafka/html/md_CONFIGURATION.html for full configuration documentation
		&kafka.ConfigMap{
			"bootstrap.servers": PushApiUrl,
			"group.id":          GroupId,
			"security.protocol": "SASL_SSL",
			"sasl.mechanism":    "OAUTHBEARER",
			"ssl.ca.location":   CertLocation,
			//"debug":             "all",
		},
	)

	if err != nil {
		fmt.Printf("Failed creating consumer: %v\n", err)
		os.Exit(1)
	}

	go consumeMessages(kafkaConsumer)

	// set up a channel for interrupting execution
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func consumeMessages(kafkaConsumer *kafka.Consumer) {
	err := kafkaConsumer.Subscribe(TopicName, nil)

	if err != nil {
		fmt.Printf("Error during subsribing topic: %v\n", err)
		return
	}

	for {
		result := kafkaConsumer.Poll(10000)

		switch event := result.(type) {
		case kafka.OAuthBearerTokenRefresh:
			refreshBearerToken(kafkaConsumer)
		case *kafka.Message:
			// actual event processing possible
			fmt.Printf("New message recieved: \n%v\n", string(event.Value))
		case kafka.Error:
			fmt.Printf("Kafka consumer error: %v\n", event)
		default:
			fmt.Printf("No new message: %v\n", event)
		}
	}
}

func refreshBearerToken(kafkaConsumer *kafka.Consumer) {
	var authResponse AuthResponse
	httpClient := http.Client{}

	data := url.Values{}
	data.Add("grant_type", "client_credentials")

	httpAuthRequest, err := http.NewRequest("POST", AuthUrl, strings.NewReader(data.Encode()))

	if err != nil {
		fmt.Printf("Error during request building: %v\n", err)
		return
	}

	httpAuthRequest.SetBasicAuth(ClientId, ClientSecret)
	httpAuthRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpAuthResponse, err := httpClient.Do(httpAuthRequest)

	if err != nil {
		fmt.Printf("Error during auth request: %v\n", err)
		return
	}

	err = json.NewDecoder(httpAuthResponse.Body).Decode(&authResponse)

	if err != nil {
		fmt.Printf("Error during response decoding: %v\n", err)
		return
	}

	err = kafkaConsumer.SetOAuthBearerToken(kafka.OAuthBearerToken{
		TokenValue: authResponse.AccessToken,
		Expiration: time.Now().Add(time.Second * time.Duration(authResponse.AccessExpires)),
	})

	if err != nil {
		fmt.Printf("Error during updating bearer token: %v\n", err)
		return
	}
}
