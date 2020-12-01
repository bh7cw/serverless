# serverless
This repository is used to log sns events in cloud watch.
It's running on aws lambda.

Secret:
- 'cicd_lambda' user credentials in `dev` or `prod`
- bucketname to store serverless application

Reference:
https://docs.aws.amazon.com/lambda/latest/dg/golang-logging.html

Instructions:
1. Prerequisites for building and deploying your application locally
    - Install [Golang](https://golang.org/dl/)
    - Place the codebase in `GOPATH/src/`
    - Install the dependencies listed in go.mod(optional)

2. Build and Deploy instructions for web application
    - Build: `go build`
    - Build for ubuntu: `env GOOS=linux GOARCH=amd64 go build`
    - Test: `go test ./...`
    - Run: `go run main.go`