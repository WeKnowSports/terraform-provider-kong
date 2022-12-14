resource "kong_target" "target" {
  upstream = kong_upstream.upstream.id
  target   = "google.com:80"
  weight   = 100
  tags     = ["user-level", "low-priority"]
}
