resource "kong_certificate" "certificate" {
  cert     = file("./data/certificate.crt")
  key      = file("./data/certificate.key")
  cert_alt = file("./data/ec_cert.pem")
  key_alt  = file("./data/ec_key.pem")
  tags     = ["user-level", "low-priority"]
}
