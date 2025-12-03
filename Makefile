swag:
	swag fmt -d .&& swag init --pd -g main.go -ot "json"

lint:
	go mod tidy && gofmt -w . && golangci-lint run
