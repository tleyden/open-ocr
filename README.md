
OpenOCR makes it simple to host your own OCR ReST API, powered by Tesseract OCR.  

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

The service consists of three docker images:

* RabbitMQ
* OpenOCR Worker
* OpenOCR HTTP API Server

### 

```
$ git clone https://github.com/tleyden/open-ocr.git
$ export ORCHARD_HOST=107.170.72.189 RABBITMQ_PASS=foo HTTP_PORT=8080
$ cd launcher
$ 

```


# Building Docker images

You can safely ignore these notes if you are just *using* the docker images.  But they may be useful if you need to rebuild your own images for any reason.

## Orchard "test" images

* git clone github.com/tleyden/docker
* cd docker/open-ocr
* orchard docker build -t open-ocr .
 
## Official docker.io images

* trusted build pointed to github.com/tleyden/docker/open-ocr