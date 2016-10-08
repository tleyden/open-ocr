# Python Script Developed over the upload-local-file.sh
# to use the requests library provided in python. The only
# dependency to run this script is requests library.

# Author: Arpit Goyal
# Linkedin: https://in.linkedin.com/in/arpit-goyal-ba599839


import json
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.mime.image import MIMEImage
import requests


class OpenOCRService(object):
    """
    This is the python script that can be incorporated
    in your python applications to connect with the openocr
    services run by open-ocr docker. This is the interpretation
    of the upload-local-file.sh which is a shell script to 
    interact with the service using curl. The service expects
    a multipart/related Content-Type.

    This service just expects the image path in your local machine
    """

    TESSERACT_OCR_SERVICE_URL = 'http://openocr:9292/ocr-file-upload'

    def __init__(self, image_path):
        self.image_path = image_path
        self.related = MIMEMultipart('related')
        self.image_format = self.image_path.split('/')[-1].split('.')[-1]

        self.set_image_payload()
        self.set_data_payload()

    def set_image_payload(self):
    	submission = MIMEImage('image', '{}'.format(self.image_format))
    	submission.set_payload(open(self.image_path, 'rb').read())
    	self.related.attach(submission)

    def set_data_payload(self):
    	data = MIMEText('application', 'json')
    	data.set_payload(json.dumps({"engine": "tesseract", "preprocessors" : ["stroke-width-transform"]}))
    	self.related.attach(data)

    def get_body(self):
    	return self.related.as_string().split('\n\n', 1)[1]

    def get_headers(self):
    	return dict(self.related.items())

    def perform_ocr(self):
    	body = self.get_body()
    	headers = self.get_headers()

    	response = requests.post(self.TESSERACT_OCR_SERVICE_URL, data = body, headers = headers)
    	return response

    def get_ocr_text(self):
    	response = self.perform_ocr()

    	return response.content
