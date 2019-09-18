# Update chain
**IMPORTANT BEFORE YOU START**: If you are a new validator you need follow ["Getting Started"](##getting-started) procedure. __DON'T USE THESE UPDATE PROCEDURES__    
      
This section describe the procedures to update chain from a version to another.    
Every update have a specific produre type.   
First type will be used only for testnet chains, while the second will be used to update mainnet chain.  
Every chain starting `commercio-testnet3000` contains type procedure adopted.   
In     

`https://github.com/commercionetwork/chains/blob/master/commercio-<chain-version>/.data`      

You should find `Update type`.


## 1. Update with "getting started" procedure
This type is similar to the "getting started" procedure.   
You need to delete or move the `~/.cnd` folder and start as fresh.    
You can mantain your wallet that is installed in `~/.cncli` folder or in your `ledger` device, or recreate it with mnemonic.   
You can create new wallet if you prefered and use a new account to become a validator.   

```bash
systemctl stop cnd
pkill cnd #We want be sure that chain process was stopped ;)
```

Delete `~/.cnd` folder

```bash
rm -rf ~/.cnd
```

or move it (if you want keep the old testnet state for your porpouses). 
Use `<previous-chain-version>` name for copy name for example

```bash
cp -r ~/.cnd ~/.cnd.<previous-chain-version>
```

Now you can start follow "getting started" procedure.
**ATTENTION**: before go haed with "getting started" procedure read follow changes about some steps

In [`step 1`](##1-Setup) in order to update the OS so that you can work properly, execute the following commands:

```bash
apt update && sudo apt upgrade -y
snap refresh --classic go # You need to update golang to last version
```


 In [`step 4`](##4-install-binaries-genesis-file-and-setup-configuration) you don't need to change the follow rows of your `~/.profile` file

```
export GOPATH="\$HOME/go"
export PATH="\$GOPATH/bin:\$PATH"
```

You need clean up your file from previous chain configurations

```bash
sed -i \
 -e '/export NODENAME=.*/d' \
 -e '/export CHAINID=.*/d' ~/.profile
```

and add new chain configs

```bash
export NODENAME="<your-moniker>"
export CHAINID=commercio-$(cat .data | grep -oP 'Name\s+\K\S+')

cat <<EOF >> ~/.profile
export NODENAME="$NODENAME"
export CHAINID="$CHAINID"
EOF
```

   

## 2. Update with cnd commands **(WIP 🛠)**
This procedure will be applied to mainnet chain and to some specific testnet update

### A. Preliminary/Risks **(WIP 🛠)**
To this type of procedure will be assigned a height of block, informations about checksum of geneis file and software version and a deadline expressed in UTC format.    
There is some risks about double signature: to avoid every sort of risks verify software and hash of `genesis.json` and specific configuration in `config.toml`.
The deadline of update must be respected: every validator that will not update just in time will be slashed.

### B. Recovery **(WIP 🛠)**
Is recommended to take a full data snapshot at the export height before update.   
This procedure is quite simple using commands below

```bash
systemctl stop cnd
cp -r ~/.cnd ~/.cnd.[OLD VERSION]
cp -r ~/.cncli ~/.cncli.[OLD VERSION]
# Save binaries also
cp -r /root/go/bin/cnd /root/go/bin/cnd.[OLD VERSION]
cp -r /root/go/bin/cncli /root/go/bin/cncli.[OLD VERSION]
```
### C. Upgrade Procedure **(WIP 🛠)**

```bash
rm -rf commercio-chains
mkdir commercio-chains && cd commercio-chains
git clone https://github.com/commercionetwork/chains.git .
cd commercio-<chain-version> # eg. cd commercio-testnet1001 
```

Compile binaries 

```bash
git init . 
git remote add origin https://github.com/commercionetwork/commercionetwork.git
git pull
git checkout tags/$(cat .data | grep -oP 'Release\s+\K\S+')
make install
```

Test if you have the correct binaries version:

```bash
cnd version
# Should output the same version written inside the .data file
```

Get height from update info

```bash
export BLOCKHEIGHT=$(cat .data | grep -oP 'Height\s+\K\S+')
```

Export state from 