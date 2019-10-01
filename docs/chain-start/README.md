# Starting a local chain
Inside the following page you will learn how to start a new Commercio.network chain that might be useful to you 
in order to perform some tests without connecting to the testnet or mainnet. 

## Installation
In order to start a local test chain you will need to install the latest `cnd` and `cncli` binaries. 
To do so, please execute the following commands:

```bash
git clone https://github.com/commercionetwork/commercionetwork
cd commercionetwork
make install
``` 

The output should look like the following: 

```
GO111MODULE=on go install -tags " ledger" ./cmd/cnd
GO111MODULE=on go install -tags " ledger" ./cmd/cncli
``` 

Now, you should be able to execute the following command: 

```
cnd version
```

If the version number is printed properly, you are ready to go.

## Chain starting
In order to start the chain, the following steps must be performed: 

1. [Resetting previous instances](#1-resetting-previous-instances)
2. [Init a new chain](#2-init-a-new-chain)
3. [Setup the genesis data](#3-setup-the-genesis-data)
4. [Collect the genesis transactions](#4-collect-the-genesis-transactions)
5. [Start the chain](#5-start-the-chain)

### 1. Resetting previous instances
In order to start a chain without any problem, you will need to reset everything. 
To do so, execute the following commands: 

```bash
rm -r ~/.cnd
cnd unsafe-reset-all
```

:::warning  
This will remove all the previous chain data so please make sure to backup 
the `~/.cnd` folder just in case you need the data back later.   
:::

### 2. Init a new chain
To initialize a new chain, please execute the following command: 

```
cnd init testchain --overwrite
``` 

### 3. Setup the genesis data
Now that you have initialized the new chain, you need to set some genesis data.  
To do so we will use some commands that require you to have a local account key name and password. 
If you haven't create one yet, please do it know by executing

```bash
cncli keys add jack
``` 

After this command, please insert a password that will be later used.

:::warning  
While creating a local key please use a password that you will remember easily as it will be used
often later during the procedure.  
:::

The output to the previous command should look something like the following:

```
- name: jack
  type: local
  address: did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke
  pubkey: did:com:pub1addwnpepqgkyyqvz2e3um89luc34wt4rlhv63jlgky6eyvc4x57ee8hngl8z2h3d3zn
  mnemonic: ""
  threshold: 0
  pubkeys: []


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

middle vanish genre gorilla label insane east need holiday fringe odor ice permit hen art benefit amazing worry evoke pigeon project van unfold fantasy
```

Once you have create a local key, you can execute the following commands: 

```shell
# Add some funds to the account
cnd add-genesis-account $(cncli keys show jack --address) 10000000000000ucommercio

# Set the account to be the government
cnd set-genesis-government-address $(cncli keys show jack --address)

# Set the initial TBR pool amount
cnd set-genesis-tbr-pool-amount 10000000000ucommercio

# Optional - Set the account to be a membership minter
cnd add-genesis-membership-minter $(cncli keys show jack --address)
```

After executing those commands, make sure your genesis file is valid by executing

```shell
cnd validate-genesis
```

This should output something similar to the following text:

```
validating genesis file at /home/user/.cnd/config/genesis.json
File at /home/user/.cnd/config/genesis.json is a valid genesis file
```

### 4. Collect the genesis transactions
Once you've setup the genesis file, you can create the genesis transaction and collect it.
To do so, please run

```shell
cnd gentx --name jack --amount 100000000ucommercio
cnd collect-gentxs
``` 

### 5. Start the chain
Once all the genesis transactions have been created, you can start the chain by running

```shell
cnd start
``` 

You should now be able to see an output that looks something like the following:

```
I[2019-09-19|10:26:06.651] Starting ABCI with Tendermint                module=main 
I[2019-09-19|10:26:12.034] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
I[2019-09-19|10:26:12.046] Committed state                              module=state height=1 txs=0 appHash=522AF70477C8C53361489DB2D592BF66C37E76C52A42DC7AE8230AD76EF3B54F
I[2019-09-19|10:26:17.128] Executed block                               module=state height=2 validTxs=0 invalidTxs=0
I[2019-09-19|10:26:17.140] Committed state                              module=state height=2 txs=0 appHash=8BD8E4D3D66A60C37B1AE721E2C7B259C36A65209575A548CB4D09BEF0B0E42E
...
```