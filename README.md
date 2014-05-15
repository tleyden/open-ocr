
OpenOCR makes it host your own OCR ReST API.  

![screenshot](http://tleyden-misc.s3.amazonaws.com/blog_images/openocr-architecture.png)

# Supported OCR engines

* [Tesseract](https://code.google.com/p/tesseract-ocr/)

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

# Launching OpenOCR on Orchard

There are several [docker](http://www.docker.io) PAAS platforms available, and OpenOCR should work on all of them.  The following instructions are geared towards [Orchard](http://www.orchardup.com), but should be easily adaptable to other platforms.

## Install Orchard CLI tool

See the [Orchard Getting Started Guide](https://www.orchardup.com/docs)
for instructions on signing up and installing their CLI management tool.

## Launch docker images

### RabbitMQ

```
orchard docker run -d -p 5672:5672 -p 15672:15672 tutum/rabbitmq
```

You will need to use `orchard hosts` and `orchard docker logs <container_id>` to get the amqp uri to use in later steps, eg `amqp://admin:8Sd7safsdafaukg@107.170.72.189:5672/`

### OpenOCR Worker

```
orchard docker run -d tleyden5iwx/open-ocr open-ocr-worker -amqp_uri "amqp://admin:8Sd7safsdafaukg@107.170.72.189:5672/"
```

### OpenOCR HTTP API Server

```
orchard docker run -d -p 8081:8081 tleyden5iwx/open-ocr open-ocr-httpd -amqp_uri "amqp://admin:8Sd7safsdafaukg@107.170.72.189:5672/"
```

# Building Docker images

You can safely ignore these notes if you are just *using* the docker images.  But they may be useful if you need to rebuild your own images for any reason.

## Orchard "test" images

* git clone github.com/tleyden/docker
* cd docker/open-ocr
* orchard docker build -t open-ocr .
 
## Official docker.io images

* trusted build pointed to github.com/tleyden/docker/open-ocr