# Getting the current membership

# CLI

## Get membership invites for user

```sh
$ cncli query accreditations get-invites-for-user [address]
```

## Get all the membership invites

```sh
$ cncli query accreditations get-invites
```

# REST

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