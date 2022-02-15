# Commercio.network nodes
Commercio.network blockchain has two different nodes types: *full* nodes and *validator* nodes.  
Inside this page you will learn the difference between them, which role they play inside the whole blockchain, 
and how you can create either one of them. 

:::tip Contribute

If you have questions or comments, please create an [issue](https://github.com/commercionetwork/commercionetwork/issues).
This is a community guide  any contribute is important! Please help us to improve by opening an [issue](https://github.com/commercionetwork/commercionetwork/issues).

:::


## Full nodes
**Full** nodes are nodes that simply store the whole transactions history. To do so, they connect to the blockchain and
each time a new transaction is made, they write it on the hard disk. This means that being a full node you will be able 
to read the whole chain transaction from its beginning, but you will need to have a large store space if you want to 
keep it running for a long period of time. 

:::tip Become a full node  
If you are interested in becoming a full node, you can read the 
[*full node installation* guide](full-node-installation.md)
:::
      
## Validator nodes
**Validator** nodes are simply [full nodes](#full-nodes) that also have the ability of validating new transactions
that should be added to the chain. In order to do so, they possess a private key with which they sign the transactions 
marking them as valid. In exchange of their work, they get back a *reward* that is given to them each time a new block 
is created.  
Without validator nodes, the whole chain couldn't exist and by becoming one of them, you will contribute to the working
of the whole Commercio.network ecosystem. 

:::tip Become a validator node  
If you wish to become a validator node and start earning rewards you can read the
[*validator node installation* guide](validator-node-installation.md)
:::   
<!--
3.0 ok f1
-->