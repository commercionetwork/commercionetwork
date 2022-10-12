# Installing a full node
After you've setup your hardware following the [hardware requirements](hardware-requirements.md) you are now ready to
become a Commercio.network full node. 

:::warning  
Make sure you have read the [hardware requirements](hardware-requirements.md) before starting  
:::

## 1. Installing the software requirements
In order to update the OS so that you can work properly, execute the following commands:

```bash
apt update && apt upgrade -y
apt install -y git gcc make unzip
snap install --classic go

export NODENAME="<your-moniker>"

echo 'export GOPATH="$HOME/go"' >> ~/.profile
echo 'export PATH="$GOPATH/bin:$PATH"' >> ~/.profile
echo 'export PATH="$PATH:/snap/bin"' >> ~/.profile
echo "export NODENAME=\"$NODENAME\"" >> ~/.profile

source ~/.profile
```

## 2. Chain selection
Before installing the node, please select which chain you would like to connect to (for example **testent11k**)

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
commercionetworkd unsafe-reset-all
# If you get a error because .commercionetwork folder is not present don't worry 

commercionetworkd init $NODENAME
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
2. [Using the state sync future](#using-the-state-sync-feature)
3. [Using the quicksync dump](#using-the-quicksync-dump)
### From the start

If you intend to syncronize eveything from the start you can skip this part and continue with the configuration.     
100.000 blocks should synchronize within an hour.     

### Using the state sync feature

Under the state sync section in `~/.commercionetwork/config/config.toml` you will find multiple settings that need to be configured in order for your node to use state sync.
You need get information from chain about trusted block using

Select open rpc services of chains

* Testnet: 
  * With name: rpc-testnet.commercio.network, rpc2-testnet.commercio.network
  * With ip: 157.230.110.18:26657, 46.101.146.48:26657
* Mainnet: (WIP)
  * https://rpc-mainnet.commercio.network, https://rpc2-mainnet.commercio.network (WIP)



```bash
TRUST_RPC1="157.230.110.18:26657"
TRUST_RPC2="46.101.146.48:26657"
curl -s "http://$TRUST_RPC1/block" | jq -r '.result.block.header.height + "\n" + .result.block_id.hash'
```

The command should be return block height and hash of block as follow:
```
5075021
EB1032C6DFC9F2708B16DF8163CAB2258B0F1E1452AEF031CA3F32004F54C9D1
```

Edit these settings accordingly:

```
[statesync]

enable = true

rpc_servers = "$TRUST_RPC1,$TRUST_RPC2"
trust_height = 5075021
trust_hash = "EB1032C6DFC9F2708B16DF8163CAB2258B0F1E1452AEF031CA3F32004F54C9D1"
```


### Using the quicksync dump:

```bash
wget "https://quicksync.commercio.network/$CHAINID.latest.tgz" -P ~/.commercionetwork/
# Check if the checksum matches the one present inside https://quicksync.commercio.network
cd ~/.commercionetwork/
tar -zxf $(echo $CHAINID).latest.tgz
```




Now you can start you full node. Enable the newly created server and try starting it using:

```bash
# Start the node  
systemctl enable commercionetworkd  
systemctl start commercionetworkd
```

Control if the sync was started. Use `Ctrl + C` to interrupt the `journalctl` command

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
This will allow it to expose some endpoints that can be used in order to query the chain state at any moment. 

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


This will start up the REST server and make it reachable using the port `1317`.     
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
git checkout cosmovisor/v0.1.0
# you can use version v1.0.0 or v1.1.0
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
│   └── node_id.json
│   └── priv_validator_key.json
├── data
│   └── priv_validator_state.json
└── cosmovisor
    └── current
    └── genesis
    └── bin
    │   └── commercionetword
    └── upgrades
    └── <name>
        └── bin
            └── commercionetword
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

ExecStart=/root/go/bin/cosmovisor start --home="/root/.commercionetwork" 

[Install]
WantedBy=multi-user.target
EOF
```

Now you can start your full node. Enable the newly created server and try to start it using:
```bash
systemctl enable commercionetworkd  
systemctl start commercionetworkd
```

Control if the sync was started. Use `Ctrl + C` to interrupt the `journalctl` command:
```bash
journalctl -u commercionetworkd -f
```

Set env of cosmovisor for you convenience

```bash
echo 'export DAEMON_NAME=commercionetworkd' >> ~/.profile
echo 'export DAEMON_HOME=/root/.commercionetwork' >> ~/.profile
echo 'export DAEMON_RESTART_AFTER_UPGRADE=true' >> ~/.profile
echo 'export DAEMON_ALLOW_DOWNLOAD_BINARIES=false' >> ~/.profile
echo 'export DAEMON_LOG_BUFFER_SIZE=512' >> ~/.profile
echo 'export UNSAFE_SKIP_BACKUP=true' >> ~/.profile
```


## Next step
Now that you are a Commercio.network full node, if you want you can become a validator.
If you wish to do so, please read the [*Becoming a validator* guide](validator-node-installation.md).
