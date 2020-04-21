# CommercioMint

The `commerciomint` module is the one that allows you to create Collateralized Debt Positions (*CDPs*) using your 
Commercio.network tokens (*ucommercio*) in order to get Commercio Cash Credits (*uccc*) in return.

A *Collateralized Debt Position* (*CDP*) is a core component of the Commercio Network blockchain whose purpose is to
create Commercio Cash Credits (`uccc`) in exchange for Commercio Tokens (`ucommercio`) which it then holds in
escrow until the borrowed Commercio Cash Credits are returned.

In simple words, opening a CDP allows you to exchange any amount of `ucommercio` to get half the amount of `uccc`. 
For example, if you open a CDP lending `100 ucommercio` will result in you receiving `50 uccc`.    

## Transactions

### Open a CDP

#### Transaction message
To open a new CDP you need to create and sign the following message:
  
```json
{
  "type": "commercio/MsgOpenCdp",
  "value": {
    "deposited_amount": [
      {
        "amount": "<amount to be deposited>",
        "denom": "<token denom to be deposited>"
      }
    ],
    "depositor": "<user address>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
openCdp
```  

### Close a CDP

#### Transaction message

To close a previously opened CDP you need to create and sign the following message:

```json
{
  "type": "commercio/MsgCloseCdp",
  "value": {
    "signer": "<user address>",
    "timestamp": "<block height at which the CDP is being inserted into the chain>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
closeCdp
```  

### Set CDP collateral rate

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

#### Transaction message

To set the CDP collateral rate you need to create and sign the following message:

```json
{
  "type": "commercio/MsgSetCdpCollateralRate",
  "value": {
    "signer": "<user address>",
    "cdp_collateral_rate": "<floating-point collateral rate>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setCdpCollateralRate
```  

## Queries
- [Read a CDP by its owner's address and timestamp](query/read-cdp.md)
- [List all CDPs of a user](query/read-cdps.md)
