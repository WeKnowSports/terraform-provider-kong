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
