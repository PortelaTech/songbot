#!/usr/bin/env bash

# CircleCI setup if running in vagrant
if [ -d /wrk ] ; then
    export CIRCLE_BRANCH=$(git branch | grep '\*' | awk '{print $2}')
    export CIRCLE_SHA1=`git log --pretty=format:'%H' -n 1`
fi

# AWS
export AWS_ACCESS_KEY_ID=`awk '$1~/aws_access_key_id/ {print $2}' ~/.aws/credentials`
export AWS_SECRET_ACCESS_KEY=`awk '$1~/aws_secret_access_key/ {print $2}' ~/.aws/credentials`
export AWS_REGION=`awk '$1~/default_region:/ {print $2}' ~/.aws/config`
export APP_NAME=`awk '$1~/application_name:/ {print $2}' .elasticbeanstalk/config.yml`
export ENV_NAME=`perl -ne 'if(m/$ENV{CIRCLE_BRANCH}:/){ $m = 1 } elsif($m && /environment:\s+(\S+)/){ print "$1\n"; exit 0; }' .elasticbeanstalk/config.yml`
export SHORT_HASH=$(echo ${CIRCLE_SHA1} | cut -c -7)
export VERSION_NAME=${CIRCLE_BRANCH}-${SHORT_HASH}
export S3_BUCKET=${APP_NAME}
# Docker
export DOCKER_USER=portela
export DOCKER_EMAIL=portela.tech@gmail.com
export DOCKER_PASS=$SECRET
# Database
export DB_NAME=songs
export DB_USER=ptuser
export DB_PASS=$SECRET
export DB_PORT=5432
export DB_HOST=localhost

# Telegram
export TELEGRAM_BOTID="660697313"
export TELEGRAM_APIKEY=$SECRET

