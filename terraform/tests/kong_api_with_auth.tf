resource "kong_api" "basic" {
  name               = "basic"
  upstream_url       = "http://www.google.com"
  hosts              = ["yo.com"]
  strip_uri          = true
  retries            = 5
  preserve_host      = true
  https_only         = false
  http_if_terminated = true
}