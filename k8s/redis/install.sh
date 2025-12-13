#/bin/bash

PASSWORD=$1

if [ $# != 1 ]; then
  echo "You should pass the password: $0 {password}"
  exit 1
fi

TEMP_VALUES=./tmp-values.yaml
ORIG_VALUES=./values.yaml

sed "s|{PASSWORD}|${PASSWORD}|g" "$ORIG_VALUES" > "$TEMP_VALUES"

helm upgrade redis -n redis -f $TEMP_VALUES --install --create-namespace bitnami/redis --version 23.2.12

rm "$TEMP_VALUES"