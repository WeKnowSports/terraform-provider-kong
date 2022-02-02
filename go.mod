module github.com/WeKnowSports/terraform-provider-kong

go 1.16

replace github.com/WeKnowSports/terraform-provider-kong/kong => ./kong

require (
	github.com/dghubble/sling v1.2.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
)
