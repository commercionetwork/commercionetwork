# Invalidating a Did power up request
Once a [Did power up request has been created](./request-did-power-up.md), it can be invalidated in two different cases: 
 
1. The user wants to cancel it. 
2. The external centralized identity invalidates it for different reasons (insufficient funds, errors, etc.). 

## Transaction message
```json
{
  "type": "commercio/MsgInvalidateDidPowerUpRequest",
  "value": {
    "editor": "<Address of the user that is invalidating the request>",
    "power_up_proof": "<Value of the proof that has been used while creating the deposit request>",
    "status": {
      "type": "<Type of the new status>",
      "message": "<Optional additional message explaining why the status has changed>"
    }
  }
}
```

### Fields requirements
| Field | Required |
| :---: | :------: |
| `editor` | Yes |
| `deposit_proof` | 
| `status` | Yes | 

#### `status`
| Field | Required | 
| :---: | :------: |
| `type` | Yes *<sup>1</sup> |
| `message` | No |

*<sup>1</sup> The `type` field value must be either one of the following strings:
- `canceled` if the user is canceling the request spontaneously
- `rejected` if the external centralized authority is rejecting the request for different reasons
  

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
invalidateDidPowerUpRequest
```