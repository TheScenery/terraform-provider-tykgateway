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
  
}
