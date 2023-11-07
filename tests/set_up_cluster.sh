#!/bin/bash

set -euxo pipefail

cluster_name=${1-epinio}
namespace=${2-epinio}

k3d cluster create $cluster_name -p '8080:80@loadbalancer' -p '8443:443@loadbalancer' --wait

# Ingress controller setup is not necessary, as k3d installs Traefik

helm repo add jetstack https://charts.jetstack.io
helm repo update
helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager  \
    --set installCRDs=true \
    --set extraArgs={--enable-certificate-owner-ref=true} \
    --create-namespace \
    --wait

# Dynamic storage provisioner setup is not necessary either, as k3d installs the `local-path` provisioner

./$(dirname $0)/install_epinio.sh $cluster_name $namespace
