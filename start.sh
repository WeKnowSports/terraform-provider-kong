#!/usr/bin/env bash
TF_LOG='DEBUG'

rm tests/terraform-provider-kong
go build -o tests/terraform-provider-kong
cd tests
terraform ${1}
cd ..