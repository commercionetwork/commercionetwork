# Buying a membership

## Transaction message
To buy a membership, the following message must be used. 

```json
{
  "type": "commercio/MsgBuyMembership",
  "value": {
    "membership_type": "<Membership type identifier>",
    "buyer": "<Your address>"
  }
}
```

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
buyMembership
```  