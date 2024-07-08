kafka-oauth-demo.py
===================

The Kafka OAuth demo shows how a customer can authenticate their client via OAuth2 to Kafka and how they can read data
from a Kafka topic using the Kafka consumer API.

This is only a sample without support and liability to its correctness!

Prerequisite
------------

The script is based on Python 3.9 and [confluent-kafka-python](https://github.com/confluentinc/confluent-kafka-python)

Related Links:

* https://docs.confluent.io/kafka-clients/python/current/overview.html
* https://github.com/confluentinc/confluent-kafka-python

You can use the [requirements.txt](./requirements.txt) to install the dependencies using:

```bash
pip install -r requirements.txt
```

In this example we assume to use Let’s Encrypt CA for SSL/TLS certificates. The confluent-kafka-python-library by default uses the 
operating system's default trusted root CA certificates for secure connections. Please ensure that your system has the 
Let's Encrypt root certificates installed. These certificates are usually included in the system's trusted root store by
default. If not, please install them manually.

For Debian/Ubuntu distributions, the system CA certificates can be updated as:
```bash
sudo update-ca-certificates
```

How to use
----------

To use the sample, please change the following parameters, which should have been sent to you beforehand.

```python
client_id = 'YOUR_CLIENT_ID'                 # If you are an MBCon customer, please use the client id you have received
client_secret = 'YOUR_CLIENT_SECRET'         # If you are an MBCon customer, please use the secret you have received
scope = 'SCOPE'                              # use the correct scope for your region
topic = 'YOUR_DEDICATED_TOPIC'               # If you are an MBCon customer, please use topic name as 'vehiclesignals.<client name>'
group = 'CONSUMER_GROUP'                     # If you are an MBCon customer, please use the received client name as the prefix. eg: '<client name>.GROUP_ID_POSTFIX_OF_YOUR_CHOICE'

bootstrap_url = 'BOOTSTRAP_URL'              # use the correct broker url for your region
oauth_token_api_url = 'OAUTH_TOKEN_API_URL'  # use the correct token API url for your region
```

after preparation, you can start the demo with 
```bash
python kafka-oauth-demo.py
```
To adjust log verbosity, please refer to:
```bash
python kafka-oauth-demo.py -h
```

Notes
-----

confluent-kafka-python supports refreshing your oauth token automatically starting from version 2.2.0. Please update the
confluent-kafka-python library as soon as it is available to your platform. In client versions before 2.2.0, broker
connections will be closed, as soon as the token expires. The client will automatically reconnect, but you will receive
error logs.

Copyright 2024 Mercedes-Benz Connectivity Services GmbH
