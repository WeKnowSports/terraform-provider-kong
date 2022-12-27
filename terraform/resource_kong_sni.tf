resource "kong_sni" "sni" {
  name        = "sni"
  certificate = kong_certificate.certificate.id
  tags        = ["user-level", "low-priority"]
}
