# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Marco Ruaro <marco.ruaro@gmail.com>"
LABEL creator="Gianguido Sorà <me@gsora.xyz>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN make build-linux

FROM ubuntu:20.04

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get -y install curl jq file tzdata
    
COPY --from=builder /app/build/Linux-AMD64/commercionetworkd .

VOLUME /cnd
WORKDIR /cnd
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM
ENV TZ America/New_York

COPY contrib/localnet/cndnode/wrapper.sh /usr/bin/wrapper.sh
