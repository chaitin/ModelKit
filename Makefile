swag:
	swag fmt -d .&& swag init --pd -g main.go -ot "json"

lint:
	go mod tidy && gofmt -w . && golangci-lint run

dev:
	unset HTTP_PROXY HTTPS_PROXY ALL_PROXY http_proxy https_proxy all_proxy; go run test/backend/main.go & cd ui/ModelModal && pnpm install && pnpm build && cd ../../test/ui_example && pnpm install && pnpm dev

