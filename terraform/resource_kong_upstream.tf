resource "kong_upstream" "upstream" {

  name                 = "my-upstream"
  algorithm            = "round-robin"
  slots                = 10000
  hash_on              = "none"
  hash_fallback        = null
  hash_on_header       = null
  hash_fallback_header = null
  hash_on_cookie       = null
  hash_on_cookie_path  = "/"
}
