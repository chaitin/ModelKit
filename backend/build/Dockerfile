FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

WORKDIR /src
ENV CGO_ENABLED=0

COPY go.* .
ARG GOMODCACHE
RUN --mount=type=cache,target=${GOMODCACHE} \
    go mod download

ARG TARGETOS TARGETARCH GOCACHE
ARG VERSION
ARG BUILD_TIME
ARG GIT_COMMIT
RUN --mount=type=bind,target=. \
--mount=type=cache,target=${GOMODCACHE} \
--mount=type=cache,target=${GOCACHE} \
GOOS=$TARGETOS GOARCH=$TARGETARCH \
go build \
-ldflags "-w -s -X 'github.com/chaitin/ModelKit/backend/pkg/version.Version=${VERSION}' -X 'github.com/chaitin/ModelKit/backend/pkg/version.BuildTime=${BUILD_TIME}' -X 'github.com/chaitin/ModelKit/backend/pkg/version.GitCommit=${GIT_COMMIT}'" \
-o /out/main \
cmd/server/main.go cmd/server/wire_gen.go 

FROM alpine:3.21 as binary

WORKDIR /app

ADD migration ./migration

COPY --from=builder /out/main /app/main

CMD [ "./main" ]