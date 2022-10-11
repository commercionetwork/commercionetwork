# Updating a validator/chain

:::danger  
If you are a new validator you need follow the [*"Becoming a validator"* procedure](validator-node-installation.md).   
**DO NOT USE THIS UPDATE PROCEDURES**  
:::    
      
This section describes the procedure that needs to be followed in order to update a validator node from one 
version to another.

## Upgrade minor release



## Upgrade via Proposal

The main way to upgrade chain and his node for a major release is via Upgrade Proposal.     
Everytime a new major release will be released a Upgrade Proposol will be submited in the chain and all users of chain with stake will be able to vote it.     
Every proposal will be show in the explorer in follow page

https://mainnet.commercio.network/proposals/

Clicking one proposal you can access to the proposal details and you can vote it with keplr extension.

Another way to vote is using `cli` interface of `commercionetworkd` program

```bash
commercionetworkd tx gov vote [id proposl] [option vote]
```

You can vote `yes`, `no`, `abstain`, `nowithveto`.     
Upgrade Proposal contains the details of upgrade, the halt block and the version of software.    

After the proposal accepted the chain will halt at the block of proposal and wait the new version of the software.     





**(WIP)**