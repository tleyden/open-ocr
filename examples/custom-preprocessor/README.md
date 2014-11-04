
Example of using a custom preprocessor.

## Build docker image

```
$ docker build -t local/custom-preprocessor .
```

## Run docker image

```
$ export AMQP_URI=amqp://admin:${RABBITMQ_PASS}@${RABBITMQ_HOST}/
$ docker run -d local/custom-preprocessor custom-preprocessor -amqp_uri "${AMQP_URI}" 
```

## Test

```
curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://bit.ly/ocrimage-swt","engine":"tesseract", "preprocessors":["custom-preprocessor"]}' http://${RABBITMQ_HOST}:${HTTP_PORT}/ocr
```