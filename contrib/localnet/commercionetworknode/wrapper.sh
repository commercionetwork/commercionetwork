#!/usr/bin/env sh

##
## Input parameters
##
ID=${ID:-0}
LOG=${LOG:-commercionetwork.log}

##
## Assert linux binary
##
if ! [ -f "/app/build/commercionetworkd" ]; then
	echo "The binary /app/build/commercionetworkd cannot be found. Please add the binary to the shared folder."
	exit 1
fi

#BINARY_CHECK="$(file "${BINARY}" | grep 'ELF 64-bit LSB executable, x86-64')"
#if [ -z "${BINARY_CHECK}" ]; then
#	echo "Binary needs to be OS linux, ARCH amd64"
#	exit 1
#fi

##
## Run binary with all parameters
##
export CNDHOME="/commercionetwork/node${ID}/commercionetwork"

if [ -d "$(dirname "${CNDHOME}"/"${LOG}")" ]; then
  "/app/build/commercionetworkd" --minimum-gas-prices 0.01ucommercio --home "${CNDHOME}" "$@" | tee "${CNDHOME}/${LOG}"
else
  "/app/build/commercionetworkd" --minimum-gas-prices 0.01ucommercio --home "${CNDHOME}" "$@"
fi