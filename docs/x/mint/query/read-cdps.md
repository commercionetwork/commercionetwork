# Reading all user opened CDPs

## Rest API
### Endpoint     
```
/cdp/${address}
```

### Parameters
| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read all the CDPs |

### Example 
#### Call
```
http://localhost:1317/mint/CDPs/did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke
```
#### Response
```json
{
  "height": "0",
  "result": [
    {
      "deposited_amount": {
        "denom": "ucommercio",
        "amount": "10000000"
      },
      "liquidity_amount": {
        "denom": "uccc",
        "amount": "500000"
      },
      "owner": "did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke/1570177686",
      "timestamp": "1570177686"
    }
  ]
}
```
