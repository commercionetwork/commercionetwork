# Updating a validator

:::danger  
If you are a new validator you need follow the [*"Becoming a validator"* procedure](validator-node-installation.md).   
**DO NOT USE THIS UPDATE PROCEDURES**  
:::    
      
This section describes the procedure that needs to be followed in order to update a validator node from one 
version to another.

Please note that each chain version has an update type associated to it. 
In order to know which one is associated to the chain version you are currently running, please do the following:

1. Go to our [chains repo](https://github.com/commercionetwork/chains)
2. Enter the folder of you current chain's version 
3. Open the `.data` file
4. Read the  value associated to the `Update type` field.

Once you know the update type, please follow the related procedure:

* Update type [`1 - Hard fork`](#update-type-1---hard-fork)
* Update type [`2 - Soft fork`](#update-type-2---soft-fork)


## Update type 1 - Hard fork
The procedure to follow in order to upgrade the chain using this update type is very similar to the 
[installation procedure](validator-node-installation.md).

Due to this procedure being a **hard fork**, this means that any of the current chain data will be completely wiped 
and a new chain will start from scratch. However, if you want you can create a copy of the current chain state in order
to use it for testing purposes.    

### 1. Wipe the current chain state
To start the procedure we need to kill the service running the chain:

```bash
systemctl stop cnd
pkill cnd
```

Now, we need to wipe the chain state. In order to do so, you have two options:

1. Backup your data  
   `cp -r ~/.cnd ~/.cnd.back`
2. Delete the data without a backup
   `rm -rf ~/.cnd`

### 2. Start a new chain
Once you have properly cleaned up the previous chain state, you are ready to start the new chain version.   
To do so, please refer to the [full node installation guide](full-node-installation.md) but remember to apply the 
following changes to the procedure described there:

**1.** Inside the [first step](full-node-installation.md#1-installing-the-software-requirements) 
in order to update the OS so that you can work properly execute the following commands:
   
```bash
apt update && sudo apt upgrade -y
snap refresh --classic go
```

**2.** During the [4th step](full-node-installation.md#3-install-binaries-genesis-file-and-setup-configuration) 
you don't need to change the follow rows of your `~/.profile` file

```bash
export GOPATH="\$HOME/go"
export PATH="\$GOPATH/bin:\$PATH"
```

You also need to clean up the files from the previous chain configurations

```bash
sed -i \
-e '/export NODENAME=.*/d' \
-e '/export CHAINID=.*/d' ~/.profile
```

and add the new chain ones

```bash
export NODENAME="<your-moniker>"
export CHAINID=commercio-$(cat .data | grep -oP 'Name\s+\K\S+')

cat <<EOF >> ~/.profile
export NODENAME="$NODENAME"
export CHAINID="$CHAINID"
EOF
```

   

## Update type 2 - Soft fork
The second update type is the one known as **soft fork**.  
In this case, the chain state will be preserved from its beginning to a certain point in time.  

### Preliminary/Risks
When signalling a required update that should follow this procedure, the following information will 
be let known to all validators:

1. A specific block height
2. The genesis file checksum 
3. The new chain software version
4. A deadline expressed in UTC format

If you are a validator, please make sure that you know all those information before proceeding with the update.
    
:::danger Double signing 
Due to the nature of the update operations, there is some risk of double signature. 
To avoid every sort of risks please verify the software version, the hash of the `genesis.json` file and the specific
configuration present inside the `config.toml` file.  
::: 

:::danger Update time  
The deadline of the update must be respected: every validator that will not update just in time will be slashed.  
:::  

### Backup
Before starting the update it ss recommended to take a full data snapshot of the chain state at the export height.     
To do so please run:

```bash
systemctl stop cnd
cp -r ~/.cnd ~/.cnd.back
cp -r ~/.cncli ~/.cncli.back
cp -r /root/go/bin/cnd /root/go/bin/cnd.back
cp -r /root/go/bin/cncli /root/go/bin/cncli.back
```

### Update procedure

#### 1. Updating the software
In order to properly update your validator node, first of all you need to clone the 
[`chains` repo](https://github.com/commercionetwork/chains):

```bash
rm -rf commercio-chains
mkdir commercio-chains && cd commercio-chains
git clone https://github.com/commercionetwork/chains.git .
cd commercio-<chain-version> # eg. cd commercio-testnet3000 
```

Once downloaded, you need to compile the binaries:

```bash
git init . 
git remote add origin https://github.com/commercionetwork/commercionetwork.git
git pull
git checkout tags/$(cat .data | grep -oP 'Release\s+\K\S+')
make install
```

After the compilation has finished successfully, please make sure you are running the correct software version: 

```bash
cnd version
# Should output the same version written inside the .data file
```


#### 2. Export the chain state (**WIP**)
Once the software has properly been updated, we need to export the current chain state and later import it.  
In order to do so, first of all you need to get the export height from the `.data` file:

```bash
export BLOCKHEIGHT=$(cat .data | grep -oP 'Height\s+\K\S+')
```

**WIP**