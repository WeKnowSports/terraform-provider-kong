// HS256 
resource "kong_consumer_jwt_credential" "jwt_auth_credential" {

  consumer  = kong_consumer.consumer.id
  key       = var.consumer_custom_id
  secret    = "HCS8Ap6G9cK37H1C696RHjhoZ8ZupF7z"
  algorithm = "HS256"
  tags      = ["user-level", "low-priority"]
}


// RS256
resource "kong_consumer" "consumer_jwt_rs256" {
  username  = "user_jwt_rs256"
  custom_id = "123456_jwt_rs256"
  tags      = ["user-level", "low-priority"]
}

resource "kong_consumer_jwt_credential" "jwt_auth_asymmetric_credential" {

  consumer       = kong_consumer.consumer_jwt_rs256.id
  key            = kong_consumer.consumer_jwt_rs256.custom_id
  rsa_public_key = file("./data/pubkey.pem")
  algorithm      = "RS256"
  tags           = ["user-level", "low-priority"]
}
