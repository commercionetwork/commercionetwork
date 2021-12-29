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

export CHAIN_DIR="/root/chain"
export GENESIS_DIR="/root/genesis"
export BINARY=${BINARY:-commercionetworkd}


CND_FLAGS="--home=$CHAIN_DIR $CND_EXTRA_FLAGS"
CND_START_FLAGS="$CND_START_FLAGS"
if [ -z "$(ls -A $CHAIN_DIR)" ]; then
	# chain directory empty, let's build a new chain
	./$BINARY unsafe-reset-all $CND_FLAGS
	./$BINARY init $NODENAME $CND_FLAGS
	cp $GENESIS_DIR/genesis.json $CHAIN_DIR/config/genesis.json
	sed -e "s|persistent_peers = \".*\"|persistent_peers = \"$(cat $GENESIS_DIR/.data | grep -oP 'Persistent peers\s+\K\S+')\"|g" $CHAIN_DIR/config/config.toml > $CHAIN_DIR/config/config.toml.tmp
	mv $CHAIN_DIR/config/config.toml.tmp  $CHAIN_DIR/config/config.toml
	sed -e "s|enable = false|enable = true|g" $CHAIN_DIR/config/app.toml | \ $CHAIN_DIR/config/app.toml.tmp
	sed -e "s|swagger = false|swagger = true|g" > $CHAIN_DIR/config/app.toml.tmp
	mv $CHAIN_DIR/config/app.toml.tmp  $CHAIN_DIR/config/app.toml
fi

./$BINARY start $CND_FLAGS $CND_START_FLAGS
