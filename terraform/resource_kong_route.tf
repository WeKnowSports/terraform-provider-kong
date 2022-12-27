resource "kong_route" "route" {

  service = kong_service.service.id

  name      = "my-route"
  protocols = ["http", "https"]
  methods   = ["GET", "POST"]
  hosts     = ["example.com", "foo.test"]
  paths     = ["/foo", "/bar"]

  # header {
  #    name   = "x-my-header"
  #    values = ["foo", "bar"]
  #  }
  #  header {
  #    name   = "x-another-header"
  #    values = ["bla"]
  #  }

  https_redirect_status_code = 426
  regex_priority             = 1
  strip_path                 = true
  path_handling              = "v0"
  preserve_host              = false
  request_buffering          = false
  response_buffering         = false
  tags                       = ["user-level", "low-priority"]

}
