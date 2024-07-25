#FROM golang:latest 
FROM commercionetwork/libraries

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

#uncomment to build with static
# RUN set -eux; \    
#     export ARCH=$(uname -m); \
#     WASM_VERSION=$(go list -m all | grep github.com/CosmWasm/wasmvm | awk '{print $2}'); \
#     if [ ! -z "${WASM_VERSION}" ]; then \
#       wget -O /lib/libwasmvm_muslc.a https://github.com/CosmWasm/wasmvm/releases/download/${WASM_VERSION}/libwasmvm_muslc.${ARCH}.a; \      
#     fi;

# RUN go mod download

# Build the Go app
RUN make build

#uncomment to build with static
# RUN CGO_ENABLED=1 go build -mod=readonly -tags "netgo,muslc" \
#     -ldflags "-X github.com/cosmos/cosmos-sdk/version.Name=commercionetwork \
#     -X github.com/cosmos/cosmos-sdk/version.AppName=commercionetworkd \
#     -X github.com/cosmos/cosmos-sdk/version.Version=${GIT_VERSION} \
#     -X github.com/cosmos/cosmos-sdk/version.Commit=${GIT_COMMIT} \
#     -X github.com/cosmos/cosmos-sdk/version.BuildTags=netgo,ledger,muslc \
#     -w -s -linkmode=external -extldflags '-Wl,-z,muldefs -lm'" \
#     -trimpath -o ./build/commercionetworkd ./cmd/commercionetworkd/main.go

VOLUME /commercionetwork
WORKDIR /commercionetwork
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM
ENV TZ=America/New_York

COPY ./contrib/localnet/commercionetworknode/wrapper.sh /usr/bin/wrapper.sh