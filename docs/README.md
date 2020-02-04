# Commercio.network Documentation



## What is `cn`
`cn` is the name of the Commercio.network application for the [Cosmos Hub](https://hub.cosmos.network/). It is shipped
with two different entrypoints: 

* `cnd`: The Commercio.network Daemon, runs a full-node of the `cn` application
* `cndcli`: The Commercio.network command-line interface, which enables interaction with a Commercio.network full-node.

`cn` is built on top of the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) using the following modules:

* `x/auth`: Accounts and signatures.
* `x/bank`: Token transfers.
* `x/staking`: Staking logic.
* `x/mint`: Inflation logic.
* `x/distribution`: Fee distribution logic.
* `x/slashing`: Slashing logic.
* `x/gov`: Governance logic.
* `x/params`: Handles app-level parameters.

A part from these modules, `cn` comes with the following custom modules: 

* [`x/docs`](x/docs/README.md): Documents storing and sharing. 
* [`x/id`](x/id/README.md): Pseudonymous identities creation.
* [`x/government`](x/government/README.md): On-chain government. 

## Nodes
If you wish to learn about the different node types that are present inside the Commercio.network chain or you 
wish to setup a new node, please refer to our [nodes section](nodes/README.md).  

## App Developers
If you're an App  developer and would like to integrate to Commercio.network, please refer to our 
[App developers guide](app_developers/README.md). 


## SDK Developers
If you're a developer and would like to help us build a SDK in your favourite language please refer to our 
[SDK developers guide](developers/README.md). 
