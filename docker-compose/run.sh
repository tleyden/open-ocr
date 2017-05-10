#!/bin/bash

#===================================================================
# clean docker containers
#===================================================================
cleanContainer()
{
	echo
	echo "Removing all previous instance"
	echo

	DOCKER_CONTAINER=$(docker ps -a -q)

	if [ "$DOCKER_CONTAINER" != "" ]
	then
		echo "Cleaning all docker container"
		docker rm $DOCKER_CONTAINER
	fi
}

#===================================================================
# clean docker images
#===================================================================
cleanImage()
{
	IMAGE_TO_CLEAN=$(docker images | grep "$1")
	if [ "$IMAGE_TO_CLEAN" != "" ]
	then
		echo "Cleaning $IMAGE_TO_CLEAN image"
		docker rmi $IMAGE_TO_CLEAN
	else
        echo "$1 already cleaned"
	fi

}


# first remove all the docker container and docker images related to the project
cleanContainer
cleanImage "tleyden5iwx/open-ocr"
cleanImage "tleyden5iwx/open-ocr-2"
cleanImage "ubuntu"
cleanImage "tleyden5iwx/open-ocr-preprocessor"
cleanImage "tleyden5iwx/stroke-width-transform" 

echo
echo "Which version of the OCR do you want to deploy: "
echo "[1] V1 (using tesseract 3.X): low memory consumption, faster but result less precise"
echo "[2] V2 (using tesseract 4.X): High accuracy but slower and moderate to high memory consumption"
echo

read -p "Choose 1 or 2: " OPEN_OCR_VERSION

OPEN_OCR_INSTANCE_NAME=""

if [ "$OPEN_OCR_VERSION" == 1 ]
then
	echo "Open ocr instance name will be open-ocr-1"
	OPEN_OCR_INSTANCE_NAME="open-ocr"
elif [ "$OPEN_OCR_VERSION" == 2 ]
then
	echo "Open ocr instance name will be open-ocr-2"
	OPEN_OCR_INSTANCE_NAME="open-ocr-2"
else
	echo "ERROR: No correct version specified (please choose between 1 and 2)"
	exit
fi

export OPEN_OCR_INSTANCE=$OPEN_OCR_INSTANCE_NAME

cd docker-compose/

sudo docker-compose up