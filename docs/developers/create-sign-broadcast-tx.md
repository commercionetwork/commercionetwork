# Create, sign and send a transaction
This guide has been made due to the lacking documentation about the offline creation, signing and broadcasting of 
transactions on a Cosmos.network chain.  
All the references about how to perform the actions described below can also be found across the different 
Cosmos.network documentation pages, such as (but not limited to): 

- [Cosmos SDK transaction signing](https://cosmos-staging.interblock.io/docs/clients/service-providers.html#cosmos-sdk-transaction-signing)
- [Cosmos.network RPC APIs](https://cosmos.network/rpc/#/)
- [Cosmos SDK Documentation](https://cosmos.network/docs/)

Please note that the above links **will not** be kept in sync with the frequent updates that the Cosmos developers 
do to their documentation structure. In order to preserve my mental sanity, I will only update this page when 
necessary, so please refer to it in order to always know how to create, sign and broadcast transactions. 

## 1. Message creation
A transaction can contain one or multiple messages. While usually a single message is sent per transaction, 
none denies of sending multiple messages inside the same transaction. Making it short, a message is a simple JSON 
object with some specific fields inside it. 
An example of message object is the following. 

```json
{
  "type": "commercio/MsgSendDocument",
  "value": {
    "sender": "<Your address>",
    "recipient": "<Recipient address>",
    "uuid": "<Document UUID>",
    "content_uri": "<Document content URI>",
    "metadata": {
      "content_uri": "<Metadata content URI>",
      "schema_type": "<Officially recognized schema type>",
      "schema": {
        "uri": "<Metadata schema URI>",
        "version": "<Metadata schema version>"
      },
      "proof": "<Metadata validation proof>"
    },
    "checksum": {
      "value": "<Document content checksum value>",
      "algorithm": "<Checksum algorithm>"
    }
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
The first thing that must be done in order to sign a message, is to retrieve all the needed data. 

#### Account number and sequence
To avoid repetition attacks, all the transactions that are sent inside a Cosmos blockchain must contain a given `account_number` and `sequence value`. Those numbers can be retrieved from the blockchain itself with

```bash
curl https://<NODE_URL>/auth/accounts/<SIGNER_ADDRESS>
```

Supposing we want to read those values for the address `did:com:10k2fswy7ce8mvfdlkt2x0mht326svvuwkxqp30` from our local node, we will then use

```bash
curl http://localhost:1317/auth/accounts/did:com:10k2fswy7ce8mvfdlkt2x0mht326svvuwkxqp30
```

This should print a JSON object similar to this:
```json
{
  "height": "1647",
  "result": {
    "type": "cosmos-sdk/Account",
    "value": {
      "address": "did:com:10k2fswy7ce8mvfdlkt2x0mht326svvuwkxqp30",
      "coins": [
        {
          "denom": "ucommercio",
          "amount": "9999900000000"
        }
      ],
      "public_key": {
        "type": "tendermint/PubKeySecp256k1",
        "value": "A5nUN8qB2iTgQp29lbiZUb6CUoqaxRUNPP3IA9TOKa9J"
      },
      "account_number": "0",
      "sequence": "1"
    }
  }
}
```

Save the `account_number` and `sequence` value in one place, they will be useful later on. 

#### Chain id
The next thing we need to create a message signature, is the chain id. This value can be retrieved using the following HTTP endpoit. 

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
  "protocol_version": {
    "p2p": "7",
    "block": "10",
    "app": "0"
  },
  "id": "233d4e3a9c68bc724f3deb41197e07052e75ab95",
  "listen_addr": "tcp://0.0.0.0:26656",
  "network": "test-chain-eDSCSs",
  "version": "0.30.2",
  "channels": "4020212223303800",
  "moniker": "testchain",
  "other": {
    "tx_index": "on",
    "rpc_address": "tcp://0.0.0.0:26657"
  }
}
```

Save the `network` value for later.

### 2.2 Signature data creation
In order to create the data that has to be signed, you will need the following values: 
1. JSON representation of the message, created inside the [message creation section](#1-message-creation)
2. `account_number` and `sequence` values, obtained inside the [account number and sequence section](#account-number-and-sequence)
3. Id of the chain, obtained inside the [chain id section](#chain-id)

After you've retrieved (or created) all those values, you are ready to start. 

The signature data is a JSON object formed as follows: 

```json
{
  "account_number": "<ACCOUNT_NUMBER>",
  "chain_id": "<CHAIN_ID>",
  "fee": {
    FEE_OBJECT
  },
  "memo": "<TX_MEMO>",
  "msgs": [
    MESSAGE_OBJECT
  ],
  "sequence": "<SEQUENCE>"
}
```

Using the same example data of the previous sections, a valid signature data will then look like the following. 

```json
{
  "account_number": "0",
  "chain_id": "test-chain-eDSCSs",
  "fee": {
    "amount": [],
    "gas": "20000"
  },
  "memo": "",
  "msgs": [
    {
      "type": "commercio/MsgSendDocument",
      "value": {
        "sender": "<Your address>",
        "recipient": "<Recipient address>",
        "uuid": "<Document UUID>",
        "content_uri": "<Document content URI>",
        "metadata": {
          "content_uri": "<Metadata content URI>",
          "schema_type": "<Officially recognized schema type>",
          "schema": {
            "uri": "<Metadata schema URI>",
            "version": "<Metadata schema version>"
          },
          "proof": "<Metadata validation proof>"
        },
        "checksum": {
          "value": "<Document content checksum value>",
          "algorithm": "<Checksum algorithm>"
        }
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
The `fee` object contains the details of the fee that the transaction creator and signer will pay when broadcasting the transaction to the blockchain. Usually this object's content is set as the default `gas` value and no `amount` content, i.e

```json
{
  "amount": [],
  "gas": "20000"
}
```

### 2.3 Signing the data
Once you've create the JSON object containing the data to sign, it is now time to sign them. In order to do so, the following steps should be made:

1. Convert the JSON object to it's compact and alphabetically ordered representation.  
   This means that the keys of the object should be alphabetically sorted, and any white space should be trimmed. The above JSON should then look like this:  
   ```
   {"account_number":"0","chain_id":"test-chain-eDSCSs","fee":{"amount":[],"gas":"20000"},"memo":"","msgs":[{"type":"commercioid/SetIdentity","value":{"ddo_reference":"hkbmhbmbmbmnbmb","did":"lkh,mjhmjhmj,hmjh","owner":"did:com:13jckgxmj3v8jpqdeq8zxwcyhv7gc3dzmrqqger"}}],"sequence":"1"}
   ```

2. Compute the SHA-256 hash of the JSON content's byte array representation.  
   ```
   sha256([]byte(compact_json))
   ```

3. Sign the hash bytes with the signer's private key.  
   ```
   sign([]byte(has))
   ```

4. Encode the resulting signature as a Base64 string.  
   ```
   base64([]byte(signature))
   ```

All in one, the operations should be

```
let compact_json = compact(json)
let hashed_json = sha256([]byte(compact_json))
let signature_bytes = sign([]byte(hashed_json))
let base64_signature = base64([]byte(signature_bytes))
```

### 2.4 Signature object creation
Once we have the Base64 signature representation, we can finally create the signature object that we will later use during the transaction creation. In order to do so, a JSON object with the following fields should be created: 

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
| `value` | string | yes | See [Public key value retrieving](#public-key-encoding) for more details |
| `signature` | string | yes | Base64 encoded value of the signature, as obtained inside the [signing the data section](#23-signing-the-data) |

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

All in one, the operations should be 
```
let compressed = encode(public_key.point, compress = true)
let base64 = Base64(compressed)
```

More on how compressed keys can be obtained is described on the following StackExchange question:  
[How are compressed PubKeys generated?](https://bitcoin.stackexchange.com/questions/69315/how-are-compressed-pubkeys-generated/69322)

---

## 3. Transaction broadcasting
The last step to correctly insert a transaction into the blockchain is broadcasting it. In order to do so, we need to create the transaction body and later sending it with an HTTP request to a full node. 

### 3.1. Creating the JSON object representation of the transaction
The first thing that needs to be done in order to create a transaction, is to create a JSON object formed as follows.

```json
{
  "msg": [
    { MESSAGE_OBJECT }
  ],
  "fee": {
    FEE_OBJECT
  },
  "signatures": [
    { SIGNATURE_OBJECT }
  ],
  "memo": ""
}
```

**Fields**

| Field | Type | Required | Description | 
| :---- | :--- | :------- | :---------- |
| `msg` | array | yes | Contains all the messages that have been signed, each one as a JSON object. These objects are the same that has been created inside the [message creation section](#1-message-creation) | 
| `fee` | object | yes | Contains the fees that the transaction signer will pay when sending the transaction itself. This is the same object that is put inside the signature data we've seen inside the [signature data creation section](#22-signature-data-creation) |
| `signatures` | array | yes | Contains all the signatures of the messages that will be sent along with the request. The object definition is the one we've seen in the previous [signature object creation section](#24-signature-object-creation) |
| `memo` | string | no | This must contain the same value of the `memo` field that is present inside the signature data we've seen on the previous [signature object creation section](#24-signature-object-creation) | 

#### Example
```json
{
  "msg": [
    {
      "account_number": "0",
      "chain_id": "test-chain-eDSCSs",
      "fee": {
        "amount": [],
        "gas": "20000"
      },
      "memo": "",
      "msgs": [
        {
          "type": "commercio/MsgSendDocument",
          "value": {
            "sender": "<Your address>",
            "recipient": "<Recipient address>",
            "uuid": "<Document UUID>",
            "content_uri": "<Document content URI>",
            "metadata": {
              "content_uri": "<Metadata content URI>",
              "schema_type": "<Officially recognized schema type>",
              "schema": {
                "uri": "<Metadata schema URI>",
                "version": "<Metadata schema version>"
              },
              "proof": "<Metadata validation proof>"
            },
            "checksum": {
              "value": "<Document content checksum value>",
              "algorithm": "<Checksum algorithm>"
            }
          }
        }
      ],
      "sequence": "1"
    }
  ],
  "fee": {
    "amount": [],
    "gas": "200000"
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
Once the transaction JSON has been created, we can now broadcast it sending it to a full node. In order to do so, a POST HTTP request should be made. 

#### 3.2.1 Creating the request body
The request body must be a JSON object formed as follows: 

```json
{ 
  "tx": {
    <TX_BODY>
  },
  "mode": "block/sync/async"
}
```

**Fields**

| Field | Type | Required | Descrizione | 
| :---- | :--- | :------- | :---------- | 
| `tx` | object | yes | Contains the data of the transaction that has been created inside the [section 3.1](#31-creating-the-json-object-representation-of-the-transaction) | 
| `mode` | string | yes | Tells when the node should return the answer. The `async` one tells the node to return immediately, `sync` tells to return after the transaction has been validated, while `block` waits till the transaction has been successfully broadcasted and returns the information of the block containing it. | 


#### 3.2.2 Performing the request
Once the JSON body has been created, the request can be made. The endpoint is the following:

```
curl -X POST https://<NODE_URL>/txs -d @<JSON_BODY_FILE>
```

Where `<JSON_BODY_FILE>` represents the path to the file containing the request body.

If we are running a local node, it should then be

```
curl -X POST http://localhost:1317/txs -d @~/request_body.json
```


#### Example request body
```json
{
  "tx": {
    "msg": [
      {
        "account_number": "0",
        "chain_id": "test-chain-eDSCSs",
        "fee": {
          "amount": [],
          "gas": "20000"
        },
        "memo": "",
        "msgs": [
          {
            "type": "commercio/MsgSendDocument",
            "value": {
              "sender": "<Your address>",
              "recipient": "<Recipient address>",
              "uuid": "<Document UUID>",
              "content_uri": "<Document content URI>",
              "metadata": {
                "content_uri": "<Metadata content URI>",
                "schema_type": "<Officially recognized schema type>",
                "schema": {
                  "uri": "<Metadata schema URI>",
                  "version": "<Metadata schema version>"
                },
                "proof": "<Metadata validation proof>"
              },
              "checksum": {
                "value": "<Document content checksum value>",
                "algorithm": "<Checksum algorithm>"
              }
            }
          }
        ],
        "sequence": "1"
      }
    ],
    "fee": {
      "amount": [],
      "gas": "200000"
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
  },
  "mode": "sync"
}
```