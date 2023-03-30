<!--
order: 3
-->

# State

The `contract` keeps state of the address whitelist, the gov module address and the ibc module address.


## Type definitions

### Whitelist
The whitelist is an item that contains the list of whitelisted addresses.
 ```rust
pub const ADDRS_WHITELIST: Item<Whitelist> = Item::new("addrs_whitelist");

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq, JsonSchema)]
pub struct Whitelist {
    pub wl: Vec<Addr>,
}
 ```
### Gov module and Ibc module 
```rust
pub const GOVMODULE: Item<Addr> = Item::new("gov_module");

pub const IBCMODULE: Item<Addr> = Item::new("ibc_module");
```

This addresses are set only once and there is no command to change them after the contract is instantiated.