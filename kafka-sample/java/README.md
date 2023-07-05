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

How to use
----------

In order to use the sample please populate at least following configurations of the [consumer.properties file](consumer.properties).

```properties
# use the clientId and clientSecret you have received:
sasl.jaas.config=org.apache.kafka.common.security.oauthbearer.OAuthBearerLoginModule required \
clientId="YOUR_CLIENT_ID" \
clientSecret="YOUR_CLIENT_SECRET";
# use the correct token API url for your region:
sasl.oauthbearer.token.endpoint.url=OAUTH_TOKEN_API_URL
# you can change the postfix of your consumer group:
group.id=YOUR_GROUP_ID
# use the correct broker url for your region
bootstrap.servers=bootstrap.msg-stream-dev.connect-business.net:443
# file path of the .p12 truststore you received:
ssl.truststore.location=PATH_TO_TRUSTSTORE
# password of the .p12 truststore you received:
ssl.truststore.password=TRUSTSTORE_PASSWORD
```

After preparation, you can start the demo with

```bash
java -jar build/libs/java-0.1.0-all.jar -t <YOUR_TOPIC_NAME>
```

Copyright 2023 Mercedes-Benz Connectivity Services GmbH
