PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
GOBIN ?= $(GOPATH)/bin
GOSUM := $(shell which gosum)

include Makefile.ledger

export GO111MODULE = on

all: test build

########################################
### Install

install: go.sum
	GO111MODULE=on go install -tags "$(build_tags)" ./cmd/cnd
	GO111MODULE=on go install -tags "$(build_tags)" ./cmd/cncli

########################################
### Build

build: go.sum
	GO111MODULE=on go build -o "build/cnd" -tags "$(build_tags)" ./cmd/cnd/main.go
	GO111MODULE=on go build -o "build/cncli" -tags "$(build_tags)" ./cmd/cncli/main.go

# TODO: Note:386 builds are disabled due to a bug inside the Cosmos SDK:
# github.com/cosmos/cosmos-sdk@v0.28.2-0.20190826165445-eeb847c8455b/simapp/state.go:75:46: constant 1000000000000 overflows

build-darwin: go.sum
	env GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ./build/Darwin-AMD64/cncli -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ./build/Darwin-AMD64/cnd -tags "$(build_tags)" ./cmd/cnd/main.go
# 	env GO111MODULE=on GOOS=darwin GOARCH=386 go build -o ./build/Darwin-386/cncli -tags "$(build_tags)" ./cmd/cncli/main.go
# 	env GO111MODULE=on GOOS=darwin GOARCH=386 go build -o ./build/Darwin-386/cnd -tags "$(build_tags)" ./cmd/cnd/main.go

build-linux: go.sum
	env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ./build/Linux-AMD64/cncli -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ./build/Linux-AMD64/cnd -tags "$(build_tags)" ./cmd/cnd/main.go
# 	env GO111MODULE=on GOOS=linux GOARCH=386 go build -o ./build/Linux-386/cncli -tags "$(build_tags)" ./cmd/cncli/main.go
# 	env GO111MODULE=on GOOS=linux GOARCH=386 go build -o ./build/Linux-386/cnd -tags "$(build_tags)" ./cmd/cnd/main.go

build-windows: go.sum
	env GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ./build/Windows-AMD64/cncli.exe -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ./build/Windows-AMD64/cnd.exe -tags "$(build_tags)" ./cmd/cnd/main.go
# 	env GO111MODULE=on GOOS=windows GOARCH=386 go build -o ./build/Windows-386/cncli.exe -tags "$(build_tags)" ./cmd/cncli/main.go
# 	env GO111MODULE=on GOOS=windows GOARCH=386 go build -o ./build/Windows-386/cnd.exe -tags "$(build_tags)" ./cmd/cnd/main.go

build-all: go.sum
	make build-darwin
	make build-linux
	make build-windows

prepare-release: go.sum build-all
	rm -f ./build/Darwin-386.zip ./build/Darwin-AMD64.zip
	rm -f ./build/Linux-386.zip ./build/Linux-AMD64.zip
	rm -f ./build/Windows-386.zip ./build/Windows-AMD64.zip
	zip -jr ./build/Darwin-AMD64.zip ./build/Darwin-AMD64/cncli ./build/Darwin-AMD64/cnd
	zip -jr ./build/Linux-AMD64.zip ./build/Linux-AMD64/cncli ./build/Linux-AMD64/cnd
	zip -jr ./build/Windows-AMD64.zip ./build/Windows-AMD64/cncli.exe ./build/Windows-AMD64/cnd.exe
# 	zip -jr ./build/Darwin-386.zip ./build/Darwin-386/cncli ./build/Darwin-386/cnd
# 	zip -jr ./build/Linux-386.zip ./build/Linux-386/cncli ./build/Linux-386/cnd
# 	zip -jr ./build/Windows-386.zip ./build/Windows-386/cncli.exe ./build/Windows-386/cnd.exe

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

########################################
### Testing

test: test_unit

test_unit:
	@VERSION=$(VERSION) go test -mod=readonly $(PACKAGES_NOSIMULATION) -tags='ledger test_ledger_mock'