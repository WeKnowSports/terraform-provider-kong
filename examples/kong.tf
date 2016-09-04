provider "kong" {
    address = "http://192.168.99.100:8001"
}

resource "kong_api" "api0" {
    name               = "test1"
    upstream_url       = "http://api.local"
    request_path       = "/api0"
    strip_request_path = true
}

resource "kong_api" "api1" {
    name               = "test2"
    upstream_url       = "http://api.local"
    request_path       = "/api1"
    strip_request_path = true
}

resource "kong_api" "api2" {
    name               = "test3"
    upstream_url       = "http://api.local"
    request_path       = "/api2"
    strip_request_path = true
}

resource "kong_api" "api3" {
    name               = "test4"
    upstream_url       = "http://api.local"
    request_path       = "/api3"
    strip_request_path = true
}

resource "kong_consumer" "consumer1" {
    username = "consumer1"
}

resource "kong_consumer" "consumer2" {
    username = "consumer2"
}

resource "kong_consumer" "consumer3" {
    username  = "asdff"
    custom_id = "123456"
}