<!--
order: 2
-->

# Contract

The contract is responsible for checking if the sender of a transfer packet has the permission to send tokens or not. To do this, the contract maintains a whitelist with authorized address that can be managed by gov proposal.

## Instantiating

Instantiate the contract taking track of the ibc module address, the gov module address and the list of authorized addresses.
Only the Gov address and the ibc module address can perform actions on the whitelist.

```rust
pub struct InstantiateMsg {
    pub gov_module: Addr,
    pub ibc_module: Addr,
    pub addrs_whitelist: Vec<Addr>,
}
```
* gov_module &rarr; address of the government module. Only this address can perform the sudo messages of the contract. It can also perform the `ManageWhitelist` command.
* ibc_module &rarr; address of the ibc module. This address can perform the `ManageWhitelist` command.
* addrs_whitelist &rarr; list of the addresses that have permission to send tokens through the IBC transfer. Addresses that are not in this list trying to send tokens over the chain ,will trigger an `unauthorized` error.
