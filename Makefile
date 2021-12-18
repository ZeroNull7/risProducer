.PHONY: clean

BINARY=ris-producer
BINARY_DIR=bin
GOARCH=amd64
IMAGE_NAME=risProducer
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

images:
	podman build -t ${IMAGE_NAME}-client:${IMAGE_TAG} -f Dockerfile.client .
	podman build -t ${IMAGE_NAME}-server:${IMAGE_TAG} -f Dockerfile.server .

push: images
	podman tag localhost/${IMAGE_NAME}-server simplefxn/${IMAGE_NAME}-server:${IMAGE_TAG}
	podman tag localhost/${IMAGE_NAME}-client simplefxn/${IMAGE_NAME}-client:${IMAGE_TAG}
	podman push simplefxn/${IMAGE_NAME}-client:${IMAGE_TAG} ${REGISTRY}/${REGISTRY_PATH}/${IMAGE_NAME}-client:${IMAGE_TAG}
	podman push simplefxn/${IMAGE_NAME}-server:${IMAGE_TAG} ${REGISTRY}/${REGISTRY_PATH}/${IMAGE_NAME}-server:${IMAGE_TAG}

deploy:
	kubectl apply -f k8s/server_deployment.yaml
	kubectl apply -f k8s/client_deployment.yaml

undeploy:
	kubectl delete -f k8s/server_deployment.yaml
	kubectl delete -f k8s/client_deployment.yaml

clean:
	${RM} ${BINARY_DIR}/${BINARY}-*
	${RM} ./{demo.log,heap.out}
