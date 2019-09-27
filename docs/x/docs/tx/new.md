# Sharing a document
In order to send a document you are required to have an identity with some tokens inside it.   
In order to know what an identity is, how to create it and how to get tokens, please refer to the 
[*"Creating an identity"* section](../../id/README.md#creating-an-identity).  

## Transaction message
In order to properly send a transaction to share a document, you will need to create and sign the
following message:

```json
{
  "type": "commercio/MsgShareDocument",
  "value": {
    "sender": "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0",
    "recipients": [
      "<Recipient address>"
    ],
    "document": {
      "uuid": "<Document UUID>",
      "content_uri": "<Document content URI>",
      "metadata": {
        "content_uri": "<Metadata content URI>",
        "schema": {
          "uri": "<Metadata schema definition URI>",
          "version": "<Metadata schema version>"
        },
        "schema_type": "<Metadata schema type>",
        "proof": "<Metadata verification proof>"
      },
      "checksum": {
        "value": "<Document content checksum value>",
        "algorithm": "<Document content checksum algorithm>"
      },
      "encryption_data": {
        "keys": [
          {
            "recipient": "<Recipient address>",
            "value": "<Encrypted and encoded symmetric key value>",
            "encoding": "<Encoding algorithm>"
          }
        ],
        "encrypted_data": [
          "<Encrypted field identifier>"
        ]
      }
    }
  }
}

```

### Notes
1. The `metatada.schema_type` and `metadata.schema` fields are mutually exclusive.
   This means that if the first one exists the second will not be used.
2. You can read which `metadata.schema_type` values are supported inside 
   the [supported metadata schemes section](../metadata-schemes.md#supported-metadata-schemes)
3. You can read which `checksum.algorithm` values are supported inside the
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

## Encrypting the data
In order to properly encrypting the data that you want to avoid being shared publicly, 
the following procedure must be followed. 

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
           "recipient": "<recipient_address>",
           "value": "<encoded_encryption_key>",
           "encoding": "hex"
         }
       ],
       "encrypted_data": [
         "<Your encrypted data identifier>"
       ]
     }
   }
   ```
   
### Supported encoding methods
Currently only the following encoding methods are supported then encoding the encrypted AES-256 key:

* `hex`
* `baes64`  
   
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

## Using the CLI

:::danger  
Please note that the following procedure is completely outdated and should not be used  
:::

In order to send a document using the CLI you can use the following command 

```bash
cncli tx commerciodocs share \
  [recipient] \
  [document-uuid] \ 
  [document-metadata-uri] \
  [metadata-schema-uri] \
  [metadata-schema-version] \
  [metadata-verification-proof] \
  [document-content-uri] \
  [checksum-value] \
  [checksum-algorithm]
```

### Parameters
| Parameter | Type | Required | Description |  
| :-------- | :---: | :-----: | :---------- |
| `recipient` |  String | Yes | Bech32 address of the document recipient | 
| `document-uuid` | String | Yes | UUID of the document that is being sent |
| `document-metadata-uri` | String | Yes | Uri to the file containing the document metadata |
| `metadata-schema-uri` | String | Yes | Uri to the file containing the definition of the metadata schema used to define the metadata of this document |
| `metadata-schema-version` | String | Yes | Version of the schema used to define the document's metadata |
| `metadata-verification-proof` | String | Yes | Proof that the client has correctly verified the validity of the metadata based on the associated schema |
| `document-content-uri` | String | No | Uri of the file containing the document's data |
| `checksum-value` | String | Yes if `document-content-uri` is present | Value of the document content's checksum | 
| `checksum-algorithm` | String | Yes if `checksum-value` is present | Algorithm used to compute the document's content checksum |

### Example usage 

```bash
cncli tx commerciodocs share \
  did:com:1d63vn76znf6fxumfdgx4rmc4wlnppv5evqnx84 \
  490e835d-dbb3-445c-aedc-df477f41d8a2 \
  http://example.com/document/metadata \
  http://example.com/document/metadata/schema \
  1.0.0 \
  353438683534383534726835347268353472 \
  https://example/document \
  7815696ecbf1c96e6894b779456d330e \
  md5 \
  --from jack
```
