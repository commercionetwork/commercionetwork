# Commercio.network Documentation



## What is `cn`
`cn` is the name of the [Commercio.network](https://commercio.network) blockchain application. It is shipped with two different entrypoints:  


* `cnd`: The Commercio.network Daemon, runs a full-node of the `cn` application
* `cndcli`: The Commercio.network command-line interface, which enables interaction with a Commercio.network full-node.

`cn` is built on top of the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) using the following modules:

* `x/auth`: Accounts and signatures.
* `x/bank`: Tokens transfers.
* `x/staking`: Staking logic.
* `x/distribution`: Fee distribution logic.
* `x/slashing`: Slashing logic.
* `x/params`: Handles app-level parameters.

Apart from these modules, `cn` comes with the following custom modules: 

* [`x/docs`](x/docs/README.md): Documents storing and sharing. 
* [`x/id`](x/id/README.md): Pseudonymous identities creation.
* [`x/government`](x/government/README.md): On-chain government. 
* [`x/memberships`](x/memberships/README.md): Allows you to invite new user and buy a membership. 
* [`x/commerciomint`](x/commerciomint/README.md): Handle Collateralized Debt Positions. 
* [`x/pricefeed`](x/pricefeed/README.md): Pricefeed to setup prices. 
* [`x/vbr`](x/vbr/README.md): Validator Block Rewards. 


## Running Nodes
If you wish to learn about the different node types that are present inside the Commercio.network chain or you 
wish to setup a new node, please refer to our [nodes section](nodes/README.md).  

## App Developers
If you're an App  developer and would like to integrate to Commercio.network, please refer to our 
[App developers guide](app_developers/README.md). 


## SDK Developers
If you're a developer and would like to help us build a SDK in your favourite language please refer to our 
[SDK developers guide](developers/README.md). 
