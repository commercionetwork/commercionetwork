# Becoming a validor in the Mainnet

To becoming a validator in the Mainnet is quite different against mainnet



## Requirements


## Undestanding command and folder structure

When you create a full node for the first time you will use this command

```bash
cnd init moniker
```

where `moniker` is the name of your node.
This command create by default a folder in home dir that named `.cnd`.    
This folder contains two sub folder 

* config: contains configurations like genesis file, toml file, a file that contains private keys and node info.
* data: contains database with the status of the chain
  
If use command with specific flag you can create the cnd home folder in other place. 

```bash
cnd init moniker --home /home/cnd-user/cnd
```

So you can choose a specific disk to register all data of the chain
You can stop you service and move all your data to another place and modify the service script to update with new position, or create with simbolic link. For example

```bash
systemctl cnd stop
sleep 7 #<-- wait complete stop the service
mv /home/cnd-user/cnd /mnt/largedisk/.
cd /home/cnd-user/
ln -s /mnt/largedisk/cnd .
systemctl cnd start
```

When you want reset hard the state of chain you need use

```bash
systemctl cnd stop
cnd unsafe-reset-all
```

or
```bash
cnd unsafe-reset-all --home /home/cnd-user/cnd
```

if you choose other place for your cnd home.    

The command `cnd`  have a set of commands that you can retrive with `-h` flag. Every subcommand could have other subcommands or specific flags. At each level you can get informations about command and flags with `-h` flag.    
For example 

```bahs
cnd -h
```

print out all command at the first level. One of this is `init`. You can get help of `init` command using

```bash
cnd init -h
```



## Undestanding node identity

Your node has some peculiar data

* Node id
* Address
* Validator address
* Ip address

### Node id

You can get it using the command

```bash
cnd tendermint show-node-id [--home /home/cnd-user/cnd]
```

A alphanumeric string like `81dd15669daea0e6d3dacbcfdcc5ffd32b56c767` should be print. This the node id identifier of your node.
This value is suitable to used in toml configurations like `persistant_peers`, `private_peers` etc. etc.


### Address

```bash
cnd tendermint show-address [--home /home/cnd-user/cnd]
```


A alphanumeric string like `did:com:valcons1scc7lnhf3pqtwtjkcy57k9xrgp2kyu7grp24nw` should be print.

### Validator address

```bash
cnd tendermint show-validator [--home /home/cnd-user/cnd]
```


A alphanumeric string like `did:com:valconspub1zcjduepqre74dapyqd76zelkp0rxhpsc34uqdas8l64dfyzxqqzfxj5s8qwqaka2y8` should be print. This is your public key of your validator, and should be used when you performing the create validator transaction


### Ip address

Ip address of your server. Your server could have multiple ip. If your server has only one ip you can get it with command 

```bash

```







## Undestanding configurations

When you create cnd home folder you get sub folder config.
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
* seeds: define a set of 
* persistant_peers: define a set of 


## How to configure sentry node



## How to configure validator node



## 