#!/usr/bin/env sh

##
## Input parameters
##
BINARYCLI=/cnd/${BINARYCLI:-commercionetworkd}
ID=${ID:-0}
LOGCLI=${LOGCLI:-cnd.cncli.log}

RPCNODE="cndnode$ID"

##
## Assert linux binary
##
if ! [ -f "${BINARYCLI}" ]; then
	echo "The binary $(basename "${BINARYCLI}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'cnd' E.g.: -e BINARY=cnd_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARYCLI" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export CNDHOME="/cnd/node${ID}/cncli"


if [ -d "$(dirname "${CNDHOME}"/"${LOGCLI}")" ]; then
  "${BINARYCLI}" rest-server --trust-node --node tcp://${RPCNODE}:26657 --laddr tcp://0.0.0.0:1317 | tee "${CNDHOME}/${LOGCLI}"
else
  "${BINARYCLI}" rest-server --trust-node --node tcp://${RPCNODE}:26657 --laddr tcp://0.0.0.0:1317
fi

