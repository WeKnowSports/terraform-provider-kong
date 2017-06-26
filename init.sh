#!/usr/bin/env bash
TF_LOG='DEBUG'
echo ${1}

export GOPATH=$HOME/go
mkdir -p $GOPATH/src/github.com/localz/terraform-provider-kong
cp -R $PWD $GOPATH/src/github.com/localz/terraform-provider-kong/
ls

go build -o terraform/init/terraform-provider-kong
cd terraform/init
terraform ${1}
cd ../..

rm -fr $GOPATH/src/github.com/localz/terraform-provider-kong
rm -f ./terraform/init/terraform-provider-kong