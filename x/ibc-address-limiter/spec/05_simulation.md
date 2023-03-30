<!--
order: 5
-->

# Simulation

This document aims to illustrate the stack creation or setup of two chains with `ignite` for testing the behaviour we may expect from the module and the contract.

## Setup chains
For this simulation we need two running chains connected by ignite relayer. We will use two local chains. Chain A will be configured by `config.yml` and chain B by `commercio2.yml`.

### Run chains
Start chain A on a terminal window
```bash
ignite chain serve -r 
```
endpoints: 
  - Tendermint node: http://localhost:26657
  - API: http://localhost:1317
  - faucet: http://localhost:4500

On another terminal window, start chain B
```bash
ignite chain serve -c commercio2.yml -r 
```
endpoints: 
  - Tendermint node: http://localhost:26659
  - API: http://localhost:1318
  - faucet: http://localhost:4501

### Relay
On a new terminal window, remove existing relayers and configure a new one.

```bash
rm -rf ~/.ignite/relayer

ignite relayer configure -a   \
--source-rpc "http://127.0.0.1:26657"   \
--source-faucet "http://127.0.0.1:4500"  \
--source-port "transfer"   \
--source-version "ics20-1"  \
--source-gasprice "0.01ucommercio"   \
--source-prefix "did:com:" \
--source-gaslimit 300000   \
--target-rpc "http://127.0.0.1:26659"   \
--target-faucet "http://127.0.0.1:4501"  \
--target-port "transfer"   \
--target-version "ics20-1"  \
--target-gasprice "0.01ucommercio"   \
--target-prefix "did:com:"   \
--target-gaslimit 300000
```
When prompted, accept the default values for Source Account and Target Account.

>If this generates the same account address for Source and target, and this would lead to some errors afterwards, it's possible to define an account by: 
>```bash
>ignite account create [name] 
>```
>Then when the terminal prompt for the account values one chain can use the default value while the other one use the previously defined account, typing it before hitting enter.

### Connect
```bash
ignite relayer connect
```
This connect the two chains and will display the ports and channels of the connection.

## Contract
The contract implement the middleware logic that permits to contol the set of address that can send token via IBC transfer, so it must be stored and instantiated on the chain. 

Since there is no custom logic for receiving tokens from other chains, it's not necessary to set the contract in both chain, one is sufficient.

### Compile & Store 
A .wasm of the contract is already available in the [`bytecode`](./bytecode/) folder. If you want to use, it sufficient to skip the compilation and just change the path of the local variable.

```bash
RUSTFLAGS="-C link-arg=-s" cargo build --release --target=wasm32-unknown-unknown --locked

WALLET_CREATOR="did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"
ADDRESS_LIMITER="$HOME/commercionetwork/x/ibc-address-limiter/target/wasm32-unknown-unknown/release/address_limiter.wasm"

commercionetworkd tx wasm store $ADDRESS_LIMITER --from $WALLET_CREATOR --fees 100000000ucommercio -o text  --gas 50000000 -y
```
Save the contract ID, for later instantition. The ID can be retrieved from the output of the previous Store.
```bash
LIMITER_CODE_ID=1
```
### Instantiate
In this step we have to instantiate the contract passing the ibc module address, the gov module address and the actual whitelist of permissioned addresses.

In this case we set only one address in the whitelist.
```bash
INIT='{"gov_module": "did:com:10d07y265gmmuvt4z0w9aw880jnsr700jdfgwkq","ibc_module": "did:com:1yl6hdjhmkf37639730gffanpzndzdpmhe5dzhs","addrs_whitelist": ["did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"]}'

commercionetworkd tx wasm instantiate $LIMITER_CODE_ID "$INIT" --label "IBC ADDRESS LIMITER" --admin $WALLET_CREATOR   --from $WALLET_CREATOR   --fees 10000ucommercio -o text --gas 50000000 -y
```
Save the contract address for setting module params and futher operations.

```bash
LIMITER_ADDRESS="did:com:14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sh7yll8"
```

### Set Params
As the contract to be used by the module it must be set as params, we submit a param change proposal and vote.

N.B. If the chain has multiple instantiation of this contract, only the one whose address is set as param of the module will be used and the others will be ignored.

```bash
commercionetworkd tx gov submit-proposal param-change $HOME/proposal.json --from gov --fees 100000000ucommercio -o text  --gas 50000000 -y

commercionetworkd tx gov vote 1 yes --from gov -y
commercionetworkd tx gov vote 1 yes --from alice -y
```

proposal.json is:
```json
{
  "title": "Ibc contract address",
  "description": "Add contract address limiter",
  "changes": [
    {
      "subspace": "address-limited-ibc",
      "key": "contract",
      "value": "did:com:14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sh7yll8"
    }
  ],
  "deposit": "10000000ucommercio"
}	
```

## Transfer Tokens

At this stage, the contract is instantiated and set as param in chain A. The whitelist contains only one the address `"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"` and for simplicity it will be `alice`. Chain A and chain B are relayed and connected by ignite.

