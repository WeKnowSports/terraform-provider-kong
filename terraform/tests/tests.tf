/**
 * Create Test API
 */
resource "kong_api" "api" {
  name               = "test"
  upstream_url       = "http://requestb.in"
  hosts              = ""
  methods            = ""
  uris               = "/localz,/test"
  strip_uri          = true
  methods            = "GET"
}

/**
 * Create Test Consumer
 */
resource "kong_consumer" "consumer" {
  username  = "user"
  custom_id = "123456"
}

/************************************************/

/**
 * Enable Basic Auth and Create Basic Auth Credentials
 */
resource "kong_api_plugin" "basic_auth" {
    api = "${kong_api.api.id}"
    name = "basic-auth"
}

/**
 * Create Basic Auth Credentials for Consumer
 */
resource "kong_consumer_basic_auth_credential" "basic_auth_credential" {
    consumer = "${kong_consumer.consumer.id}"
    username = "user123"
    password = "password"
}

/************************************************/

/**
 * Enable JWT and Create JWT Auth Credentials
 */
resource "kong_api_plugin" "jwt" {
    api = "${kong_api.api.id}"
    name = "jwt"
}

/**
 * Create JWT Crednetials for Consumer
 */
resource "kong_consumer_jwt_credential" "jwt_auth_credential" {
  consumer = "${kong_consumer.consumer.id}"
  secret = "secret"
}

/************************************************/

/**
 * Enable Key Authentication
 */
resource "kong_api_plugin" "key_authentication" {
  api = "${kong_api.api.id}"
  name = "key-auth"

  config {
    key_names = "apikey"
  }
}

resource "kong_consumer_key_auth_credential" "key_auth_credential" {
  consumer = "${kong_consumer.consumer.id}"
  key = "test"
}

/************************************************/

/**
 * Enable Rate Limiting
 */
resource "kong_api_plugin" "rate_limiting" {
  api = "${kong_api.api.id}"
  name = "rate-limiting"

  config = {
    day = "1000"
  }
}