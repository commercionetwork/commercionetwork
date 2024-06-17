# Installing a full node with statesync

This guide  provides instructions for installing a node in mainnet with statesync and using a dump snapshot from https://quicksync.commercio.network/ as a node data.

After you've setup your hardware following the [hardware requirements](hardware-requirements.md) you are now ready to
become a Commercio.network full node. 

:::warning  
Make sure you have read the [hardware requirements](hardware-requirements.md) before starting  
:::

## 1. Installing the software requirements
Choose the name of your node. In the guide you will need to replace the variable `<your-moniker>` with the name you choose
In order to update the OS so that you can work properly, execute the following commands:

```bash
apt update && apt upgrade -y
apt install -y git gcc make unzip jq
snap install --classic go

export NODENAME="<your-moniker>"

echo 'export GOPATH="$HOME/go"' >> ~/.profile
echo 'export PATH="$GOPATH/bin:$PATH"' >> ~/.profile
echo 'export PATH="$PATH:/snap/bin"' >> ~/.profile
echo "export NODENAME=\"$NODENAME\"" >> ~/.profile

source ~/.profile
```

## 2. Chain selection

Select the mainnet chain and clone its repository

```bash
CHAINVERSION=commercio-3
rm -rf commercio-chains
mkdir commercio-chains && cd commercio-chains
git clone https://github.com/commercionetwork/chains.git .
cd $CHAINVERSION
```
## 3. Install binaries and setup configuration

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

Test if you have the correct binaries version. Get version from chain data

```bash
cat .data | grep -oP 'Release\s+\K\S+'
```

Compare the result with the output of command

```bash
commercionetworkd version
```

Setup the validator node name. We will use the same name for node as well as the wallet key:

```bash
export CHAINID=commercio-$(cat .data | grep -oP 'Name\s+\K\S+')
cat <<EOF >> ~/.profile
export CHAINID="$CHAINID"
EOF
```

Init the `.commercionetwork` folder with the basic configuration

```bash
commercionetworkd unsafe-reset-all --home ~/.commercionetwork
# If you get a error because .commercionetwork folder is not present don't worry 

commercionetworkd init $NODENAME --home ~/.commercionetwork
# If you get a error because .commercionetwork folder is present don't worry 
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


:::warning  
At this point there may be some differences if you are using `KMS` with `HSM`.
:::

If you are using a kms you have to configure the node so that it can accept the connection from the sign provider.         
Usually the connection is within a private lan network, or via vpn. Assuming that the validator node has an internal ip such as 10.1.1.1 the configuration to set up will be


```toml
priv_validator_laddr = "tcp://10.1.1.1:26658"
```

Edit `~/.commercionetwork/config/config.toml`, search `priv_validator_laddr`, and modify the configuration.

**WARN**: the configuration of KMS must already be done. 


## 4. Set the statesync configuration

Download the statesyync and extract the data folder
```
wget "https://quicksync.commercio.network/$(echo $CHAINID)-statesync.latest.tgz" -P ~/.commercionetwork/
# Check if the checksum matches the one present inside https://quicksync.commercio.network
cd ~/.commercionetwork/
tar -zxf $(echo $CHAINID).latest.tgz
rm $(echo $CHAINID)-statesync.latest.tgz
```

## 5. Install and config cosmovisor

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

Environment="DAEMON_NAME=commercionetworkd"
Environment="DAEMON_HOME=/root/.commercionetwork"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_LOG_BUFFER_SIZE=512"
Environment="UNSAFE_SKIP_BACKUP=true"

ExecStart=/root/go/bin/cosmovisor run start --home="/root/.commercionetwork" 

[Install]
WantedBy=multi-user.target
EOF
```


Set env of cosmovisor for you convenience

```bash
echo 'export DAEMON_NAME=commercionetworkd' >> ~/.profile
echo 'export DAEMON_HOME=/root/.commercionetwork' >> ~/.profile
echo 'export DAEMON_RESTART_AFTER_UPGRADE=true' >> ~/.profile
echo 'export DAEMON_ALLOW_DOWNLOAD_BINARIES=false' >> ~/.profile
echo 'export DAEMON_LOG_BUFFER_SIZE=512' >> ~/.profile
echo 'export UNSAFE_SKIP_BACKUP=true' >> ~/.profile
source ~/.profile
```

## 6. Start service and sync the node


Now you can start your full node. Enable the newly created service and try to start it using:
```bash
systemctl enable commercionetworkd  
systemctl start commercionetworkd
```

Check if the sync has been started. Use `Ctrl + C` to interrupt the `journalctl` command:
```bash
journalctl -u commercionetworkd -f
```
