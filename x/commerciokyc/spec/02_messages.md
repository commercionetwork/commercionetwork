<!--
order: 2
-->

# Messages
To invite a user, the following message must be used: 





## MsgInviteUser


### Protobuf message

```protobuf
message MsgInviteUser {
  string recipient = 1;
  string sender = 2;
}
```
### Transaction message

```json
{
  "type": "commercio/MsgInviteUser",
  "value": {
    "recipient": "<address of the user to be invited>",
    "sender": "<your address>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
inviteUser
```

## MsgBuyMembership

Buying a membership



### Protobuf message

```protobuf
message MsgBuyMembership {
  string membership_type = 1;
  string buyer = 2;
  string tsp = 3;
}
```

### Transaction message


```json
{
  "type": "commercio/MsgBuyMembership",
  "value": {
    "membership_type": "<membership type identifier>",
    "buyer": "<buyer address>",
    "tsp": "<tsp address>"
  }
}
```
#### Action type

If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
buyMembership
```




## MsgSetMembership

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

To arbitrarily set a user's membership, the following message must be used:


### Protobuf message

```protobuf
message MsgSetMembership {
  string government = 1;
  string subscriber = 2;
  string new_membership = 3;
}
```

### Transaction message

```json
{
  "type": "commercio/MsgSetMembership",
  "value": {
    "government": "<address of the government that sends this message>",
    "subscriber": "<address which membership will change based on this message>",
    "new_membership": "<membership type identifier>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setMembership
```  



## MsgRemoveMembership


:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::


To arbitrarily set a user's membership, the following message must be used:

### Protobuf message

```protobuf
message MsgRemoveMembership {
  string government = 1;
  string subscriber = 2;
}
```

### Transaction message


```json
{
  "type": "commercio/MsgRemoveMembership",
  "value": {
    "government": "<address of the government that sends this message>",
    "subscriber": "<address which membership will change based on this message>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
removeMembership
```  

## MsgAddTsp

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

### Protobuf message

```protobuf
message MsgAddTsp {
  string tsp = 1;
  string government = 2;
}
```

### Transaction message


```json
{
  "type": "commercio/MsgAddTsp",
  "value": {
    "tsp": "<address of the user to be recognized as a TSP>",
    "government": "<government address>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
addTsp
```  

## MsgRemoveTsp

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

### Protobuf message

```protobuf
message MsgRemoveTsp {
  string tsp = 1;
  string government = 2;
}
```

### Transaction message

```json
{
  "type": "commercio/MsgRemoveTsp",
  "value": {
    "tsp": "<address of the user to be recognized as a TSP>",
    "government": "<government address>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
removeTsp
```  



## MsgDepositIntoLiquidityPool

### Protobuf message

```protobuf
message MsgDepositIntoLiquidityPool {
  string depositor = 1;
  repeated cosmos.base.v1beta1.Coin amount = 2
    [(gogoproto.nullable) = false, (gogoproto.castrepeated) = "github.com/cosmos/cosmos-sdk/types.Coins"];
}
```

### Transaction message

To deposit a given amount into the Memberships reward pool, the following message must be used:

```json
{
  "type": "commercio/MsgDepositIntoLiquidityPool",
  "value": {
    "depositor": "<address that deposits into the pool>",
    "amount": [
      {
        "amount": "<amount to be deposited>",
        "denom": "<token denom to be deposited>"        
      }
    ]
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
depositIntoLiquidityPool
```  

