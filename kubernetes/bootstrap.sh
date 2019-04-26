#!/bin/sh -
#
# bootstrap
# automate the deployment of Open-OCR on a Kubernetes Cluster
# based on https://github.com/tleyden/open-ocr/tree/master/kubernetes
#
# Copyright (c) 2019 diego casati <dcasati@dcasati.net>
#
# Permission to use, copy, modify, and distribute this software for any
# purpose with or without fee is hereby granted, provided that the above
# copyright notice and this permission notice appear in all copies.
#
# THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
# WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
# MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
# ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
# WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
# ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
# OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
#

# usage: -c - clean up the demo.
cleanup() {
    echo "cleaning up the Kubernetes environment"

    kubectl delete -f \
    pods/rabbitmq.yaml,services/rabbitmq.yml,replication-controllers/open-ocr-worker.yaml,pods/open_ocr_httpd.yml,services/open_ocr_httpd.yml

    kubectl delete secrets rabbit-mq-password
    exit 0
}

# show usage and exit.
__usage="usage: `basename $0` [-cit]

Options:
    -c     clean up the demo.
    -i     install the demo.
    -t     run a cURL test against the API.
"
usage() {
    echo "$__usage"
    exit 1
}

# Attempt to create a random password for RabitMQ.
create_rabbitmq_secret(){
    echo "Creating a random RabbitMQ password"
    echo "You will want to replace the YOUR_RABBITMQ_PASS below with something more secure.\n"

    date | md5sum | awk '{print $1}' > ./password
    kubectl create secret generic rabbit-mq-password --from-file=./password
}

# Clone the repo if not cloned.
first_run() {
    local LOCAL_REPO="open-ocr"

    if [ -d $LOCAL_REPO ]; then
        echo Clone OpenOCR repo
        git clone https://github.com/tleyden/open-ocr.git
    fi
}

# Launch RabbitMQ 
launch_rabbitmq() {
    local RABBITMQ_STATUS
    kubectl create -f pods/rabbitmq.yaml
    kubectl create -f services/rabbitmq.yml

    printf "%s" "waiting until RabitMQ is ready"

    RABBITMQ_STATUS=`kubectl get po -o=jsonpath='{.items[?(@.metadata.labels.name=="rabbitmq")].status.phase}'`
    
    while [ $RABBITMQ_STATUS != "Running" ]; do
        RABBITMQ_STATUS=`kubectl get po -o=jsonpath='{.items[?(@.metadata.labels.name=="rabbitmq")].status.phase}'`
        printf "%s" "."
        sleep 2
    done
    echo "RabbitMQ is ready."
}

# Launch REST API Server
launch_rest_api(){
    echo "creating the REST API Server\n"
    kubectl create -f pods/open_ocr_httpd.yml
    kubectl create -f services/open_ocr_httpd.yml
}

# Launch OCR Worker
launch_ocr_worker(){
    echo "creating the Open-OCR workers\n"
    kubectl create -f replication-controllers/open-ocr-worker.yaml
}

# usage: -t - checks if the LoadBalancer IP address is up and running
test_rest_api() {
    LOADBALANCER_IP=`kubectl get service -o jsonpath='{.items[?(@.metadata.name=="open-ocr-httpd-service")].status.loadBalancer.ingress[].ip}'`
    
    echo "running curl against the REST API\n"
    curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://bit.ly/ocrimage","engine":"tesseract"}' http://$LOADBALANCER_IP/ocr
    
    # bail out if curl can't get to the REST API server
    if [ $? > 0 ]; then
        exit 1
    fi
            
    exit 0
}

# usage: -i - installs the entire demo
run_install() {
    first_run
    create_rabbitmq_secret
    launch_rabbitmq
    launch_rest_api
    launch_ocr_worker
    test_rest_api

    exit 0
}

while getopts "cit" opt; do
	case $opt in
	c)	cleanup
        ;;
    i)  run_install
        ;;
	t)	test_rest_api
        ;;
	*)	usage 
        exit 1
        ;;
	esac
done
shift $(( $OPTIND - 1 ))

if [ $OPTIND = 1 ]; then
	usage
	exit 0
fi