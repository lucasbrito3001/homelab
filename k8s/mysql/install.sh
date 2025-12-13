#/bin/bash

ROOT_PASSWORD=$1
NAMESPACE=$2

if [ $# != 2 ]; then
  echo "You should pass the password: $0 {password} {namespace}"
  exit 1
fi

ROOT_PASSWORD_BASE64=$(echo -n "$ROOT_PASSWORD" | base64 -w 0)

STATEFULSET_PATH=./manifests/statefulset.yaml
SECRET_PATH=./manifests/secret.yaml
SERVICE_PATH=./manifests/service.yaml

TMP_SECRET_PATH=./manifests/tmp-secret.yaml

sed "s|{PASSWORD}|${ROOT_PASSWORD_BASE64}|g" "$SECRET_PATH" > "$TMP_SECRET_PATH"

k create ns $NAMESPACE

k apply -n $NAMESPACE -f $TMP_SECRET_PATH

rm "$TMP_SECRET_PATH"

k apply -n $NAMESPACE -f $STATEFULSET_PATH
k apply -n $NAMESPACE -f $SERVICE_PATH