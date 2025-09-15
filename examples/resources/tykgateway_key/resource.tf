# Copyright (c) HashiCorp, Inc.

resource "tykgateway_key" "key" {
  key_config = jsonencode(
    {
      "allowance" : 1000,
      "rate" : 1000,
      "per" : 1,
      "quota_max" : 10000,
      "quota_renews" : 1688131200,
      "org_id" : "default",
      "access_rights" : {
        "httpbin-api" : {
          "api_id" : "httpbin-api",
          "api_name" : "Httpbin API"
        }
      }
  })
}