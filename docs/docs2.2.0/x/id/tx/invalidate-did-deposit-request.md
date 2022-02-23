# Invalidating a Did deposit request
Once a [Did deposit request has been created](./request-did-deposit.md), it can be invalidated in two different cases: 
 
1. The user wants to cancel it. 
2. The external centralized identity invalidates it for different reasons (insufficient funds, errors, etc.). 

## Transaction message
```json
{
  "type": "commercio/MsgInvalidateDidDepositRequest",
  "value": {
    "editor": "<Address of the user that is invalidating the request>",
    "deposit_proof": "<Value of the proof that has been used while creating the deposit request>",
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
invalidateDidDepositRequest
```