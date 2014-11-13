
FROM ubuntu:14.04

ENV GOPATH /opt/go

# Get dependencies
RUN apt-get update && apt-get install -y \
  git \
  golang \
  gcc \
  libleptonica-dev \
  libtesseract3 \
  libtesseract-dev \
  tesseract-ocr \
  tesseract-ocr-eng \
  tesseract-ocr-ara \
  tesseract-ocr-bel \
  tesseract-ocr-ben \
  tesseract-ocr-bul \
  tesseract-ocr-ces \
  tesseract-ocr-dan \
  tesseract-ocr-deu \
  tesseract-ocr-ell \
  tesseract-ocr-fin \
  tesseract-ocr-fra \
  tesseract-ocr-heb \
  tesseract-ocr-hin \
  tesseract-ocr-ind \
  tesseract-ocr-isl \
  tesseract-ocr-ita \
  tesseract-ocr-jpn \
  tesseract-ocr-kor \
  tesseract-ocr-nld \
  tesseract-ocr-nor \
  tesseract-ocr-pol \
  tesseract-ocr-por \
  tesseract-ocr-ron \
  tesseract-ocr-rus \
  tesseract-ocr-spa \
  tesseract-ocr-swe \
  tesseract-ocr-tha \
  tesseract-ocr-tur \
  tesseract-ocr-ukr \
  tesseract-ocr-vie \
  tesseract-ocr-chi-sim \
  tesseract-ocr-chi-tra 

# In theory, should be able to set export TESSDATA_PREFIX=/usr/share/tesseract-ocr/, 
# but when I tried I still got error: Error opening data file /usr/local/share/tessdata/eng.traineddata
# Workaround: just copy the data to where it expects
RUN mkdir -p /usr/local/share/tessdata/ && \
    cp -R /usr/share/tesseract-ocr/tessdata/* /usr/local/share/tessdata/

RUN mkdir -p $GOPATH

# go get open-ocr
RUN go get -u -v -t github.com/tleyden/open-ocr

# build open-ocr-httpd binary and copy it to /usr/bin
RUN cd $GOPATH/src/github.com/tleyden/open-ocr/cli-httpd && go build -v -o open-ocr-httpd && cp open-ocr-httpd /usr/bin

# build open-ocr-worker binary and copy it to /usr/bin
RUN cd $GOPATH/src/github.com/tleyden/open-ocr/cli-worker && go build -v -o open-ocr-worker && cp open-ocr-worker /usr/bin

