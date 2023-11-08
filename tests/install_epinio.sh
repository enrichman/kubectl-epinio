#!/bin/bash

set -euxo pipefail

cluster_name=${1-epinio}
namespace=${2-epinio}

external_ip="$(kubectl get svc -n kube-system traefik -o jsonpath={@.status.loadBalancer.ingress})"
external_ip=$(echo $external_ip | grep -Eo --color=never '([0-9]{1,3}.){3}[0-9]{1,3}')

helm repo add epinio https://epinio.github.io/helm-charts
helm repo update
helm upgrade --install epinio epinio/epinio --namespace $namespace --create-namespace \
    --set global.domain=$external_ip.omg.howdoi.website \
    --wait

k3d kubeconfig get $cluster_name