On new terminal window, alice transfer 200000ucommercio to a user in chain B:
```bash
commercionetworkd tx ibc-transfer \
transfer \
transfer channel-0 did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd 200000ucommercio \
--from did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd \
--fees 100000000ucommercio -o text  \
--gas 50000000 -y
```

Verify on chain B
```bash
commercionetworkd query bank balances did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd --node tcp://localhost:26659
```
Response:
```bash
balances:
- amount: "200000"
  denom: ibc/A0D0A020FC536B3F6E0F36B036E3C1941035D3428C7B76D4C2A55D0B2277C7E9
- amount: "20000000000"
  denom: uccc
- amount: "199899980000"
  denom: ucommercio
pagination:
  next_key: null
  total: "0"
```

Someone who has not the permission, try transfer:
```bash
commercionetworkd tx ibc-transfer \
transfer \
transfer channel-0 did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd 200000ucommercio \
--from did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn \
--fees 100000000ucommercio -o text  \
--gas 50000000 -y
```

Error response:
```bash
failed to execute message; message index: 0: Unauthorized! The sender did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn
  has not authorization: execute wasm contract failed: unauthorized
```

Add address to the whitelist and retry.

```bash
ADD_ADDRESS='{"manage_whitelist":{"add_addrs":{"addresses":["did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn"]}}}'

commercionetworkd tx gov submit-proposal execute-contract \
$LIMITER_ADDRESS "$ADD_ADDRESS" \
--title "Ibc contract execute" \
--description "Execute contract"  \
--deposit "10000000ucommercio" \
--from gov \
--fees 10000ucommercio -o text \
--gas 50000000 \
--run-as did:com:10d07y265gmmuvt4z0w9aw880jnsr700jdfgwkq -y

commercionetworkd query wasm contract-state smart $LIMITER_ADDRESS '{"get_whitelist":{}}'
```

Query response:
```bash
data:
  wl:
  - did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd
  - did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn
```

Transfer and verify on chain B
```bash
commercionetworkd tx ibc-transfer \
transfer \
transfer channel-0 did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd 200000ucommercio \
--from did:com:1829s409tjju2luhudq5dfeus6je3vfdnjv9tpn \
--fees 100000000ucommercio -o text  \
--gas 50000000 -y

commercionetworkd query bank balances did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd --node tcp://localhost:26659
```

Response:
```bash
balances:
- amount: "400000"
  denom: ibc/EF505ADD9D17C91A013DC23BD975CEB7AA3CC118052940A20A3D623A75658B4C
- amount: "20000000000"
  denom: uccc
- amount: "199900000000"
  denom: ucommercio
pagination:
  next_key: null
  total: "0"
```

Remove alice from the whitelist.
```bash
REMOVE_ADDR='{"manage_whitelist":{"remove_addrs":{"addresses":["did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"]}}}'

commercionetworkd tx gov submit-proposal execute-contract \
$LIMITER_ADDRESS "$REMOVE_ADDR" \
--title "remove address" \
--description "propose remove address" \
--deposit "10000000ucommercio" \
--fees 10000ucommercio -o text \
--from gov
--gas 50000000 \
--run-as did:com:10d07y265gmmuvt4z0w9aw880jnsr700jdfgwkq -y
```

Alice has no more permission to transfer tokens.
```bash
commercionetworkd tx ibc-transfer \
transfer \
transfer channel-0 did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd 200000ucommercio \
--from did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd \
--fees 100000000ucommercio -o text  \
--gas 50000000 -y
```

Error response:
```bash
failed to execute message; message index: 0: Unauthorized! The sender did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd has not authorization: execute wasm contract failed: unauthorized
```

Empty whitelist to allow everyone transfer tokens.
```bash
EMPTY_WL='{"manage_whitelist":{"reset_whitelist":{}}}'

commercionetworkd tx gov submit-proposal execute-contract \
$LIMITER_ADDRESS "$EMPTY_WL" \
--title "Empty whitelist" \
--description "propose empty wl" \
--deposit "10000000ucommercio" \
--fees 10000ucommercio -o text \
--from gov
--gas 50000000 \
--run-as did:com:10d07y265gmmuvt4z0w9aw880jnsr700jdfgwkq -y

commercionetworkd tx ibc-transfer \
transfer \
transfer channel-0 did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd 200000ucommercio \
--from did:com:1zg4jreq2g57s4efrl7wnh2swtrz3jt9nfaumcm \
--fees 100000000ucommercio -o text  \
--gas 50000000 -y
```

Check balance on chain B
```bash
commercionetworkd query bank balances did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd --node tcp://localhost:26659
```

Response:
```bash
balances:
- amount: "600000"
  denom: ibc/EF505ADD9D17C91A013DC23BD975CEB7AA3CC118052940A20A3D623A75658B4C
- amount: "20000000000"
  denom: uccc
- amount: "199900000000"
  denom: ucommercio
pagination:
  next_key: null
  total: "0"
```
