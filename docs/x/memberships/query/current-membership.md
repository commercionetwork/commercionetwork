# Getting the current membership

## Endpoint     
```
/memberships/${address}
```

## Parameters
| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read the membership |

## Example 
### Call
```
http://localhost:1317/memberships/did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke
```

### Response
```json
{
  "height": "0",
  "result": {
    "user": "did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke",
    "membership_type": "bronze"
  }
}
```