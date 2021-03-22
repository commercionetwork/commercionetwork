# CommercioKYC 

Inside Commercio.network we designed a system to make sure to have a network of trusted participants (Know your customer).

To do so we've implemented a *commerciokyc* system which allows you to buy a membership by a TSP (Trust Service Provider) to display to everyone that you've been invited by an already verified members to join the network. 

## Buying a membership
Memberships can be bought **exclusively on chain** and **exclusively through a TSP**. 
To do so you are required to possess and spend an amount of Commercio Cash Credits (*CCC*) greater or 
equal to the price of the membership you wish and request to a TSP to buy it. 


### Requirements
In order to buy a membership, the following requirements must be met: 

1. You must have been invited by a user already having a Green membership or superior.  
   Please refer to the [invitation procedure page](#invitation-process) 
   for more details on invitations. 
2. You must own a sufficient amount of CCCs to buy a membership or buy it directaly from Trusted Service Provider. 
   Please refer to the [proper page](../commerciomint/README.md) to know more about how to get CCCs out of your Commercio Tokens.
3. You need to transfer the amount of CCCs to your Trusted Service Provider or otherwise let the Trusted Service Provider complete the transaction for you as a service.
4. Your Trusted Service Provider will purchase the membership for you after proper verification of your identity.
   
The system will **prevent** sending an invitation to an account already active on the chain or to an address already invited.

### Buying process

Once you've met all the [requirements](#requirements), you can buy a membership by 
performing a [buying transaction](#buying-a-membership).  
After doing so, you will be able to verify your membership status 
using the [membership query feature](#queries).

## Invitation process
The invitation process is the first su-process of the accreditation procedure. 
which allows you to invite a new user to join Commercio.network so that you can later be rewarded if he later
buys a membership.  

### Requirements
Suppose that you want to invite another user and later receive a reward once he buys a membership. 
Then, the following requirements must be met:

* You must already have [bought](#buying-a-membership) a membership. 
* You must know the user's identity that he has created by following the [specific guide](../id/tx/create-an-identity.md).

## Inviting a user
Once all the [above requirements](#requirements) have been met, you can invite the user by performing an
[invitation transaction](#sending-an-invite) 

Once done, you will be automatically rewarded when the user buys a new membership.  

## Membership types
The currently supported memberships and their prices are the following.

::: tip  
Please note that each price is expressed using *Commercio Cash Credit* as measurement unit.  
::: 

| Membership | Identifier | Price | 
| :-------: | :---: | :---- |
| Green | `green` | `5` | 
| Bronze | `bronze` | `25` | 
| Silver | `silver` | `250` | 
| Gold | `gold` | `2500` | 
| Black | `black` | `50000` |

## Rewards values
The reward value given to you is based on your membership type and the membership 
that any user you've accreditated buys. 

::: tip  
Please note that each reward is expressed using *Commercio Token* as measurement unit.  
::: 

| Invitee / Invited | Green | Bronze | Silver | Gold | Black |
| :--------------: | :----: | :----: | :----: | :---: | :---: |
| Green | 0.05 | 0.5 | 7.5 | 100 | 1'250 | 
| Bronze | 0.125 | 1.25 | 25 | 375 | 5'000 | 
| Silver | 0.5 | 5 | 75 | 1'000 | 12'500 |
| Gold | 2 | 12.5 | 150 | 1'750 | 20'000 |
| Black | 2.5 | 17.5 | 200 | 2'250 | 25'000 |  

Please note that the number of rewards is capped to a maximum of **12.5 millions tokens**.
After all the tokens have been distributed, any following invite will not be rewarded anymore.






## Transactions

### Buying a membership

#### Transaction message
To buy a membership, the following message must be used: 

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



### Sending an invite

#### Transaction message
To invite a user, the following message must be used: 

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



### Adding a TSP

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

#### Transaction message
To recognize an address as a TSP, the following message must be used: 

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




### Removing a TSP

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

#### Transaction message
To recognize an address as a TSP, the following message must be used: 

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





### Deposit into reward pool

#### Transaction message
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





### Set user membership

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

#### Transaction message
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



### Remove user membership

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

#### Transaction message
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







## Queries

### Getting current membership invites

#### CLI

Get all the all invites:

```bash
cncli query commerciokyc invites
```

#### REST

Endpoints:
     
```
/commerciokyc/invites
```

##### Example 

Getting all invites:

```
http://localhost:1317/commerciokyc/invites
```

##### Response
```json
{
  "height": "0",
  "result": [
    {
      "sender": "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
      "sender_membership": "black",
      "user": "did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf",
      "status": 0
    }
  ]
}
```

### Getting the TSP list

#### CLI

```bash
cncli query commerciokyc trusted-service-providers
```

#### REST

Endpoints:
     
```
/commerciokyc/tsps
```

##### Example 

Getting all TSPs:

```
http://localhost:1317/commerciokyc/tsps
```

##### Response
```json
{
  "height": "0",
  "result": [
    "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen"
  ]
}
```

### Getting the reward pool funds amount

#### CLI

```bash
cncli query commerciokyc pool-funds
```

#### REST

Endpoints:
     
```
/commerciokyc/pool-funds
```

##### Example 

Getting the reward pool funds amount:

```
http://localhost:1317/commerciokyc/accreditations-funds
```

##### Response
```json
{
  "height": "0",
  "result": [
    {
      "denom": "ucommercio",
      "amount": "9999899990000"
    }
  ]
}
```

### Getting user membership

#### CLI

```bash
cncli query commerciokyc membership [user]
```

#### REST

Endpoints:
     
```
/commerciokyc/membership/{address}
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read membership|

##### Example 

Getting the reward pool funds amount:

```
http://localhost:1317/commerciokyc/membership/did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen
```

##### Response
```json
{
  "height": "0",
  "result": {
    "user": "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
    "tsp_address" : "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
    "membership_type": "black",
    "expiry_at" : "2022-03-21T00:00:00Z"

}
```


### Getting all memberships

#### CLI

```bash
cncli query commerciokyc memberships
```

#### REST

Endpoints:
     
```
/commerciokyc/memberships
```

##### Example 

Getting the reward pool funds amount:

```
http://localhost:1317/commerciokyc/memberships
```

##### Response
```json
{
  "height": "0",
  "result": [
    {
      "user": "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
      "tsp_address" : "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
      "membership_type": "black",
      "expiry_at" : "2022-03-21T00:00:00Z"
    }
  ]
```


### Getting all memberships sold by tsp

#### CLI

```bash
cncli query commerciokyc sold [tsp-address]
```

#### REST

Endpoints:
     
```
/commerciokyc/sold/{address}
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read membership|

##### Example 

Getting the reward pool funds amount:

```
http://localhost:1317/commerciokyc/sold/did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen
```

##### Response
```json
{
  "height": "0",
  "result": [
    {
      "user": "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
      "tsp_address" : "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
      "membership_type": "black",
      "expiry_at" : "2022-03-21T00:00:00Z"
    }
  ]
```