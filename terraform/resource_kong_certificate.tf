resource "kong_certificate" "certificate" {
  cert = file("./data/certificate.crt")
  key = file("./data/certificate.key")
}
