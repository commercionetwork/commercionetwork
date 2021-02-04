# Requesting a Did Power Up
A *Did Power Up* is the expression we use when referring to the willingness of a user to move a specified amount of tokens 
from external centralized entity to one of his
private pairwise Did, making them able to send documents (which indeed require the user to spend some tokens as fees). 

A user who wants to execute a Did Power Up must have previously sent tokens to the public address of the centralized entity.
  
This action is the second and final step that must be done when [creating a pairwise Did](../creating-pairwise-did.md).  

:::tip  
If you wish to know more about the overall pairwise Did creation sequence, please refer to the
[pairwise Did creation specification](../creating-pairwise-did.md) page  
:::    

## Transaction message
```javascript
{
  "type": "commercio/MsgRequestDidPowerUp",
  "value": {
    "claimant": "address that sent funds to the centralized entity before",
    "amount": [
      {
        "denom": "ucommercio",
        "amount": "amount to transfer to the pairwise did, integer"
      }
    ],
    "proof": "proof string",
    "id": "randomly-generated UUID v4",
    "proof_key": "proof encryption key"
  }
}
```

### `value` fields requirements
| Field | Required |
| :---: | :------: |
| `claimant` | Yes |
| `amount` | Yes |
| `proof` | Yes | 
| `id` | Yes | 
| `proof_key` | Yes | 


### Creating the `proof` value

To create the `proof` field value, the following steps must be followed:

1. create the `signature_json` formed as follow.  
   ```javascript
   {
    "sender_did": "user who sends the power-up request",
    "pairwise_did": "pairwise address to power-up",
    "timestamp": "UNIX-formatted timestamp",
   }
   ```

2. retrive the public key of external centralized entity **Tk**, by querying the `cncli` REST API
3. calculate SHA-256 `HASH` of the concatenation of `sender_did`, `pairwise_did` and `timestamp` fields, taken from `signature_json`
4. do a PKCS1v15 signature of `HASH` with the RSA private key associated to RSA public key inserted in the `sender_did` DDO - this process yields the `SIGN(HASH)` value
5. convert `SIGN(HASH)` in **base64** `BASE64(SIGN(HASH))`, this is the value to be placed in the `signature` field 
6. add the `signature` field to the `signature_json` JSON:

   ```javascript
   {
    "sender_did": "user who sends the power-up request",
    "pairwise_did": "pairwise address to power-up",
    "timestamp": "UNIX-formatted timestamp",
    "signature": "the value BASE64(SIGN(HASH))",
   }
   ```
7. create a random 256-bit [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) key `F`
8. generate a random 96-bit nonce `N`
8. using the AES-256 key generated at point (6), encrypt the `payload`:
   1. remove all the white spaces and line ending characters
   2. encrypt the resulting string bytes using the `AES-GCM` mode, `F` as key, obtaining `CIPHERTEXT`
   3. concatenate bytes of `CIPHERTEXT` and `N` and encode the resulting bytes in **base64**, obtaining the `value` `proof` content 
9. encrypt the AES-256 key:
   1. encrypt the `F` key bytes using the centralized entity's RSA public key found in its Did Document, in PKCS1v15 mode.  
   2. encode the resulting bytes in **base64**, obtaining the `value` `proof_key` content 


## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
requestDidPowerUp
```
