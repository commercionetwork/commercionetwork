# Unlocking an account 

:::warning  
This transaction can only be performed by the [government](../../government/README.md).  
:::

## Transaction message
In order to allow a blocked user to send tokens again, you must use the following message:

```json
{
  "type": "commercio/MsgUnlockAccountSend",
  "value": {
    "address": "<Address of the user to be unlocked>",
    "signer": "<Government address>"
  }
}
```

### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `address` | Yes |
| `signer` | Yes |

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
unlockAccountSend
```  