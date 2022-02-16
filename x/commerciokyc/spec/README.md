<!--
order: 0
title: Commerciokyc Overview
parent:
  title: "commerciokyc"
-->

# CommercioKYC 

## Abstract

This document specifies the commerciokyc module of the Commercio Network.
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
| Green | 0.05 | 0.5 | 7.5 | 100 | 1,250 | 
| Bronze | 0.125 | 1.25 | 25 | 375 | 5,000 | 
| Silver | 0.5 | 5 | 75 | 1,000 | 12,500 |
| Gold | 2 | 12.5 | 150 | 1,750 | 20,000 |
| Black | 2.5 | 17.5 | 200 | 2,250 | 25,000 |  

Please note that the number of rewards is capped to a maximum of **12.5 millions tokens**.
After all the tokens have been distributed, any following invite will not be rewarded anymore.






## ABR Liquidity pool



## Module Accounts

The supply functionality introduces a new type of `auth.Account` which can be used by
modules to allocate tokens and in special cases mint or burn tokens. At a base
level these module accounts are capable of sending/receiving tokens to and from
`auth.Account`s and other module accounts. This design replaces previous
alternative designs where, to hold tokens, modules would burn the incoming
tokens from the sender account, and then track those tokens internally. Later,
in order to send tokens, the module would need to effectively mint tokens
within a destination account. The new design removes duplicate logic between
modules to perform this accounting.

The `ModuleAccount` interface is defined as follows:

```go
type ModuleAccount interface {
  auth.Account               // same methods as the Account interface

  GetName() string           // name of the module; used to obtain the address
  GetPermissions() []string  // permissions of module account
  HasPermission(string) bool
}
```

> **WARNING!**
> Any module or message handler that allows either direct or indirect sending of funds must explicitly guarantee those funds cannot be sent to module accounts (unless allowed).

The supply `Keeper` also introduces new wrapper functions for the auth `Keeper`
and the comerciokyc `Keeper` that are related to `ModuleAccount`s in order to be able
to:

- Get and set `ModuleAccount`s by providing the `Name`.
- Send coins from and to other `ModuleAccount`s or standard `Account`s
  (`BaseAccount` or `VestingAccount`) by passing only the `Name`.
- `Mint` or `Burn` coins for a `ModuleAccount` (restricted to its permissions).

### Permissions

Each `ModuleAccount` has a different set of permissions that provide different
object capabilities to perform certain actions. Permissions need to be
registered upon the creation of the supply `Keeper` so that every time a
`ModuleAccount` calls the allowed functions, the `Keeper` can lookup the
permissions to that specific account and perform or not the action.

The available permissions are:

- `Minter`: allows for a module to mint a specific amount of coins.
- `Burner`: allows for a module to burn a specific amount of coins.

## Contents

1. **[State](01_state.md)**
2. **[Keepers](02_keepers.md)**
   - [Common Types](02_keepers.md#common-types)
3. **[Messages](03_messages.md)**
   - [MsgInviteUser](03_messages.md#msginviteuser)
   - [MsgBuyMembership](03_messages.md#msgbuymembership)
   - [MsgSetMembership](03_messages.md#msgsetmembership)
   - [MsgRemoveMembership](03_messages.md#msgremovemembership)
   - [MsgAddTsp](03_messages.md#msgaddtsp)
   - [MsgRemoveTsp](03_messages.md#msgremovetsp)
   - [MsgDepositIntoLiquidityPool](03_messages.md#msgdepositintoliquiditypool)
4. **[Events](04_events.md)**
   - [Handlers](04_events.md#handlers)
5. **[Client](05_client.md)**
   - [Query](05_client.md#query)
ì