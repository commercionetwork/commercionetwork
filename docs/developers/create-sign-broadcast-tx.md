# Create, sign and send a transaction

This guide has been made due to the lacking documentation about the offline creation, signing and broadcasting of 
transactions on a Cosmos.network chain.  
All the references about how to perform the actions described below can also be found across the different 
Cosmos.network documentation pages, such as (but not limited to): 

- [Cosmos SDK transaction signing](https://cosmos-staging.interblock.io/docs/clients/service-providers.html#cosmos-sdk-transaction-signing)
- [Cosmos.network RPC APIs](https://cosmos.network/rpc/#/)
- [Cosmos SDK Documentation](https://cosmos.network/docs/)

Please note that the above links **will not** be kept in sync with the frequent updates that the Cosmos developers 
do to their documentation structure.

## 1. Message creation
A transaction can contain one or multiple messages. 

A message is a simple JSON object with some specific fields inside it. 

An example of message object is the following:

```json
{
  "type": "cosmos-sdk/MsgSend",
  "value": {
    "from_address": "<Your address>",
    "to_address": "<Recipient address>",
    "amount": [
      {
        "denom" : 10,
        "amount" : "ucommercio"
      }
    ]
  }
}
```

**Fields**

| Field | Type | Required | Description |
| :---- | :--- | :------- | :---------- |
| `type` | string | yes | Contains the type of the message that is represented inside the object. In order to see all the possible `type` values, please refer to the [proper page](message-types.md) |
| `value` | object | yes | Contains the real value of the message. For each message type, a different set of fields will be present inside this object. In order to know which fields should be sent inside this object for each message type, please refer to the [proper page](message-types.md) |

---

## 2. Message signing
Once created, messages must be signed. In order to do so, the following data will be used.

1. Private key associated to the account.  
   This can be retrieved from the 12/24 key words using the [BIP-39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) specification
2. Address of the signer.  
   This can be retrieved using the [BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki) specification, setting `did:com` as the human readable part and using the `m/44'/118'/0'/0/0` derivation path.

### 2.1 Data retrieval
Some data needs to be retrieved in order to correctly sign a transaction.

#### Account number and sequence
To avoid repetition attacks, all the transactions sent inside a Cosmos blockchain must contain a given `account_number` and `sequence value`. 

They can be retrieved from the blockchain itself:

```bash
curl https://<NODE_URL>/auth/accounts/<SIGNER_ADDRESS>
```

Supposing we want to read those values for the address `did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen` from our local node, we will then use

```bash
curl http://localhost:1317/auth/accounts/did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen
```

This should print a JSON object similar to this:
```json
{
  "height": "4866",
  "result": {
    "type": "cosmos-sdk/Account",
    "value": {
      "address": "did:com:1l9rr5ck7ed30ny3ex4uj75ezrt03gfp96z7nen",
      "coins": [
        {
          "denom": "uccc",
          "amount": "10000000000000"
        },
        {
          "denom": "ucommercio",
          "amount": "9999899990000"
        }
      ],
      "public_key": "did:com:pub1addwnpepqw6amy77xennkrkh3d32pz8ykr5kvuwx97w5ychn87ett8m2dzhzzxyynp4",
      "account_number": 0,
      "sequence": 1
    }
  }
}
```

Keep track of `account_number` and `sequence`, they'll be used later. 

#### Chain id

Next thing we need is the chain ID, which can be retrieved using the `/node_info` endpoint:

```bash
curl https://<NODE_URL>/node_info
```

If we are using a local node, we can use

```bash
curl http://localhost:1317/node_info
```

This should return a JSON object similar to the following one. 

```json
{
  "node_info": {
    "protocol_version": {
      "p2p": "7",
      "block": "10",
      "app": "0"
    },
    "id": "f5920c1e69fb917eae04ed270b447f1af3129b8b",
    "listen_addr": "tcp://0.0.0.0:26656",
    "network": "testnet",
    "version": "0.33.3",
    "channels": "4020212223303800",
    "moniker": "testchain",
    "other": {
      "tx_index": "on",
      "rpc_address": "tcp://127.0.0.1:26657"
    }
  },
  "application_version": {
    "name": "commercionetwork",
    "server_name": "cnd",
    "client_name": "cndcli",
    "version": "2.1.0-8-gc37c4748",
    "commit": "c37c474812cb80b13622462a09bf3998c634e875",
    "build_tags": "netgo,ledger",
    "go": "go version go1.14 darwin/amd64"
  }
}
```

The value assigned to the `network` value is our chain ID, which will be `testnet`.

### 2.2 Signature data creation

In order to create the payload which will be signed later, the following values are needed:
1. JSON representation of the message, created inside the [1. message creation section](#_1-message-creation)
2. `account_number` and `sequence` values, obtained inside the [2.1 account number and sequence section](#account-number-and-sequence)
3. chain ID, obtained inside the [2.1 chain id section](#chain-id)

The signature data is a JSON object formed as follows: 

```json
{
  "account_number": "<ACCOUNT_NUMBER>",
  "chain_id": "<CHAIN_ID>",
  "fee": {
    "<your fees>"
  },
  "memo": "<TX_MEMO>",
  "msgs": [
    "<your messages>"
  ],
  "sequence": "<SEQUENCE>"
}
```

Using the same example data of the previous sections, a valid signature data will then look like the following. 

```json
{
  "account_number": "0",
  "chain_id": "testnet",
  "fee": {
    "amount": [
      {
        "denom" : 10000,
        "amount" : "ucommercio"
      }
    ],
    "gas": "20000"
  },
  "memo": "",
  "msgs": [
    {
      "type": "cosmos-sdk/MsgSend",
      "value": {
        "from_address": "<Your address>",
        "to_address": "<Recipient address>",
        "amount": [
          {
            "denom" : 10,
            "amount" : "ucommercio"
          }
        ]
      }
    }
  ],
  "sequence": "1"
}
```

#### Important notes
##### Order of the fields
When serializing them to the JSON format, the fields **must be in alphabetical order**, that means that they should look exactly like the example above.

##### The fee object
The `fee` object contains the fees that the transaction creator and signer will pay when broadcasting it to the blockchain:

```json
{
  "amount": [
    "<coin>"
  ],
  "gas": "20000"
}
```

The `amount` array contains object of type `coin`:

```json
{
  "denom": "<token denom>",
  "amount": "<integer amount of denom tokens>"
}
```

### 2.3 Signing the data

Once you've create the JSON object containing the data to sign, it's time to sign them.

In order to do so, the following steps must be followed:

1. Convert the JSON object to it's compact and alphabetically ordered representation.  

   This means that the keys of the object should be alphabetically sorted, and any white space should be removed.
    
   The above JSON should then look like this:  
   
   ```json
   {"account_number":"0","chain_id":"testnet","fee":{"amount":[{"denom":10000,"amount":"ucommercio"}],"gas":"20000"},"memo":"","msgs":[{"type":"cosmos-sdk/MsgSend","value":{"amount":[{"denom":10,"amount":"ucommercio"}],"from_address":"<Youraddress>","to_address":"<Recipientaddress>"}}],"sequence":"1"}
   ```

2. Compute the SHA-256 hash of the JSON content's byte array representation.  
   
   ```
   sha256([]byte(compact_json))
   ```

3. Sign the hash bytes with the signer's private key.  
   
   ```
   sign([]byte(hash))
   ```

4. Encode the resulting signature as a Base64 string.  
   
   ```
   base64([]byte(signature))
   ```


### 2.4 Signature object creation
Once we have the base64 signature representation, we can finally create the signature object that we will later use during the transaction creation. 

In order to do so, a JSON object with the following fields should be created: 

```json
{
  "pub_key": {
    "type": "<PUB_KEY_TYPE>",
    "value": "<PUB_KEY_VALUE>"
  },
  "signature": "<BASE64_SIGNATURE_VALUE>"
}
```

**Fields**

| Field | Type | Required | Description | 
| :---- | :--- | :------- | :---------- |
| `type` | string | yes | Contains the type of the public key associated with the signing key we used before. The supported types are `tendermint/PubKeySecp256k1` and `tendermint/PubKeyEd25519` |
| `value` | string | yes | See [Public key value encoding](#public-key-encoding) for more details |
| `signature` | string | yes | Base64 encoded value of the signature, as obtained inside the [2.3 signing the data section](#23-signing-the-data) |

#### Public key encoding

In order to properly encode the public key value of the signer, the following steps should be followed. 

1. Obtain the compressed format of the public key point value.  

   ```
   let compressed = encode(public_key.point, compressed = true)
   ```

2. Encode the compressed value as a Base64 string.  

   ```
   let encoded = Base64(compressed)
   ```

More on how compressed keys can be obtained is described on the following StackExchange question:  
[How are compressed PubKeys generated?](https://bitcoin.stackexchange.com/questions/69315/how-are-compressed-pubkeys-generated/69322)

---

## 3. Transaction broadcasting

The last step to correctly insert a transaction into the blockchain is broadcasting it. 

In order to do so, we need to create the transaction body and later sending it with an HTTP request to a full node. 

### 3.1. Creating the JSON object representation of the transaction

The first thing that needs to be done in order to create a transaction, is to create a structure like the following JSON object:

```json
{
  "msg": [
    { "MESSAGE_OBJECT" }
  ],
  "fee": {
    "FEE_OBJECT"
  },
  "signatures": [
    { "SIGNATURE_OBJECT" }
  ],
  "memo": ""
}
```

**Fields**

| Field | Type | Required | Description | 
| :---- | :--- | :------- | :---------- |
| `msg` | array | yes | Contains all the messages that have been signed, each one as a JSON object. These objects are the same that has been created inside the [1. message creation section](#_1-message-creation) | 
| `fee` | object | yes | Contains the fees that the transaction signer will pay when sending the transaction itself. This is the same object that is put inside the signature data we've seen inside the [2.2 signature data creation section](#_2-2-signature-data-creation) |
| `signatures` | array | yes | Contains all the signatures of the messages that will be sent along with the request. The object definition is the one we've seen in the previous [2.4 signature object creation section](#_2-4-signature-object-creation) |
| `memo` | string | no | This must contain the same value of the `memo` field that is present inside the signature data we've seen on the previous [2.2 signature data creation section](#_2-2-signature-data-creation) | 

#### Example
```json
{
  "msg": [
    {
      "type": "cosmos-sdk/MsgSend",
      "value": {
        "from_address": "<Your address>",
        "to_address": "<Recipient address>",
        "amount": [
          {
            "denom" : 10,
            "amount" : "ucommercio"
          }
        ]
      }
    }
  ],
  "fee": {
    "amount": [
      {
        "denom" : 10000,
        "amount" : "ucommercio"
      }
    ],
    "gas": "20000"
  },
  "signatures": [
    {
      "pub_key": {
        "type": "tendermint/PubKeySecp256k1",
        "value": "AnFSuINPl9229iZdH9z2C9vVi7acnM7mM02Z9AEtsvnj"
      },
      "signature": "DTcJz2V0JxpcdtAlg/pavyB/k+0RnbgulMjIGHtl3g4LwHrG7vnZ6eYll6FkVRkjSB2VSNrLdYiWbriB1Y8KTQ=="
    }
  ],
  "memo": ""
}
```

### 3.2 Broadcasting the transaction

Once the transaction JSON has been created, we can now broadcast it sending it to a full node. 

In order to do so, a POST HTTP request should be made. 

#### 3.2.1 Creating the request body
The request body must be a JSON object formed as follows: 

```json
{ 
  "tx": {
    "<TX_BODY>"
  },
  "mode": "<broadcast mode>"
}
```

**Fields**

| Field | Type | Required | Descrizione | 
| :---- | :--- | :------- | :---------- | 
| `tx` | object | yes | Contains the data of the transaction that has been created inside the [section 3.1](#_3-1-creating-the-json-object-representation-of-the-transaction) | 
| `mode` | string | yes | Tells when the node should return the answer. | 

The `mode` field must assume one of those values:

 - `async`: the node to returns immediately, without waiting for the validation process
 - `sync`: the node returns after the transaction has been validated, without waiting for block inclusion
 - `block`: waits until the transaction has been successfully verified and included in a block, and returns the information of the block containing it 


#### 3.2.2 Performing the request

Once the JSON body has been created, it can be broadcasted to the node. The endpoint is the following:

```
curl -X POST https://<NODE_URL>/txs -d @<JSON_BODY_FILE>
```

Where `<JSON_BODY_FILE>` represents the path to the file containing the request body.

For example, to broadcast the `request_body.json` file to our local node:
 
```
curl -X POST http://localhost:1317/txs -d @request_body.json
```


#### Example request body
```json
{
  "tx": {
    "msg": [
      {
        "type": "cosmos-sdk/MsgSend",
        "value": {
          "from_address": "<Your address>",
          "to_address": "<Recipient address>",
          "amount": [
            {
              "denom" : 10,
              "amount" : "ucommercio"
            }
          ]
        }
      }
    ],
    "fee": {
      "gas": "string",
      "amount": [
        {
          "denom": "ucommerio",
          "amount": "10000"
        }
      ]
    },
    "memo": "string",
    "signature": {
      "signature": "MEUCIQD02fsDPra8MtbRsyB1w7bqTM55Wu138zQbFcWx4+CFyAIge5WNPfKIuvzBZ69MyqHsqD8S1IwiEp+iUb6VSdtlpgY=",
      "pub_key": {
        "type": "tendermint/PubKeySecp256k1",
        "value": "Avz04VhtKJh8ACCVzlI8aTosGy0ikFXKIVHQ3jKMrosH"
      }
    }
  },
  "mode": "sync"
}
```