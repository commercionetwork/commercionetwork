# Open a CDP

:::warning  
Before performing this type of transaction be sure to understand what a CDP is, and how it works.  
You could loose your token if you aren't be careful.  
:::

## Transaction message
To open a new CDP you need to create and sign the following message.
  
```json
{
  "type": "commercio/MsgOpenCDP",
  "value": {
    "cdp_request": {
      "deposited_amount": "<Token to be deposited as a collateral (supports only integers)>",
      "signer": "<User address>",
      "timestamp": "<Timestamp of when the CDP request was made>"
    }
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