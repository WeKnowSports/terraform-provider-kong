#!/usr/bin/env bash
TF_LOG='DEBUG'
echo ${1}

export GOPATH=$HOME/gocode
mkdir -p $GOPATH/github.com/localz/terraform-provider-kong
cp -R $PWD $GOPATH/src/github.com/localz/terraform-provider-kong/

go build -o terraform/tests/terraform-provider-kong
cd terraform/tests
terraform ${1}
cd ../..

rm -fr $GOPATH/github.com/localz/terraform-provider-kong
rm -f ./terraform/tests/terraform-provider-kong