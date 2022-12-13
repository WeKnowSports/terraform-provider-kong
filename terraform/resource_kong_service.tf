// Service
resource "kong_service" "service" {

  name            = "my_service"
  protocol        = "https"
  host            = "httpbin.org"
  port            = "443"
  path            = "/get"
  retries         = 5
  connect_timeout = 60000
  write_timeout   = 60000
  read_timeout    = 60000
}
