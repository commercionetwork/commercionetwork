# Sharing a document
In order to send a document you are required to have an identity with some tokens inside it.   
In order to know what an identity is, how to create it and how to get tokens, please refer to the 
[*"Creating an identity"* section](../../id/README.md#creating-an-identity).  

## Transaction message
In order to properly send a transaction to share a document, you will need to create and sign the
following message:

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

### Notes
1. The `metatada.schema_type` and `metadata.schema` fields are mutually exclusive.
   This means that if the first one exists the second will not be used.
2. You can read which `metadata.schema_type` values are supported inside 
   the [supported metadata schemes section](../metadata-schemes.md#supported-metadata-schemes)
3. You can read which `checksum.algorithm` values are supported inside the
   [supported checksum algorithms section](#supported-checksum-algorithm)  

## Using the CLI 
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

