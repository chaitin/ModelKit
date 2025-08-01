PLATFORM=linux/amd64
TAG=main
OUTPUT=type=docker,dest=${HOME}/tmp/ModelKit_server.tar
GOCACHE=${HOME}/.cache/go-build
GOMODCACHE=${HOME}/go/pkg/mod
REGISTRY=ModelKit.chaitin.cn/ModelKit
VERSION=dev-${shell git rev-parse HEAD}
BUILD_TIME=${shell date -u +"%Y-%m-%dT%H:%M:%SZ"}
GIT_COMMIT=${shell git rev-parse HEAD}

# make build PLATFORM= TAG= OUTPUT= GOCACHE=
image:
	docker buildx build \
	  -f build/Dockerfile \
	  --build-arg GOCACHE=${GOCACHE} \
	  --build-arg GOMODCACHE=${GOMODCACHE} \
	  --build-arg REPO_COMMIT=$(shell git rev-parse HEAD) \
	  --build-arg VERSION=${VERSION} \
	  --build-arg BUILD_TIME=${BUILD_TIME} \
	  --build-arg GIT_COMMIT=${GIT_COMMIT} \
	  --platform ${PLATFORM} \
	  --tag ${REGISTRY}/backend:${TAG} \
	  --output ${OUTPUT} \
	  .

image-nginx:
	docker buildx build \
	  -f build/Dockerfile.nginx \
	  --platform ${PLATFORM} \
	  --tag ${REGISTRY}/nginx:${TAG} \
	  --output ${OUTPUT} \
	  .

wire:
	wire cmd/server/wire.go cmd/server/main.go

swag:
	swag fmt -d internal && swag init --pd -g cmd/server/main.go -ot "json"

migrate_sql:
	migrate create -ext sql -dir migration -seq ${SEQ}
