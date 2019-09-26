# Sending a document
In order to send a document you are required to have an identity with some tokens inside it.   

::: tip  
To know what an identity is, how to create it and how to get tokens, please refer to the 
[*"Creating an identity"* section](../../id/README.md#creating-an-identity).  
:::  

## Transaction message
In order to properly send a transaction to share a document, you will need to create and sign the
following message.

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

### Fields requirements
| Field | Required | 
| :---: | :------: |
| `sender` | Yes | 
| `recipient` | Yes |
| `uuid` | Yes | 
| `content_uri` | No | 
| `metadata` | Yes |
| `checksum` | No | 

#### `metadata`
| Field | Required | 
| :---: | :------: |
| `content_uri` | Yes | 
| `schema_type` | No *<sup>1</sup> *<sup>2</sup>  | 
| `schema` | No *<sup>1</sup> |
| `proof` | Yes | 

- *<sup>1</sup> The `schema_type` and `schema` fields are mutually exclusive.
This means that if the first one exists the second will not be used.
   
- *<sup>2</sup> You can read which `schema_type` values are supported inside 
   the [supported metadata schemes section](../metadata-schemes.md#supported-metadata-schemes)
   
##### `metadata.schema`
| Field | Required | 
| :---: | :------: |
| `uri` | Yes | 
| `version` | Yes | 

#### `checksum`
| Field | Required | 
| :---: | :------: |
| `value` | Yes |
| `algorithm` | Yes *<sup>1</sup> |

- *<sup>1</sup> You can read which `checksum.algorithm` values are supported inside the
[supported checksum algorithms section](#supported-checksum-algorithm)  

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

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
shareDocument
```