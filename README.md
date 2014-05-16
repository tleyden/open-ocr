
OpenOCR makes it simple to host your own OCR REST API.

The heavy lifting OCR work is handled by [Tesseract OCR](https://code.google.com/p/tesseract-ocr/).

[Docker](http://www.docker.io) is used to containerize the various components of the service.

![screenshot](http://tleyden-misc.s3.amazonaws.com/blog_images/openocr-architecture.png)

# Launching OpenOCR on Orchard

There are several [docker](http://www.docker.io) PAAS platforms available, and OpenOCR should work on all of them.  The following instructions are geared towards [Orchard](http://www.orchardup.com), but should be easily adaptable to other platforms.

## Install Orchard CLI tool

See the [Orchard Getting Started Guide](https://www.orchardup.com/docs)
for instructions on signing up and installing their CLI management tool.

## Find out your orchard host

```
$ orchard hosts
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

This will start three docker instances:

* RabbitMQ
* OpenOCR Worker
* OpenOCR HTTP API Server

You are now ready to decode images -> text via your REST API.

# Test the REST API 

**Request**

```
$ curl -X POST -H "Content-Type: application/json" -d '{"img_url":"bit.ly/ocrimage","engine":"tesseract"}' http://$ORCHARD_HOST:$HTTP_PORT/ocr
```

**Response**

It will return the decoded text for the [test image](http://bit.ly/ocrimage):

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

