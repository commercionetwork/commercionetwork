DEP := $(shell command -v dep 2> /dev/null)

get_tools:
ifndef DEP
	@echo "Installing dep"
	go get -u -v github.com/golang/dep/cmd/dep
else
	@echo "Dep is already installed..."
endif

get_vendor_deps:
	@echo "--> Generating vendor directory via dep ensure"
	@dep ensure -v -vendor-only

update_vendor_deps:
	@echo "--> Running dep ensure"
	@dep ensure -v -update

build:
	go build -o "cnd" -tags "ledger" ./cmd/cnd/main.go
	go build -o "cncli" -tags "ledger" ./cmd/cncli/main.go

build-darwin:
	env GOOS=darwin GOARCH=386 go build -o ./build/darwin/cncli-darwin-386 ./cmd/cncli/main.go
	env GOOS=darwin GOARCH=386 go build -o ./build/darwin/cnd-darwin-386 ./cmd/cnd/main.go
	env GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/cncli-darwin-amd64 ./cmd/cncli/main.go
	env GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/cnd-darwin-amd64 ./cmd/cnd/main.go

build-linux:
	env GOOS=linux GOARCH=386 go build -o ./build/linux/cncli-linux-386 ./cmd/cncli/main.go
	env GOOS=linux GOARCH=386 go build -o ./build/linux/cnd-linux-386 ./cmd/cnd/main.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/linux/cncli-linux-amd64 ./cmd/cncli/main.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/linux/cnd-linux-amd64 ./cmd/cnd/main.go

build-windows:
	env GOOS=windows GOARCH=386 go build -o ./build/windows/cncli-windows-386.exe ./cmd/cncli/main.go
	env GOOS=windows GOARCH=386 go build -o ./build/windows/cnd-windows-386.exe ./cmd/cnd/main.go
	env GOOS=windows GOARCH=amd64 go build -o ./build/windows/cncli-windows-amd64.exe ./cmd/cncli/main.go
	env GOOS=windows GOARCH=amd64 go build -o ./build/windows/cnd-windows-amd64.exe ./cmd/cnd/main.go

build-all:
	make build-darwin
	make build-linux
	make build-windows

install:
	go install -tags "ledger" ./cmd/cnd
	go install -tags "ledger" ./cmd/cncli

crossbuild_windows_linux:
	set GOARCH=amd64
	set GOOS=linux
	make build