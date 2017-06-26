#!/usr/bin/env bash
TF_LOG='DEBUG'

rm terraform/tests/terraform-provider-kong
go build -o terraform/tests/terraform-provider-kong

cd terraform/tests
terraform ${1}
cd ../..