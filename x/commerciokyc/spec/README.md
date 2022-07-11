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
[invitation transaction](02_messages.md#msginviteuser) 

Once done, you will be automatically rewarded when the user buys a new membership.  

## Membership types
The currently supported memberships and their prices are the following.

::: tip  
Please note that each price is expressed using *Commercio Cash Credit* as measurement unit.  
::: 

| Membership | Identifier | Price | 
| :-------: | :---: | :---- |
| Green | `green` | `1` | 
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
| Green | 0.01 | 0.1 | 1.5 | 20 | 250 | 
| Bronze | 0.025 | 1.25 | 25 | 375 | 5,000 | 
| Silver | 0.1 | 5 | 75 | 1,000 | 12,500 |
| Gold | 0.4 | 12.5 | 150 | 1,750 | 20,000 |
| Black | 0.5 | 17.5 | 200 | 2,250 | 25,000 |  

Please note that the number of rewards is capped to a maximum of **12.5 millions tokens**.
After all the tokens have been distributed, any following invite will not be rewarded anymore.

## ABR Liquidity pool



## Contents

1. **[State](01_state.md)**
2. **[Messages](02_messages.md)**
   - [MsgInviteUser](02_messages.md#msginviteuser)
   - [MsgBuyMembership](02_messages.md#msgbuymembership)
   - [MsgSetMembership](02_messages.md#msgsetmembership)
   - [MsgRemoveMembership](02_messages.md#msgremovemembership)
   - [MsgAddTsp](02_messages.md#msgaddtsp)
   - [MsgRemoveTsp](02_messages.md#msgremovetsp)
   - [MsgDepositIntoLiquidityPool](02_messages.md#msgdepositintoliquiditypool)
3. **[Events](03_events.md)**
   - [Handlers](03_events.md#handlers)
4. **[Client](04_client.md)**
   - [Query](04_client.md#query)
   - [gRPC](04_client.md#gRPC)
   - [Rest](04_client.md#rest)
