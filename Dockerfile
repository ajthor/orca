# ----------------------------------------
# Base
#
FROM scratch AS base
ADD rootfs.tar.gz /

# ----------------------------------------
# Golang
#
FROM base AS golang

ENV GOLANG_VERSION=1.8
ENV GOLANG_DOWNLOAD_URL=https://storage.googleapis.com/golang

RUN apk add --no-cache \
    ca-certificates \
 && apk add --no-cache --virtual .build-deps \
    bash \
    curl \
    gcc \
    go \
    musl-dev \
    openssl \
 && export GOROOT_BOOTSTRAP="$(go env GOROOT)" \
 && curl -SLO "$GOLANG_DOWNLOAD_URL/go$GOLANG_VERSION.src.tar.gz" \
 && echo "406865f587b44be7092f206d73fc1de252600b79b3cacc587b74b5ef5c623596  go1.8.src.tar.gz" | sha256sum -c - \
 && mkdir -p "/usr/local/go" \
 && tar -xzf "go$GOLANG_VERSION.src.tar.gz" -C "/usr/local/go" --strip-components=1 \
 && rm "go$GOLANG_VERSION.src.tar.gz" \
 && cd /usr/local/go/src && ./make.bash \
 && mkdir -p "/go/src" "/go/bin" && chmod -R 777 "/go" \
 && apk del .build-deps

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $GOPATH

# ----------------------------------------
# Builder
#
FROM golang AS builder

# Set environment variables.
ENV DISTRIBUTION_DIR=/go/src/github.com/gorobot-library/orca

# Install dependencies.
RUN apk add --no-cache \
    ca-certificates \
 && apk add --no-cache --virtual .build-deps \
    gcc \
    g++ \
    git \
    glide \
    make \
    mercurial

# Set the working directory to the distribution directory.
WORKDIR $DISTRIBUTION_DIR

# Clone the git repo into the container.
RUN git clone https://github.com/gorobot-library/orca .

# Get Go dependencies.
RUN glide install

# Build the binary.
RUN go generate ./initialize \
 && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orca ./cmd/orca

# ----------------------------------------
# Orca
#
# Creates the orca image using the binary from the builder.
FROM scratch

# Copy the ca-certificates files.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary.
COPY --from=builder /go/src/github.com/gorobot-library/orca/orca /

# Run the binary.
ENTRYPOINT ["/orca"]
