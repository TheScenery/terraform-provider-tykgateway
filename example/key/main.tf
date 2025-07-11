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
  # alias      = "example-key"
  # org_id     = "default"
  # hashed     = false
  # allowance  = 1000
  # rate       = 10
  # per        = 60
  # quota_max  = 10000
  # quota_remaining = 10000
  # quota_renewal_rate = 3600
  # quota_renews = 0
  # expires    = 0
  # tags       = ["test", "terraform"]

  access_rights = {
    "api_1" = {
      # api_id   = "api_1"
      # api_name = "Example API"
      # versions = ["v1"]
      # allow_urls = [
      #   {
      #     url     = "/test"
      #     methods = ["GET", "POST"]
      #   }
      # ]
      endpoints = [
        {
          # path = "/test"
          methods = [
            {
              # name = "GET"
              limit = {
                per = {value = 60}
                rate = 10
                # smoothing = {
                #   delay     = 1
                #   enabled   = true
                #   step      = 1
                #   threshold = 10
                #   trigger   = 0.5
                # }
              }
            }
          ]
        }
      ]
      # limit = {
      #   max_query_depth = 5
      #   rate            = 10
      #   per             = 60
      #   quota_max       = 10000
      #   quota_remaining = 10000
      #   quota_renewal_rate = 3600
      #   quota_renews    = 0
      #   throttle_interval = 1
      #   throttle_retry_limit = 3
      #   smoothing = {
      #     delay     = 1
      #     enabled   = true
      #     step      = 1
      #     threshold = 10
      #     trigger   = 0.5
      #   }
      # }
    }
  }
}
