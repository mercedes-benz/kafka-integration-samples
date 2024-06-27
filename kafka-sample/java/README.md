Java oauth kafka consumer
===================

The Kafka OAuth demo shows how a customer can authenticate their client via OAuth2 to Kafka and how they can read data
from a Kafka topic API.

This is only a sample without support and liability to its correctness!

Prerequisite
------------

The code is based on java version 17 and gradle. Required dependencies:

* [com.github.johnrengelman.shadow](https://github.com/johnrengelman/shadow)
* [apache kafka-clients](https://kafka.apache.org/documentation/)
* [jackson-databind](https://mvnrepository.com/artifact/com.fasterxml.jackson.core/jackson-databind) (runtimeOnly)
* [log4j2](https://logging.apache.org/log4j/2.x/)
* [apache commons-cli](https://commons.apache.org/proper/commons-cli/index.html)

you can use the gradle `shadowJar` task to build an executable [jarfile](build/libs/java-0.1.0-all.jar):

```bash
gradle shadowJar
```

Please ensure that your Java environment includes Let's Encrypt CA certificates in its truststore to securely connect
to our server. If not, please ensure that your system has the Let's Encrypt root certificates installed. These
certificates are usually included in the system's trusted root store by default. If not, please install them manually.

For Debian/Ubuntu distributions, the CA certificates can be updated as:
```bash
sudo update-ca-certificates
```

How to use
----------

In order to use the sample please populate at least following configurations of the [consumer.properties file](consumer.properties).

```properties
# use the clientId and clientSecret you have received:
sasl.jaas.config=org.apache.kafka.common.security.oauthbearer.OAuthBearerLoginModule required \
clientId="YOUR_CLIENT_ID" \
clientSecret="YOUR_CLIENT_SECRET" \
scope="SCOPE";
# use the correct token API url for your region:
sasl.oauthbearer.token.endpoint.url=OAUTH_TOKEN_API_URL
# you can change the postfix of your consumer group:
group.id=YOUR_CLIENT_NAME.GROUP_ID_POSTFIX
# use the correct bootstrap url for your region
bootstrap.servers=BOOTSTRAP_URL
```

After preparation, you can start the demo with

```bash
java -jar build/libs/java-0.1.0-all.jar -t <YOUR_TOPIC_NAME>
```

Copyright 2023 Mercedes-Benz Connectivity Services GmbH
