# Join the Public Testnet
:::tip Current Testnet.   
See the [testnet repo](https://github.com/commercionetwork/testnets) for information on the latest testnet, including the correct
version of the Commercio.network executable to use and the details about the genesis file. 
:::

::: warning 
You need to [install Commercio.network](./installation.md) before you go further.
:::

## Starting a new node
> NOTE: If you ran a full node on previous testnet, please skip to [Upgrading from Previous Testnet](#upgrading-your-node).

To start a new node, the mainnet instructions apply: 
* [Join the mainnet](./join-mainnet.md)
* [Deploy a validator](./validator-setup.md)

The only difference is the executable version and the genesis file.
See the [testnet repo](https://github.com/commercionetwork/testnets) for information on testnets, including
the correct version of the Commercio.network executable to use and the details about the genesis file. 


## Upgrading your node  
These instructions are for full nodes that have ran on previous versions of Commercio.network and would like to upgrade
to the latest testnet. 

### Reset the data
First, remove the outdated files and reset the current chain data. 

```bash
rm $HOME/.cnd/config/addrbook.json $HOME/.cnd/config/genesis.json
cnd unsafe-reset-all
```

Your node is now in a pristine state while keeping the original `priv_validator.json` and `config.toml`. 
If you had any sentry nodes or full nodes setup before, your node will still try to connect to them, but may fail if 
they haven't also been upgraded.

::: warning 
Make sure that every node has a unique `priv_validator.json`. Do not copy the `priv_validator.json` from an old node to 
multiple new nodes. Running two nodes with the same `priv_validator.json` will cause you to double sign.
:::

## Software upgrade
Now it is time to upgrade the software:

```bash
cd $GOPATH/src/github.com/commercionetwork/commercionetwork
git fetch --all && git checkout master
make update_tools install
```

> NOTE: If you have issues at this step, please check that you have the latest stable version of Go installed. 

Note we use `master` here since it contains the latest stable release. 
See the [chains repo](https://github.com/commercionetwork/chains) for details on which version is needed for which 
testnet, and the [chain release page](https://github.com/commercionetwork/commercionetwork/releases) for details on 
each release.

Your full node has been cleanly upgraded!