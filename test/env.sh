#!/bin/bash

: "${SPLUNK_OPERATOR_IMAGE:=splunk/splunk-operator:latest}"
: "${SPLUNK_ENTERPRISE_IMAGE:=splunk/splunk:latest}"
: "${CLUSTER_PROVIDER:=eks}"
: "${CLUSTER_NAME:=julien-int-test-2}"
: "${NUM_WORKERS:=3}"
: "${NUM_NODES:=2}"
: "${COMMIT_HASH:=}"
#: "${ECR_REGISTRY:=}"
: "${VPC_PUBLIC_SUBNET_STRING:=subnet-0921cea9bcffd7b77,subnet-0dbbc27abdf4a416e,subnet-0dec6ad34b32e791f}" // Public VPC for cluster
: "${VPC_PRIVATE_SUBNET_STRING:=subnet-0c068ca7d9c468e09,subnet-0b9a43cb73e3f9799,subnet-0fa22fc6046b4591f}" // Private VPC for cluster
: "${ECR_REGISTRY:=667741767953.dkr.ecr.us-west-2.amazonaws.com/julien}}"

# Docker registry to use to push the test images to and pull from in the cluster
if [ -z "${PRIVATE_REGISTRY}" ]; then
    case ${CLUSTER_PROVIDER} in
      kind)
        PRIVATE_REGISTRY=localhost:5000
        ;;
      eks)
        if [ -z "${ECR_REGISTRY}" ]; then
          echo "Please define ECR_REGISTRY that specified where images are pushed and pulled from."
          exit 1
        fi
        PRIVATE_REGISTRY="${ECR_REGISTRY}"
        ;;
    esac
fi
