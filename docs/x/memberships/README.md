# Memberships 
Inside Commercio.network we've designed a system to make sure that the network that will be build is made
of trusted participants.

To do so we've implemented a *membership* system which allows you to buy a membership (represented by an *NFT token*) 
to display to everyone that you've been invited by an already verified members to join the network. 

## Buying a membership
Memberships can be bought **exclusively on chain**. 
To do so you are required to possess an amount of Commercio Cash Credits (*CCC*) greater or 
equals to price of the membership you wish to buy.  

### Requirements
In order to buy a membership, the following requirements must be met: 

1. You must have been invited by a user already having a Bronze membership or superior.  
   Please refer to the [invitation procedure page](../accreditations/invitation-process.md) 
   to know more about invitations. 
2. You must have been verified by a Trusted Service Provider.  
   Please refer to the [verification procedure](../accreditations/verification-process.md)
   to know more about the verification process.
3. You must possess a sufficient amount of CCCs to buy a membership.  
   Please refer to the [proper page](../mint/README.md) to know more about how to get CCCs out of your Commercio Tokens.
   
### Buying process
Once you've met all the [requirements](#requirements), you can buy a membership by 
performing a [buying transaction](./tx/buy-membership.md).  
After doing so, you will be able to verify your membership status 
using the [membership query feature](./query/current-membership.md).  

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
| Black | `black` | `25000` |

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