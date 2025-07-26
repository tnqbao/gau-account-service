#!/usr/bin/bash

set -e
source .env
sh apply_envsubst.sh
kubectl --kubeconfig kubeconfig.yaml delete -k ./