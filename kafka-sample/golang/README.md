kafka-oauth-consumer.go
===================

The Kafka OAuth demo shows how a customer can authenticate their client via OAuth2 to Kafka and how they can read data
from a Kafka topic using the Kafka consumer API.

This is only a sample without support and liability to its correctness!

Prerequisite
------------

The code is based on Go 1.20 and the confluent-kafka-go v2.1.1.

Links:

* https://github.com/confluentinc/confluent-kafka-go

You can use the [go.mod](./go.mod) to install the dependencies and build the application using:

```bash
go install
go build kafka-oauth-consumer.go
```

In this example we assume to use Let’s Encrypt CA for SSL/TLS certificates. The confluent-kafka-go library uses the operating 
system's default trusted root CA certificates for secure connections. Please ensure that your system has the 
Let's Encrypt root certificates installed. These certificates are usually included in the system's trusted root store 
by default. If not, please install them manually.

For Debian/Ubuntu distributions, the system CA certificates can be updated as:
```bash
sudo update-ca-certificates
```

How to use
----------

To use the sample, please change the following parameters, which should have been sent to you beforehand.

```go
package main

const (
	ClientId         = "YOUR_CLIENT_ID"                 // If you are an MBCon customer, please use the client id you have received
	ClientSecret     = "YOUR_CLIENT_SECRET"             // If you are an MBCon customer, please use the secret you have received
	Scope            = "SCOPE"                          // use the correct scope for your region
	TopicName        = "YOUR_DEDICATED_TOPIC"           // If you are an MBCon customer, please use topic name as 'vehiclesignals.<client name>'
	GroupId          = "CONSUMER_GROUP"                 // If you are an MBCon customer, please use the received client name as the prefix. eg: '<client name>.GROUP_ID_POSTFIX_OF_YOUR_CHOICE'
	BootstrapUrl     = "BOOTSTRAP_URL"                  // use the correct broker url for your region
	OauthTokenApiUrl = "OAUTH_TOKEN_API_URL"            // use the correct token API url for your region
)
```

after preparation, you can start the demo with

```bash
kafka-oauth-consumer
```

To adjust log verbosity, please refer to:

```bash
kafka-oauth-consumer -h
```

Notes
-----

confluent-kafka-go supports refreshing your oauth token automatically starting from version 2.2.0. Please update the
confluent-kafka-go library as soon as it is available to your platform. In client versions before 2.2.0, broker
connections will be closed, as soon as the token expires. The client will automatically reconnect, but you will receive
error logs.

The bundled librdkafka of confluent-kafka-go does not support OIDC on all systems. If you get error logs mentioning that
some Configuration properties are not supported on this build, you either need to install the
underlying [librdkafka](https://github.com/confluentinc/librdkafka) library on your system and build the project using:

```bash
go build -tags dynamic kafka-oauth-consumer.go 
```

or implement a custom login callback. Details for installing librdkafka can be
found [here](https://github.com/confluentinc/confluent-kafka-go#installing-librdkafka).

Copyright 2024 Mercedes-Benz Connectivity Services GmbH
