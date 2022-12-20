resource "kong_consumer_key_auth_credential" "key-auth-credential" {

  consumer = kong_consumer.consumer.id
  key      = "KkRLVM9E8xWf7E0M9jIyWv8MXSMdDB2J"
  tags     = ["user-level", "low-priority"]
  ttl      = 3600
}
