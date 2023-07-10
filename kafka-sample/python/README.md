kafka-oauth-demo.py
===================

The Kafka OAuth demo shows how a customer can authenticate their client via OAuth2 to Kafka and how they can read data
from a Kafka topic API and.

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

In order to use the sample please change the following parameters. Those parameters should previously send to you.

```python
client_id = 'YOUR_CLIENT_ID'  # use the client you have received
client_secret = 'YOUR_CLIENT_SECRET'  # use the secret you have received
topic = f'vehiclesignals.{client_id}'  # use topic for the client you have received
group = f'{client_id}.GROUP_ID_POSTFIX'  # you can change the postfix of your consumer group

bootstrap_url = 'BOOTSTRAP_URL'  # use the correct broker url for your region
oauth_token_api_url = 'OAUTH_TOKEN_API_URL'  # use the correct token API url for your region
```

You also need to provide the CA certificate location.

```python
root_ca_file = './cluster-ca.crt'  # file path of your CA certificate (must be a PEM file)
```

The certificate file must be in PEM format. You can extract the PEM from a .p12 using `openssl` and `sed`:

```bash
openssl pkcs12 -in cluster-ca.p12 -nokeys | \
sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > cluster-ca.crt
```

Or export it, using a GUI-based tool like [keystore explorer](https://keystore-explorer.org/)

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

Copyright 2023 Mercedes-Benz Connectivity Services GmbH
