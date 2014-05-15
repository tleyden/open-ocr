
OpenOCR makes it host your own OCR ReST API.  

![screenshot](http://tleyden-misc.s3.amazonaws.com/blog_images/openocr-architecture.png)

Currently [Tesseract](https://code.google.com/p/tesseract-ocr/) is the only supported OCR engine.

# REST API call example

**Request**

```
curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://cl.ly/image/132b2C0T1S3q/Screen%20Shot%202014-05-10%20at%2012.32.18%20PM.png","engine":0}' http://107.170.72.189:8081/ocr
```

**Response**

```
< HTTP/1.1 200 OK
< Date: Tue, 13 May 2014 16:18:50 GMT
< Content-Length: 283
< Content-Type: text/plain; charset=utf-8
<
You can create local variables for the pipelines within the template by
preﬁxing the variable name with a “$" sign. Variable names have to be
composed of alphanumeric characters and the underscore. In the example
below I have used a few variations that work for variable names.

```

# Launch OpenOCR on Orchard

There are several [docker](http://www.docker.io) PAAS platforms available, and OpenOCR should work on all of them.  The following instructions are geared towards [Orchard](http://www.orchardup.com), but should be easily adaptable to other platforms.

## Install Orchard CLI tool

See the [Orchard Getting Started Guide](https://www.orchardup.com/docs)
for instructions on signing up and installing their CLI management tool.

## Launch docker images

### RabbitMQ

```
orchard docker run -d -p 5672:5672 -p 15672:15672 tutum/rabbitmq
```

### OpenOCR Worker

```
orchard docker run  open-ocr-worker open-ocr-worker -amqp_uri "amqp://admin:8Sd7safsdafaukg@107.170.72.189:5672/"
```

### OpenOCR HTTP API Server

```
orchard docker run -p 8081:8081 open-ocr-worker open-ocr-httpd -amqp_uri "amqp://admin:8Sd7safsdafaukg@107.170.72.189:5672/"
```

# How to build docker images:

* git clone github.com/tleyden/docker
* cd docker/ocr-worker
* orchard docker build -t open-ocr-worker .
 