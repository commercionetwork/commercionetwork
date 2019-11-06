# Setting a user as verified

:::warning  
This transaction type is accessible only to Trusted Service Providers.  
Trying to perform this transaction without being a TSP will result in an error.  
:::

## Transaction message
To mark a user as verified, the following message should be used: 

```json
{
  "type": "commercio/MsgSetUserVerified",
  "value": {
    "user": "<Address of the user to mark as verified>",
    "verifier": "<Trusted Service Provider address>"
  }
}
```

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setUserVerified
```  