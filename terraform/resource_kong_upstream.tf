resource "kong_upstream" "upstream" {

  name                 = "sample_upstream"
  algorithm            = "consistent-hashing"
  slots                = 100
  hash_on              = "header"
  hash_fallback        = "cookie"
  hash_on_header       = "HeaderName"
  hash_fallback_header = "FallbackHeaderName"
  hash_on_cookie       = "CookieName"
  hash_on_cookie_path  = "/path"
  host_header          = "x-host"
  tags                 = ["a", "b"]
  client_certificate   = kong_certificate.certificate.id

  healthchecks {
    active {
      type                     = "http"
      http_path                = "/status"
      timeout                  = 10
      concurrency              = 20
      https_verify_certificate = false
      https_sni                = "some.domain.com"
      healthy {
        successes     = 0
        interval      = 0
        http_statuses = [200, 302]
      }
      unhealthy {
        timeouts      = 0
        interval      = 0
        tcp_failures  = 0
        http_failures = 0
        http_statuses = [429, 404, 500, 501, 502, 503, 504, 505]
      }
    }
    passive {
      type = "https"
      healthy {
        successes     = 0
        http_statuses = [200, 201, 202, 203, 204, 205, 206, 207, 208, 226, 300, 301, 302, 303, 304, 305, 306, 307, 308]

      }
      unhealthy {
        timeouts      = 0
        tcp_failures  = 0
        http_failures = 0
        http_statuses = [429, 500, 503]
      }
    }
  }
}
