# Open a CDP
A *Collateralized Debt Position* (*CDP*) is a core component of the Commercio Network blockchain whose purpose is to
create Commercio Cash Credits (`uccc`) in exchange for Commercio Tokens (`ucommercio`) which it then holds in
escrow until the borrowed Commercio Cash Credits are returned.

In simple words, opening a CDP allows you to exchange any amount of `ucommercio` to get half the amount of `uccc`. 
For example, if you open a CDP lending `100 ucommercio` will result in you receiving `50 uccc`.    

## Transaction message
To open a new CDP you need to create and sign the following message.
  
```json
{
  "type": "commercio/MsgOpenCdp",
  "value": {
    "deposited_amount": [
      {
        "amount": "<Amount to be deposited>",
        "denom": "<Token denom to be deposited>"
      }
    ],
    "signer": "<User address>",
    "timestamp": "<Timestamp of when the CDP request was made>"
  }
}
```

### About the Timestamp
The timestamp is a way to track time as a running total of seconds.  
In this case the timestamp has to be the running total of seconds from when the CDP request was made.

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
openCdp
```  