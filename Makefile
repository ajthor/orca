.PHONY: all

ORCA_IMAGE := "gorobot/orca"

ALPINE_ARCH := armhf
ALPINE_DIST := v3.5
ALPINE_VERSION := 3.5.0
ALPINE_MIRROR := http://nl.alpinelinux.org/alpine
ALPINE_FILENAME := alpine-minirootfs-$(ALPINE_VERSION)-$(ALPINE_ARCH).tar.gz
ALPINE_DOWNLOAD_URL := $(ALPINE_MIRROR)/$(ALPINE_DIST)/releases/$(ALPINE_ARCH)/$(ALPINE_FILENAME)

DOCKER_DEV_OPTS := --rm -it -v "$PWD:/go/src/github.com/gorobot-library/orca"

DOCKER_RUN_ORCA_OPTS := --rm -it -v "/var/run/docker.sock:/var/run/docker.sock"
DOCKER_RUN_ORCA := docker run $(DOCKER_RUN_ORCA_OPTS) $(ORCA_IMAGE)

default: build

all: build

alpine:
	if [ ! -e "rootfs.tar.gz" ]; then \
		curl -sSL $(ALPINE_DOWNLOAD_URL) -o rootfs.tar.gz ; \
	fi

build: alpine
	docker build -t $(ORCA_IMAGE) .

dev:
	docker run $(DOCKER_DEV_OPTS) base/golang:1.8 sh

clean:
	rm -r $(TEMPDIR)

shell: build
	$(DOCKER_RUN_ORCA) sh
