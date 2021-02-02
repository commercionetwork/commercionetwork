# Memberships

Inside Commercio.network we designed a system to make sure to have a network of trusted participants.

To do so we've implemented a *membership* system which allows you to buy a membership to display to everyone that 
you've been invited by an already verified members to join the network. 

## Buying a membership
Memberships can be bought **exclusively on chain**. 
To do so you are required to possess and spend an amount of Commercio Cash Credits (*CCC*) greater or 
equal to the price of the membership you wish to buy.  

### Requirements
In order to buy a membership, the following requirements must be met: 

1. You must have been invited by a user already having a Bronze membership or superior.  
   Please refer to the [invitation procedure page](#invitation-process) 
   to know more about invitations. 
2. You must possess a sufficient amount of CCCs to buy a membership.  
   Please refer to the [proper page](../commerciomint/README.md) to know more about how to get CCCs out of your Commercio Tokens.
   
The system will **prevent** sending an invite to an account with a balance greater than zero across CCCs and/or Commercio token.

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
| Bronze | `bronze` | `25` | 
| Silver | `silver` | `250` | 
| Gold | `gold` | `2500` | 
| Black | `black` | `50000` |

## Rewards values
The reward value given to you is based on your membership type and the membership 
that any user you've accreditated buys. 

| Invitee / Invited | Bronze | Silver | Gold | Black |
| :--------------: | :----: | :----: | :---: | :---: |
| Bronze | 1.25 | 25 | 375 | 5'000 | 
| Silver | 5 | 75 | 1'000 | 12'500 |
| Gold | 12.5 | 150 | 1'750 | 20'000 |
| Black | 17.5 | 200 | 2'250 | 25'000 |  

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
    "buyer": "<your address>"
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

### Deposit into reward pool

#### Transaction message
To deposit a given amount into the Memberships reward pool, the following message must be used:

```json
{
  "type": "commercio/MsgDepositIntoLiquidityPool",
  "value": {
    "depositor": "<address that deposits into the pool>",
    "amount": [
      "<amount of coins to deposit into the pool>"
    ]
  }
}
```

where components of `amount` must be objects of the following type:

```json
{
    "denom": "<token denom>",
    "amount": "<integer amount of tokens of denom>"
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
    "government_address": "<address of the government that sends this message>",
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

## Queries

### Getting current membership invites

#### CLI

Get membership invites for user:

```bash
cncli query accreditations get-invites-for-user [address]
```

Get all the membership invites:

```bash
cncli query accreditations get-invites
```

#### REST

Endpoints:
     
```
/invites/${address}
/invites
```

Parameters (if any):

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current invite if any |

##### Example 

Getting invites for `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/invites/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf
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

##### Example 

Getting all invites:

```
http://localhost:1317/invites
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
cncli query accreditations trusted-service-providers
```

#### REST

Endpoints:
     
```
/tsps
```

##### Example 

Getting all TSPs:

```
http://localhost:1317/tsps
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
cncli query accreditations pool-funds
```

#### REST

Endpoints:
     
```
/accreditations-funds
```

##### Example 

Getting the reward pool funds amount:

```
http://localhost:1317/accreditations-funds
```

##### Response
```json
{
  "height": "0",
  "result": {
    "denom": "ucommercio",
    "amount": "9999899990000"
  }
}
```

### Getting user membership

#### CLI

```bash
cncli query accreditations user-membership
```

#### REST

Endpoints:
     
```
/membership/{address}
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read membership|

##### Example 

Getting the reward pool funds amount:

```
http://localhost:1317/membership/did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen
```

##### Response
```json
{
  "height": "0",
  "result": {
    "user": "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
    "membership_type": "black"
  }
}
```
