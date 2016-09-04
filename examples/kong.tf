provider "kong" {
    address = "http://localhost:8001/"
}

resource "kong_api" "api" {
    name               = "api"
    request_host       = "api.local"
    request_path       = "api"
    strip_request_path = true
    preserver_host     = false
    upstream_url       = "http://api.local/"
}