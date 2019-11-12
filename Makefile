# VERSION
# If no version is specified as a parameter of make, the last git hash
# value is taken.
# VERSION?=$(shell git describe --abbrev=0)+hash.$(shell git rev-parse --short HEAD)
VERSION?=$(shell git describe --abbrev=0)+hash.$(shell git rev-parse --short HEAD)

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

# CONTAINER_RUNTIME / BUILD_IMAGE
# The CONTAINER_RUNTIME variable will be used to specified the path to a
# container runtime. This is needed to start and run a container image defined
# by the BUILD_IMAGE variable. The BUILD_IMAGE container serve as build
# environment to execute the different make steps inside. Therefore, the bulid
# environment requires all necessary dependancies to build this project.
CONTAINER_RUNTIME?=$(shell which docker)
BUILD_IMAGE:=volkerraschek/build-image:latest

# REGISTRY_MIRROR / REGISTRY_NAMESPACE
# The REGISTRY_MIRROR variable contains the name of the registry server to push
# on or pull from container images. The REGISTRY_NAMESPACE defines the Namespace
# where the CONTAINER_RUNTIME will be search for container images or push them
# onto. The most time it's the same as REGISTRY_USER.
REGISTRY_USER:=volkerraschek
REGISTRY_MIRROR=docker.io
REGISTRY_NAMESPACE:=${REGISTRY_USER}

# CONTAINER_IMAGE_VERSION / CONTAINER_IMAGE_NAME / CONTAINER_IMAGE
# Defines the name of the new container to be built using several variables.
BASE_IMAGE=busybox:latest
CONTAINER_IMAGE_NAME=${EXECUTABLE}
CONTAINER_IMAGE_VERSION?=latest
CONTAINER_IMAGE=${REGISTRY_NAMESPACE}/${CONTAINER_IMAGE_NAME}:${CONTAINER_IMAGE_VERSION}

README_FILE:=README.md

# BINARIES
# ==============================================================================
PHONY:=all

${EXECUTABLE}: bin/tmp/${EXECUTABLE}

all: ${EXECUTABLE_TARGETS}

bin/linux/amd64/${EXECUTABLE}: bindata
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o ${@}

bin/linux/arm/5/${EXECUTABLE}: bindata
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o ${@}

bin/linux/arm/7/${EXECUTABLE}: bindata
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o ${@}

bin/tmp/${EXECUTABLE}: bindata
	go build -ldflags "-X main.version=${VERSION}" -o ${@}

# BINDATA
# ==============================================================================
BINDATA_TARGETS:=\
	pkg/hub/bindata_test.go

PHONY+=bindata
bindata: ${BINDATA_TARGETS}

pkg/hub/bindata_test.go:
	go-bindata -pkg hub_test -o ${@} README.md

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
container-image/build:
	${CONTAINER_RUNTIME} build \
		--build-arg BASE_IMAGE=${BASE_IMAGE} \
		--build-arg BUILD_IMAGE=${BUILD_IMAGE} \
		--build-arg EXECUTABLE_TARGET=bin/linux/amd64/${EXECUTABLE} \
		--build-arg GOPROXY=${GOPROXY} \
		--build-arg VERSION=${VERSION} \
		--no-cache \
		--tag ${CONTAINER_IMAGE} \
		--tag ${REGISTRY_MIRROR}/${CONTAINER_IMAGE} \
		.

	if [ -f $(shell which docker) ] && [ "${CONTAINER_RUNTIME}" == "$(shell which podman)" ]; then \
		podman push ${REGISTRY_MIRROR}/${CONTAINER_IMAGE} docker-daemon:${CONTAINER_IMAGE}; \
	fi

container-image/push: container-image/build
	${CONTAINER_RUNTIME} login ${REGISTRY_MIRROR} --username ${REGISTRY_USER} --password ${REGISTRY_PASSWORD}
	${CONTAINER_RUNTIME} push ${REGISTRY_MIRROR}/${CONTAINER_IMAGE}

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
		${BUILD_IMAGE} \
			make ${COMMAND} \
				VERSION=${VERSION}

# PHONY
# ==============================================================================
.PHONY: ${PHONY}