resource "kong_route" "route" {

  service = kong_service.service.id

  #name          = local.route_name
  protocols     = ["http", "https"]
  methods       = []
  paths         = ["/my-path"]
  strip_path    = true
  hosts         = ["localhost"]
  preserve_host = false

  header {
    name   = "X-Custom"
    values = ["hello", "hi"]
  }
}
