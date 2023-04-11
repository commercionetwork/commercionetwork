<!--
order: 4
-->

# Messages

The contract exposes different types of messages: 'ExecuteMsg', 'SudoMsg' and 'QueryMsg' and for each type of message there is a different level of restriction.


## ExecuteMsg
ManageWhitelist: 

```rust
pub enum ExecuteMsg {
    ManageWhitelist(ExecWhitelist),
}

pub enum ExecWhitelist {
    AddAddrs{
        addresses: Vec<Addr>,
    },
    RemoveAddrs{
        addresses: Vec<Addr>,
    },
    ResetWhitelist{},
    New{
        whitelist: Vec<Addr>,
    }
}
```
* AddAddrs - add one or n addresses to whitelist
* RemoveAddrs - remove one or n addresses from whitelist
* ResetWhitelist - empty whitelist
* New - Remove previous list if set and create new one with given list

This messages are runnable only by `gov_module` and `ibc_module` addresses given at the instantiation stage. Execute with other addresses will produce an `Unauthorized` error.

## SudoMsg
Sudo messages can only be executed by the chain.

```rust
pub enum SudoMsg {
    SendPacket {
        packet: Packet,
    },
    RecvPacket {
        packet: Packet,
    },
}
```

* SendPacket - Check if sender is in the whitelist. If it isn't, it will return an Unauthorized error
* RecvPacket - Do not perform any special control

These sudo messages receive the packet from the chain and extract the necessary information to process the packet and determine if it should be address limited. 

## QueryMsg
```rust
GetWhitelist { }
```
Returns the list of addresses that currently have permission to send tokens through IBC transfer. This queryMsg can be executed by any addresse. 
