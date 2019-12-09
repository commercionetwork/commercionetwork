# Close a CDP
## Transaction message
To close a previously opened CDP you need to create and sign the following message.

```json
{
  "type": "commercio/MsgCloseCdp",
  "value": {
    "signer": "<User address>",
    "timestamp": "<Block height at which the CDP is being inserted into the chain>"
  }
}
```

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
closeCdp
```  