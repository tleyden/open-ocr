
[![Build Status](https://drone.io/github.com/tleyden/open-ocr/status.png)](https://drone.io/github.com/tleyden/open-ocr/latest) [![GoDoc](http://godoc.org/github.com/tleyden/open-ocr?status.png)](http://godoc.org/github.com/tleyden/open-ocr) 


OpenOCR makes it simple to host your own OCR REST API.

The heavy lifting OCR work is handled by [Tesseract OCR](https://code.google.com/p/tesseract-ocr/).

[Docker](http://www.docker.io) is used to containerize the various components of the service.

![screenshot](http://tleyden-misc.s3.amazonaws.com/blog_images/openocr-architecture.png)

# Features

* Scalable message passing architecture via RabbitMQ.
* Platform independence via Docker containers.
* Supports 31 languages in addition to English 
* Ability to use an image pre-processing chain.  An example using [Stroke Width Transform](https://github.com/tleyden/open-ocr/wiki/Stroke-Width-Transform) is provided.
* Pass arguments to Tesseract such as character whitelist and page segment mode.
* [REST API docs](http://docs.openocr.apiary.io/)
* A [Go REST client](http://github.com/tleyden/open-ocr-client) is available.

# Launching OpenOCR on a Docker PAAS

OpenOCR can easily run on any PAAS that supports Docker containers.  Here are the instructions for a few that have already been tested:

* [Google Compute Engine](https://github.com/tleyden/open-ocr/wiki/Installation-on-Google-Compute-Engine)
* [AWS](https://github.com/tleyden/open-ocr/wiki/Installation-on-CoreOS-Fleet)
* [Tutum](https://github.com/tleyden/open-ocr/wiki/Installation-on-Tutum)

If your preferred PAAS isn't listed, please open a [Github issue](https://github.com/tleyden/open-ocr/issues) to request instructions.

# Launching OpenOCR on Ubuntu 14.04

OpenOCR can be launched on anything that supports Docker, such as Ubuntu 14.04.  

Here's how to install it from scratch and verify that it's working correctly.

## Install Docker

See [Installing Docker on Ubuntu](https://docs.docker.com/installation/ubuntulinux/) instructions.

## Find out your host address

```
$ ifconfig
eth0      Link encap:Ethernet  HWaddr 08:00:27:43:40:c7
          inet addr:10.0.2.15  Bcast:10.0.2.255  Mask:255.255.255.0
          ...
```

The ip address `10.0.2.15` will be used as the `RABBITMQ_HOST` env variable below.

## Launch docker images

Here's how to launch the docker images needed for OpenOCR.

```
$ curl -O https://raw.githubusercontent.com/tleyden/open-ocr/master/launcher/launcher.sh
$ export RABBITMQ_HOST=10.0.2.15 RABBITMQ_PASS=supersecret2 HTTP_PORT=8080
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
$ curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://bit.ly/ocrimage","engine":"tesseract"}' http://10.0.2.15:$HTTP_PORT/ocr
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

The REST API also supports:

* Uploading the image content via `multipart/related`, rather than passing an image URL.  (example client code provided in the [Go REST client](http://github.com/tleyden/open-ocr-client))
* Tesseract config vars (eg, equivalent of -c arguments when using Tesseract via the command line) and Page Seg Mode 
* Ability to use an image pre-processing chain, eg [Stroke Width Transform](https://github.com/tleyden/open-ocr/wiki/Stroke-Width-Transform).
* Non-English languages

See the [REST API docs](http://docs.openocr.apiary.io/) and the [Go REST client](http://github.com/tleyden/open-ocr-client) for details.


# Uploading local files using curl

The supplied `docs/upload-local-file.sh` provides an example of how to upload a local file using curl with `multipart/related` encoding of the json and image data:
* usage: `docs/upload-local-file.sh <urlendpoint> <file> [mimetype]`
* download the example ocr image `wget http://bit.ly/ocrimage`
* example: `docs/upload-local-file.sh http://10.0.2.15:$HTTP_PORT/ocr-file-upload ocrimage` 


# Community

* Follow [@OpenOCR](https://twitter.com/openocr) on Twitter
* Checkout the [Github issue tracker](https://github.com/tleyden/open-ocr/issues)

# License

OpenOCR is Open Source and available under the Apache 2 License.
