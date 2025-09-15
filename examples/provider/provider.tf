# Copyright (c) HashiCorp, Inc.

# Configutation-based provider

provider "tykgateway" {
  gateway_url = "https://host/tyk-gateway"
  api_key     = "your-api-key"
}