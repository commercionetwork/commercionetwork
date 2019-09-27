# Adding an oracle

:::warning  
This transaction type is reserved to the [*government*](../../government/README.md).  
If you want to become an oracle please contact a Commercio.network administrator.  
:::


## Transaction message
To add a new oracle into the set of current oracles you need to create and sign the following message. 

```json
{
  "type": "commercio/MsgAddOracle",
  "value": {
    "oracle": "<Did of the oracle to be added>",
    "signer": "<Government address>"
  }
}
```