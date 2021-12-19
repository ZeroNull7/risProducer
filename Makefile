.PHONY: clean

BINARY=ris-producer
BINARY_DIR=bin
GOARCH=amd64
IMAGE_NAME=ris-producer
IMAGE_TAG=latest
REGISTRY=harbor.simplefxn.com
REGISTRY_PATH=library

all: linux darwin windows

linux: mod
	GOOS=linux GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-linux-${GOARCH} main.go

darwin: mod
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-darwin-${GOARCH} main.go

windows: mod
	GOOS=windows GOARCH=${GOARCH} go build -o ${BINARY_DIR}/${BINARY}-windows-${GOARCH}.exe main.go

mod: 
	go mod tidy

image:
	podman build -t ${IMAGE_NAME}:${IMAGE_TAG} -f Dockerfile .

push: image
	podman tag localhost/${IMAGE_NAME} simplefxn/${IMAGE_NAME}:${IMAGE_TAG}
	podman push simplefxn/${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY}/${REGISTRY_PATH}/${IMAGE_NAME}:${IMAGE_TAG}

deploy:
	kubectl apply -f k8s/deployment.yaml

undeploy:
	kubectl delete -f k8s/deployment.yaml

clean:
	${RM} ${BINARY_DIR}/${BINARY}-*
	${RM} ./{demo.log,heap.out}
