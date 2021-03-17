#!/bin/bash

# TODO:
# - allow cli argument passthrough for key name
# - check to see if keyname already exists

# load up variables
source ./config.sh

NOW=$(date '+%Y%m%d%H%M%S')
LOGFILE="${LOGDIR}/create_keypair-${NOW}.log"
echo "Removing Full AWS infrastructure for ${APP_TAG_NAME}:${APP_TAG_VALUE}" | tee ${LOGFILE}

echo "Running create_keypair.sh at ${NOW}" | tee -a ${LOGFILE}

aws ec2 create-key-pair ${PREAMBLE} --key-name ${KEY_NAME} --query 'KeyMaterial' --output text > ${KEY_FILE}
chmod 400 ${KEY_FILE}

echo "Done." | tee -a ${LOGFILE}
exit 0