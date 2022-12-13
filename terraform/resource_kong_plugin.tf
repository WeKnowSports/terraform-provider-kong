locals {
  plugin_rate_limiting_config = {
    minute = 1
  }
}

// Rate-Limiting plugin with configuration
resource "kong_plugin" "rate_limiting" {

  name    = "rate-limiting"
  service = kong_service.service.id

  config_json = jsonencode(local.plugin_rate_limiting_config)
}


// Prometheus plugin with default configuration
resource "kong_plugin" "prometheus" {
  name = "prometheus"
}
