# Join the mainnet

> See the [chains repo](https://github.com/commercionetwork/chains) for information on the mainnet, including the 
details about the genesis file. 

:::warning 
Please make you you have [installed commercio.network](./installation.md) before you go further.
::: 

## Setting up a new node
These instructions are for setting up a brand new full node from scratch. 

First, initialize the node and create the necessary config file: 

```bash
cnd init <your_custom_moniker>
```

> Note. Monikers can contains only ASCII characters. Using Unicode characters will render your node unreachable.

You can edit this `moniker` later, in the `~/.cnd/config/config.toml` file: 

```toml
# A custom human readable name for this node
moniker = "<your_custom_moniker>"
```

You can edit the `~/.cnd/config/cnd.toml` file in order to enable the anti spam mechanism and reject incoming 
transactions with less than the minimum gas prices: 

```toml
# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### main base config options #####

# The minimum gas prices a validator is willing to accept for processing a
# transaction. A transaction's fees must meet the minimum of any denomination
# specified in this config (e.g. 10uatom).

minimum-gas-prices = ""
```

You full-node has been initialized! 

## Genesis & Seeds
### Copy the Genesis File
Fetch the mainnet's `genesis.json` file into the `cnd`'s config directory. 

```bash
mkdir -p $HOME/.cnd/config
curl https://raw.githubusercontent.com/commercionetwork/launch/master/genesis.json > $HOME/.cnd/config/genesis.json
```

Note that we use the `latest` directory inside the [chains repo](https://github.com/commercionetwork/chains) which 
contains details for the mainnet like the latest version and genesis file. 

> If you want to connect to the public testnet instead, click [here](./join-testnet.md).

To verify the correctness of the configuration run: 

```bash
cnd start
```

### Add Seed Nodes
Your node needs to know how to find peers. You'll need to add healthy seed nodes to `$HOME/.cnd/config/config.toml`. 
The [chains repo](https://github.com/commercionetwork/chains) contains the links to some seed nodes.

If those seeds aren't working, you can ask for peers on the [Commercio.network Telegram Group](https://t.me/CommercioNetwork).

For more information on seeds and peers, you can [read this](https://github.com/tendermint/tendermint/blob/develop/docs/tendermint-core/using-tendermint.md#peers).

## A Note on Gas and Fees
:::warning 
On Commercio.network mainnet, the accepted denom is ucommercio, where 1commercio = 1.000.000ucommercio
:::

Transactions on the Cosmos Hub network need to include a transaction fee in order to be processed. 
This fee pays for the gas required to run the transaction. The formula is the following:

```
fees = ceil(gas * gasPrices)
```

The `gas` is dependent on the transaction. Different transaction require different amount of gas. 
The `gas` amount for a transaction is calculated as it is being processed, but there is a way to estimate it beforehand 
by using the `auto` value for the `gas` flag. Of course, this only gives an estimate. You can adjust this estimate with the 
flag `--gas-adjustment` (default `1.0`) if you want to be sure you provide enough `gas` for the transaction.

The `gasPrice` is the price of each unit of gas. Each validator sets a `min-gas-price` value, and will only include 
transactions that have a `gasPrice` greater than their `min-gas-price`.

The transaction `fees` are the product of `gas` and `gasPrice`. As a user, you have to input 2 out of 3. 
The higher the `gasPrice/fees`, the higher the chance that your transaction will get included in a block.

:::tip
For the mainnet, the recommended gas-prices is `0.025ucommercio`.
:::

## Set `minimum-gas-prices`
Your full-node keeps unconfirmed transactions in its mempool. In order to protect it from spam, it is better to set a 
`minimum-gas-prices` that the transaction must meet in order to be accepted in your node's mempool. 
This parameter can be set inside `~/.cnd/config/gaiad.toml`.

The initial recommended `min-gas-prices` is `0.025ucommercio`, but you might want to change it later.

## Run a full node
Start the full node with this command: 

```bash
cnd start
```

Check that everything is running smoothly: 

```bash
cncli status
```

View the status of the network with the [Commercio.network Explorer](https://explorer.commercio.network).

## Export the current state
Gaia can dump the entire application state to a JSON file, which could be useful for manual analysis and can also 
be used as the genesis file of a new network.

Export the state with:
```bash
cnd export > [filename].json
```

You can also export the state from a particular height (at the end of processing the block of that height):

```bash
cnd export --height [height] > [filename].json
```

If you plan to start a new network from the exported state, export with the `--for-zero-height` flag:

```bash
cnd export --height [height] --for-zero-height > [filename].json
```

## Verify the mainnet state
Verifying the mainnet state helps to prevent a catastrophe by running invariants on each block on your full node.
In essence, by running invariants you ensure that the state of mainnet is the correct expected state. 
One vital invariant check is that no atoms are being created or destroyed outside of expected protocol, however there 
are many other invariant checks each unique to their respective module. 
Because invariant checks are computationally expensive, they are not enabled by default. 
To run a node with these checks start your node with the `assert-invariants-blockly` flag:

```bash
cnd start --assert-invariants-blockly
```

If an invariant is broken on your node, your node will panic and prompt you to send a transaction which will halt 
mainnet. For example the provided message may look like:

```bash
invariant broken:
    loose token invariance:
        pool.NotBondedTokens: 100
        sum of account tokens: 101
    CRITICAL please submit the following transaction:
        gaiacli tx crisis invariant-broken staking supply
```

When submitting a invariant-broken transaction, transaction fee tokens are not deducted as the blockchain will halt 
(aka. this is a free transaction).

## Upgrade to Validator Node
You now have an active full node. What's the next step? You can upgrade your full node to become a 
Commercio.network Validator. 
The top 100 validators have the ability to propose new blocks to Commercio.network. 
Continue onto the [Validator Setup](./validator-setup.md).