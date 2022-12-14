locals {
  plugin_rate_limiting_config = {
    minute = 1
  }
}

// Enable plugin on service
resource "kong_plugin" "rate_limiting_on_service" {

  name      = "rate-limiting"
  service   = kong_service.service.id
  protocols = ["grpc", "grpcs", "http", "https"]

  config_json = jsonencode(local.plugin_rate_limiting_config)
}

// Enable plugin on route
resource "kong_plugin" "rate_limiting_on_route" {

  name      = "rate-limiting"
  route     = kong_route.route.id
  protocols = ["grpc", "grpcs", "http", "https"]

  config_json = jsonencode(local.plugin_rate_limiting_config)
}

// Enable plugin on consumer
resource "kong_plugin" "rate_limiting_on_consumer" {

  name      = "rate-limiting"
  consumer  = kong_consumer.consumer.id
  protocols = ["grpc", "grpcs", "http", "https"]
  enabled   = false
  tags      = ["user-level", "low-priority"]

  config_json = jsonencode(local.plugin_rate_limiting_config)
}


// Prometheus plugin with default configuration
resource "kong_plugin" "prometheus" {
  name      = "prometheus"
  protocols = ["grpc", "grpcs", "http", "https"]
}
