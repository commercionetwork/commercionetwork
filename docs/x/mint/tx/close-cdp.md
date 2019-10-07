# Close a CDP
## Transaction message
To close a previously opened CDP you need to create and sign the following message.

```json
{
  "type": "commercio/MsgCloseCDP",
  "value": {
    "signer": "<User address>",
    "timestamp": "<Timestamp of when the CDP request was made>"
  }
}
```

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
closeCdp
```  