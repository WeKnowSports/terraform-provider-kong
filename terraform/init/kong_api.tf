//resource "kong_api" "admin" {
//  name               = "kong-admin"
//  upstream_url       = "http://localhost:8001"
//  hosts              = "localhost, 127.0.0.1"
//  strip_uri          = true
//  retries            = 5
//  preserve_host      = true
//  https_only         = false
//  http_if_terminated = true
//}
//
//resource "kong_consumer" "admin" {
//  username  = "username"
//}
//
//resource "kong_api_plugin" "admin_basic_auth" {
//  api = "${kong_api.admin.id}"
//  name = "basic-auth"
//}
//
//resource "kong_consumer_basic_auth_credential" "basic_auth_credential" {
//  consumer = "${kong_consumer.admin.id}"
//  username = "username"
//  password = "password"
//}

resource "kong_upstream" "upstream_new" {
  name = "test_upstream_new"
  healthchecks = {
        active = {
            healthy = {
                http_statuses = [ 200 ],
                successes = 5
           }
        },
        passive = {
            healthy = {
                http_statuses = [ 200 ],
                successes = 5
            },
            unhealthy = {
                http_failures = 10,
                http_statuses = [ 404, 500, 503 ],
                tcp_failures = 10,
                timeouts = 5
            }
        }
    }
}

resource "kong_target" "target_new" {
  upstream = "${kong_upstream.upstream_new.id}"
  target = "www.google.com"
}
