# Updating a validator

:::danger  
If you are a new validator you need follow the [*"Becoming a validator"* procedure](validator-node-installation.md).   
**DO NOT USE THIS UPDATE PROCEDURES**  
:::    
      
This section describes the procedure that needs to be followed in order to update a validator node from one 
version to another.




For updating will be sent to the validators the type and a brief explanation of the steps to be taken.    
This guide contains only a general basic explanation of the procedures to be applied.     

In general there are 3 types of updates

* Update type [`1 - Hard fork`](#update-type-1-hard-fork)
* Update type [`2 - Soft fork`](#update-type-2-soft-fork)
* Update type [`3 - Reset/reinstall chain`](#update-type-3-reset-reinstall-chain)


## Update type 1 - Hard fork

In this case, the chain state will be preserved from its beginning to a certain point in time.  

### Preliminary/Risks
When signalling a required update that should follow this procedure, the following information will 
be let known to all validators:

1. A specific block height to stop chain: `<stop_block_height>` 
2. The new chain software version: `<new_version>` 
3. A deadline expressed in UTC format
4. A new genesis_time in UTC format (something like 2020-06-12T08:55:00Z): `<genesis_time>` 
5. The new chain_id: `<new_chain_id>` 
6. The old chain_id: `<old_chain_id>` 

If you are a validator, please make sure that you know all those information before proceeding with the update.

:::tip
You should reach [Setup height to stop the chain](#setup-height-to-stop-the-chain) point before update event.   
You prepare all settings as soon as you become aware of the update data
:::
    
:::danger Double signing 
Due to the nature of the update operations, there is some risk of double signature. 
To avoid every sort of risks please verify the software version, the hash of the `genesis.json` file and the specific
configuration present inside the `config.toml` file.  
::: 

:::danger 
Update time  
The deadline of the update must be respected: every validator that will not update just in time will be slashed as offline node.  
:::  

### Install support software

```bash
sudo apt install jq -y
```


### Setup environment
Each validator has its own specific environment and therefore will have to modify the variables according to its setup.   
**The following setups are general guidelines**



```bash
cd $HOME
export UPDATE_PROFILE="$HOME/.profile_update_cnd"
echo 'export HOME_CND="'$HOME'/.cnd"' > $UPDATE_PROFILE
echo 'export HOME_CND_CONFIG="$HOME_CND/config"' >> $UPDATE_PROFILE
echo 'export HOME_CND_DATA="$HOME_CND/data"' >> $UPDATE_PROFILE
echo 'export HOME_CNCLI="'$HOME'/.cncli"' > $UPDATE_PROFILE
echo 'export APP_TOML="$HOME_CND_CONFIG/app.toml"' >> $UPDATE_PROFILE
echo 'export BIN_DIR="'$HOME'/go/bin"' >> $UPDATE_PROFILE
echo 'export SRC_GIT_DIR="'$HOME'/commercionetwork"' >> $UPDATE_PROFILE
echo 'export BUILD_DIR="$SRC_GIT_DIR/build"' >> $UPDATE_PROFILE
echo 'export NEW_CHAIN_ID="<new_chain_id>"' >> $UPDATE_PROFILE
echo 'export CHAIN_ID="<old_chain_id>"' >> $UPDATE_PROFILE
echo 'export BUILD_VERSION="<new_version>"' >> $UPDATE_PROFILE
echo 'export ALT_BLOCK=<stop_block_height>' >> $UPDATE_PROFILE
echo 'export NEW_GENESIS_TIME="<genesis_time>"' >> $UPDATE_PROFILE

source $UPDATE_PROFILE

echo ". $UPDATE_PROFILE" >> $HOME/.profile
```

### Compile new executables

In order to properly update your validator node, first of all you need to clone the 
[`commercionetwork` repo](https://github.com/commercionetwork/commercionetwork):


```bash
git clone https://github.com/commercionetwork/commercionetwork.git $SRC_GIT_DIR
```

Once downloaded, you need to compile the binaries with correct version:

```bash
cd $SRC_GIT_DIR
git pull
git checkout $BUILD_VERSION
git pull
make build
```

Verify your version
```bash
./build/cnd version
```
should be the same value of `<new_version>`

### Setup height to stop the chain

```bash
sed -e "s|halt-height = .*|halt-height = $ALT_BLOCK|g" \
  $APP_TOML > $APP_TOML.tmp
mv $APP_TOML.tmp $APP_TOML
sudo service cnd restart
```

### Control if chain is stopped

After the block `<stop_block_height>` a chain stop message should appear at the indicated block.
To verify that the chain is actually locked check with reading the logs

```bash
sudo journalctl -u cnd -f
```

### Stop services
Make sure you have stopped services and export the chain

```bash
sudo systemctl stop cnd
sudo systemctl stop cncli # if you running rest server
sudo pkill cnd
sudo pkill cncli
```

and control it reading the logs

```bash
sudo journalctl -u cnd -f
```

### Export chain state

```bash
cnd export --for-zero-height > $HOME/export_for_$NEW_CHAIN_ID.json
```

### Backup executables
```bash
cp -r $BIN_DIR/cnd $BIN_DIR/cnd.back
cp -r $BIN_DIR/cncli $BIN_DIR/cncli.back
```

### Changing executables

```bash
cp $BUILD_DIR/cn* $BIN_DIR/.
```
### Migrate chain 

Try to find out migrate version from list of current migration

```bash
cnd migrations-list 
```

If you find `<new_version>` in the ouput command you ca use the follow command

```bash
cnd migrate \
  $BUILD_VERSION \
  $HOME/export_for_$NEW_CHAIN_ID.json \
  --chain-id $NEW_CHAIN_ID \     
  --genesis-time $NEW_GENESIS_TIME > \
  $HOME/genesis_for_$NEW_CHAIN_ID.json
  
```

If you don't find `<new_version>` in the valid migrations list use 

```bash
cat $HOME/export_for_$NEW_CHAIN_ID.json | \
  jq '.genesis_time="'$NEW_GENESIS_TIME'"' | \
  jq '.chain_id="'$NEW_CHAIN_ID'"' > \
  $HOME/genesis_for_$NEW_CHAIN_ID.json
```

### Check hash of genesis
To confirm that all validator nodes have produced the same genesis, a comparison must be made between all update participants to demonstrate the consistency of the data.

```bash
jq -S -c -M '' $HOME/genesis_for_$NEW_CHAIN_ID.json | shasum -a 256
```

Compare the value of the hash obtained with all other participants in the update


### Backup
Before starting the update it is recommended to take a full data snapshot of the chain state at the export height.     
To do so please run:

```bash
systemctl stop cnd
systemctl stop cncli # if you running rest server
mkdir $HOME_CND.$CHAIN_ID
mv $HOME_CND_DATA $HOME_CND.$CHAIN_ID/.
cp -r $HOME_CND_CONFIG $HOME_CND.$CHAIN_ID/.
cp -r $HOME_CNCLI $HOME_CNCLI.$CHAIN_ID
```


### Reset chain and change geneis file

```bash
cnd unsafe-reset-all
```


```bash
cp $HOME/genesis_for_$NEW_CHAIN_ID.json $HOME_CND_CONFIG/genesis.json
```

### Restart chain

Before restarting the chain we have to reset the value for height for the stop to zero.

```bash
sed -e "s|halt-height = .*|halt-height = 0|g" \
  $APP_TOML > $APP_TOML.tmp
mv $APP_TOML.tmp $APP_TOML
systemctl start cnd
```

### Control chain restart

The new chain should restart at the time set in the `genesis_time`.         
Control if the sync starts. Use `Ctrl + C` to interrupt the `journalctl` command    

```bash
journalctl -u cnd -f
# OUTPUT SHOULD BE LIKE BELOW AFTER GENESI TIME
#
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.722] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.728] Committed state                              module=state height=1 txs=0 appHash=9815044185EB222CE9084AA467A156DFE6B4A0B1BAAC6751DE86BB31C83C4B08
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.745] Executed block                               module=state height=2 validTxs=0 invalidTxs=0
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.751] Committed state                              module=state height=2 txs=0 appHash=96BFD9C8714A79193A7913E5F091470691B195E1E6F028BC46D6B1423F7508A5
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.771] Executed block                               module=state height=3 validTxs=0 invalidTxs=0
```


### Restart the REST API if you need
```
cncli config chain-id $NEW_CHAIN_ID
cncli rest-server --chain-id $NEW_CHAIN_ID --trust-node
``` 






## Update type 2 - Soft fork

This type of upgrade is generally optional: it is released to fix some minor bugs in the software and does not change the status of the chain


### Preliminary/Risks
When signalling a optional update that should follow this procedure, the following information will 
be let known to all validators:

1. The new chain software version: `<new_version>` 

If you are a validator, please make sure that you know all those information before proceeding with the update.



### Setup environment
Each validator has its own specific environment and therefore will have to modify the variables according to its setup.   
**The following setups are general guidelines**


```bash
cd $HOME
export UPDATE_PROFILE="$HOME/.profile_update_cnd"
echo 'export HOME_CND="'$HOME'/.cnd"' > $UPDATE_PROFILE
echo 'export HOME_CND_CONFIG="$HOME_CND/config"' >> $UPDATE_PROFILE
echo 'export HOME_CND_DATA="$HOME_CND/data"' >> $UPDATE_PROFILE
echo 'export BIN_DIR="'$HOME'/go/bin"' >> $UPDATE_PROFILE
echo 'export SRC_GIT_DIR="'$HOME'/commercionetwork"' >> $UPDATE_PROFILE
echo 'export BUILD_DIR="$SRC_GIT_DIR/build"' >> $UPDATE_PROFILE
echo 'export BUILD_VERSION="<new_version>"' >> $UPDATE_PROFILE

source $UPDATE_PROFILE

echo ". $UPDATE_PROFILE" >> $HOME/.profile
```



### Compile new executables

In order to properly update your validator node, first of all you need to clone the 
[`commercionetwork` repo](https://github.com/commercionetwork/commercionetwork):


```bash
git clone https://github.com/commercionetwork/commercionetwork.git $SRC_GIT_DIR
```

Once downloaded, you need to compile the binaries with correct version:

```bash
cd $SRC_GIT_DIR
git pull
git checkout $BUILD_VERSION
git pull
make build
```

Verify your version
```bash
./build/cnd version
```
should be the same value of `<new_version>`

### Backup executables
```bash
cp -r $BIN_DIR/cnd $BIN_DIR/cnd.back
cp -r $BIN_DIR/cncli $BIN_DIR/cncli.back
```

### Restart chain changing executables

```bash
sudo systemctl stop cnd
sudo systemctl stop cncli
cp $BUILD_DIR/cn* $BIN_DIR/.
sudo systemctl start cnd
sudo systemctl start cncli

```

### Control chain restart

The new chain should restart at the time set in the `genesis_time`.         
Control if the sync starts. Use `Ctrl + C` to interrupt the `journalctl` command    

```bash
journalctl -u cnd -f
# OUTPUT SHOULD BE LIKE BELOW AFTER GENESI TIME
#
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.722] Executed block                               module=state height=1 validTxs=0 invalidTxs=0
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.728] Committed state                              module=state height=1 txs=0 appHash=9815044185EB222CE9084AA467A156DFE6B4A0B1BAAC6751DE86BB31C83C4B08
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.745] Executed block                               module=state height=2 validTxs=0 invalidTxs=0
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.751] Committed state                              module=state height=2 txs=0 appHash=96BFD9C8714A79193A7913E5F091470691B195E1E6F028BC46D6B1423F7508A5
# Aug 13 16:30:20 commerciotestnet-node4 cnd[351]: I[2019-08-13|16:30:20.771] Executed block                               module=state height=3 validTxs=0 invalidTxs=0
```




## Update type 3 - Reset/reinstall chain
This procedure is very rare and is only performed in borderline cases.     
With this process the chain is completely reinitialized and a new genesis is produced by the chain management system.    
Usually this procedure is applied on development and test chains.    



### Backup executables
```bash
cp -r $BIN_DIR/cnd $BIN_DIR/cnd.back
cp -r $BIN_DIR/cncli $BIN_DIR/cncli.back
```

### Backup chain database
Before starting the update it is recommended to take a full data snapshot of the chain state at the export height.     
To do so please run:

```bash
systemctl stop cnd
systemctl stop cncli # if you running rest server
mkdir $HOME_CND.$CHAIN_ID
mv $HOME_CND_DATA $HOME_CND.$CHAIN_ID/.
cp -r $HOME_CND_CONFIG $HOME_CND.$CHAIN_ID/.
cp -r $HOME_CNCLI $HOME_CNCLI.$CHAIN_ID
```

### Reinstall node
Follow all steps in [Installing a full node](full-node-installation.md)      
Follow all steps in [Becoming a validator](validator-node-installation.md)
