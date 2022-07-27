# VERSION
# If no version is specified as a parameter of make, the last git hash
# value is taken.
# VERSION?=$(shell git describe --abbrev=0)+hash.$(shell git rev-parse --short HEAD)
VERSION?=0.0.0+hash.$(shell git rev-parse --short HEAD)

# EXECUTABLE
# Executable binary which should be compiled for different architecures
EXECUTABLE:=dhdu

# LINUX_EXECUTABLES AND TARGETS
LINUX_EXECUTABLES:=\
	linux/amd64/${EXECUTABLE} \
	linux/arm/5/${EXECUTABLE} \
	linux/arm/7/${EXECUTABLE}

LINUX_EXECUTABLE_TARGETS:=${LINUX_EXECUTABLES:%=bin/%}

# UNIX_EXECUTABLES AND TARGETS
# Define all executables for different architectures and operation systems
UNIX_EXECUTABLES:=\
	${LINUX_EXECUTABLES}

UNIX_EXECUTABLE_TARGETS:=\
	${LINUX_EXECUTABLE_TARGETS}

# EXECUTABLE_TARGETS
# Include all UNIX and Windows targets.
EXECUTABLES:=\
	${UNIX_EXECUTABLES}

EXECUTABLE_TARGETS:=\
	${UNIX_EXECUTABLE_TARGETS}

# CONTAINER_RUNTIME
# The CONTAINER_RUNTIME variable will be used to specified the path to a
# container runtime. This is needed to start and run a container images.
CONTAINER_RUNTIME?=$(shell which docker)

# BUILD_IMAGE
# Definition of the container build image, in which the Binary are compiled from
# source code
BUILD_IMAGE_REGISTRY:=docker.io
BUILD_IMAGE_NAMESPACE:=volkerraschek
BUILD_IMAGE_NAME:=build-image
BUILD_IMAGE_VERSION:=latest
BUILD_IMAGE_FULL=${BUILD_IMAGE_REGISTRY}/${BUILD_IMAGE_NAMESPACE}/${BUILD_IMAGE_NAME}:${BUILD_IMAGE_VERSION:v%=%}
BUILD_IMAGE_SHORT=${BUILD_IMAGE_NAMESPACE}/${BUILD_IMAGE_NAME}:${BUILD_IMAGE_VERSION:v%=%}

# BASE_IMAGE
# Definition of the base container image
BASE_IMAGE_REGISTRY:=docker.io
BASE_IMAGE_NAMESPACE:=library
BASE_IMAGE_NAME:=alpine
BASE_IMAGE_VERSION:=3.11.2
BASE_IMAGE_FULL=${BASE_IMAGE_REGISTRY}/${BASE_IMAGE_NAMESPACE}/${BASE_IMAGE_NAME}:${BASE_IMAGE_VERSION:v%=%}
BASE_IMAGE_SHORT=${BASE_IMAGE_NAMESPACE}/${BASE_IMAGE_NAME}:${BASE_IMAGE_VERSION:v%=%}

# CONTAINER_IMAGE
# Definition of the container image
CONTAINER_IMAGE_REGISTRY:=docker.io
CONTAINER_IMAGE_NAMESPACE:=volkerraschek
CONTAINER_IMAGE_NAME:=${EXECUTABLE}
CONTAINER_IMAGE_VERSION:=latest
CONTAINER_IMAGE_FULL=${CONTAINER_IMAGE_REGISTRY}/${CONTAINER_IMAGE_NAMESPACE}/${CONTAINER_IMAGE_NAME}:${CONTAINER_IMAGE_VERSION:v%=%}
CONTAINER_IMAGE_SHORT=${CONTAINER_IMAGE_NAMESPACE}/${CONTAINER_IMAGE_NAME}:${CONTAINER_IMAGE_VERSION:v%=%}

README_FILE:=README.md

# BINARIES
# ==============================================================================
PHONY:=all

${EXECUTABLE}: bin/tmp/${EXECUTABLE}

all: ${EXECUTABLE_TARGETS}

bin/linux/amd64/${EXECUTABLE}: bindata
	CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  GOPROXY=${GOPROXY} \
  GOPRIVATE=${GOPRIVATE} \
	go build -ldflags "-X main.version=${VERSION}" -o ${@}

