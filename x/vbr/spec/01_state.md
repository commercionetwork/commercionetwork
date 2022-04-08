<!--
order: 1
-->

# State

The `vbr` module keeps state of `poolAmount`.
Pool amount is used to set the total reward pool and the module account balance in the module init if the module account is zero. It's needed also in the export of chain to hold the module account balance.
When migrate occurs the balance of the module account is set using reward pool amount.


## Store

### Pool

The module register in the store the pool amount.

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `vbr:pool:` | &rarr; | _pollAmount_ |

## Type definitions

#### `poolAmount` definition


```protobuf
repeated cosmos.base.v1beta1.DecCoin poolAmount = 1 [(gogoproto.moretags) = "yaml:\"pool_amount\"",
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.DecCoins",
    (gogoproto.nullable) = false] ;
```
