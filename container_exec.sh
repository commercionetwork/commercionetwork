#!/bin/bash

# commercio.network blockchain container startup script
#
# This script spin up a commercio.network blockchain if needed,
# i.e. it checks if a chain already exists before creating a new one.
#
# Environment variables needed:
# CHAINID: your chain ID
# NODENAME: your node name
# CHAIN_DIR: chain directory
# GENESIS_DIR: folder with genesis files informations
#
# If you're gonna deploy a new chain, make sure to pass a "genesis.json"
# and ".data" files by mounting a Docker volume on /root/genesis.

# if $CHAIN_DIR is empty, assume we need to spin up a new chain

printf "This script will be replaced by the script for version 3.0\n"

exit 0

export CHAIN_DIR="/app/chain"
#export GENESIS_DIR="/app/genesis"

CND_FLAGS="--home=$CHAIN_DIR $CND_EXTRA_FLAGS"
CND_START_FLAGS="$CND_START_FLAGS"
if [ -z "$(ls -A $CHAIN_DIR)" ]; then
	# chain directory empty, let's build a new chain
	./commercionetworkd unsafe-reset-all $CND_FLAGS
	./commercionetworkd init $NODENAME $CND_FLAGS
	cp $GENESIS_DIR/genesis.json $CHAIN_DIR/config/genesis.json
	sed -e "s|persistent_peers = \".*\"|persistent_peers = \"$(cat $GENESIS_DIR/.data | grep -oP 'Persistent peers\s+\K\S+')\"|g" $CHAIN_DIR/config/config.toml > $CHAIN_DIR/config/config.toml.tmp
	mv $CHAIN_DIR/config/config.toml.tmp  $CHAIN_DIR/config/config.toml
fi

./commercionetworkd start $CND_FLAGS $CND_START_FLAGS &
sleep 3 # let cnd start first before running cncli rest server
./commercionetworkd rest-server --chain-id=$CHAINID --laddr $CNCLI_LISTEN_ADDR 
