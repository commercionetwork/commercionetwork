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
fund your identity. Each and every transaction on the blockchain has a cost, and to pay for it you must have some tokens.  
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
        "type": "RsaSignatureKey2018",
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
    }
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
 - key with suffix `#keys-2` must be of type `RsaSignatureKey2018`, and must be a valid RSA PKIX public key.
 
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

1. Create a `value` JSON as specified earlier, excluding only `proof` field. In example include `@context`, `id` and `publicKey` fields:
```javascript
{
 "@context": "https://www.w3.org/ns/did/v1",
 "id": "your DID",
 "publicKey": "your public keys",
}
```
and we will call this json `did_document_unsigned`. 
:::warning
**Note**: There may be fields other than those used in this example such as `service` and many others, and they should always be included in the `did_document_unsigned`. 
:::
1. alphabetically sort the `did_document_unsigned` and remove all the white spaces and line endings characters.
2. obtain hash of resulting string bytes using **SHA-256**. 
3. sign the result of the hashing process using your DID's private key, which you assigned to the `verificationMethod` `proof` JSON field
4. encode the result in **base64** obtaining `signatureValue`.

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

A user who wants to execute a Did Power Up must have previously sent tokens to the public address of the centralized entity **Tk**.
Retriving Did of **Tk** using public endpoint [/government/tumbler](../government/#retrieving-the-tumbler-address) or by command `cncli query government tumbler-address`
  
This action is the second and final step that must be done when [creating a pairwise Did](creating-pairwise-did.md).  

:::tip  
If you wish to know more about the overall pairwise Did creation sequence, please refer to the
[pairwise Did creation specification](creating-pairwise-did.md) page  
:::    

#### Transaction message
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

2. Retrive the public key of external centralized entity **Tk** [resolving its DDO](../id/#reading-a-user-did-document). Retriving Did of **Tk** using public endpoint [/government/tumbler](../government/#retrieving-the-tumbler-address) or by command `cncli query government tumbler-address`
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
9.  using the AES-256 key generated at point (7), encrypt the `payload`:
   1. remove all the white spaces and line ending characters
   2. encrypt the resulting string bytes using the `AES-GCM` mode, `F` as key, obtaining `CIPHERTEXT`
   3. concatenate bytes of `CIPHERTEXT` and `N` and encode the resulting bytes in **base64**, obtaining the `value` `proof` content 
10. encrypt the AES-256 key:
   4. encrypt the `F` key bytes using the centralized entity's RSA public key found in its Did Document, in PKCS1v15 mode.  
   5. encode the resulting bytes in **base64**, obtaining the `value` `proof_key` content 


#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
requestDidPowerUp
```


-------------
## Accessible to the tumbler service

**The following explanations are specific to the tumbler service and should not be taken into account for user transactions**


### Did power up
Once a user has properly **create a Did power up request**, the external centralized entity can
accept such request and fund the pairwise Did specified. 

#### Transaction message
```json
{
  "type": "commercio/MsgRequestDidPowerUp",
  "value": {
    "claimant": "<Address of the Did to fund>",
    "amount": [
      {
        "amount": "<Amount of coins to be sent>",
        "denom": "<Denom of the coin to send>"
      }   
    ],
    "proof": "<proof>",
    "id": "<uuid>",
    "proof_key": "<proof_key>"
  }
}
```

### Fields requirements
| Field | Required |
| :---: | :------: |
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
    "timestamp": "<Timestamp>",
   }
   ```

2. Retrive the public key of external centralized entity **Tk** [resolving its DDO](../id/#reading-a-user-did-document). Retriving Did of **Tk** using public endpoint [/government/tumbler](../government/#retrieving-the-tumbler-address) or by command `cncli query government tumbler-address`
3. Calculate SHA-256 `HASH` of `sender_did`, `pairwise_did` and `timestamp` concatenation
4. Sign in format PKCS1v15 the `HASH` with the RSA private key associated to RSA public key inserted in the DDO. Now we have `SIGN(HASH)`
5. Convert `SIGN(HASH)` in Base64 notation `BASE64(SIGN(HASH))` and use it to add `signature` field 

6. Create a `payload` JSON made as follow:

```json
{
"sender_did": "<User did>",
"pairwise_did": "<Pairwise Did to power up>",
"timestamp": "<Timestamp>",
"signature": "BASE64(SIGN(HASH))",
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



#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
requestDidPowerUp
```
------------------

## Queries


### Reading a user Did Document

#### REST

Endpoint:
```
/identities/{did}
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
  "height":"0",
  "result":{
    "owner": "did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx",
    "did_document": {
        "@context": "https://www.w3.org/ns/did/v1",
        "id": "did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx",
        "publicKey": [
        {
            "id": "did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx#keys-1",
            "type": "RsaVerificationKey2018",
            "controller": "did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx",
            "publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvaM5rNKqd5sl1flSqRHg\nkKdGJzVcktZs0O1IO5A7TauzAtn0vRMr4moWYTn5nUCCiDFbTPoMyPp6tsaZScAD\nG9I7g4vK+/FcImcrdDdv9rjh1aGwkGK3AXUNEG+hkP+QsIBl5ORNSKn+EcdFmnUc\nzhNulA74zQ3xnz9cUtsPC464AWW0Yrlw40rJ/NmDYfepjYjikMVvJbKGzbN3Xwv7\nZzF4bPTi7giZlJuKbNUNTccPY/nPr5EkwZ5/cOZnAJGtmTtj0e0mrFTX8sMPyQx0\nO2uYM97z0SRkf8oeNQm+tyYbwGWY2TlCEXbvhP34xMaBTzWNF5+Z+FZi+UfPfVfK\nHQIDAQAB\n-----END PUBLIC KEY-----\n"
        },
        {
            "id": "did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx#keys-2",
            "type": "RsaSignatureKey2018",
            "controller": "did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx",
            "publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuk6XjcPA9Zjpo3dgtHOz\n97cbDi6sRmoGZOFqBBaLVvGT1Cgi4Hp00I5z7WP13OCeaV6dkZLCTRyuLeMTxbXn\nRVBOMnfjMek0xjX4X3DkRaDXstk6OOlJJa8BbBkrs3xb4xXInyoTYyS//F+Hkzjg\nleZPdUw7Fa1/NMtMBoUDSb93IrrO1RBtvOE7/I+85q8khkL3zT8MfD9X9li+fidD\n/xpyikMt3ZsmYJs417FeB0v6chrQjlcMrqYyKmfEwze1tgh6fVrxYUAdCDtclL9x\n/HsSzCOSwgIMNPwArV5v6lsoeyRy1ufA8MpwBWPKN4St0a6DoTizQmqaLZ/kgwZb\nqQIDAQAB\n-----END PUBLIC KEY-----\n"
        }
        ],
        "proof": {
            "type": "EcdsaSecp256k1VerificationKey2019",
            "created": "2020-04-22T04:23:50.73112321Z",
            "proofPurpose": "authentication",
            "controller": "did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx",
            "verificationMethod": "did:com:pub1addwnpepqt6lnn5v0c3rys49v5v9f4kvcchehnu7kyk8t8vce5lsxfy7e2pxwyvmf6t",
            "signatureValue": "nIgRvObXlF2OIbktZcQJw0UU7zDEku8cEBq7194YOjhEvD5wBZ+TcNu9GNRZucC6OyuplHfK6uo57+3lVQbpgA=="
        },
        "service": null
    }
    }
}
```


### Reading a user Did power up request

#### CLI

```bash
cncli query id power-up-request [Request-id] [flags]
```


#### REST

Endpoint:
```
/powerUpRequest/{id}
```

Parameters  
| Parameter | Description |
| :-------: | :---------- | 
| `id` | Request id |


##### Example 

```
http://localhost:1317/powerUpRequest/28b6f9d1-347a-432b-a3a4-e018617a63d8
```

#### Response
```json
{
  "height":"0",
  "result":{
    "status": null,
    "claimant": "did:com:150jp3tx96frukqg6v870etf02q0cp7em78wu48",
    "amount": [
        {
        "denom": "ucommercio",
        "amount": "100"
        }
    ],
    "proof": "bgdNnyQrh4vKUl7BeZPaWRdmzXMRvckjYJEll95OSx7Qf6IkS+NXE4zhtuWi6lG/dBkZUbZEbTAxqjb3qTMlOm2J5fQH5LqCg66aSZBISczXfNMPjjMKH0+F+WEZ1GoUWlta6PBMCfWn2UdltsAbm+GJia3QUPwx4UNheIejjOYyd8b0dAlXTwms+NaBZ5K4nLVC9nrDA6u7EL6tHU5KlP50XQ6/1mD9IHR9GdexNMw0OTchk4mEWWdgnfdKNKQt8qOQoSOItqfsu0I6jD5w+sXK2tg1Zgc8XDIvERGT85G3qxapddyK87nfxEKJ975mkA9yGxYJIk5xFwofbJL2S4cPDmJBLiWUDuBiL8XQldEm4bQBtk3HCbV7eDCbziB431cZw3ThHJcgWMGJT6WMFs+Hsv6UVtWULiAq4n517AIoaQvw754VswZi/1nYlkozLZCJdjhFZh9WoEVuC84iB8zBCcsKCC7TcbLJJc9Fev/5nqXNBWxMW2Fm1IrVcXb+MhqgwkdQgh/tDpksjxHixOJ830I1gSgpRYI7ig75qQIGQu2mk4ZGlcSTst4v+ksub8I1DOSz9kXb7cxnf9CRW3I+Vj9J7ll4MO2ayIrnK5OajhkBHFwhVuqXz2QtODW3ruEL4ySMC43ManHHXC9IWmI/w6sp5qzJ3+j4DuKJzf1iFW2yCbOMlMuGsFzoHiPEAnI5ohLstufJd4lnLg==",
    "id": "28b6f9d1-347a-432b-a3a4-e018617a63d8",
    "proof_key": "Fti1Z+NVPhZGcTdpTDrQXO1bOlf2FWW5EE61VgzJDJMDn/KwvJ9xrHrEPKtnGRjZX1oRhxDI7BIv7i03hoyG7GDhbymEDQrY3Hia00y52opD8k+9B/WqMa35t0j0lBpEgp2ZeyA8QTxDmQ9kro0uCMsvFv9XBkW34cIfAUJBRFhk/yA37u6wUGve5AG9DUO86yysaPa8Y5c+vohNdynbYgtSc1maDag9E4E2w57YjqiXussyl4bG7l4j7weSKEHYy07Bv3G+VKGJ8T28HKHeVZTO2TWDPPWAJG5+HdQC9D7ME02dPvHzvrApSVIXjT/Sfx0G5YW3flw7T05UOuOmsA=="
    }
}
```



### Reading a user Did power up pending request

#### REST

Endpoint:
```
/pendingPowerUpRequests
```

##### Example 

```
http://localhost:1317/pendingPowerUpRequests
```

#### Response
```json
{
  "height":"0",
  "result": [
    {
        "status": {
          "type": "approved",
          "message": "request approved"
        },
        "claimant": "did:com:150jp3tx96frukqg6v870etf02q0cp7em78wu48",
        "amount": [
        {
            "denom": "ucommercio",
            "amount": "100"
        }
        ],
        "proof": "qlD0VR/nNf2AHwECFmhXPu/U3lQjLbzcimZ59nuSw5TSBPGewUKzh95l2LkDW/uBAJffyUWuwBYTi+eQ+ZKH+w2XbX8HQ7jo6E16END8+vDGzYmuz8bdbHFavk+8r4Ps9cPkpwAeN97xDN478KlVNHdoL5ZR+B4eMYKcmJRGNV2hD6X/DoTD8Lt9HUU3WTa67O7rwhtTgFdRssjJJpCHGkO6TKMG3qDAp0K8OGiAZ7eE9a3cFYX9umuZVhDo/Yen5SBtAmMXpj7cUEKUpsIlRdoTy5vmsNtmaSQxoXQlJqctH19hqC4GCSxMyyJCwIUD5zMO38Udn7MsTg2wI6BDxgrGfiHgiyAuCTNLjZ6FNL9NtEo2l/JoSbEkydzayJvkzTCl8Sl5nQu10XkapZIp7z1eRAh00W8yxxEOHLWoP6y2H5heT8EPfIrV25Y+osTtCbyTnl7z/2WDuwBpKXJ0muHlzElSocfvh0U8Q5HhTXsa0hGKVWuIUcsuBL/goRI6mf+aO8CZ9KKakR3F0brqkJU5yBh8t3v3qYTyFwj4ZB6FQFrFDZvNwwXv1k47NnNTVwUDc9VCHVAYzlftgEIczH2lanzGlURVnp/IutCb8jUx7P5qIM55wG2TUdjsWTQ8xg7QZQQfZ5SnIfWswF7+vG2eFywMjmCbdBalFrhim2qQ2UNNtg9xhBf6wZrYhUYRPdxBFcZyCexldMy9RQ==",
        "id": "3ba73f16-9241-4e13-9879-413becb0818a",
        "proof_key": "WN3zWFva7ZQNUZ5cnO1fumY9SlGhclQzX8lQf9FrsQnrgn39UJx0KApJDyV1nd03C2+ZXuQT6+79D+8RJX7MQw4SXUvQIu8A2Ta3C/KvzvSbUlfVXv2Y5OuBpYHbCZI6GRmXbg0m/bf2CfJLMcKDRPYlgCH86yHPeJ4BR7UQKLwhtNZeuGoYwjgoqmwCs95gbPycmDHCp4nAdrrvehHss6uj8v6JMxny9cQzMk0FDNK+vjOy+ULI/SeVEPex0DOFW3hYJjpvTQv9PpiwfinDyvTAPm1ahtG6A/b9ujYoAINcjAaCAo9e4pY679mWTv70Ii9PTr94U6tQwNOMp64sZQ=="
    }
  ]
}
```

### Reading a user Did power up approved request

#### REST

Endpoint:
```
/pendingPowerUpRequests
```

##### Example 

```
http://localhost:1317/approvedPowerUpRequests
```

#### Response
```json
{
  "height":"0",
  "result": [
    {
        "status": null,
        "claimant": "did:com:150jp3tx96frukqg6v870etf02q0cp7em78wu48",
        "amount": [
        {
            "denom": "ucommercio",
            "amount": "100"
        }
        ],
        "proof": "bgdNnyQrh4vKUl7BeZPaWRdmzXMRvckjYJEll95OSx7Qf6IkS+NXE4zhtuWi6lG/dBkZUbZEbTAxqjb3qTMlOm2J5fQH5LqCg66aSZBISczXfNMPjjMKH0+F+WEZ1GoUWlta6PBMCfWn2UdltsAbm+GJia3QUPwx4UNheIejjOYyd8b0dAlXTwms+NaBZ5K4nLVC9nrDA6u7EL6tHU5KlP50XQ6/1mD9IHR9GdexNMw0OTchk4mEWWdgnfdKNKQt8qOQoSOItqfsu0I6jD5w+sXK2tg1Zgc8XDIvERGT85G3qxapddyK87nfxEKJ975mkA9yGxYJIk5xFwofbJL2S4cPDmJBLiWUDuBiL8XQldEm4bQBtk3HCbV7eDCbziB431cZw3ThHJcgWMGJT6WMFs+Hsv6UVtWULiAq4n517AIoaQvw754VswZi/1nYlkozLZCJdjhFZh9WoEVuC84iB8zBCcsKCC7TcbLJJc9Fev/5nqXNBWxMW2Fm1IrVcXb+MhqgwkdQgh/tDpksjxHixOJ830I1gSgpRYI7ig75qQIGQu2mk4ZGlcSTst4v+ksub8I1DOSz9kXb7cxnf9CRW3I+Vj9J7ll4MO2ayIrnK5OajhkBHFwhVuqXz2QtODW3ruEL4ySMC43ManHHXC9IWmI/w6sp5qzJ3+j4DuKJzf1iFW2yCbOMlMuGsFzoHiPEAnI5ohLstufJd4lnLg==",
        "id": "28b6f9d1-347a-432b-a3a4-e018617a63d8",
        "proof_key": "Fti1Z+NVPhZGcTdpTDrQXO1bOlf2FWW5EE61VgzJDJMDn/KwvJ9xrHrEPKtnGRjZX1oRhxDI7BIv7i03hoyG7GDhbymEDQrY3Hia00y52opD8k+9B/WqMa35t0j0lBpEgp2ZeyA8QTxDmQ9kro0uCMsvFv9XBkW34cIfAUJBRFhk/yA37u6wUGve5AG9DUO86yysaPa8Y5c+vohNdynbYgtSc1maDag9E4E2w57YjqiXussyl4bG7l4j7weSKEHYy07Bv3G+VKGJ8T28HKHeVZTO2TWDPPWAJG5+HdQC9D7ME02dPvHzvrApSVIXjT/Sfx0G5YW3flw7T05UOuOmsA=="
    }
  ]
}
```


### Reading a user Did power up rejected request

#### REST

Endpoint:
```
/rejectedPowerUpRequests
```

##### Example 

```
http://localhost:1317/rejectedPowerUpRequests
```

#### Response
```json
{
  "height":"0",
  "result": [
    {
        "status": {
          "type": "rejected",
          "message": "insufficient fund"
        },
        "claimant": "did:com:150jp3tx96frukqg6v870etf02q0cp7em78wu48",
        "amount": [
        {
            "denom": "ucommercio",
            "amount": "100000000000000"
        }
        ],
        "proof": "S5hyg4slMxm9fK8PTNDs8tHmQcBfWXG0vqrNHLXY5K1qUz3QwZYjR9nzJoNDJh18aPsXper7rNBbyZPOm5K//x8Bqm2EJkdnHd7woa5eFqpziGaHxqvgPaLGspH47tnVilARTeF23L2NVHWcEWuo9U5cWg52l1lOixOG+DehT3vC9KjLqg0YqBoL2u0LTLqQMON4UUjC8JwzT/RMs30OYGsWuLc9s48RtJCQJZ+yAg3U6jZn3OokGwWWjYxF9tAsMR48KilHsPigsa9WPnaAyCMSJ05hOqjBxWiSHYiH1nAefFqHtNFXhJF3LRUCJ2xnSHxJC5Ndj4HFzUjyK4aiV1mtRlRcsqmXU80HEk7IzI74HYpW74F8LzXNsh8Pbl7HXoIzEiOHB5XStFnrxkIL3sYAJGH/pGbX3SxeyfoZhY4ikEyqX3OB7Pat2yHh/63XSPThRVpD7g0gy5N2aKBz3vrHCPhe3QQTzWmKlJOcg1FE5ZtSUEHdVQbm1GD9zP6KZDfbekh9+xU0EFczW9JF/we61LTvMF1KoxaBpL46O/J6ROEOQsb03hLEMadBKxZ+XaqAHiQWKu6G5YH2opNTGKcvSyNfDInOvAygUOfzLgTCWp7JOU09hWBKW1ya2yJNJMZ6q9giEAlqS/qqYy4gAqZKjt7nF0siOb3Vz6zEaXdhCcqrfnNN6n/kFXWz24yAucW+/EHt+hsygEVUZQ==",
        "id": "d423c645-fd50-4841-8138-192ee8e23dde",
        "proof_key": "L0QIWxtHeWeUQhmfWqB2n+MZXFqEYctltilM0j69tBd1drUoUSz/vUkaPadQAdKqtQOD43Py7/JZt5IFyx7iDdphzJEX7bqq+B6nC2DQUeISEiXwtDmJYMp20/N23DY2T7L/Z/dzbxRZDWoUhtr9fRPeJL8NHtPqU9YZw2f1tgMk2t/ZMKtBhYzO5BnF8Crmshjw6b6KA3fK+j7YrmF8fVpVFCdz5jd7cprf5RIqwVjt4w1cYZWeKvGLWeGVX3oiCB67EzXZVUCsD03evr90GDY9qGLfUaWJdBkNjByDotLY0OhrKpcZ+O0IZyZv1+YKx7ZDoPAsEJqpqw4M9bGQRg=="
    }
  ]
}
```


