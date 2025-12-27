#/bin/bash

REDIS_HOST=$1
REDIS_PORT=$2
REDIS_PASS=$3
MONGODB_URI=$4
NAMESPACE=$5

if [ $# != 5 ]; then
  echo "You should pass the password: $0 {redis_host} {redis_port} {redis_pass} {mongodb_uri} {namespace}"
  exit 1
fi

REDIS_HOST_BASE64=$(echo -n "$REDIS_HOST" | base64 -w 0)
REDIS_PORT_BASE64=$(echo -n "$REDIS_PORT" | base64 -w 0)
REDIS_PASS_BASE64=$(echo -n "$REDIS_PASS" | base64 -w 0)
MONGODB_URI_BASE64=$(echo -n "$MONGODB_URI" | base64 -w 0)

DEPLOYMENT_PATH=./manifests/deployment.yaml
SECRET_PATH=./manifests/secret.yaml
SERVICE_PATH=./manifests/service.yaml

TMP_SECRET_PATH=./manifests/tmp-secret.yaml

sed "s|{REDIS_HOST}|${REDIS_HOST_BASE64}|g;s|{REDIS_PORT}|${REDIS_PORT_BASE64}|g;s|{REDIS_PASS}|${REDIS_PASS_BASE64}|g;s|{MONGODB_URI}|${MONGODB_URI_BASE64}|g" "$SECRET_PATH" >> "$TMP_SECRET_PATH"

k create ns $NAMESPACE
k label ns $NAMESPACE istio-injection=enabled

k apply -n $NAMESPACE -f $TMP_SECRET_PATH

rm "$TMP_SECRET_PATH"

k apply -n $NAMESPACE -f $DEPLOYMENT_PATH
k apply -n $NAMESPACE -f $SERVICE_PATH