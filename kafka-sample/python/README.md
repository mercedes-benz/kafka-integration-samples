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

Our Kafka server uses Let’s Encrypt CA for SSL/TLS certificates. The confluent-kafka-python-library by default uses the 
operating system's default trusted root CA certificates for secure connections. Please ensure that your system has the 
Let's Encrypt root certificates installed. These certificates are usually included in the system's trusted root store by
default. If not, please install them manually.

For Debian/Ubuntu distributions, the system CA certificates can be updated as:
```bash
sudo update-ca-certificates
```

How to use
----------

In order to use the sample, please substitute the following parameters with those you have received via secure email.

```python
client_name = 'YOUR_CLIENT_NAME'             # use the client name you have received
client_id = 'YOUR_CLIENT_ID'                 # use the client you have received
client_secret = 'YOUR_CLIENT_SECRET'         # use the secret you have received
scope = 'SCOPE'                              # use the correct scope for your region
topic = f'vehiclesignals.{client_name}'      # use topic for the client you have received
group = f'{client_name}.GROUP_ID_POSTFIX'    # you can change the postfix of your consumer group

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
