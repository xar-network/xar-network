# Simple usage with a mounted data directory:
# > docker build -t zar .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.zard:/root/.zard -v ~/.zarcli:/root/.zarcli zar zard init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.zard:/root/.zard -v ~/.zarcli:/root/.zarcli zar zard start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

# Set working directory for the build
WORKDIR /go/src/github.com/zar-network/zar-network

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make tools && \
    make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/zard /usr/bin/zard
COPY --from=build-env /go/bin/zarcli /usr/bin/zarcli

# Run zard by default, omit entrypoint to ease using container with zarcli
CMD ["zard"]
