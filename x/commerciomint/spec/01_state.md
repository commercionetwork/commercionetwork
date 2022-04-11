<!--
order: 1
-->

# State

The `commerciomint` module keeps state of the Exchange Trade Positions


## Store

### Positions 


| Key |  | Value |
| ------- | ---------- | ---------- | 
| `commerciomint:etp:[owner][ID]` | &rarr; | _Position_ |

## Type definitions

### Positions
Positions are objects that are created when a user deposits an amount of Commercio Cash Credit (CCC). However the holded collateral is proportional to the position's exchange rate.

 ```protobuf
 message Position {
  string owner = 1;
  int64 collateral = 2;
  cosmos.base.v1beta1.Coin credits = 3;
  google.protobuf.Timestamp created_at = 4 [ (gogoproto.stdtime) = true ];
  string ID = 5;
  string exchange_rate = 6 [
    (gogoproto.nullable) = false,
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec"
  ];
}
 ```
