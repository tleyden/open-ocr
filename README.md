
OpenOCR makes it simple to host your own OCR ReST API, powered by [Tesseract OCR](https://code.google.com/p/tesseract-ocr/)

![screenshot](http://tleyden-misc.s3.amazonaws.com/blog_images/openocr-architecture.png)

# Launching OpenOCR on Orchard

There are several [docker](http://www.docker.io) PAAS platforms available, and OpenOCR should work on all of them.  The following instructions are geared towards [Orchard](http://www.orchardup.com), but should be easily adaptable to other platforms.

## Install Orchard CLI tool

See the [Orchard Getting Started Guide](https://www.orchardup.com/docs)
for instructions on signing up and installing their CLI management tool.

## Find out your orchard host

```
$ orchard hosts
```

and you should see a result like:

```
NAME                SIZE                IP
default             512M                107.170.72.189
```

The ip address `107.170.72.189` will be used as the `ORCHARD_HOST` env variable below.

## Launch docker images

Here's how to launch the docker images needed for OpenOCR.

```
$ curl -O https://raw.githubusercontent.com/tleyden/open-ocr/master/launcher/launcher.sh
$ export ORCHARD_HOST=107.170.72.189 RABBITMQ_PASS=supersecret2 HTTP_PORT=8080
$ chmod +x launcher.sh
$ ./launcher.sh
```

*Note: as you may have noticed, you don't even need to clone this git repo!*

This will start three docker instances:

* RabbitMQ
* OpenOCR Worker
* OpenOCR HTTP API Server

You are now ready to decode images -> text via your REST API.

# Test the REST API 

**Request**

```
$ TEST_IMAGE=http://tleyden-misc.s3.amazonaws.com/blog_images/ocr_test.png
$ curl -X POST -H "Content-Type: application/json" -d '{"img_url":"$TEST_IMAGE","engine":"tesseract"}' http://$ORCHARD_HOST:$HTTP_PORT/ocr
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

