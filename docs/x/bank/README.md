# Bank
The `bank` module allows you to send any token amount you possess to any other user inside the Commercio.network chain. 


## Sending tokens

### Transaction message
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

#### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `amount` | Yes |
| `to_address` | Yes |
| `from_address` | Yes |

### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
send
```  

:::warning    
Please note that you might not be able to perform this transaction if you are marked as a blocked account.  
To know more about blocked accounts, please read [Blocked accounts](#blocked-accounts).  
:::


## Blocked accounts
In order to preserve the token price stability, we've decided to block some users from sending the token too early.  
These users are not common users but are people that have contributed to the Commercio.network early stage development 
and which tokens selling might compromise the whole token stability. 

Such blocked accounts **are not able to send any kind of token to any user**. 
If they try to do so, the transaction will simply fail and will not be considered valid from the whole system. 

## Adding a blocked account
In order to add an account as blocked you need to be the [government](../government/README.md).  


### Blocking an account 

:::warning  
This transaction can only be performed by the [government](../../government/README.md).  
:::

#### Transaction message
In order to prevent a specific user from being able to send any token, you must use the following message: 

```json
{
  "type": "commercio/MsgBlockAccountSend",
  "value": {
    "address": "<Address of the user to be blocked>",
    "signer": "<Government address>"
  }
}
```

##### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `address` | Yes |
| `signer` | Yes |

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
blockAccountSend
```  


### Unlocking a blocked account
In order to unlock a previously blocked account you need to be the [government](../government/README.md).  

:::warning  
This transaction can only be performed by the [government](../government/README.md).  
:::

#### Transaction message
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

##### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `address` | Yes |
| `signer` | Yes |

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
unlockAccountSend
```  

