
# download all unit files
wget https://raw.githubusercontent.com/tleyden/open-ocr/master/docs/coreos-fleet/httpd.service
wget https://raw.githubusercontent.com/tleyden/open-ocr/master/docs/coreos-fleet/rabbitmq.service
wget https://raw.githubusercontent.com/tleyden/open-ocr/master/docs/coreos-fleet/rabbitmq_announce.service
wget https://raw.githubusercontent.com/tleyden/open-ocr/master/docs/coreos-fleet/worker.service

# start all unit files
echo "Starting units .."
fleetctl start rabbitmq.service rabbitmq_announce.service httpd.service worker.service

