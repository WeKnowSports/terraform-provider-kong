resource "kong_target" "target" {
  upstream = kong_upstream.upstream.id
  target   = "google.com:80"
}
