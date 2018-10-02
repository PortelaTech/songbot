#!/bin/bash


rm -rf /tmp/ebversion
mkdir -p /tmp/ebversion

cat > /tmp/ebversion/Dockerrun.aws.json << EOF
{
  "AWSEBDockerrunVersion": "1",
  "Image": {
    "Name": "${DOCKER_USER}/${APP_NAME}:${VERSION_NAME}",
    "Update": "true"
  },
  "Ports": [
    {
      "ContainerPort": 8080
    }
  ]
}
EOF

mkdir .ebextensions 2>/dev/null
cp -r .ebextensions /tmp/ebversion
cat > /tmp/ebversion/.ebextensions/09_env.config << EOF
option_settings:
  aws:elasticbeanstalk:application:environment:
    AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
    AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}
    VERSION_NAME: ${VERSION_NAME}
    CIRCLE_BRANCH: ${CIRCLE_BRANCH}
    CIRCLE_SHA1: ${CIRCLE_SHA1}
    DB_USER: ${DB_USER}
    DB_NAME: ${DB_NAME}
    DB_PASS: ${DB_PASS}
    DB_HOST: ${DB_HOST}
    DB_PORT: ${DB_PORT}
EOF

cd /tmp/ebversion

zip -v -r ${VERSION_NAME}.zip Dockerrun.aws.json .ebextensions

aws s3 cp ${VERSION_NAME}.zip s3://${S3_BUCKET} \
    --region ${AWS_REGION}

aws elasticbeanstalk delete-application-version \
    --application-name ${APP_NAME} \
    --version-label ${VERSION_NAME} \
    --region ${AWS_REGION}


aws elasticbeanstalk create-application-version \
    --application-name ${APP_NAME} \
    --version-label ${VERSION_NAME} \
    --source-bundle S3Bucket="${S3_BUCKET}",S3Key="${VERSION_NAME}.zip" \
    --region ${AWS_REGION}

cd -

echo
echo "Deploy to docker"
echo "================"
docker login -u${DOCKER_USER} -p${DOCKER_PASS}
docker tag ${DOCKER_USER}/${APP_NAME}:latest ${DOCKER_USER}/${APP_NAME}:${VERSION_NAME}
docker push ${DOCKER_USER}/${APP_NAME}:${VERSION_NAME}

echo
echo "Update EB environment"
echo "====================="
aws elasticbeanstalk update-environment \
    --application-name ${APP_NAME} \
    --environment-name ${ENV_NAME} \
    --version-label ${VERSION_NAME} \
    --region ${AWS_REGION}
