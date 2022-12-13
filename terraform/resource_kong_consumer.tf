resource "kong_consumer" "consumer" {
  username  = var.consumer_username
  custom_id = var.consumer_custom_id
}
