#!/usr/bin/env bash
TF_LOG='DEBUG'
echo ${1}

go build -ldflags -w -o terraform/${1}/terraform-provider-kong
cd terraform/${1}
terraform ${2}
