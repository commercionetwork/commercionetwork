# Sending an invite

## Transaction message
To invite a user, the following message must be used. 

```json
{
  "type": "commercio/MsgInviteUser",
  "value": {
    "recipient": "<Address of the user to be invited>",
    "sender": "<Your address>"
  }
}
```

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
inviteUser
```  