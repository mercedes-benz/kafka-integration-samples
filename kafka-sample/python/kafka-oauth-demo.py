#!/usr/bin/env python
# coding: utf-8

#
# Copyright 2021 Mercedes-Benz Connectivity Services GmbH
#
# SPDX-License-Identifier: MIT

# Load libraries, please replace with the way you integrate your libs.
import pip

def install(package):
    if hasattr(pip, 'main'):
        pip.main(['install', package])
    else:
        pip._internal.main(['install', package])



# Declarations and Parameters

import ssl

clientid = 'YOUR_CLIENT_ID' # CHANGE HERE
pwd = 'YOUR_CLIENT_PWD'     # CHANGE HERE
topic = 'TOPIC_NAME'        # CHANGE HERE
group = 'GROUP_NAME'        # you can change the postfix
root_ca_file = './ca.pem'   # CHANGE HERE, put here the location of your CA certificate (must be PEM file)

broker_url = 'BROKER_URL'
oauth_token_api_url = 'OAUTH_TOKEN_API_URL'

security_protocol = 'SASL_SSL'
sasl_mechanism = 'OAUTHBEARER'
sasl_user = clientid
sasl_pwd = pwd


# Create a new context using system defaults, disable all but TLS1.2
context = ssl.create_default_context(cafile=root_ca_file)
context.options &= ssl.OP_NO_TLSv1
context.options &= ssl.OP_NO_TLSv1_1


# The TokenProvider required for the OAuth authentication to retrieve the bearer token.

from aiokafka.abc import AbstractTokenProvider
from oauthlib.oauth2 import BackendApplicationClient
from requests_oauthlib import OAuth2Session

class CustomTokenProvider(AbstractTokenProvider):
        async def token(self):
            return await asyncio.get_running_loop().run_in_executor(
                None, self._token)
        
        
        def _token(self):
            client = BackendApplicationClient(client_id=sasl_user)
            oauth = OAuth2Session(client=client)            
            token_json = oauth.fetch_token(token_url=oauth_token_api_url, client_secret=sasl_pwd)
            print("Got token")
            token = token_json['access_token']

            return token


# Simple loop to read messages from topic

import asyncio
from aiokafka import AIOKafkaConsumer
from json import loads
from aiokafka.errors import KafkaError


print("Init consumer")

consumer = AIOKafkaConsumer(
    topic,
    bootstrap_servers=[broker_url],
    security_protocol = security_protocol,
    ssl_context = context,
    sasl_plain_username = sasl_user,
    sasl_plain_password = sasl_pwd,
    sasl_oauth_token_provider = CustomTokenProvider(),
    sasl_mechanism = sasl_mechanism,
    auto_offset_reset = 'earliest',
    enable_auto_commit = True,
    group_id=group,    
    value_deserializer=lambda x: loads(x.decode('utf-8')))


async def main():
    install('aiokafka')
    install('requests_oauthlib')

    await consumer.start()
    print("Got consumer")

    try:
        async for message in consumer:
            print("read msg ", message)
            print('o {}'.format(message.value))
    except Exception as ex:
        print(ex)

if __name__ == "__main__":
    main()


