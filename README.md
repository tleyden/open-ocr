[![GoDoc](http://godoc.org/github.com/tleyden/open-ocr?status.png)](http://godoc.org/github.com/tleyden/open-ocr) 
[![Join the chat at https://gitter.im/tleyden/open-ocr](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/tleyden/open-ocr?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


OpenOCR makes it simple to host your own OCR REST API.

The heavy lifting OCR work is handled by [Tesseract OCR](https://code.google.com/p/tesseract-ocr/).

[Docker](http://www.docker.io) is used to containerize the various components of the service.

![screenshot](http://tleyden-misc.s3.amazonaws.com/blog_images/openocr-architecture.png)

# Features

* Scalable message passing architecture via RabbitMQ.
* Platform independence via Docker containers.
* [Kubernetes support](https://github.com/tleyden/open-ocr/tree/master/kubernetes): workers can run in a Kubernetes Replication Controller
* Supports 31 languages in addition to English 
* Ability to use an image pre-processing chain.  An example using [Stroke Width Transform](https://github.com/tleyden/open-ocr/wiki/Stroke-Width-Transform) is provided.
* Pass arguments to Tesseract such as character whitelist and page segment mode.
* [REST API docs](http://docs.openocr.apiary.io/)
* A [Go REST client](http://github.com/tleyden/open-ocr-client) is available.


# Launching OpenOCR on a Docker PAAS

OpenOCR can easily run on any PAAS that supports Docker containers.  Here are the instructions for a few that have already been tested:

* [Launch on Google Container Engine GKE - Kubernetes](https://github.com/tleyden/open-ocr/wiki/Installation-on-Google-Container-Engine)
* [Launch on AWS with CoreOS](https://github.com/tleyden/open-ocr/wiki/Installation-on-CoreOS-Fleet)
* [Launch on Google Compute Engine](https://github.com/tleyden/open-ocr/wiki/Installation-on-Google-Compute-Engine)

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


# Launching OpenOCR command run.sh

 * [Install docker](https://docs.docker.com/installation/)
 * [Install docker-compose](https://docs.docker.com/compose/)
 * `git clone https://github.com/tleyden/open-ocr.git`
 * `cd open-ocr/docker-compose`
 * Type ```./run.sh ``` (in case you don't have execute right type ```sudo chmod +x run.sh```
 * The runner will ask you if you want to delete the images (choose y or n for each)
 * The runner will ask you to choose between version 1 and 2
   * Version 1 is using the ocr Tesseract 3.04. The memory usage is light. It is pretty fast and not costly in term of size (a simple aws instance with 1GB of ram and 8GB of storage is sufficiant). Result are acceptable
   * Version 2 is using the ocr Tesseract 4.00. The memory usage is light. It is less fast than tesseract 3 and more costly in term of size (an simple aws instance with 1GB of ram is sufficient but with an EBS of 16GB of storage). Result are really better compared to version 3.04.
   * To see a comparative you can have a look to the [official page of tesseract](https://github.com/tesseract-ocr/tesseract/wiki/4.0-Accuracy-and-Performance)


**You can use the docker-compose without the run.sh. For this just do:**

```
# for v1
export OPEN_OCR_INSTANCE=open-ocr

# for v2
export OPEN_OCR_INSTANCE=open-ocr-2

# then up (with -d to start it as deamon)
docker-compose up

```

Docker Compose will start four docker instances

* [RabbitMQ](https://index.docker.io/u/tutum/rabbitmq/)
* [OpenOCR Worker](https://index.docker.io/u/tleyden5iwx/open-ocr/)
* [OpenOCR HTTP API Server](https://index.docker.io/u/tleyden5iwx/open-ocr/)
* [OpenOCR Transform Worker](https://registry.hub.docker.com/u/tleyden5iwx/open-ocr-preprocessor/)

You are now ready to decode images → text via your REST API.

# Launching OpenOCR with Docker Compose on OSX

 * [Install docker](https://docs.docker.com/installation/)
 * [Install docker toolbox](https://www.docker.com/products/docker-toolbox)
 * Checkout OpenOCR repository 
 * `cd docker-compose directory`
 * `docker-machine start default`
 * `docker-machine env` 
 * Look at the Docker host IP address
 * Run  `docker-compose up -d` to run containers as daemons or `docker-compose up` to see the log in console
 

## How to test the REST API after turning on the docker-compose up

Where `IP_ADDRESS_OF_DOCKER_HOST` is what you saw when you run `docker-machine env` (e.g. 192.168.99.100)
and where `HTTP_POST` is the port number inside the `.yml` file inside the docker-compose directory presuming it should be the same 9292.

**Request**

```
$ curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://bit.ly/ocrimage","engine":"tesseract"}' http://IP_ADDRESS_OF_DOCKER_HOST:HTTP_PORT/ocr
```

Assuming the values are (192.168.99.100 and 9292 respectively)

```
$ curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://bit.ly/ocrimage","engine":"tesseract"}' http://192.168.99.100:9292/ocr
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
 
# Test the REST API 

## With image url

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

## With image base64


**Request**

```
$ curl -X POST -H "Content-Type: application/json" -d '{"img_base64":"<YOUR BASE 64 HERE>","engine":"tesseract"}' http://10.0.2.15:$HTTP_PORT/ocr
```


## The REST API also supports:

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

# Client Libraries

* **Go** [open-ocr-client](https://github.com/tleyden/open-ocr-client)
* **C#** [open-ocr-dotnet](https://github.com/alex-doe/open-ocr-dotnet)

# License

OpenOCR is Open Source and available under the Apache 2 License.
