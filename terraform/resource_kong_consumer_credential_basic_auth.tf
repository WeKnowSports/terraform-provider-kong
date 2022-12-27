resource "kong_consumer_basic_auth_credential" "basic_auth_credential" {
  consumer = kong_consumer.consumer.id
  username = "user123"
  password = "password"
  tags     = ["user-level", "low-priority"]
}
