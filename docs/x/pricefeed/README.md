# Pricefeed
The `pricefeed` module allows to external actors (called *oracles*) to insert into the Commercio.network
blockchain the current prices of different assets.

## Transactions

### Adding an oracle

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::


#### Transaction message
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

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
addOracle
```  

### Set a price for an asset

::: tip  
In order to perform this transaction you need to be an *oracle*.  
If you wish to become one, please take a look at the [*"Adding an oracle"*](#adding-an-oracle) section.  
:::

If you have been added as an oracle you can set a price into the blockchain **once every block**.  
This means that trying to set a price for the same asset more than once every 6 seconds (the estimated block time) 
will result in an error when trying to send the transactions. 

#### Transaction message
In order to set a raw price for a specific asset, you need to create and sign the following message. 

```json
{
  "type": "commercio/MsgSetPrice",
  "value": {
    "oracle": "<Did of the oracle>",
    "price": {
      "asset_name": "<Name of the asset>",
      "value": "<Price of the asset (supports decimal numbers)>",
      "expiry": "<Block height after which the price should be considered invalid>"
    }
  }
}
```

##### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `oracle` | Yes |
| `price` | Yes |

##### `price` fields requirements
| Field | Required | 
| :---: | :------: | 
| `asset_name` | Yes |
| `value` | Yes[^1] |
| `expiry` | Yes |

[^1]: The `value` field value must have a 18 decimal digits, for example, to set price `1` you must use `1.000000000000000000`. 

#### About the expiration
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

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setPrice
```    