
# download all unit files
wget https://raw.githubusercontent.com/tleyden/open-ocr/master/docs/coreos-fleet/httpd.service
wget https://raw.githubusercontent.com/tleyden/open-ocr/master/docs/coreos-fleet/rabbitmq.service
wget https://raw.githubusercontent.com/tleyden/open-ocr/master/docs/coreos-fleet/rabbitmq_announce.service
wget https://raw.githubusercontent.com/tleyden/open-ocr/master/docs/coreos-fleet/worker.service

# run 3 workers by default
cp worker.service worker.1.service; cp worker.service worker.2.service; cp worker.service worker.3.service

# start all unit files
echo "Starting units .."
fleetctl start rabbitmq.service rabbitmq_announce.service httpd.service worker.[0-9].service

