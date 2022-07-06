<!--
order: 5
-->

# Client

## Transactions

### Sharing a document

#### CLI

```bash
commercionetworkd tx docs share \
  [recipient] \
  [document-uuid] \
  [document-metadata-uri] \
  [metadata-schema-uri] \
  [metadata-schema-version] \
  [document-content-uri] \
  [checksum-value] \
  [checksum-algorithm] 
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `recipient`               | Address of the recipient for the document  |
| `document-uuid`           | Document ID following the UUID format |
| `document-metadata-uri`   | Metadata content URI |
| `metadata-schema-uri`     | Metadata schema definition URI |
| `metadata-schema-version` | Metadata schema version |
| `document-content-uri`    | **Optional.** Document content URI |
| `checksum-value`          | **Optional.** Document content checksum value |
| `checksum-algorithm`      | **Optional.** Document content checksum algorithm |

**Flags:**

| Parameter              | Type         | Default | Description |
| :-------              | :----------  | :---------- | :---------- |
| `sign`                 | `bool`       | _false_ | specifies that we want to sign the document |
| `sign-storage-uri`     | `string`     | `""`    | the storage URI to sign |
| `sign-signer-instance` | `string`     | `""`    | the signer instance needed to sign |
| `sign-vcr-id`          | `string`     | `""`    | the vcr id needed to sign |
| `sign-certificate-profile` | `string` | `""`    | the certificate profile needed to sign |
| `sign-sdn-data`        | `string`     | `""`    | the sdn data needed to sign |

### Sending a receipt

#### CLI

```bash
commercionetworkd tx docs send-receipt [recipient] [tx-hash] [document-uuid] [proof]
```

**Parameters:**

| Parameter | Description |
| :-------: | :---------- | 
| `recipient`     | Address of the user who initially shared the associated document  |
| `tx-hash`       | Transaction hash in which the document has been shared |
| `document-uuid` | ID of the associated document |
| `proof` | **Optional.** Reading proof | 

This command generates a random UUID for the receipt.

## Queries

### List sent documents

#### CLI

```bash
commercionetworkd query docs sent-documents [address]
```

#### gRPC

```
grpcurl -plaintext \
    -d '{"address":"did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"}' \
    lcd-mainnet.commercio.network:9090 \
    commercionetwork.commercionetwork.documents.Query/SentDocuments
```

#### REST

```
/commercionetwork/documents/document/{address}/sent
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current sent documents |

##### Example 

Getting sent docs from `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/commercionetwork/documents/document/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/sent
```

### List received documents

#### CLI

```bash
commercionetworkd query docs received-documents [address]
```

#### gRPC

```
grpcurl -plaintext \
    -d '{"address":"did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"}' \
    lcd-mainnet.commercio.network:9090 \
    commercionetwork.commercionetwork.documents.Query/ReceivedDocument
```

#### REST

```
/commercionetwork/documents/document/{address}/received
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current received documents |


##### Example 

Getting docs for `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/commercionetwork/documents/document/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/received
```

### List sent receipts

#### CLI

```bash
commercionetworkd query docs sent-receipts [address]
```

#### gRPC

```
grpcurl -plaintext \
    -d '{"address":"did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"}' \
    lcd-mainnet.commercio.network:9090 \
    commercionetwork.commercionetwork.documents.Query/SentDocumentsReceipts
```

#### REST

```
/commercionetwork/documents/receipts/{address}/sent
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current sent receipts |

##### Example 

Getting sent receipts from `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/commercionetwork/documents/receipts/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/sent
```

### List received receipts

#### CLI

```bash
commercionetworkd query docs received-receipts [address]
```

#### gRPC

```
grpcurl -plaintext \
    -d '{"address":"did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"}' \
    lcd-mainnet.commercio.network:9090 \
    commercionetwork.commercionetwork.documents.Query/ReceivedDocumentsReceipts
```

#### REST

```
/commercionetwork/documents/receipts/{address}/received
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current received receipts |



##### Example 

Getting receipts for `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/commercionetwork/documents/receipts/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/received
```

### List receipts associated to a certain document

#### CLI

```bash
commercionetworkd query docs documents-receipts [documentUUID]
```

#### gRPC

```
grpcurl -plaintext \
    -d '{"UUID":"3469ca3e-8fe6-4d1f-9713-11418bb9a8f4"}' \
    lcd-mainnet.commercio.network:9090 \
    commercionetwork.commercionetwork.documents.Query/DocumentsReceipts
```

#### REST

```
/commercionetwork/documents/document/{UUID}/receipts
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `UUID` | Document ID of the document for which to read current received receipts |

##### Example 

Getting receipts associated to the document with ID `d83422c6-6e79-4a99-9767-fcae46dfa371`:

```
http://localhost:1317/commercionetwork/documents/document/d83422c6-6e79-4a99-9767-fcae46dfa371/receipts
```

### Get document with specific `documentUUID`

#### CLI

```bash
commercionetworkd query docs show-document [documentUUID]
```

#### gRPC

```
grpcurl -plaintext \
    -d '{"UUID":"3469ca3e-8fe6-4d1f-9713-11418bb9a8f4"}' \
    lcd-mainnet.commercio.network:9090 \
    commercionetwork.commercionetwork.documents.Query/Document
```

#### REST

```
/commercionetwork/documents/document/{UUID}
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `UUID` | Document ID of the document |

##### Example 

Getting receipts associated to the document with ID `d83422c6-6e79-4a99-9767-fcae46dfa371`:

```
http://localhost:1317/commercionetwork/documents/document/d83422c6-6e79-4a99-9767-fcae46dfa371
```
