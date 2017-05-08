# Terraform provider for Kong

Uses [Terraform](http://www.terraform.io) to configure APIs in [Kong](http://www.getkong.org). It fully supports creating APIs and consumers, but plugins and credentials are not complete (most plugins will work though).

## Example usage

```Terraform
provider "kong" {
    address = "http://192.168.99.100:8001"
}

resource "kong_api" "api" {
    name               = "api"
    upstream_url       = "http://api.local"
    uris       = "/api"
    strip_uris = false
}

resource "kong_consumer" "consumer" {
    username  = "user"
    custom_id = "123456"
}

resource "kong_api_plugin" "basic_auth" {
    api = "${kong_api.api.id}"
    name = "basic-auth"
}

resource "kong_api_plugin" "jwt" {
    api = "${kong_api.api.id}"
    name = "jwt"
}

resource "kong_api_plugin" "rate_limiting" {
    api  = "${kong_api.api.id}"
    name = "rate-limiting"

    config {
        minute = "100"
    }
}

resource "kong_consumer_basic_auth_credential" "basic_auth_credential" {
    consumer = "${kong_consumer.consumer.id}"
    username = "user123"
    password = "password"
}

resource "kong_consumer_jwt_credential" "jwt_credential" {
    consumer = "${kong_consumer.consumer.id}"
    secret   = "secret"
}
```
