# Updating a validator/chain

:::danger  
If you are a new validator you need follow the [*"Becoming a validator"* procedure](validator-node-installation.md).   
**DO NOT USE THIS UPDATE INSTRUCTIONS**  
:::    
      
This section describes the procedure that needs to be followed in order to update a validator node from one 
version to another.     
The instructions below are only to decribe the general behavior when a upgrade event occours.     
Every upgrade event will be associated to the peculiar instructions and the validators will be warned through standard comunication channels (Discord, Telegram, e-mail etc. etc.)

## Upgrade minor release

When a minor release is publish the steps that the node should be performed are

1. Enter in your node and clone the repository in your server with `git clone https://github.com/commercionetwork/commercionetwork.git` and move in the new folder `cd commercionetwork`. **If you already have the repository folder use** `cd commercionetwork; git pull`
1. Checkout the code to new versione with `git checkout <version target>; git pull` where `<version target>` is the new version or tag to install
2. Compile the new versione with the command `make build`. Control if your version is equal `<version target>` using `./build/commercionetworkd version`
3. Stop the validator/node service with `systemctl stop commercionetworkd`
4. Put new version of chain in the right folder:
   1. **Cosmovisor**: copy the binary in cosmovisor current folder `cp build/commercionetworkd ~/.commercionetwork/cosmovisor/current/bin/commercionetworkd`
   2. **Old style installation**: copy the binary in `GOBIN` path. In the most cases `cp build/commercionetworkd ~/go/bin/commercionetworkd`
   3. **Other enviroment**: if you have a custom env copy the new binary replacing the one on which the chain service is based
5. Start the validator/node service with `systemctl start commercionetworkd`
6. Check the validator/node blocks issuance with `journalctl -u commercionetworkd -f`
7. Also check the validator signature blocks by looking at [Commercio Network explorer](https://mainnent.commercio.network)



## Upgrade via Proposal

The main way to upgrade chain and its nodes for a major release is via the Upgrade Proposal.     
Everytime a new major release will be released a Upgrade Proposol will be submited in the chain and all users of chain with stake will be able to vote it.     
Each proposal will be displayed in the explorer in this [page](https://mainnet.commercio.network/proposals/)

Clicking one proposal you can access to the proposal details and you can vote it with keplr extension.

Another way to vote is using `cli` interface of `commercionetworkd` program

```bash
commercionetworkd tx gov vote [id proposl] [option vote]
```

You can vote `yes`, `no`, `abstain`, `nowithveto`.     
Upgrade Proposal contains the details of upgrade, the halt block and the version of software.    

After acceptance of the proposal the chain will halt at the block of proposal and wait the new version of the software.     

Using cosmovisor you can put new software in the `upgrades` folder, and wait that the cosmovisor deamon perform all job by itself.    
An example of these procedure could be find [here](https://github.com/commercionetwork/commercio-consortium/tree/master/upgrade/3.1.0-4.0.0/en) 


## Emergency Upgrade

In some cases an emergency update may be necessary. In this case a special message will be sent to all validators, sometimes in the group channel and in other case with a 1-1 comunication.    
These upgrades is quite rare and are performed only when a real danger issue is found in the chain software.


## Dump Upgrade

This type of upgrade is a special case. In epochal cases, the chain could stop because consensus issues. Sometimes, in such cases, the only way to restart the chain is dump the whole state, correct the errors, install new software in all validators, import the state, install the new genesis, and start the chain.    
An example of these procedure could be find [here](https://github.com/commercionetwork/commercio-consortium/blob/master/upgrade/2.2.0-3.0.0/en/README.md)



**(WIP)**