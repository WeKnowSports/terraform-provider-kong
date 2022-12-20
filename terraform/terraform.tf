terraform {
  required_version = ">= 0.13"
  required_providers {
    kong = {
      source  = "local/paas/kong"
      version = "= 1.0.1"
    }
  }
}
