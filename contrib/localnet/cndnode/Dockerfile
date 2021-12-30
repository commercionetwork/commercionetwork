# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM ubuntu:20.04
RUN apt update && \
    apt upgrade -y && \
    DEBIAN_FRONTEND=noninteractive apt install -y git gcc make unzip jq wget tar gzip zip curl file tzdata

# Add Maintainer Info
LABEL maintainer="Marco Ruaro <marco.ruaro@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
#COPY go1.17.5.linux-amd64.tar.gz ./
RUN wget https://go.dev/dl/go1.17.5.linux-amd64.tar.gz && \
     tar -xzf go1.17.5.linux-amd64.tar.gz -C /usr && \
     chmod +x /usr/go/bin/go && \
     export PATH=/usr/go/bin:$PATH && \
     go mod download


# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed


# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN export PATH=/usr/go/bin:$PATH && \
    make install

VOLUME /cnd
WORKDIR /cnd
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM
ENV TZ America/New_York

COPY contrib/localnet/cndnode/wrapper.sh /usr/bin/wrapper.sh
