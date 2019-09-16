# docker hub description updater

[![Build Status](https://travis-ci.com/volker-raschek/docker-hub-description-updater.svg?branch=master)](https://travis-ci.com/volker-raschek/docker-hub-description-updater)
[![Go Report Card](https://goreportcard.com/badge/github.com/volker-raschek/docker-hub-description-updater)](https://goreportcard.com/report/github.com/volker-raschek/docker-hub-description-updater)
[![GoDoc Reference](https://godoc.org/github.com/volker-raschek/docker-hub-description-updater?status.svg)](http://godoc.org/github.com/volker-raschek/docker-hub-description-updater)
[![Docker Pulls](https://img.shields.io/docker/pulls/volkerraschek/docker-hub-description-updater)](https://hub.docker.com/r/volkerraschek/docker-hub-description-updater)

By specifying the login data for hub.docker.com you can update the short and
long description of a docker repository.

## Usage

Several options are available to update the descriptions. Either based on
Markdown files or as a normal string, which is passed as argument when calling.
The examples below describe two ways, the binary and container based way.

### Example 1: Update full description of the repository with a Markdown file

```bash
dhdu \
  -user=<username> \
  -password=<password> \
  -namespace=<namespace> \
  -repository=<repository> \
  -full-description-file=./README.md
```

```bash
docker run \
  --rm \
  --volume $(pwd):/workspace \
    volkerraschek/dhdu \
      -user=<username> \
      -password=<password> \
      -namespace=<namespace> \
      -repository=<repository> \
      -full-description-file=./README.md
```

### Example 2: Update full description of the repository over an argument

```bash
dhdu -user=<username> \
     -password=<password> \
     -namespace=<namespace> \
     -repository=<repository> \
     -full-description="My awesome description"
```

```bash
docker run \
  --rm \
  --volume $(pwd):/workspace \
    volkerraschek/dhdu \
      -user=<username> \
      -password=<password> \
      -namespace=<namespace> \
      -repository=<repository> \
      -full-description="My awesome description"
```

## Compiling the source code

There are two different ways to compile dhdu from scratch. The easier ways is
to use the pre-defined container image in the Makefile, which has included all
dependancies to compile dhdu. Alternatively, if all dependencies are met,
dhdu can also be compiled without the container image. Both variants are
briefly described.

### Compiling the source code via container image

To compile dhdu via container image it's necessary, that a container runtime
is installed. In the Makefile is predefined docker, but it's can be also used
podman. Execute `make container-run/dhdu` to start the compiling process.

```bash
make container-run/dhdu
```

#### Compiling the source code without container image

Make sure you have installed go >= v1.12. Execute `make dhdu` to compile
dhdu without a container-image. There should be a similar output as when
compiling dhdu via the container image.
