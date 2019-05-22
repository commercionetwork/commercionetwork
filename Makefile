GOBIN ?= $(GOPATH)/bin
GOSUM := $(shell which gosum)

include Makefile.ledger

all: install

install: go.sum
	GO111MODULE=on go install -tags "$(build_tags)" ./cmd/cnd
	GO111MODULE=on go install -tags "$(build_tags)" ./cmd/cncli

build: go.sum
	GO111MODULE=on go build -o "cnd" -tags "$(build_tags)" ./cmd/cnd/main.go
	GO111MODULE=on go build -o "cncli" -tags "$(build_tags)" ./cmd/cncli/main.go

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

build-darwin:
	env GO111MODULE=on GOOS=darwin GOARCH=386 go build -o ./build/darwin/cncli-darwin-386 -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=darwin GOARCH=386 go build -o ./build/darwin/cnd-darwin-386 -tags "$(build_tags)" ./cmd/cnd/main.go
	env GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/cncli-darwin-amd64 -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/cnd-darwin-amd64 -tags "$(build_tags)" ./cmd/cnd/main.go

build-linux:
	env GO111MODULE=on GOOS=linux GOARCH=386 go build -o ./build/linux/cncli-linux-386 -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=linux GOARCH=386 go build -o ./build/linux/cnd-linux-386 -tags "$(build_tags)" ./cmd/cnd/main.go
	env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ./build/linux/cncli-linux-amd64 -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ./build/linux/cnd-linux-amd64 -tags "$(build_tags)" ./cmd/cnd/main.go

build-windows:
	env GO111MODULE=on GOOS=windows GOARCH=386 go build -o ./build/windows/cncli-windows-386.exe -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=windows GOARCH=386 go build -o ./build/windows/cnd-windows-386.exe -tags "$(build_tags)" ./cmd/cnd/main.go
	env GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ./build/windows/cncli-windows-amd64.exe -tags "$(build_tags)" ./cmd/cncli/main.go
	env GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o ./build/windows/cnd-windows-amd64.exe -tags "$(build_tags)" ./cmd/cnd/main.go

build-all:
	make build-darwin
	make build-linux
	make build-windows
