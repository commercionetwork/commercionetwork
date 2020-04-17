# Requesting a Did power up
A *Did power up* is a term we use when referring to the willingness of a user to move a specified amount of tokens 
from external centralized entity to one of his
private pairwise Did, making them able to send documents (which indeed require the user to spend some tokens as fees).  
User has previously send tokens to public address of centralized entity.
  
This action is the second and final step that must be done when [creating a pairwise Did](../creating-pairwise-did.md).  

:::tip  
If you wish to know more about the overall pairwise Did creation sequence, please refer to the
[pairwise Did creation specification](../creating-pairwise-did.md) page  
:::    

## Transaction message
```json
{
  "type": "commercio/MsgRequestDidPowerUp",
  "value": {
    "status": "",
    "claimant": "<Address that is able to spend the funds (the recipient used during the deposit procedure)>",
    "amount": "int64",
    "proof": "string",
    "id": "string",
    "proof_key": "string"
  }
}
```

### Fields requirements
| Field | Required |
| :---: | :------: |
| `status` | No (Don't use) |
| `claimant` | Yes |
| `amount` | Yes |
| `proof` | Yes | 
| `id` | Yes | 
| `proof_key` | Yes | 


### Creating the `proof ` value
When creating the `proof ` field value, the following steps must be followed. 


1. Create the `signature_json` formed as follow.  
   ```json
   {
    "sender_did": "<User did>",
    "pairwise_did": "<Pairwise Did to power up>",
    "timestamp": <Timestamp>,
   }
   ```

2. Retrive the public key of external centralized entity **Tk**
3. Calculate SHA-256 `HASH` of `sender_did`, `pairwise_did` and `timestamp` concatenation
4. Sign in format PKCS1v15 the `HASH` with the RSA private key associated to RSA public key inserted in the DDO. Now we have `SIGN(HASH)`
5. Convert `SIGN(HASH)` in Base64 notation `BASE64(SIGN(HASH))` and use it to add `signature` field 

6. Create a `payload` JSON made as follow:

   ```json
   {
    "sender_did": "<User did>",
    "pairwise_did": "<Pairwise Did to power up>",
    "timestamp": <Timestamp>,
    "signature": `BASE64(SIGN(HASH))`,
   }
   ```
7. Create a random [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) key `F`

8. Generate a random 96-bit nonce `N`

8. Using the AES-256 key generated at point (6), encrypt the `payload`.
   1. Remove all the white spaces and line ending characters. 
   2. Encrypt the resulting string bytes using `F`, obtaining `CIPHERTEXT`  
      Note that the AES encryption method must be `AES-GCM`.
   3. Concatenate bytes of `CIPHERTEXT` and `N` and encode the resulting bytes using the Base64 encoding method, obtaining `proof` 

   
9. Encrypt the AES-256 key.
   1. Encrypt the `F`key bytes using the centralized system's RSA public key using PKCS1v15 mode.  
   2. Encode the resulting bytes using the Base64 encoding method, obtaining `proof_key` 


## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
requestDidPowerUp
```