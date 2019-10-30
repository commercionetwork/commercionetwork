# Requesting a Did deposit
A *Did deposit* is a term we use when referring to the willingness of a user to deposit a specific amount into the 
system's liquidity pool taking such funds from his account balance.  
This action is the first step that must be done when [creating a pairwise Did](../creating-pairwise-did.md). 
After the funds have been properly sent to the liquidity pool, you can request the 
[*power up* of a pairwise Did](./request-did-power-up.md) which will properly send the deposited funds 
(or part of them) from the liquidity pool to the specified private Did, without making impossible to correlate
your public identity to the private one. 

:::tip  
If you wish to know more about the overall pairwise Did creation sequence, please refer to the
[pairwise Did creation specification](../creating-pairwise-did.md) page  
:::    

## Transaction message
```json
{
  "type": "commercio/MsgRequestDidDeposit",
  "value": {
    "recipient": "<Address user that will be able to spend the funds to active his pairwise Did(s)>",
    "amount": [
      {
        "amount": "<Amount to send>",
        "denom": "<Denom of the coin to send>"
      }
    ],
    "proof": "<Deposit proof>",
    "encryption_key": "<Encrypted AES-256 key used to symmetrically encrypt the proof, hex encoded>",
    "from_address": "<Your address>"
  }
}
```

### Fields requirements
| Field | Required |
| :---: | :------: |
| `recipient` | Yes |
| `amount` | Yes |
| `proof` | Yes | 
| `encryption_key` | Yes | 
| `from_address` | Yes | 

### Creating the `proof` and `encryption_key` values
When creating the `proof` and `encryption_key` fields' values, the following steps must be followed. 

1. Create the `signature_json` formed as follow.  
   ```json
   {
     "recipient": "<Did address of the recipient>",
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
     "recipient": "<Did address of the recipient>",
     "timestamp": "<Timestamp>",
     "signature": "<Previously signed data, hex encoded>"
   }
   ```
   
4. Create a random [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) key.

5. Using the AES-256 key generated at point (4), encrypt the `payload`.
   1. Remove all the white spaces and line ending characters. 
   2. Encrypt the resulting string bytes using the AES-256 key.  
      Note that the AES encryption method must be `AES`.
   3. Encode the resulting bytes using the HEX encoding method.  
   
   **Doing such, you will obtain the value of the `proof` field.**
   
6. Encrypt the AES-256 key.
   1. Encrypt the key bytes using the centralized system's public key.  
      Note that the RSA encryption method must be `RSA/ECB/PKCS1Padding`.
   2. Encode the resulting bytes using the HEX encoding method. 

   **Doing this, you will obtain the value of the `encrypted_key` field.**
   

#### Pseudo-coding
```
let signature_json = {
  "recipient": "<Did address of the recipient>",
  "timestamp": "<Timestamp>"
}
let json_signature = sign(removeWhiteSpaces(sort(signature_json)));

let payload = {
  "recipient": signature_json.recipient,
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
requestDidDeposit
```