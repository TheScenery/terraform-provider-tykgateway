terraform {
  required_providers {
    tykgateway = {
      source = "github.com/thescenery/tykgateway"
    }
  }
}

provider "tykgateway" {
  gateway_url = "http://192.168.5.119/tyk-gateway"
  api_key     = "foo"
}

resource "tykgateway_key" "key1" {
  allowance   = 1000
  rate         = 1000
  per          = 1
  quota_max    = 10000
  quota_renews = 1688131200
  org_id       = "default"
  access_rights = {
    httpbin-api = {
      api_id   = "httpbin-api"
      api_name = "Httpbin API"
    }
  }
}
