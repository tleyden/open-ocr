#!/usr/bin/env bash

# Launches all components of the OpenOCR service
# 
# How to run this script for generic docker
# 
#   https://github.com/tleyden/open-ocr/blob/master/README.md
# 
# How to run this script for Orchard docker PAAS
# 
#   https://github.com/tleyden/open-ocr/wiki/Installation-on-Orchard
#

#if [ ! -n "$RABBITMQ_HOST" ] ; then
#  echo "You must define RABBITMQ_HOST"
#  exit
#fi

# let's assume we want the standard linked containers approach
# if RABBITMQ_HOST is not set, set it to a default value
[ -z ${RABBITMQ_HOST+x} ] && RABBITMQ_HOST="openocr_rabbitmq"

if [ ! -n "$RABBITMQ_PASS" ] ; then
  echo "You must define RABBITMQ_PASS"
  exit
fi

if [ ! -n "$HTTP_PORT" ] ; then
  echo "You must define HTTP_PORT"
  exit
fi

if [ "$1" == "orchard" ] ; then
      export DOCKER="orchard docker" 
elif [ "$1" == "gce" ] ; then
      export DOCKER="sudo docker" 
else
      export DOCKER="docker"
fi

export AMQP_URI=amqp://admin:${RABBITMQ_PASS}@${RABBITMQ_HOST}/

$DOCKER run -d -p 5672:5672 -p 15672:15672 -e RABBITMQ_PASS=${RABBITMQ_PASS} --name "${RABBITMQ_HOST}" tutum/rabbitmq

echo "Waiting 30s for rabbit MQ to startup .."
sleep 30 # workaround for startup race condition issue

$DOCKER run -d -p ${HTTP_PORT}:${HTTP_PORT} --link "${RABBITMQ_HOST}" tleyden5iwx/open-ocr open-ocr-httpd -amqp_uri "${AMQP_URI}" -http_port ${HTTP_PORT}

$DOCKER run -d --link "${RABBITMQ_HOST}" tleyden5iwx/open-ocr open-ocr-worker -amqp_uri "${AMQP_URI}"



