resource "kong_service" "service" {

  name            = "my_service"
  retries         = 5
  protocol        = "http"
  host            = "example.com"
  port            = "80"
  path            = "/some_api"
  connect_timeout = 60000
  write_timeout   = 60000
  read_timeout    = 60000
  tags            = ["user-level", "low-priority"]

  client_certificate = null
  tls_verify         = false
  tls_verify_depth   = null
  ca_certificates    = null

  //  Works only with HTTPS protocol
  //  client_certificate = kong_certificate.certificate.id
  //  tls_verify = false
  //  tls_verify_depth = null
  //  ca_certificates = ["4e3ad2e4-0bc4-4638-8e34-c84a417ba39b", "51e77dc2-8f3e-4afa-9d0e-0e3bbbcfd515"]

}
