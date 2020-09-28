OWNER=bimoyong
ALIAS=file
TYPE=srv
IMAGE_NAME=${OWNER}/${ALIAS}-${TYPE}
GIT_COMMIT=${shell git rev-parse --short HEAD}
GIT_TAG=${shell git describe --abbrev=0 --tags --always --match "v*"}
IMAGE_TAG=${GIT_TAG}-${GIT_COMMIT}

SERVER_NAME=go.${TYPE}.${ALIAS}

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

ARCHS := arm arm64 amd64
build: proto
	for item in ${ARCHS}; do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$item go build -o ./bin/$$item *.go; \
		chmod +x ./bin/$$item; \
	done

test:
	go test -v ./... -cover

docker:
	docker build --build-arg NAME=${SERVER_NAME} -t ${IMAGE_NAME}:${IMAGE_TAG} .
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${IMAGE_NAME}:latest
	docker push ${IMAGE_NAME}:${IMAGE_TAG}
	docker push ${IMAGE_NAME}:latest

docker_multiarch:
	docker buildx build \
		--platform linux/arm,linux/arm64,linux/amd64 \
		--build-arg NAME=${SERVER_NAME} \
		--build-arg VER=${IMAGE_TAG} \
		--tag ${IMAGE_NAME}:${IMAGE_TAG} \
		--tag ${IMAGE_NAME}:latest \
		--push \
		.

.PHONY: run vendor proto build test docker docker_multiarch
