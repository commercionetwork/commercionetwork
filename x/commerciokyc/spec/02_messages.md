<!--
order: 2
-->

# Messages
To invite a user, the following message must be used: 

## MsgInviteUser

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

