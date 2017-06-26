#!/usr/bin/env bash
TF_LOG='DEBUG'
echo ${1}

rm terraform/init/terraform-provider-kong
go build -o terraform/init/terraform-provider-kong

cd terraform/init
terraform ${1}
cd ../..