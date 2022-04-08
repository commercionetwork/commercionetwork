<!--
order: 2
-->

# Messages

## Increment Block Rewards Pool

### Protobuf message

```protobuf
message MsgIncrementBlockRewardsPool {
  string funder = 1 [(gogoproto.moretags) = "yaml:\"funder\""];
  repeated cosmos.base.v1beta1.Coin amount = 2 [(gogoproto.moretags) = "yaml:\"amount\"",
  (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins",
  (gogoproto.nullable) = false] ;
}
```


### Transaction message
To increment the block rewards pool you need to create and sign the following message:
  
```json
{
  "type": "commercio/MsgIncrementBlockRewardsPool",
  "value": {
    "funder": "<user address>",
    "amount": [
      {
        "denom": "<token denom>",
        "amount": "<amount to be incremented>"
      }
    ],
  }
}
```


#### Fields requirements
| Field | Required | Limit/Format |
| :---: | :------: | :------: |
| `funder` | Yes | bech32 | 
| `amount` | Yes |  []coins | 

### Action type
If you want to [list past transactions](../../../docs/developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
incrementBlockRewardsPool
```  


## Set Parameters

:::warning  
This transaction type is accessible only to the [government](../../government/spec/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

### Protobuf message

```protobuf
message MsgSetParams{
  string Government = 1 [(gogoproto.moretags) = "yaml:\"government\""];
  string distr_epoch_identifier = 2 [(gogoproto.moretags) = "yaml:\"distr_epoch_identifier\""];
  string earn_rate = 3 [(gogoproto.moretags) = "yaml:\"earn_rate\"",
  (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
  (gogoproto.nullable) = false];
}
```


#### Transaction message

To set the module's params you need to create and sign the following message:

```json
{
  "type": "commercio/MsgSetParams",
  "value": {
    "government": "<government address>",
    "distr_epoch_identifier": "<distribution epoch identifier>",
    "earn_rate": "<floating-point earn rate>"
  }
}
```

##### Fields requirements
| Field | Required | Limit/Format |
| :---: | :------: | :------: |
| `government` | Yes | bech32 | 
| `distr_epoch_identifier` | Yes | existing epoch identifier|
| `earn_rate` | Yes | Dec |



#### Action type
If you want to [list past transactions](../../../docs/developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setParams
```