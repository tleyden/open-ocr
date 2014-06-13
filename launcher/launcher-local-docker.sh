#!/usr/bin/env bash

# Launches all components of the OpenOCR service on docker running locally
# 
# NOTE: For DOCKER_HOST, you can run ifconfig and use the eth0 interface address
# 
# To run this script:
#
# $ export DOCKER_HOST=10.0.2.15
# $ export RABBITMQ_PASS=supersecret2
# $ export HTTP_PORT=8080  
# $ ./launcher-local-docker.sh

if [ ! -n "$DOCKER_HOST" ] ; then
  echo "You must define DOCKER_HOST"
  exit
fi

if [ ! -n "$RABBITMQ_PASS" ] ; then
  echo "You must define RABBITMQ_PASS"
  exit
fi

if [ ! -n "$HTTP_PORT" ] ; then
  echo "You must define HTTP_PORT"
  exit
fi

export AMQP_URI=amqp://admin:${RABBITMQ_PASS}@${DOCKER_HOST}/

docker run -d -p 5672:5672 -p 15672:15672 -e RABBITMQ_PASS=${RABBITMQ_PASS} tutum/rabbitmq

echo "Waiting 30s for rabbit MQ to startup .."
sleep 30 # workaround for startup race condition issue

docker run -d -p ${HTTP_PORT}:${HTTP_PORT} tleyden5iwx/open-ocr open-ocr-httpd -amqp_uri "${AMQP_URI}" -http_port ${HTTP_PORT}

docker run -d tleyden5iwx/open-ocr open-ocr-worker -amqp_uri "${AMQP_URI}"



