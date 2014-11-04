FROM ubuntu:14.04

ENV GOPATH /opt/go
RUN mkdir -p $GOPATH

# since we are gonna install open-ocr, we need all the tesseract / leptonica deps
RUN apt-get -q -y update && \
    apt-get -q -y install git && \
    apt-get -q -y install golang && \
    apt-get -q -y install gcc && \
    apt-get -q -y install libleptonica-dev && \
    apt-get -q -y install libtesseract3 libtesseract-dev && \
    apt-get -q -y install tesseract-ocr-eng

# get open-ocr
RUN go get -u -v github.com/tleyden/open-ocr

# build open-ocr-preprocessor binary and copy it to /usr/bin
RUN cd $GOPATH/src/github.com/tleyden/open-ocr/examples/custom-preprocessor && go build -v -o custom-preprocessor && cp custom-preprocessor /usr/bin