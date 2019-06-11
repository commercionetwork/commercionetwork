# Commercio.network Documentation

# What is `cn`
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

* `x/commercioauth`: Easy account management.
* `x/commerciodocs`: Documents storing and sharing. 
* `x/commercioid`: Pseudonymous identities creation.

Next, learn how to [install Commercio.network](./installation.md) 