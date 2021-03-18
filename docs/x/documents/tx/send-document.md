# Sending a document
In order to send a document you are required to have an identity with some tokens inside it.   

::: tip  
To know what an identity is, how to create it and how to get tokens, please refer to the 
[*"Creating an identity"* section](../../id/tx/create-an-identity.md).  
:::  

## Transaction message
In order to properly send a transaction to share a document, you will need to create and sign the
following message.

```json
{
  "type": "commercio/MsgShareDocument",
  "value": {
    "sender": "<Sender Did>",
    "recipients": [
      "<Recipient address>"
    ],
    "uuid": "<Document UUID>",
    "content_uri": "<Document content URI>",
    "metadata": {
      "content_uri": "<Metadata content URI>",
      "schema": {
        "uri": "<Metadata schema definition URI>",
        "version": "<Metadata schema version>"
      },
      "schema_type": "<Metadata schema type>"
    },
    "checksum": {
      "value": "<Document content checksum value>",
      "algorithm": "<Document content checksum algorithm>"
    },
    "encryption_data": {
      "keys": [
        {
          "recipient": "<Recipient address>",
          "value": "<Encrypted and encoded symmetric key value>"
        }
      ],
      "encrypted_data": [
        "<Encrypted field identifier>"
      ]
    },
    "do_sign": {
        "storage_uri": "uri://storage",
        "signer_instance": "did S",
        "sdn_data": [
                "firstName",                
                "lastName",                
                "sin",                
                "email"
        ],
        "vcrId": "<identity VCR Identifier"
        "certificateProfile": "<one of the profiles supported by S>"
    }
  }
}
```

### Fields requirements
| Field | Required | 
| :---: | :------: |
| `sender` | Yes | 
| `recipients` | Yes |
| `uuid` | Yes | 
| `content_uri` | No | 
| `metadata` | Yes |
| `checksum` | No | 
| `encryption_data` | No | 

##### `metadata`
| Field | Required | 
| :---: | :------: |
| `content_uri` | Yes | 
| `schema_type` | No *<sup>1</sup> *<sup>2</sup>  | 
| `schema` | No *<sup>1</sup> |

- *<sup>1</sup> The `schema_type` and `schema` fields are mutually exclusive.
This means that if the first one exists the second will not be used.
   
- *<sup>2</sup> You can read which `schema_type` values are supported inside 
   the [supported metadata schemes section](../metadata-schemes.md#supported-metadata-schemes)
   
###### `metadata.schema`
| Field | Required | 
| :---: | :------: |
| `uri` | Yes | 
| `version` | Yes | 

##### `checksum`
| Field | Required | 
| :---: | :------: |
| `value` | Yes |
| `algorithm` | Yes *<sup>1</sup> |

- *<sup>1</sup> You can read which `checksum.algorithm` values are supported inside the
[supported checksum algorithms section](#supported-checksum-algorithm)  

##### `encryption_data`
| Field | Required | 
| :---: | :------: |
| `key` | Yes |
| `encrypted_data` | Yes |

##### `do_sign`
| Field | Required | 
| :---: | :------: |
| `storage_uri` | Yes ?? |
| `signer_instance` | Yes |
| `sdn_data` | No |
| `vcrId` | Yes |
| `certificateProfile` | Yes/No ?? |


* storage_uri
* signer_instance
* sdn_data: contains an array with a list of required fields for Subject Distinguish Name. The names of fields are x509 standard compliant


## Supported checksum algorithm
When computing the checksum of a document's contents, you must use one of the following supported checksum algorithms.  
Not using one of these will result in your transaction being rejected or mishandled by recipients. 

| Algorithm | Specification |
| :-------: | :-----------: |
| `md5` | [MD5](https://www.ietf.org/rfc/rfc1321.txt) |
| `sha-1`| [SHA-1](https://tools.ietf.org/html/rfc3174) |
| `sha-224` | [RFC 4634](https://tools.ietf.org/html/rfc4634) |
| `sha-256` | [RFC 4634](https://tools.ietf.org/html/rfc4634) |
| `sha-384` | [RFC 4634](https://tools.ietf.org/html/rfc4634) |
| `sha-512` | [RFC 4634](https://tools.ietf.org/html/rfc4634) |

#### Checksum validity check
Please note that, when sending a document that has an associated checksum, the validity of the checksum itself is
checked only formally. This means that we only check that the hash value has a valid length, but we do not check 
if the given has is indeed the hash of the document's content. It should be the client responsibility to perform this 
check.  

## Encrypting the data

::: tip

The following is just an example on how to do file encryption, you're free to use any other algorithm!

:::

In order to properly encrypting the data that you want to avoid being shared publicly, 
the following procedure should be followed.

We'll use AES-256 in CBC mode to encrypt a file, and let the recipient decrypt it by sharing with
it the AES encryption key.

1. Generate a safe AES-256 encryption key. A key size of 256 bits is recommended.
   ```
   aes_key = get_random_aes_key(key_size = 256)
   ```

2. Use the AES key to encrypt the data you desire using the AES-256 CBC method.  
   ```
   encrypted_data = aes_encrypt_cbc(
     key = aes_key, 
     initialization_vector = null
   )
   ```
   
3. Encrypt the AES-256 key using the recipient's public encryption key  
   ```
   encrypted_aes_key = rsa_encrypt(
     key = recipient.public_rsa_encryption_key,
     value = aes_key
   )    
   ```
   
4. Encode the encrypted AES-256 key  
   ```
   encoded_encryption_key = hex_encode(encrypted_aes_key)
   ```
   
4. Compose the encryption data  
   ```json
   {
     "encryption_data": {
       "keys": [
         {
           "recipient": "<Recipient address>",
           "value": "<Hex encoded encryption key>"
         }
       ],
       "encrypted_data": [
         "<Your encrypted data identifier>"
       ]
     }
   }
   ```

The `encrypted_data` field does not contain the encrypted payload itself, but rather denotes what message property is encrypted with `aes_key`. 

`encrypted_data` only accepts the following identifiers:
 - `content_uri`
 - `metadata.content_uri`
 - `metadata.schema.uri`

A special identifier, `content`, can be used to specify that `aes_key` has been used to encrypt a file exchanged by other means of communication.

### Supported encrypted data
Please note that when specifying which data you have encrypted for the document recipient, you need to use one or 
more of the following identifiers inside the `encryption_data.encrypted_data` field.  
Inserting other non supported values inside such a field will result in the transactions being rejected as not valid.   

| Identifier | Referenced data | 
| :--------: | :-------------- |
| `content` | Document's file contents |
| `content_uri` | Value of the `content_uri` field |
| `metadata.content_uri` | Value of the `content_uri` field inside the `metadata` object |
| `metadata.schema.uri` | Value of the `uri` field inside the `metadata`'s `schema` sub-object |

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
shareDocument
```

## Transaction cost
To make sure every person on the network has predictable costs when sending a document, you are required to set 
a minimum `fee.amount` value inside the transaction when using the `commercio/MsgShareDocument` message type.  

This minimum fee value is always equivalent at 0.01 euro and can be paid in two different methods: 

1. Using `10000uccc`, which are millionth of Commercio Cash Credits.  
   In order to know how to get the Commercio Cash Credits, please refer to the [*"Opening a Etp"* page](../../commerciomint/#mint-commercio-cash-credit-ccc) 
2. Using the 0.01 euro equivalent of Commercio Tokens.  
   Please note that the usage of such tokens is discouraged as it is more difficult to properly have an estimate on how
   many tokens are needed. When using it you should set a 10% more of the minimum tokens required.     