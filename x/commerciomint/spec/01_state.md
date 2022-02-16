<!--
order: 1
-->

# State

The `x/commerciomint` module keeps state of the follow objects:

1. Positions
2. Pool amount

## Positions
Positions are objects that are created when an user deposit an amount of Commercio Cash Credit(CCC). However the holded collateral is proportional to the position's exchange rate.

 ```proto
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

## Pool amount
Pool amount is used to set the liquidity pool to the module account in the module init if the module account is zero. It's needed also in the of chain export to hold the module account balance.
When migrate occurs the balance of the module account is set using pool amount liquidity.

```proto
repeated cosmos.base.v1beta1.Coin pool_amount = 2 [
    (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins",
    (gogoproto.nullable) = false
  ];
```