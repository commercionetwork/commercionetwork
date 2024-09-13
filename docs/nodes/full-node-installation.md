# Installing a full node
After you've setup your hardware following the [hardware requirements](hardware-requirements.md) you are now ready to
become a Commercio.network full node. 

:::warning  
Make sure you have read the [hardware requirements](hardware-requirements.md) before starting  
:::

## 1. Installing the software requirements

:::tip
You can try to install the node with our installation tool [commercionetwork-installer](https://github.com/commercionetwork/chain-installer)
:::

In order to update the OS so that you can work properly, execute the following commands:

```bash
apt update && apt upgrade -y
apt install -y git gcc make unzip jq tree 
snap install --classic go

export NODENAME="<your-moniker>"

echo 'export GOPATH="$HOME/go"' >> ~/.profile
echo 'export PATH="$GOPATH/bin:$PATH"' >> ~/.profile
echo 'export PATH="$PATH:/snap/bin"' >> ~/.profile
echo "export NODENAME=\"$NODENAME\"" >> ~/.profile

source ~/.profile
```

## 2. Chain selection
Before installing the node, please select which chain you would like to connect to
- **Mainnet**: commercio-3
- **Testnet**: commercio-testnet11k


```bash
rm -rf commercio-chains
mkdir commercio-chains && cd commercio-chains
git clone https://github.com/commercionetwork/chains.git .
cd commercio-<chain-version>
```

:::tip
Always remember to pick the latest chain version listed inside [chains repo](https://github.com/commercionetwork/chains) 
::: 

## 3. Install binaries, genesis file and setup configuration

Compile binaries 

```bash
pkill commercionetworkd
git init . 
git remote add origin https://github.com/commercionetwork/commercionetwork.git
git pull
git checkout tags/$(cat .data | grep -oP 'Release\s+\K\S+')
go mod verify
make install
```

Test if you have the correct binaries version:

```bash
commercionetworkd version
# Should output the same version written inside the .data file.
# cat .data | grep -oP 'Release\s+\K\S+'
```

Setup the validator node name. We will use the same name for node as well as the wallet key:

```bash
export CHAINID=commercio-$(cat .data | grep -oP 'Name\s+\K\S+')
cat <<EOF >> ~/.profile
export CHAINID="$CHAINID"
EOF
```

Init the `.commercionetwork` folder with the basic configuration

:::warning  
At this point there may be some differences if you are using `KMS` with `HSM`. Specifications will be published shortly.
:::

```bash
commercionetworkd unsafe-reset-all --home ~/.commercionetwork
# If you get a error because .commercionetwork folder is not present don't worry 

commercionetworkd init $NODENAME --home ~/.commercionetwork
# If you get a error because .commercionetwork folder is present don't worry 
```

Install `genesis.json` file

```bash
rm -rf ~/.commercionetwork/config/genesis.json
cp genesis.json ~/.commercionetwork/config
```

Change the persistent peers inside `config.toml` file

```bash
sed -e "s|persistent_peers = \".*\"|persistent_peers = \"$(cat .data | grep -oP 'Persistent peers\s+\K\S+')\"|g" ~/.commercionetwork/config/config.toml > ~/.commercionetwork/config/config.toml.tmp
mv ~/.commercionetwork/config/config.toml.tmp  ~/.commercionetwork/config/config.toml
```

Change the seeds inside the `config.toml` file
```bash
sed -e "s|seeds = \".*\"|seeds = \"$(cat .data | grep -oP 'Seeds\s+\K\S+')\"|g" ~/.commercionetwork/config/config.toml > ~/.commercionetwork/config/config.toml.tmp
mv ~/.commercionetwork/config/config.toml.tmp  ~/.commercionetwork/config/config.toml
```

Change `external_address` value to contact your node using public ip of your node:
```bash
PUB_IP=`curl -s -4 icanhazip.com`
sed -e "s|external_address = \".*\"|external_address = \"$PUB_IP:26656\"|g" ~/.commercionetwork/config/config.toml > ~/.commercionetwork/config/config.toml.tmp
mv ~/.commercionetwork/config/config.toml.tmp  ~/.commercionetwork/config/config.toml
```
## 4. Configure the service

```bash
tee /etc/systemd/system/commercionetworkd.service > /dev/null <<EOF  
[Unit]
Description=Commercio Node
After=network-online.target

[Service]
User=root
ExecStart=/root/go/bin/commercionetworkd start
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF
```

## 5. Sync node


Choose 1 of these 3 ways to syncronize your node to the blockchain:
1. [From the start](#from-the-start)
2. [Using the state sync feature](#using-the-state-sync-feature)
3. [Using the quicksync dump](#using-the-quicksync-dump)
### From the start

If you intend to syncronize eveything from the start you can skip this part and continue with the configuration.     
100.000 blocks should synchronize within an hour.     

### Using the state sync feature

State sync is a Cosmos SDK module that allows validators to quickly join the network by syncing their node with a snapshot-enabled RPC from a trusted block height. This reduces the sync time from days to minutes but provides only the most recent state, not the full transaction history.

Under the state sync section in `~/.commercionetwork/config/config.toml`, you will find multiple settings that need to be configured for your node to use state sync. The following pieces of information must be obtained for light client verification:

- At least 2 available RPC servers
- A trusted height.
- The block ID hash of the trusted height.

**RPC servers:**

- Testnet:
  https://rpc-testnet.commercio.network, http://rpc2-testnet.commercio.network
- Mainnet:
  https://rpc-mainnet.commercio.network, https://rpc2-mainnet.commercio.network

:::tip
You can get informations about rpc services at [chain data](https://github.com/commercionetwork/chains) repository
:::


**Testnet**
```bash
TRUST_RPC1="rpc-testnet.commercio.network:80"
TRUST_RPC2="rpc2-testnet.commercio.network:80"
CURR_HEIGHT=$(curl -s "http://$TRUST_RPC1/block" | jq -r '.result.block.header.height')
TRUST_HEIGHT=$((CURR_HEIGHT-(CURR_HEIGHT%10000)))
TRUST_HASH=$(curl -s "http://$TRUST_RPC1/block?height=$TRUST_HEIGHT" | jq -r '.result.block_id.hash')
# CONFIG OUTPUT
printf "\n\nenable = true\nrpc_servers = \"$TRUST_RPC1,$TRUST_RPC2\"\ntrust_height = $TRUST_HEIGHT\ntrust_hash = \"$TRUST_HASH\"\n\n"
```

**Mainnet**
```bash
TRUST_RPC1="rpc-mainnet.commercio.network:80"
TRUST_RPC2="rpc2-mainnet.commercio.network:80"
CURR_HEIGHT=$(curl -s "http://$TRUST_RPC1/block" | jq -r '.result.block.header.height')
TRUST_HEIGHT=$((CURR_HEIGHT-(CURR_HEIGHT%10000)))
TRUST_HASH=$(curl -s "http://$TRUST_RPC1/block?height=$TRUST_HEIGHT" | jq -r '.result.block_id.hash')
# CONFIG OUTPUT
printf "\n\nenable = true\nrpc_servers = \"$TRUST_RPC1,$TRUST_RPC2\"\ntrust_height = $TRUST_HEIGHT\ntrust_hash = \"$TRUST_HASH\"\n\n"
```

The command should be return something like the following:
```
enable = true
rpc_servers = "rpc-mainnet.commercio.network:80,rpc2-mainnet.commercio.network:80"
trust_height = 6240000
trust_hash = "FCA27CBCAC3EECAEEBC3FFBB5B5433A421EF4EA873EB2A573719B0AA5093EF4C"
```

Edit `~/.commercionetwork/config/config.toml` with these settings accordingly:

```toml
[statesync]

enable = true
rpc_servers = "rpc-mainnet.commercio.network:80,rpc2-mainnet.commercio.network:80"
trust_height = 6240000
trust_hash = "FCA27CBCAC3EECAEEBC3FFBB5B5433A421EF4EA873EB2A573719B0AA5093EF4C"
```


### Using the quicksync dump

```bash
wget "https://quicksync.commercio.network/$CHAINID.latest.tgz" -P ~/.commercionetwork/
# Check if the checksum matches the one present inside https://quicksync.commercio.network
cd ~/.commercionetwork/
tar -zxf $CHAINID.latest.tgz
rm $CHAINID.latest.tgz
```




Now you can start you full node. Enable the newly created server and try starting it using:

```bash
# Start the node  
systemctl enable commercionetworkd  
systemctl start commercionetworkd
```

Check if the sync has been started. Use `Ctrl + C` to interrupt the `journalctl` command

```bash
journalctl -u commercionetworkd -f
# OUTPUT SHOULD BE LIKE BELOW
#
# Aug 13 16:30:20 commerciotestnet-node4 commercionetworkd[351]: I[2019-08-13|16:30:20.722] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
# Aug 13 16:30:20 commerciotestnet-node4 commercionetworkd[351]: I[2019-08-13|16:30:20.728] Committed state                              module=state height=1 txs=0 appHash=9815044185EB222CE9084AA467A156DFE6B4A0B1BAAC6751DE86BB31C83C4B08
# Aug 13 16:30:20 commerciotestnet-node4 commercionetworkd[351]: I[2019-08-13|16:30:20.745] Executed block                               module=state height=2 validTxs=0 invalidTxs=0
# Aug 13 16:30:20 commerciotestnet-node4 commercionetworkd[351]: I[2019-08-13|16:30:20.751] Committed state                              module=state height=2 txs=0 appHash=96BFD9C8714A79193A7913E5F091470691B195E1E6F028BC46D6B1423F7508A5
# Aug 13 16:30:20 commerciotestnet-node4 commercionetworkd[351]: I[2019-08-13|16:30:20.771] Executed block                               module=state height=3 validTxs=0 invalidTxs=0
```

## 6. Start the REST API
Each full node can start up its own REST API service. 
This allows it to expose some endpoints that can be used in order to query the chain state at any moment. 

If you want to start such a service, you need to change the parameters of your `~/.commercionetwork/config/app.toml` as follow

```toml
...
[api]

# Enable defines if the API server should be enabled.
enable = true

# Swagger defines if swagger documentation should automatically be registered.
swagger = true

# Address defines the API server to listen on.
address = "tcp://0.0.0.0:1317"
...

``` 

Apply the configuration using
```bash
systemctl restart commercionetworkd
```


This starts up the REST server making it reachable on port `1317`.     
**From here, if you want you can use services such as [Nginx](https://www.nginx.com/) in order to make it available to other devices.**


## Install and config cosmovisor

### What is cosmovisor?

`cosmovisor` is a small process manager for Cosmos SDK application binaries that monitors the governance module for incoming chain upgrade proposals. 
If it sees a proposal that gets approved, cosmovisor can automatically download the new binary, stop the current binary, switch from the old binary to the new one, and finally restart the node with the new binary.    
**As in all other explanations, it is assumed that you are acting as root. Change the parameters and variables depending on your installation environment**
### Installation

Download and compile cosmovisor:
```bash
cd $HOME
git clone https://github.com/cosmos/cosmos-sdk.git
cd cosmos-sdk
git checkout cosmovisor/v1.3.0
cd cosmovisor
make cosmovisor
cp cosmovisor $HOME/go/bin
```

Make cosmovisor folder structure:
```bash
mkdir -p $HOME/.commercionetwork/cosmovisor/genesis/bin
mkdir -p $HOME/.commercionetwork/cosmovisor/upgrades
```

Copy `commercionetwork` to cosmovisor folder
```bash
cp $HOME/go/bin/commercionetworkd $HOME/.commercionetwork/cosmovisor/genesis/bin
```

After installation your `.commercionetwork` folder should be structured like below

```
.commercionetwork
├── config
│   └── app.toml
│   └── config.toml
│   └── genesis.json
│   └── node_key.json
│   └── priv_validator_key.json
├── data
│   └── priv_validator_state.json
└── cosmovisor
    └── current -> /path/to/the/current/version/of/commercionetworkd
    └── genesis
    │   └── bin
    │      └── commercionetworkd
    └── upgrades
        └── <name>
           └── bin
               └── commercionetworkd
```


Configure the service:
```bash

tee /etc/systemd/system/commercionetworkd.service > /dev/null <<EOF  
[Unit]
Description=Commercio Network Node
After=network.target

[Service]
User=root
LimitNOFILE=4096

Restart=always
RestartSec=3

Environment="LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/root/bin/go" # <-- set this only if you compiled "commercionetworkd" locally
Environment="DAEMON_NAME=commercionetworkd"
Environment="DAEMON_HOME=/root/.commercionetwork"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_LOG_BUFFER_SIZE=512"
Environment="UNSAFE_SKIP_BACKUP=true" # Set to false if you want make backup during the upgrade

ExecStart=/root/go/bin/cosmovisor run start --home="/root/.commercionetwork" 

[Install]
WantedBy=multi-user.target
EOF
```

Now you can start your full node. Enable the newly created server and try to start it using:
```bash
systemctl enable commercionetworkd  
systemctl start commercionetworkd
```

Check if the sync has been started. Use `Ctrl + C` to interrupt the `journalctl` command:
```bash
journalctl -u commercionetworkd -f
```

Set the enviorment variables of cosmovisor for you convenience

```bash
echo 'export DAEMON_NAME=commercionetworkd' >> ~/.profile
echo 'export DAEMON_HOME=/root/.commercionetwork' >> ~/.profile
echo 'export DAEMON_RESTART_AFTER_UPGRADE=true' >> ~/.profile
echo 'export DAEMON_ALLOW_DOWNLOAD_BINARIES=false' >> ~/.profile
echo 'export DAEMON_LOG_BUFFER_SIZE=512' >> ~/.profile
echo 'export UNSAFE_SKIP_BACKUP=true' >> ~/.profile
```

## Verify Node Synchronization with Blockchain 

To verify if your node is properly synchronized with the blockchain, you can compare its height with the current blockchain height by following steps below:

First, ensure you have `curl` and `jq` installed on your server. While `curl` is usually installed by default, you may need to install `jq` using the command:

```bash
sudo apt install jq
```

Then, run these two commands:
```bash
curl -s 127.0.0.1:26657/block | jq -r '.result.block.header.height'
```
> The number returned represents the current height of your node.

```bash
curl -s https://rpc-mainnet.commercio.network/block | jq -r '.result.block.header.height'
```
> The number returned represents the height of the blockchain.

The two values can differ by a few blocks; this depends on when the commands are executed. For example, if you run the second command a minute later, you should expect a difference of approximately 10 blocks. However, if the difference is significantly larger, it means your node is not aligned.

You can also check the node's operation through logs, using the following command:

```bash
journalctl -u commercionetworkd.service -f
```
You should see the node logs and block production. To stop the log output, press Control + C.

## Next step
Now that you are a Commercio.network full node, if you want you can become a validator.
If you wish to do so, please read the [*Becoming a validator* guide](validator-node-installation.md).
