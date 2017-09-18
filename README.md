# terraform-provider-kong (rapid7 fork)
**NOTE**: This repository was forked from https://github.com/WeKnowSports/terraform-provider-kong.

This fork's master branch updates main.go to point at this fork's kong package for use with `go get`/`go install`.

The master branch includes several patches contributed (or in review) from the rapid7 fork back to the upstream project.

## Installation
Install terraform and terraform-provider-kong to your `GOPATH` (e.g. under your terraform project):

```sh
export GOPATH=$(pwd)/vendor
export GOBIN=$(pwd)/vendor/bin
export PATH=$(pwd)/vendor/bin

go get github.com/hashicorp/terraform
go install github.com/hashicorp/terraform

go get github.com/rapid7/terraform-provider-kong
go install github.com/rapid7/terraform-provider-kong

terraform init
terraform plan
```

## Development
1. Use the [installation steps](#installation) above.
2. Update the terraform provider plugin under `${GOPATH}/src/github.com/rapid7/terraform-provider-kong`.
3. Use `go build` in the provider plugin directory.
4. Use `go install` in your terraform project.
5. Use `terraform init` in your terraform project.
6. Use the updated provider.

---

# Terraform provider for Kong
Uses [Terraform](http://www.terraform.io) to configure APIs in [Kong](http://www.getkong.org). It fully supports creating APIs and consumers, but plugins and credentials are not complete (most plugins will work though).

```
go build -o tests/terraform-provider-kong
```

## Compile and terraform plan / apply

### Start kong


```Shell
docker-compose up -d
```

## Run plan
```Shell
./start plan
```

## Run apply
```Shell
./start apply
```

## Example usage

Please refer to terraform/tests
