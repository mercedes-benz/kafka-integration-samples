#!/usr/bin/env python3
# coding: utf-8

#
# Copyright 2021 Mercedes-Benz Connectivity Services GmbH
#
# SPDX-License-Identifier: MIT

import argparse
import logging

from confluent_kafka import Consumer
from confluent_kafka.cimpl import KafkaException

client_name = 'YOUR_CLIENT_NAME'             # use the client name you have received
client_id = 'YOUR_CLIENT_ID'                 # use the client you have received
client_secret = 'YOUR_CLIENT_SECRET'         # use the secret you have received
scope = 'SCOPE'                              # use the correct scope for your region
topic = f'vehiclesignals.{client_name}'      # use topic for the client you have received
group = f'{client_name}.GROUP_ID_POSTFIX'    # you can change the postfix of your consumer group

bootstrap_url = 'BOOTSTRAP_URL'              # use the correct broker url for your region
oauth_token_api_url = 'OAUTH_TOKEN_API_URL'  # use the correct token API url for your region

security_protocol = 'SASL_SSL'
sasl_mechanism = 'OAUTHBEARER'


def consume_kafka():
    consumer = configure_consumer()
    app_logger.info("Consumer created")
    consumer.subscribe([topic])
    app_logger.info('Subscribed to topic %s', topic)

    app_logger.info("Starting poll loop")
    try:
        while True:
            msg = consumer.poll(timeout=1.0)
            if msg is None:
                app_logger.debug("No new messages")
                continue
            if msg.error():
                raise KafkaException(msg.error())
            else:
                # Process the received message
                print(msg.value().decode('utf-8'))
                # commit errors can be handled with
                # on_commit: Callable[[KafkaError, list[TopicPartition]], None] = commit_cb
                # config['on_commit']=on_commit
                consumer.commit(asynchronous=True)
    except KeyboardInterrupt:
        pass
    finally:
        app_logger.info("shutting down consumer.")
        consumer.close()


def configure_consumer():
    conf = {
        'bootstrap.servers': bootstrap_url,
        'security.protocol': security_protocol,
        'sasl.mechanism': sasl_mechanism,
        'sasl.oauthbearer.method': 'OIDC',
        'sasl.oauthbearer.client.id': client_id,
        'sasl.oauthbearer.client.secret': client_secret,
        'sasl.oauthbearer.scope': scope,
        'sasl.oauthbearer.token.endpoint.url': oauth_token_api_url,
        'group.id': group,
        'auto.offset.reset': 'earliest',
        'enable.auto.commit': True,
    }
    if args.kafka_debug != ["none"]:
        conf['debug'] = 'consumer,security' if args.kafka_debug is None else ','.join(args.kafka_debug)

    consumer_logger = logging.getLogger('consumer')
    handler = logging.StreamHandler()
    handler.setFormatter(logging.Formatter(LOG_FORMAT))
    consumer_logger.addHandler(handler)

    consumer = Consumer(conf, logger=consumer_logger)
    return consumer


def parse_arguments():
    parser = argparse.ArgumentParser(
        prog="oidc-demo",
        description="The Kafka OAuth demo shows how a customer can authenticate their client via OAuth2 to Kafka \
        and how they can read data from a Kafka topic API and using oidc.",
    )
    parser.add_argument(
        '-log',
        '--loglevel',
        default='info',
        help='set root loglevel, default=info'
    )
    parser.add_argument(
        '-kafka-debug',
        nargs='*',
        help='set verbosity of kafka debug logs, disable with -kafka-debug=none. \
        Only useful together with -log=debug. \
        See CONFIGURATION.md in https://github.com/confluentinc/librdkafka for valid `debug` options, \n\
        default=consumer security'
    )
    return parser.parse_args()


if __name__ == '__main__':
    args = parse_arguments()

    LOG_FORMAT = '[%(name)s]: %(asctime)s - %(levelname)s - %(message)s'
    logging.basicConfig(encoding='utf-8', level=args.loglevel.upper(), format=LOG_FORMAT)
    app_logger = logging.getLogger("kafka-oauth-demo")

    consume_kafka()
