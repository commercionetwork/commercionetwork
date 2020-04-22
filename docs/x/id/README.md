# Id
The `id` module is the one that allows you to create a new identity and associate to it a 
Did Document.

## Transactions
Using the `id` module you can perform the following transactions.

**Accessible to everyone**

### Creating an identity
First of all, let's define what an **identity** is inside the Commercio Network blockchain.  

> An identity is the method used inside the Commercio Network blockchain in order to identify documents' sender.

In order to create an identity, you simply have to create a Commercio Network address, which will have the 
following form: 

```
did:com:<unique part>
```

In order to do so, you can use the CLI and execute the following command: 

```bash
cncli keys add <key-name>
``` 

You will be required to set a password in order to safely store the key on your computer.  

:::warning
Please note that password will be later asked you when signing the transactions so be sure you remember it.
:::  

After inserting the password, you will be shown the mnemonic that can be used in order to import your account 
(and identity) into a wallet. 

```
- name: jack
  type: local
  address: did:com:13jckgxmj3v8jpqdeq8zxwcyhv7gc3dzmrqqger
  pubkey: did:com:pub1addwnpepqfdl6s8hdwdya9zvn5wtx8ty3qsqqqd2ddvygc5zutnrryh5x9ju73jdfg8
  mnemonic: ""
  threshold: 0
  pubkeys: []


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

scorpion what indoor keen topic cricket uphold inch cactus six suffer coin popular honey vendor clown day twin during vague midnight emerge man inform
```

#### Using an identity
Once you have created it, in order to start performing a transaction with your identity you firstly have to 
fund your identity. Each and every transaction on the blockchain has a cost, and to pay for it you have to have some 
tokens.  
If you want to receive some tokens in **Test-net**, please use faucet service or tell us inside our [official Telegram group](https://t.me/commercionetwork) 
and we will send you some as soon as possible.





### Associating a Did Document to your identity 
Being your account address a Did, using the Commercio Network blockchain you can associate to it a Did Document
containing the information that are related to your public (or private) identity.  
In order to do so you will need to perform a transaction and so your account must have first received some tokens. To
know how to get them, please take a look at the [*"Using an identity"* section](create-an-identity.md#using-an-identity). 

#### Transaction message
In order to properly send a transaction to set a DID Document associating it to your identity, you will need
to create and sign the following message:

```javascript
{
  "type": "commercio/MsgSetIdentity",
  "value": {
    "@context": "https://www.w3.org/ns/did/v1",
    "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
    "publicKey": [
      {
        "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc#keys-1",
        "type": "RsaVerificationKey2018",
        "controller": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
        "publicKeyPem": "-----BEGIN PUBLIC KEY----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB-----END PUBLIC KEY-----\r\n"
      },
      {
        "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc#keys-2",
        "type": "RsaSignature2018",
        "controller": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
        "publicKeyPem": "-----BEGIN PUBLIC KEY----MIGfM3TvO3Ku3PJgZ9PO4qRw7+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCQVvTkCbc9A0GCSqGSIbqd4pNXtgbfbwJGviZ6kQIDAQAB-----END PUBLIC KEY-----\r\n"
      }
    ],
    "proof": {
      "type": "EcdsaSecp256k1VerificationKey2019",
      "created": "2019-02-08T16:02:20Z",
      "proofPurpose":"authentication",
      "controller": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
      "verificationMethod": "<did bech32 pubkey>",
      "signatureValue": "QNB13Y7Q91tzjn4w=="
    },
  }
}
```

##### `value` fields requirements

| Field | Required | 
| :---: | :------: | 
| `@context` | Yes (Must be `https://www.w3.org/ns/did/v1`) |
| `id` | Yes |
| `publicKey` | Yes |
| `proof` | Yes |

The `id` field represents the DID you want to associate the provided identity to.

The `publicKey` field represents the public keys users can use to communicate safely with you.

Each key **must** have an `id` field defined by the concatenation of the content of the `id` field, along with a `#keys-NUMBER` suffix, where `NUMBER` must be an integer.

The `controller` key field must be equal to the `id` field content.

The commercio.network blockchain requires at least two keys, defined in the following way:

 - key with suffix `#keys-1` must be of type `RsaVerificationKey2018`, and must be a valid RSA PKIX public key;
 - key with suffix `#keys-2` must be of type `RsaSignature2018`, and must be a valid RSA PKIX public key.
 
A `commercio/MsgSetIdentity` transaction that **doesn't** meet these requirements will be discarded.

##### Proof fields requirements

| Field | Required | Value | 
| :---: | :------: | :------: | 
| `type` | Yes | must always be `EcdsaSecp256k1VerificationKey2019` |
| `created` | Yes | creation date in UTC format |
| `proofPurpose` | Yes | must always be `authentication` |
| `controller` | Yes | same value specified in the `id` field |
| `verificationMethod` | Yes | bech32-encoded public key associated with the address specified in the `id` field |
| `signatureValue` | Yes | see explaination below |

##### Creating the `signatureValue` value

In order to create `signatureValue`, the following steps must be followed

1. Create a `value` JSON as specified earlier, including only the `@context`, `id` and `publicKey` fields:
```javascript
{
 "@context": "https://www.w3.org/ns/did/v1",
 "id": "your DID",
 "publicKey": "your public keys",
}
```
2. alphabetically sort the `did_document_unsigned` and remove all the white spaces and line endings characters.
3. obtain hash of resulting string bytes using **SHA-256**. 
4. sign the result of the hashing process using your DID's public key, which you assigned to the `verificationMethod` `proof` JSON field
5. encode the result in **base64** obtaining `signatureValue`.

The signature commercio.network accepts is `EcdsaSecp256k1VerificationKey2019`, which is a type of elliptic-curve signature scheme.

The signature format produced in step 4, must be of the `r || s` kind, otherwise the identity creation **will** fail.

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setIdentity
```  




### Requesting a Did Power Up
A *Did Power Up* is the expression we use when referring to the willingness of a user to move a specified amount of tokens 
from external centralized entity to one of his
private pairwise Did, making them able to send documents (which indeed require the user to spend some tokens as fees). 

A user who wants to execute a Did Power Up must have previously sent tokens to the public address of the centralized entity.
  
This action is the second and final step that must be done when [creating a pairwise Did](../creating-pairwise-did.md).  

:::tip  
If you wish to know more about the overall pairwise Did creation sequence, please refer to the
[pairwise Did creation specification](../creating-pairwise-did.md) page  
:::    

### Transaction message
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

#### `value` fields requirements
| Field | Required |
| :---: | :------: |
| `claimant` | Yes |
| `amount` | Yes |
| `proof` | Yes | 
| `id` | Yes | 
| `proof_key` | Yes | 


#### Creating the `proof` value

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


#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
requestDidPowerUp
```



**Accessible to the tumbler service**


### Did power up
Once a user has properly [create a Did power up request](./request-did-power-up.md), the external centralized entity can
accept such request and fund the pairwise Did specified. 

#### Transaction message
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
    "signer": "<Tumbler address>"
  }
}
```

##### Fields requirements
| Field | Required |
| :---: | :------: |
| `recipient` | Yes | 
| `amount` | Yes |
| `activation_reference` | Yes |
| `signer` | Yes |

###### `amount`
| Field | Required |
| :---: | :------: |
| `amount` | Yes |
| `denom` | Yes | 

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
powerUpDid
```

