OWNER=bimoyong
ALIAS=file
TYPE=srv
IMAGE_NAME=${OWNER}/${ALIAS}-${TYPE}
COMMIT=${shell git rev-parse --short HEAD}
TAG=${shell git describe --abbrev=0 --tags --always --match "v*"}
IMAGE_TAG=${TAG}-${COMMIT}

SERVER_NAME=go.${TYPE}.${ALIAS}

all: build

run:
	# MICRO_REGISTRY=consul \
	# MICRO_REGISTRY_ADDRESS=localhost:8500 \
	# MICRO_BROKER=kafka \
	# MICRO_BROKER_ADDRESS=localhost:9092 \

	MICRO_SERVER_VERSION=latest \
	MICRO_SERVER_NAME=${SERVER_NAME} \
	go run *.go

vendor:
	go mod vendor

proto:
	# protoc --proto_path=${GOPATH}/pkg/mod:. --micro_out=. --go_out=. proto/file/file.proto

build: proto
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/app *.go;
	chmod +x ./bin/app;

test:
	go test -v ./... -cover

docker:
	docker build \
		--build-arg NAME=${SERVER_NAME} \
		--build-arg VER=${IMAGE_TAG} \
		--build-arg GIT_AUTH=${GIT_AUTH} \
		--build-arg GOPRIVATE=github.com/bimoyong/* \
		--tag ${IMAGE_NAME}:${IMAGE_TAG} \
		.

	docker tag \
		${IMAGE_NAME}:${IMAGE_TAG} \
		${IMAGE_NAME}:latest

	docker push ${IMAGE_NAME}:${IMAGE_TAG}
	docker push ${IMAGE_NAME}:latest

docker_multiarch:
	docker buildx build \
		--platform linux/arm,linux/arm64,linux/amd64 \
		--build-arg NAME=${SERVER_NAME} \
		--build-arg VER=${IMAGE_TAG} \
		--build-arg GIT_AUTH=${GIT_AUTH} \
		--build-arg GOPRIVATE=github.com/bimoyong/* \
		--tag ${IMAGE_NAME}:${IMAGE_TAG} \
		--tag ${IMAGE_NAME}:latest \
		--push \
		.

.PHONY: run vendor proto build test docker docker_multiarch
