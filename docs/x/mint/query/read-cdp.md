# Reading a user opened CDP

## Rest API
### Endpoint     
```
/cdp/${address}/${timestamp}
```

### Parameters
| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read the CDP |
| `timestamp`| Timestamp of when the CDP request was made |

### Example 
#### Call
```
http://localhost:1317/mint/CDP/did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke/1570177686
```
#### Response
```json
{
  "height": "0",
  "result": {
    "deposited_amount": {
      "amount": "10000000",
      "denom": "ucommercio"
    },
    "liquidity_amount": {
      "amount": "500000",
      "denom": "uccc"
    },
    "owner": "did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke/1570177686",
    "timestamp": "1570177686"
  }
}
```
