#!/usr/bin/env bash
TF_LOG='DEBUG'
echo ${1}

rm init/terraform-provider-kong
go build -o init/terraform-provider-kong

cd init
terraform ${1}
cd ..