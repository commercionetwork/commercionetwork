# Sending tokens

## Transaction message
In order to send any token amount to a user, you need to use the following message:

```json
{
  "type": "cosmos-sdk/MsgSend",
  "value": {
    "amount": [
      {
        "denom": "<Token denomination to be sent>",
        "amount": "<Token amount to be sent>"
      }
    ],
    "to_address": "<Address of the recipient>",
    "from_address": "<Your address>"
  }
}
```

### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `amount` | Yes |
| `to_address` | Yes |
| `from_address` | Yes |

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
send
```  