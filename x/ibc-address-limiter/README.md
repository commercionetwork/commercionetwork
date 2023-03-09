# IBC Address Limiter

The IBC Address Limiter module is responsible for adding a governance-configurable address limit to IBC transfers.
This is intended to handle a list of addresses that have permission to send tokens via IBC transfers.

The architecture of this package is a minimal go package which implements an [IBC Middleware](https://github.com/cosmos/ibc-go/blob/f57170b1d4dd202a3c6c1c61dcf302b6a9546405/docs/ibc/middleware/develop.md) that wraps the [ICS20 transfer](https://ibc.cosmos.network/main/apps/transfer/overview.html) app, and calls into a cosmwasm contract.
All the actual IBC address limiting logic is then implemented in the cosmwasm contract. 
The Cosmwasm code can be found in the [`contracts`](./contracts/) package, with bytecode findable in the [`bytecode`](./bytecode/) folder.

## Code structure

As mentioned at the beginning, the Go code is a relatively minimal ICS 20 wrapper, that dispatches relevant calls to a cosmwasm contract that implements the address limiting functionality.

### Go Middleware

To achieve this, the middleware  needs to implement  the `porttypes.Middleware` interface and the
`porttypes.ICS4Wrapper` interface. This allows the middleware to send and receive IBC messages by wrapping 
any IBC module, and be used as an ICS4 wrapper by a transfer module (for sending packets or writing acknowledgements).

#### Parameters

The middleware uses the following parameters:

| Key             | Type   |
|-----------------|--------|
| ContractAddress | string |

1. **ContractAddress** -
   The contract address is the address of an instantiated version of the contract provided under `./contracts/`

### Cosmwasm Contract
#### Instantiating

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

The contract specifies the following messages:

#### Query
```rust
GetWhitelist { }
```
Returns the list of addresses that currently have permission to send tokens through IBC transfer 

#### Exec

ManageWhitelist(ExecWhitelist):

```rust
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

#### Sudo
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