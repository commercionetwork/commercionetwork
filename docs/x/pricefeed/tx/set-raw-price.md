# Setting a raw price for an asset

::: tip  
In order to perform this transaction you need to be an *oracle*.  
If you wish to become one, please take a look at the [*"Adding an oracle"* page](add-oracle.md)  
:::

If you have been added as an oracle you can set a raw price into the blockchain **once every block**.  
This means that trying to set a price for the same asset more than once every 6 seconds (the estimated block time) 
will result in an error when trying to send the transactions. 

## Transaction message
In order to set a raw price for a specific asset, you need to create and sign the following message. 

```json
{
  "type": "commercio/MsgSetPrice",
  "value": {
    "oracle": "<Did of the oracle>",
    "price": {
      "asset_name": "<Name of the asset>",
      "price": "<Price of the asset (supports decimal numbers)>",
      "expiry": "<Block height after which the price should be considered invalid>"
    }
  }
}
```

### About the expiration
When sending the above message, you will need to specify a valid `expiry` value.  
The best way to do so is to:

1. Query the latest block height.  
   You can do so by using the following REST API:  
   ```
   /blocks/latest
   ```
   
2. Once you have the latest block height, add to it $k = \frac{n}{6}$ where $n$ is the number of seconds you
   want the price to be valid for. 
   
Please note that adding a price with a block height **equals or lower** than the current block height will result 
in the price never be considered valid and thus never be taken into consideration.

We suggest you to choose a value of $k = 10$ which will make the price valid for more or less 1 minute.  