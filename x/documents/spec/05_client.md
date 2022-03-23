<!--
order: 5
-->

# Client

## Transactions

### Sharing a document

#### CLI

```bash
commercionetworkd tx docs share [recipient] [document-uuid] [document-metadata-uri] [metadata-schema-uri] [metadata-schema-version] [document-content-uri] [checksum-value] [checksum-algorithm] 
```

**Parameters:**

| Parameter | Description |
| :-------: | :---------- | 
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
| :-------:              | :----------  | :---------- | :---------- |
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

#### REST

```
/docs/{address}/sent
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current sent documents |

##### Example 

Getting sent docs from `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/docs/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/sent
```

### List received documents

#### CLI

```bash
commercionetworkd query docs received-documents [address]
```

#### REST

```
/docs/{address}/received
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current received documents |


##### Example 

Getting docs for `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/docs/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/received
```

### List sent receipts

#### CLI

```bash
commercionetworkd query docs sent-receipts [address]
```

#### REST

```
/receipts/{address}/sent
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current sent receipts |

##### Example 

Getting sent receipts from `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/receipts/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/sent
```

### List received receipts

#### CLI

```bash
commercionetworkd query docs received-receipts [address]
```

#### REST

```
/receipts/{address}/received
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current received receipts |


##### Example 

Getting receipts for `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/receipts/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/received
```

### List receipts associated to a certain document

#### CLI

```bash
commercionetworkd query docs documents-receipts [documentUUID]
```

#### REST

```
documents/document/{UUID}/receipts
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `UUID` | Document ID of the document for which to read current received receipts |

##### Example 

Getting receipts associated to the document with ID `d83422c6-6e79-4a99-9767-fcae46dfa371`:

```
http://localhost:1317/document/d83422c6-6e79-4a99-9767-fcae46dfa371/receipts
```
