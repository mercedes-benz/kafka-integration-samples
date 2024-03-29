kafka-oauth-consumer.go
===================

The Kafka OAuth demo shows how a customer can authenticate their client via OAuth2 to Kafka and how they can read data
from a Kafka topic API.

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

On linux systems, the client expects the hosts' ca-certificate at `/etc/pki/tls/ca-bundle.crt` (fedora/redhat defaults).
If your systems default location differs from that (e.g. Debian/Ubuntu), please make sure to copy or symlink the
certificate.

Example for Debian/Ubuntu distributions:
```bash
mkdir -p /etc/pki/tls/certs
ln -s /etc/ssl/certs/ca-certificates.crt /etc/pki/tls/certs/ca-bundle.crt
```

How to use
----------

In order to use the sample please change the following parameters.

```go
package main

const (
	ClientId         = "YOUR_CLIENT_ID"               // use the client you have received
	ClientSecret     = "YOUR_CLIENT_SECRET"           // use the secret you have received
	TopicName        = "vehiclesignals." + ClientId   // use topic for the client you have received
	GroupId          = ClientId + ".GROUP_ID_POSTFIX" // you can change the postfix of your consumer group
	RootCaFile       = "PATH_TO_CERT"                 // file path of your CA certificate (must be a PEM file)
	BootstrapUrl     = "BOOTSTRAP_URL"                // use the correct broker url for your region
	OauthTokenApiUrl = "OAUTH_TOKEN_API_URL"          // use the correct token API url for your region

)
```

You also need to provide the CA certificate location. The certificate file must be in PEM format. You can extract the
PEM from a .p12 using `openssl` and `sed`:

```bash
openssl pkcs12 -in cluster-ca.p12 -nokeys | \
sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > cluster-ca.crt
```

Or export it, using a GUI-based tool like [keystore explorer](https://keystore-explorer.org/)

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

Copyright 2023 Mercedes-Benz Connectivity Services GmbH
