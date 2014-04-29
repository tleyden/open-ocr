
ocr-worker does the following:

* Looks for ocr-pending-job messages on rabbitmq
* Runs tesseract (or other ocr engine) on the given image url
* Adds an ocr-finished-job to rabbitmq with the decoded text.

## ocr-pending-job

Message Fields:

* uuid
* image_url
* engine ( tesseract | fake )

## ocr-finished-job

Message Fields:

* pending_job_uuid
* result 

