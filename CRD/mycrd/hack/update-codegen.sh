#!/usr/bin/env bash

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.

bash $(cd $(dirname $0) && pwd)/../vendor/k8s.io/code-generator/generate-groups.sh all \
    mycrd/pkg/client mycrd/pkg/apis mycrdgroup:v1 \
    --output-base $(cd $(dirname $0) && pwd)/../../ \
    --go-header-file $(cd $(dirname $0) && pwd)/boilerplate.go.txt
