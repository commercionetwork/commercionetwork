# Requesting a Did power up
A *Did power up* is a term we use when referring to the willingness of a user to move a specified amount of tokens 
from the liquidity pool (in which he has previously [deposited something](./request-did-deposit.md)) to one of his
private pairwise Did, making them able to send documents (which indeed require the user to spend some tokens as fees).  
  
This action is the second and final step that must be done when [creating a pairwise Did](../creating-pairwise-did.md).  

:::tip  
If you wish to know more about the overall pairwise Did creation sequence, please refer to the
[pairwise Did creation specification](../creating-pairwise-did.md) page  
:::    

## Transaction message
```json
{
  "type": "commercio/MsgRequirePowerUpDid",
  "value": {
    "claimant": "<Address that is able to spend the funds (the recipient used during the deposit procedure)>",
    "amount": [
      {
        "amount": "<Amount to send to the pairwise Did>",
        "denom": "<Denom of the coin to send>"
      }
    ],
    "proof": "<Power up proof>",
    "encryption_key": "<Encrypted AES-256 key used to symmetrically encrypt the proof, hex encoded>"
  }
}
```

### Fields requirements
| Field | Required |
| :---: | :------: |
| `claimant` | Yes |
| `amount` | Yes |
| `proof` | Yes | 
| `encryption_key` | Yes | 


### Creating the `proof` and `encryption_key` values
When creating the `proof` and `encryption_key` fields' values, the following steps must be followed. 

1. Create the `signature_json` formed as follow.  
   ```json
   {
     "pairwise_did": "<Pairwise Did to power up>",
     "timestamp": "<Timestamp>"
   }
   ```
   
2. Sign the `signature_json`. 
   1. Alphabetically sort the `signature_json`
   2. Remove all the white spaces and line endings characters. 
   3. Sign the resulting string bytes using the private signature key. 
   
3. Create a `payload` JSON made as follow:
   ```json
   {
     "pairwise_did": "<Pairwise Did to power up>",
     "timestamp": "<Timestamp>",
     "signature": "<Previously signed data, hex encoded>"
   }
   ```
   
4. Create a random [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) key.

5. Using the AES-256 key generated at point (4), encrypt the `payload`.
   1. Remove all the white spaces and line ending characters. 
   2. Encrypt the resulting string bytes using the AES-256 key.
   3. Encode the resulting bytes using the HEX encoding method.  
   
   **Doing such, you will obtain the value of the `proof` field.**
   
6. Encrypt the AES-256 key.
   1. Encrypt the key bytes using the centralized system's public key.
   2. Encode the resulting bytes using the HEX encoding method. 

   **Doing this, you will obtain the value of the `encrypted_key` field.**
   

#### Pseudo-coding
```
let signature_json = {
  "pairwise_did": "<Pairwise Did to power up>",
  "timestamp": "<Timestamp>"
}
let json_signature = sign(removeWhiteSpaces(sort(signature_json)));

let payload = {
  "pairwise_did": signature_json.pairwise_did,
  "timestamp": signature_json.timestamp,
  "signature": hex.encode(json_signature)
}

let aes_key = AES256.random()
let encrypted_payload = aes_key.encrypt(removeWhiteSpaces(payload))

let proof = hex.encode(encrypted_payload)

let encrypted_key = hex.encode(rsa_pub_key.encrypt(aes_key)) 
```

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
requestDidPowerUp
```