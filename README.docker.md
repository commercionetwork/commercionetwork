# commercio.network blockchain in Docker

## WiP


This repository contains a Dockerfile with all the needed settings to run a complete commercio.network full node and LCD REST server.

You need to export the following environment variables:
 - `CHAINID`, your chain ID
 - `NODENAME`, your full node name

This container makes use of two volumes, `/root/chain` and `/root/genesis`.

While the genesis volume is only needed the first time you set up the container, mounting it every time is not an issue.

This container will setup a new blockchain node in the chain volume only if the mounted directory is empty.

Before starting up the container, if you want to deploy a completely new chain, place the `genesis.json` and `.data` files in your genesis volume, then start up the chain:

```sh
$ docker run -it -p "5000:5000" -e CNCLI_LISTEN_ADDR="tcp://0.0.0.0:5000" -e CHAINID="commercio-testnet5000" NODENAME="name" -v $(pwd)/genesis:/root/genesis -v $(pwd)/chain:/root/chain commercionetwork:latest
```

This command starts a local chain, and exposes the `cncli` REST server on the host 5000 port.
