# CommercioDOCS
CommercioDOCS is the module that allows you to send a document to another user, and retrieve the list of documents
that you have received. 

{:toc}

## Sending a document
In order to send a document you are required to have an identity with some tokens inside it.   
In order to know what an identity is, how to create it and how to get tokens, please refer to the 
[*"Creating an identity"* section](commercioid.md#creating-an-identity).  

### Using the CLI 
In order to send a document using the CLI you can use the following command 

```shell
cncli tx commerciodocs send-document
```

#### Parameters
| Parameter | Type | Required | Description |  
| :-------: | :---: | :-----: | :---------- |
| `recipient` |  

#### Example usage 

```shell
cncli tx commerciodocs send-document \
  [recipient] \
  [document-content-uri] \
  [metadata-content-uri] \
  [metadata-schema-uri] \
  [metadata-schema-version] \
  [computation-proof] \
  [checksum-value] \
  [checksum-algorithm]
```

### Creating a transaction offline
In order to properly send a `commerciodocs/SendDocument` transaction, you will need to create and sign the
following message:

```json
{
  "sender": "<Your address>",
  "recipient": "<Recipient address>",
  "content_uri": "<Document content URI>",
  "metadata": {
    "content_uri": "<Metadata content URI>",
    "schema": {
      "uri": "<Metadata schema URI>",
      "version": "<Metadata schema version>"
    },
    "proof": "<Metadata validation proof>"
  },
  "checksum": {
    "value": "<Document content checksum value>",
    "algorithm": "<Checsum algorithm>"
  }
}
```

## 