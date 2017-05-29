.PHONY: all

ORCA_IMAGE := "gorobot/orca"
ORCA_DEV_IMAGE := gorobot/orca:dev

ALPINE_ARCH := armhf
ALPINE_DIST := v3.5
ALPINE_VERSION := 3.5.0
ALPINE_MIRROR := http://nl.alpinelinux.org/alpine
ALPINE_FILENAME := alpine-minirootfs-$(ALPINE_VERSION)-$(ALPINE_ARCH).tar.gz
ALPINE_DOWNLOAD_URL := $(ALPINE_MIRROR)/$(ALPINE_DIST)/releases/$(ALPINE_ARCH)/$(ALPINE_FILENAME)

DOCKER_DEV_OPTS := --rm -it -v "$$PWD:/go/src/github.com/gorobot-library/orca" -v "/var/run/docker.sock:/var/run/docker.sock" --name dev

DOCKER_RUN_OPTS := --rm -it -v "/var/run/docker.sock:/var/run/docker.sock"
DOCKER_RUN := docker run $(DOCKER_RUN_ORCA_OPTS) $(ORCA_IMAGE)

default: build

all: build

alpine:
	if [ ! -e "rootfs.tar.gz" ]; then \
		curl -sSL $(ALPINE_DOWNLOAD_URL) -o rootfs.tar.gz ; \
	fi

build: alpine
	docker build -t $(ORCA_IMAGE) .

build-dev:
	docker build -t $(ORCA_DEV_IMAGE) -f Dockerfile.dev .

dev: build-dev
	docker run $(DOCKER_DEV_OPTS) $(ORCA_DEV_IMAGE) sh

clean:
	rm -f rootfs.tar.gz

shell: build
	$(DOCKER_RUN) sh