### Change Did power up status (wip)





## Queries


### Reading a user Did Document

#### REST

Endpoint:
```
/identities/${did}
```

Parameters  
| Parameter | Description |
| :-------: | :---------- | 
| `did` | Address of the user for which to read the Did Document |

##### Example 

```
http://localhost:1317/identities/did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke
```

#### Response
```json
{
  "height": "0",
  "result": {
    "owner": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
    "did_document": {
      "@context": "https://www.w3.org/2019/did/v1",
      "id": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
      "publicKey": [
        {
          "id": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-1",
          "type": "Secp256k1VerificationKey2018",
          "controller": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
          "publicKeyHex": "028b722575fe90167ccae99ab06a9f155fbb11f2045b35a452e635efc57182624a"
        },
        {
          "id": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-2",
          "type": "RsaVerificationKey2018",
          "controller": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
          "publicKeyHex": "3082010a0282010100a365a9d0ac3bc66e256268ef8126b1c9acbb977ecaa140f0f738f28e6645e038bf84ccce0e53052726c8cad0cd3eeacfd2036959917355765ecf43ebc487889dad4e388787b231c8351cafc5394572046942642f6062566a90dc309f4fe910707ed6bbb310e0fe879ee31d4ed3eb74ffff0eda3b8f0cfcfb70392ce936143c13cdcb6a11bd997d0405e7d1bcd043315e7851c30bacce8985d006f794bcb50b861b90c580fee6958a668983c0ba06bd70d6165b1b73b6666ecc0818cbd69bc09aab6d497fe5c58e46bb1b4a795bb99a40d5793fd23588d8d804e5473569bfd1454d1003c2bc74de8ef9db35a00911446df32e2071a964c7b606ffc665a5d879bd0203010001"
        }
      ],
      "authentication": [
        "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-1"
      ],
      "proof": {
        "type": "LinkedDataSignature2015",
        "created": "2019-11-11T13:44:48.829363Z",
        "creator": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-1",
        "signatureValue": "3045022100d2070318a640077c202137c0ac4c64ea2a9274baf3b8b20f7d2d526d881b9a2602204ca3d0141fb1ebe4a87d9efdd4f03f89862aeadb192e87183671d1a2ec9dec11"
      },
      "service": null
    }
  }
}
```


### Reading a user Did power up request

#### CLI

```bash
cncli query docs received-documents [address]
```

