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

######## Start a new stage from scratch #######
FROM golang:latest  

WORKDIR /root/

ARG LOG_DIR=/root/logs
ARG CHAIN_DIR=/root/chain
ARG GENESIS_DIR=/root/genesis

# Create Log Directory
RUN mkdir -p ${LOG_DIR}
RUN mkdir -p ${CHAIN_DIR}
RUN mkdir -p ${GENESIS_DIR}

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/build/Linux-AMD64/commercionetworkd .
COPY container_exec.sh .
RUN chmod +x container_exec.sh

# Declare volumes to mount
VOLUME [${LOG_DIR}]
VOLUME [${CHAIN_DIR}]
VOLUME [${GENESIS_DIR}]

# Command to run the executable
CMD ./container_exec.sh
