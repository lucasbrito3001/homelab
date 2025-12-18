#/bin/bash

ROOT_USER=$1
ROOT_PASS=$2
NAMESPACE=$3

if [ $# != 3 ]; then
  echo "You should pass the password: $0 {username} {password} {namespace}"
  exit 1
fi

ROOT_USER_BASE64=$(echo -n "$ROOT_USER" | base64 -w 0)
ROOT_PASS_BASE64=$(echo -n "$ROOT_PASS" | base64 -w 0)

STATEFULSET_PATH=./manifests/statefulset.yaml
SECRET_PATH=./manifests/secret.yaml
SERVICE_PATH=./manifests/service.yaml

TMP_SECRET_PATH=./manifests/tmp-secret.yaml

sed "s|{USERNAME}|${ROOT_USER_BASE64}|g;s|{PASSWORD}|${ROOT_PASS_BASE64}|g" "$SECRET_PATH" >> "$TMP_SECRET_PATH"

k create ns $NAMESPACE

k apply -n $NAMESPACE -f $TMP_SECRET_PATH

rm "$TMP_SECRET_PATH"

k apply -n $NAMESPACE -f $STATEFULSET_PATH
k apply -n $NAMESPACE -f $SERVICE_PATH