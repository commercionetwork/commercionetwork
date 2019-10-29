# Did power up
Once a user has properly [create a Did power up request](./request-did-power-up.md), the external centralized entity can
accept such request and fund the pairwise Did specified. 

## Transaction message
```json
{
  "type": "commercio/MsgPowerUpDid",
  "value": {
    "recipient": "<Address of the Did to fund>",
    "amount": [
      {
        "amount": "<Amount of coins to be sent>",
        "denom": "<Denom of the coin to send>"
      }   
    ],
    "activation_reference": "<Encrypted power up proof used inside the request>",
    "signer": "<Government address>"
  }
}
```

### Fields requirements
| Field | Required |
| :---: | :------: |
| `recipient` | Yes | 
| `amount` | Yes |
| `activation_reference` | Yes |
| `signer` | Yes |

#### `amount`
| Field | Required |
| :---: | :------: |
| `amount` | Yes |
| `denom` | Yes | 

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
powerUpDid
```
