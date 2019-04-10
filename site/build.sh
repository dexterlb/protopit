#!/bin/bash

cd "$(dirname "$(readlink -f "${0}")")"

go run 'github.com/DexterLB/protopit/site/build' "${@}"
