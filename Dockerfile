# Dockerfile basato su Debian con glibc
FROM golang:1.21 as builder

# Installazione delle dipendenze necessarie
RUN apt-get update && apt-get install -y \
    build-essential \
    wget \
    git \
    && apt-get clean

WORKDIR /commercionetwork

# Copia i file del progetto
COPY go.mod go.sum ./

# Scarica le dipendenze Go
RUN go mod download

# Scarica libwasmvm_muslc.a
RUN  set -eux; \    
    export ARCH=$(uname -m); \
    WASM_VERSION=$(go list -m all | grep github.com/CosmWasm/wasmvm | awk '{print $2}') \
    && wget -O /lib/libwasmvm_muslc.a https://github.com/CosmWasm/wasmvm/releases/download/${WASM_VERSION}/libwasmvm_muslc.${ARCH}.a

COPY . .

# Build dell'eseguibile
RUN go build \
    -mod=readonly \
    -tags "netgo" \
    -ldflags \
        "-X github.com/cosmos/cosmos-sdk/version.Name=commercionetwork \
         -X github.com/cosmos/cosmos-sdk/version.AppName=commercionetworkd \
         -X github.com/cosmos/cosmos-sdk/version.Version=$(git describe --tags) \
         -X github.com/cosmos/cosmos-sdk/version.Commit=$(git rev-parse HEAD) \
         -X github.com/cosmos/cosmos-sdk/version.BuildTags=netgo \
         -w -s" \
    -trimpath \
    -o /commercionetwork/build/commercionetworkd \
    ./cmd/commercionetworkd

# Compress the binary to reduce size
# RUN upx --best /commercionetwork/build/commercionetworkd

# Final image
FROM debian:bullseye-slim

# Copia l'eseguibile dal builder
COPY --from=builder /commercionetwork/build/commercionetworkd /bin/commercionetworkd

# Setup dell'ambiente
ENV HOME=/commercionetwork
WORKDIR $HOME

# Esponi le porte necessarie
EXPOSE 26656 26657 1317 9090 9091

ENTRYPOINT ["commercionetworkd"]
