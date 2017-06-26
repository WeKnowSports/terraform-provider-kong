#!/usr/bin/env bash
TF_LOG='DEBUG'
echo ${1}

export GOPATH=$HOME/go
mkdir -p $GOPATH/src/github.com/localz/terraform-provider-kong
cp -R ./ $GOPATH/src/github.com/localz/terraform-provider-kong/
ls

go build -o terraform/${1}/terraform-provider-kong
cd terraform/${1}
terraform ${2}
cd ../..

rm -fr $GOPATH/src/github.com/localz/terraform-provider-kong
rm -f ./terraform/init/terraform-provider-kong