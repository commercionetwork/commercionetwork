# Install Commercio.network
This guide will explain how to install the `cnd` and `cncli` entrypoints onto your system. 
With these installed on a server, you can participate in the mainnet as either a [Full Node](./join-mainnet.md) or 
a [Validator](./validator-setup.md).

## Install Go
Install `go` by following the [official docs](https://golang.org/doc/install). Remember to set your `$GOPATH`, `$GOBIN`, 
and `$PATH` environment variables, for example:

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bash_profile
echo "export GOBIN=$GOPATH/bin" >> ~/.bash_profile
echo "export PATH=$PATH:$GOBIN" >> ~/.bash_profile
source ~/.bash_profile
```

Remember that Go 1.12.1+ is required for Commercio.network. 

## Install the binaries
Next, let's install the latest version of Commercio.network. 
Here we'll use the master branch, which contains the latest stable release. 
If necessary, make sure you git checkout the correct [released version](https://github.com/commercionetwork/commercionetwork/releases).

```bash
mkdir -p $GOPATH/src/github.com/commercionetwork
cd $GOPATH/src/github.com/commercionetwork
git clone https://github.com/commercionetwork/commercionetwork
cd commercionetwork && git checkout master
make tools install
```

> NOTE: If you have issues at this step, please check that you have the latest stable version of Go installed.

That will install the `cnd` and `cncli` binaries. Verify that everything is OK by running the following commands: 

```bash
cnd version --long
cncli version --long
```

### Build tags
Build tags indicate special features that have been enabled in the binary. 

| Build tag | Description |
| :-------- | :---------- |
| netgo | Name resolution will use pure Go code |
| ledger | Ledger devices are supported (hardware wallets) | 

## Next 
Now you can [join the mainnet](./join-mainnet.md), [the public testnet](./join-testnet.md) 
or [create your own testnet](./deploy-testnet.md).