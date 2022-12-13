resource "kong_consumer_acl_group" "acl_group" {

  consumer = kong_consumer.consumer.id
  group    = var.consumer_username
}
