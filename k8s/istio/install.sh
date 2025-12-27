#/bin/bash

helm repo add kiali https://kiali.org/helm-charts
helm repo update

helm upgrade --install istio-base istio/base -n istio-system --set defaultRevision=default
helm upgrade --install -n istio-system istiod istio/istiod -f values.yaml
helm upgrade \
    --install \
    --namespace istio-system \
    --set auth.strategy="anonymous" \
    kiali-server \
    kiali/kiali-server

kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.28/samples/addons/prometheus.yaml
