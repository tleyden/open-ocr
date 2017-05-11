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
# check and clean
#===================================================================
checkAndCleanImage()
{
	IMAGE_TO_CHECK=$(docker images | grep "$1")
	if [ "$IMAGE_TO_CHECK" != "" ]
	then
		echo
		echo "$1 has been found. Do you want to clean it?"
		echo
		read -p "Choose y/n: " CHOICE
		
		if [ $CHOICE == "y" ]
		then
			cleanImage $1
		elif [ $CHOICE == "n" ]
		then 
			echo
			echo "Keep $1"
			echo
		else
			echo
			echo "Wrong choice please retry"
			checkAndCleanImage $1
		fi
	else
        	echo "$1 Not present"
		echo
	fi

}

#===================================================================
# clean docker images
#===================================================================
cleanImage()
{
	echo "Cleaning $1 image"
	docker rmi $1
	echo "$1 has been cleaned"
}


# first remove all the docker container and docker images related to the project
cleanContainer

# put a space as we open-ocr and open-ocr-2 will match on open-ocr grep
checkAndCleanImage "tleyden5iwx/open-ocr "
checkAndCleanImage "tleyden5iwx/open-ocr-2 "
checkAndCleanImage "tleyden5iwx/rabbitmq"
checkAndCleanImage "tleyden5iwx/open-ocr-preprocessor"

echo
echo "Which version of the OCR do you want to deploy: "
echo "[1] V1 (using tesseract 3.X): low memory consumption, faster but result less precise"
echo "[2] V2 (using tesseract 4.X): High accuracy but slower and moderate to high memory consumption"
echo

read -p "Choose 1 or 2: " OPEN_OCR_VERSION

OPEN_OCR_INSTANCE_NAME=""

if [ "$OPEN_OCR_VERSION" == 1 ]
then
	echo "Open ocr instance name will be open-ocr"
	OPEN_OCR_INSTANCE_NAME="open-ocr"
elif [ "$OPEN_OCR_VERSION" == 2 ]
then
	echo "Open ocr instance name will be open-ocr-2"
	OPEN_OCR_INSTANCE_NAME="open-ocr-2"
else
	echo "ERROR: No correct version specified (please choose between 1 and 2)"
	exit
fi

echo
echo "Do you want to start it as deamon?"
echo

read -p "Choose 'y' for yes or anyother for no: " CHOICE_DEAMON

DEAMON_OPTION=""

if [ "$CHOICE_DEAMON" == "y" ]
then
	echo "Use deamon option -d"
	DEAMON_OPTION="-d"
fi

export OPEN_OCR_INSTANCE=$OPEN_OCR_INSTANCE_NAME

sudo -E docker-compose up $DEAMON_OPTION
