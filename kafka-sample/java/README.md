Java oauth kafka consumer
===================

The Kafka OAuth demo shows how a customer can authenticate their client via OAuth2 to Kafka and how they can read data
from a Kafka topic using the Kafka consumer API.

This is only a sample without support and liability to its correctness!

Prerequisite
------------

The code is based on java version 21 and gradle. Required dependencies:

* [gradle shadow plugin](https://github.com/GradleUp/shadow)
* [apache kafka-clients](https://kafka.apache.org/documentation/)
* [jose4j](https://bitbucket.org/b_c/jose4j/wiki/Home) (see also [KIP-1139](https://cwiki.apache.org/confluence/display/KAFKA/KIP-1139%3A+Add+support+for+OAuth+jwt-bearer+grant+type))
* [jackson-databind](https://mvnrepository.com/artifact/com.fasterxml.jackson.core/jackson-databind) (runtimeOnly)
* [log4j2](https://logging.apache.org/log4j/2.x/)
* [apache commons-cli](https://commons.apache.org/proper/commons-cli/index.html)

you can use the gradle `shadowJar` task to build an executable [jarfile](build/libs/java-0.1.0-all.jar):

```bash
gradle shadowJar
```


In this example we assume to use Let’s Encrypt CA for SSL/TLS certificates. These certificates are usually included in the default 
truststore of java environment. If not, please update your java version.


How to use
----------

To use the sample please change at least the following configurations of the [consumer.properties file](consumer.properties).

```properties
# use the correct bootstrap url for your region
bootstrap.servers=BOOTSTRAP_URL
# if you are an MBCon customer, use the received client name as the prefix. eg: '<client name>.GROUP_ID_POSTFIX_OF_YOUR_CHOICE':
group.id=CONSUMER_GROUP
# use the correct token API url for your region:
sasl.oauthbearer.token.endpoint.url=OAUTH_TOKEN_API_URL

# if you are an MBCon customer, use the clientId and clientSecret you have received along with the correct scope for your region:
sasl.oauthbearer.client.credentials.client.id=YOUR_CLIENT_ID
sasl.oauthbearer.client.credentials.client.secret=YOUR_CLIENT_SECRET
sasl.oauthbearer.scope=SCOPE
```

After preparation, you can start the demo with

```bash
java '-Dorg.apache.kafka.sasl.oauthbearer.allowed.urls=<OAUTH_TOKEN_API_URL>' -jar build/libs/java-0.1.0-all.jar -t <YOUR_TOPIC_NAME>
```

And note that, for the Kafka client version 4 and newer, the token URL must be explicitly repeated in the shown property.

Copyright 2024-2025 Mercedes-Benz Connectivity Services GmbH
