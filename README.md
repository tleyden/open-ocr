
OpenOCR makes it simple to host your own OCR REST API.

The heavy lifting OCR work is handled by [Tesseract OCR](https://code.google.com/p/tesseract-ocr/).

[Docker](http://www.docker.io) is used to containerize the various components of the service.

![screenshot](http://tleyden-misc.s3.amazonaws.com/blog_images/openocr-architecture.png)

# Launching OpenOCR on Ubuntu 14.04

## Install Docker

See [Installing Docker on Ubuntu](https://docs.docker.com/installation/ubuntulinux/) instructions.

## Find out your host address

```
$ ifconfig
eth0      Link encap:Ethernet  HWaddr 08:00:27:43:40:c7
          inet addr:10.0.2.15  Bcast:10.0.2.255  Mask:255.255.255.0
          ...
```

The ip address `10.0.2.15` will be used as the `DOCKER_HOST` env variable below.

## Launch docker images

Here's how to launch the docker images needed for OpenOCR.

```
$ curl -O https://raw.githubusercontent.com/tleyden/open-ocr/master/launcher/launcher.sh
$ export DOCKER_HOST=10.0.2.15 RABBITMQ_PASS=supersecret2 HTTP_PORT=8080
$ chmod +x launcher.sh
$ ./launcher.sh
```

This will start three docker instances:

* [RabbitMQ](https://index.docker.io/u/tutum/rabbitmq/)
* [OpenOCR Worker](https://index.docker.io/u/tleyden5iwx/open-ocr/)
* [OpenOCR HTTP API Server](https://index.docker.io/u/tleyden5iwx/open-ocr/)

You are now ready to decode images → text via your REST API.

# Test the REST API 

**Request**

```
$ curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://bit.ly/ocrimage","engine":"tesseract"}' http://$DOCKER_HOST:$HTTP_PORT/ocr
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

# Launching OpenOCR on a Docker PAAS

You can also run OpenOCR on any PAAS that supports Docker containers.  Here are the instructions for a few that have already been tested:

* [Google Compute Engine](https://github.com/tleyden/open-ocr/wiki/Installation-on-Google-Compute-Engine)
* [Orchard](https://github.com/tleyden/open-ocr/wiki/Installation-on-Orchard)
* More coming soon ..

# License

OpenOCR is Open Source and available under the Apache 2 License.
