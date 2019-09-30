# Trusted metadata scheme proposers

## Officially recognized metadata schemes
When [sharing a document](tx/send-document.md) with a user, you can decide to adopt a metadata schema 
that is officially recognized inside the chain. 

Those schemes are composed of few information and are proposed by **trusted proposers**. 

## Trusted proposers
*Trusted metadata scheme proposers* are specific addresses that are recognized as 
trusted by the [chain government](../government/README.md).  

These addresses can identify either companies or individuals, and are allowed to propose new metadata
schemes whenever they wish to do so, without asking any permission to anyone. 

### Adding a new proposer
If you represent the chain government, you can add a new trusted schema proposer
using the [appropriate procedure](./tx/add-trusted-metadata-schema-proposer.md).
  
Before doing so, please note that once added, you cannot remove a trusted proposer without halting the chain. 
This means that you should pay a lot of attention when deciding who is going to be a proposer, as he will be later 
allowed to propose any new metadata scheme that he wants without any further approval.     
