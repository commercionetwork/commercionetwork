# Becoming a validor in the Mainnet

To becoming a validator in the Mainnet is quite different against testnet



## Requirements


## Undestanding command and folder structure

When you create a full node for the first time you will use this command

```bash
commercionetworkd init
```

This command create by default a folder in home dir that named `.commercionetwork`.    
This folder contains two sub folder 

* config: contains configurations like genesis file, toml file, a file that contains private keys and node info.
* data: contains database with the status of the chain
  
If use command with specific flag you can create the commercionetworkd home folder in other place. 

```bash
commercionetworkd init --home /home/commercionetworkd-user/commercionetworkd
```

So you can choose a specific disc to register all data of the chain
You can stop you service and move all your data to another place and modify the service script to update with new position, or create with simbolic link. For example

```bash
systemctl commercionetworkd stop
sleep 7 #<-- wait complete stop the service
mv /home/commercionetworkd-user/commercionetworkd /mnt/largedisk/.
cd /home/commercionetworkd-user/
ln -s /mnt/largedisk/commercionetworkd .
systemctl commercionetworkd start
```




## Undestanding configurations

When you create commercionetworkd home folder you get sub folder config.
In that folder there are some file

* genesis.json: the main file that all nodes share, and the manifest of the chain
* config.toml: a toml file that contains a set of parameters for your specific node
* app.toml: contains a small set of parameters to interact with your applications
* priv_validator_keys.json: contains address, public key and private key of your full node. **ATTENTION**: 
* node_key.json: contains informations about your node

### genesis.json

Genesis file is provided from the consortium at first time and, after upgrade of chain, will be producted from the state of the chain.
The most important parameters of this file are
* genesis_time: show when the chain should be started
* chain_id: the id of chain.
  


### config.toml

Any node can configures a set of parameters to get own needs.
The most important parameters of this file are

* moniker: define the unique name of node
* persistant_peers: define a set of 


## How to configure sentry node



## How to configure validator node



## 