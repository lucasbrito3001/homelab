#/bin/bash

helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

helm upgrade tempo grafana/tempo -n monitoring -f values.yaml --install --create-namespace
