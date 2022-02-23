# Accept a deposit request
Once a user has properly [created a Did deposit request](./request-did-deposit.md), the external centralized entity 
will be able to properly accept such a request and perform a deposit withdraw.

## Transaction message
```json
{
  "type": "commercio/MsgMoveDeposit",
  "value": {
    "deposit_proof": "<Proof used inside the deposit request>",
    "signer": "<Government address>"
  }
}
```  
### Fields requirements
| Field | Required |
| :---: | :------: |
| `deposit_proof` | Yes | 
| `signer` | Yes |

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
moveDeposit
```
