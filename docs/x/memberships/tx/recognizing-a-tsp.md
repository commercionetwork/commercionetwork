# Recognizing a TSP

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

## Transaction message
To recognize an address as a TSP, the following message must be used. 

```json
{
  "type": "commercio/MsgAddTsp",
  "value": {
    "tsp": "<Address of the user to be recognized as a TSP>",
    "government": "<Government address>"
  }
}
```

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
addTsp
```  