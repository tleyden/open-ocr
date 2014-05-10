
# How to build docker images:

* git clone github.com/tleyden/docker
* cd docker/ocr-worker
* orchard docker build -t open-ocr-worker .
 
## Test the worker image:

* orchard docker run -i -t open-ocr-worker /bin/bash
* cd /opt/go/src/github.com/tleyden/open-ocr/cli-worker
* go build -v
* ./cli-worker -amqp_uri "amqp://admin:<password>@101.170.72.189:5672/

## Test the httpd image:

* orchard docker run -p 8081:8081 -i -t open-ocr-worker /bin/bash
* cd /opt/go/src/github.com/tleyden/open-ocr/cli-httpd
* go build -v
* ./cli-httpd -amqp_uri "amqp://admin:<password>@101.170.72.189:5672/

# Test curl request:

```
curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://101.170.72.189:8081/img","engine":0}' http://101.170.72.189:8081/ocr
```