bin/linux/arm/5/${EXECUTABLE}: bindata
	CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=arm \
  GOARM=5 \
  GOPROXY=${GOPROXY} \
  GOPRIVATE=${GOPRIVATE} \
	go build -ldflags "-X main.version=${VERSION}" -o ${@}

bin/linux/arm/7/${EXECUTABLE}: bindata
	CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=arm \
  GOARM=5 \
  GOPROXY=${GOPROXY} \
  GOPRIVATE=${GOPRIVATE} \
	go build -ldflags "-X main.version=${VERSION}" -o ${@}

bin/tmp/${EXECUTABLE}: bindata
	GOPROXY=${GOPROXY} \
	GOPRIVATE=${GOPRIVATE} \
	go build -ldflags "-X main.version=${VERSION}" -o ${@}

# BINDATA
# ==============================================================================
BINDATA_TARGETS := \
	pkg/hub/bindata.go

PHONY+=bindata
bindata: ${BINDATA_TARGETS}

pkg/hub/bindata.go:
	go-bindata -pkg hub -o ${@} README.md

# TEST
# ==============================================================================
PHONY+=test
test: clean bin/tmp/${EXECUTABLE}
	REGISTRY_USER=${REGISTRY_USER} \
	REGISTRY_PASSWORD=${REGISTRY_PASSWORD} \
	REGISTRY_NAMESPACE=${REGISTRY_NAMESPACE} \
	CONTAINER_IMAGE_NAME=${CONTAINER_IMAGE_NAME} \
	README_FILE=${README_FILE} \
		go test -v ./pkg/...

# CLEAN
# ==============================================================================
PHONY+=clean
clean:
	rm --force ${EXECUTABLE} || true
	rm --force --recursive bin || true
	rm --force --recursive ${BINDATA_TARGETS} || true

# CONTAINER IMAGE STEPS
# ==============================================================================
PHONY+=container-image/build/amd64
container-image/build/amd64:
	${CONTAINER_RUNTIME} build \
		--build-arg BASE_IMAGE=${BASE_IMAGE_FULL} \
		--build-arg BUILD_IMAGE=${BUILD_IMAGE_FULL} \
		--build-arg EXECUTABLE=${EXECUTABLE} \
		--build-arg EXECUTABLE_TARGET=bin/linux/amd64/${EXECUTABLE} \
		--build-arg GOPROXY \
		--build-arg GOPRIVATE \
		--build-arg VERSION=${VERSION} \
		--file Dockerfile \
		--no-cache \
		--tag ${CONTAINER_IMAGE_FULL} \
		--tag ${CONTAINER_IMAGE_SHORT} \
		.

PHONY+=container-image/push/amd64
container-image/push/amd64: container-image/build/amd64
	${CONTAINER_RUNTIME} login ${CONTAINER_IMAGE_REGISTRY} \
		--username ${CONTAINER_IMAGE_REGISTRY_USER} \
		--password ${CONTAINER_IMAGE_REGISTRY_PASSWORD}
	${CONTAINER_RUNTIME} push ${CONTAINER_IMAGE_FULL}

# CONTAINER STEPS - BINARY
# ==============================================================================
PHONY+=container-run/all
container-run/all:
	$(MAKE) container-run COMMAND=${@:container-run/%=%}

PHONY+=${UNIX_EXECUTABLE_TARGETS:%=container-run/%}
${UNIX_EXECUTABLE_TARGETS:%=container-run/%}:
	$(MAKE) container-run COMMAND=${@:container-run/%=%}

# CONTAINER STEPS - CLEAN
# ==============================================================================
PHONY+=container-run/clean
container-run/clean:
	$(MAKE) container-run COMMAND=${@:container-run/%=%}

# GENERAL CONTAINER COMMAND
# ==============================================================================
PHONY+=container-run
container-run:
	${CONTAINER_RUNTIME} run \
		--rm \
		--volume ${PWD}:/workspace \
		${BUILD_IMAGE_FULL} \
			make ${COMMAND} \
				VERSION=${VERSION} \
				GOPROXY=${GOPROXY} \
				GOPRIVATE=${GOPRIVATE} \

# PHONY
# ==============================================================================
.PHONY: ${PHONY}