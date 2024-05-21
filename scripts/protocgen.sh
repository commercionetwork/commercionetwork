#!/usr/bin/env bash

unameOut="$(uname -s)"
case "${unameOut}" in
    Darwin*)    set -eo pipefail;;
    *)          set -o pipefail
esac

# get protoc executions
# go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos 2>/dev/null
# get cosmos sdk from github
#go get github.com/cosmos/cosmos-sdk@v0.47.8 2>/dev/null

# Get the path of the cosmos-sdk repo from go/pkg/mod
#cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
#proto_dirs=$(find . -path ./third_party -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
protoc_install_proto_gen_doc() {
  echo "Installing protobuf protoc-gen-doc plugin"
  (go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest 2> /dev/null)
}

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./commercionetwork -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yaml $file
      # protoc --proto_path=. \
      #        --grpc-gateway_out=logtostderr=true,paths=source_relative:./ \
      #        --go-grpc_out=paths=source_relative:./ \
      #        $file
    fi
  done
done

protoc_install_proto_gen_doc

echo "moving the files..."
cp -r ./github.com/commercionetwork/commercionetwork/* ../
rm -rf ./github.com
go mod tidy