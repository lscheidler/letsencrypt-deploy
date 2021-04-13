module github.com/lscheidler/letsencrypt-deploy

go 1.14

//replace github.com/lscheidler/letsencrypt-lambda => ../letsencrypt-lambda

require (
	github.com/aws/aws-sdk-go v1.30.20
	github.com/lscheidler/letsencrypt-lambda v0.1.0
)
