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
Endpoint:

```
commercionetwork.commercionetwork.documents.Query/SentDocuments
```

##### Example

```bash
grpcurl -plaintext \
    -d '{"address":"did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.documents.Query/SentDocuments
```

##### Response
```json
"Document": [
    {
      "sender": "did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya",
      "recipients": [
        "did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"
      ],
      "UUID": "000132d4-77b9-4159-8456-4b8301e1c717",
      "metadata": {
        "contentURI": "4fa2d6ec7a74ae3cf7ab40f51fe6074594c471883e7546f10505f2f0797f4f80",
        "schema": {
          "URI": "foxsign.app/shareDocument/au",
          "version": "1.0.14"
        }
      },
      "contentURI": "4fa2d6ec7a74ae3cf7ab40f51fe6074594c471883e7546f10505f2f0797f4f80",
      "checksum": {
        "value": "1a766b02dd3b397d4def5e943c4f8574c6acaa98d5525491b97c19b76b35b463",
        "algorithm": "sha-256"
      }
    },
    {
      "sender": "did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya",
      "recipients": [
        "did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"
      ],
      "UUID": "0001e283-9b41-4242-a595-6a04b352868e",
      "metadata": {
        "contentURI": "7475e7aabb526a8aeb35015f8f1d4892757c12f7110edcb9429e33860fd07758",
        "schema": {
          "URI": "foxsign.app/shareDocument/cb",
          "version": "1.0.9"
        }
      },
      "contentURI": "7475e7aabb526a8aeb35015f8f1d4892757c12f7110edcb9429e33860fd07758",
      "checksum": {
        "value": "97748dc132b3fe26283a8c4a1f34963364dc2238009a881df628ed3b5e2d2511",
        "algorithm": "sha-256"
      }
    },
    ...
  ],
  "pagination": {
    "nextKey": "OjAwMjMxYzQxLTMyNDktNGUxNC05ZDVkLTIwZWVlMjEwM2Q4OQ==",
    "total": "177888"
  }
}
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
Endpoint:

```
commercionetwork.commercionetwork.documents.Query/ReceivedDocument
```

```
grpcurl -plaintext \
    -d '{"address":"did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.documents.Query/ReceivedDocument
```

##### Response
```json
{
  "Document": [
    {
      "sender": "did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya",
      "recipients": [
        "did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"
      ],
      "UUID": "000132d4-77b9-4159-8456-4b8301e1c717",
      "metadata": {
        "contentURI": "4fa2d6ec7a74ae3cf7ab40f51fe6074594c471883e7546f10505f2f0797f4f80",
        "schema": {
          "URI": "foxsign.app/shareDocument/au",
          "version": "1.0.14"
        }
      },
      "contentURI": "4fa2d6ec7a74ae3cf7ab40f51fe6074594c471883e7546f10505f2f0797f4f80",
      "checksum": {
        "value": "1a766b02dd3b397d4def5e943c4f8574c6acaa98d5525491b97c19b76b35b463",
        "algorithm": "sha-256"
      }
    },
    ...
  ],
  "pagination": {
    "nextKey": "OjAwMjMxYzQxLTMyNDktNGUxNC05ZDVkLTIwZWVlMjEwM2Q4OQ==",
    "total": "177895"
  }
}

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
Endpoint:

```
commercionetwork.commercionetwork.documents.Query/SentDocumentsReceipts
```

```
grpcurl -plaintext \
    -d '{"address":"did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.documents.Query/SentDocumentsReceipts
```

##### Response
```json
{
  "receipt": [
    {
      "UUID": "a783a39a-eabb-4e5c-b879-a27538247232",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "ECD43F1BEC153DC129FBF13F6836E527BA3F28E0651A4285B9CAF96D4E6483E3",
      "documentUUID": "4939a995-e979-41a3-9b03-bef6f1fc7044"
    },
    {
      "UUID": "8e95e73b-382a-4d26-9eb2-d4264c2ff854",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "8B1A600ABCBFA5D31C01D98BDB0905C83E89660EDB7166965021B5974EDDCA41",
      "documentUUID": "64b80f68-0050-4d87-af8e-89bb5d34dfc9"
    },
    {
      "UUID": "0a039c5f-c467-4f1e-914c-19c1eb9f4d6b",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "170678B90D0CEEC638C5EC1305B85B52AFC05E9C3209E75052AC5BBE03B94C3C",
      "documentUUID": "c0e59b0c-98ed-4eac-b0c4-c05088f73479"
    },
    {
      "UUID": "f77efdcc-981e-4c8b-9f3e-63b3af4ac212",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "A784C3CD6255D1C675D0F009059436F6FC255C87587CAA513291BCCDCE869E28",
      "documentUUID": "f77efdcc-981e-4c8b-9f3e-63b3af4ac212"
    }
  ],
  "pagination": {
    "total": "4"
  }
}
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
Endpoint:

```
commercionetwork.commercionetwork.documents.Query/ReceivedDocumentsReceipts
```

```
grpcurl -plaintext \
    -d '{"address":"did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.documents.Query/ReceivedDocumentsReceipts
```

##### Response
```json
{
  "ReceiptReceived": [
    {
      "UUID": "0a039c5f-c467-4f1e-914c-19c1eb9f4d6b",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "170678B90D0CEEC638C5EC1305B85B52AFC05E9C3209E75052AC5BBE03B94C3C",
      "documentUUID": "c0e59b0c-98ed-4eac-b0c4-c05088f73479"
    },
    {
      "UUID": "8e95e73b-382a-4d26-9eb2-d4264c2ff854",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "8B1A600ABCBFA5D31C01D98BDB0905C83E89660EDB7166965021B5974EDDCA41",
      "documentUUID": "64b80f68-0050-4d87-af8e-89bb5d34dfc9"
    },
    {
      "UUID": "a783a39a-eabb-4e5c-b879-a27538247232",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "ECD43F1BEC153DC129FBF13F6836E527BA3F28E0651A4285B9CAF96D4E6483E3",
      "documentUUID": "4939a995-e979-41a3-9b03-bef6f1fc7044"
    },
    {
      "UUID": "f77efdcc-981e-4c8b-9f3e-63b3af4ac212",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "A784C3CD6255D1C675D0F009059436F6FC255C87587CAA513291BCCDCE869E28",
      "documentUUID": "f77efdcc-981e-4c8b-9f3e-63b3af4ac212"
    }
  ],
  "pagination": {
    "total": "4"
  }
}
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
Endpoint:

```
commercionetwork.commercionetwork.documents.Query/DocumentsReceipts
```

```
grpcurl -plaintext \
    -d '{"UUID":"4939a995-e979-41a3-9b03-bef6f1fc7044"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.documents.Query/DocumentsReceipts
```

##### Response
```json
{
  "Receipts": [
    {
      "UUID": "a783a39a-eabb-4e5c-b879-a27538247232",
      "sender": "did:com:1a0v2kdjkm95gq5qq7ygvczdyuymt6hg3c2su0c",
      "recipient": "did:com:1ujh8ldcy2k737vwz8k6cw86uhfvwfe5peay8gg",
      "txHash": "ECD43F1BEC153DC129FBF13F6836E527BA3F28E0651A4285B9CAF96D4E6483E3",
      "documentUUID": "4939a995-e979-41a3-9b03-bef6f1fc7044"
    }
  ],
  "pagination": {
    "total": "1"
  }
}

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
Endpoint:

```
commercionetwork.commercionetwork.documents.Query/Document
```

##### Example

```bash
grpcurl -plaintext \
    -d '{"UUID":"3469ca3e-8fe6-4d1f-9713-11418bb9a8f4"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.documents.Query/Document
```

##### Response
```json
{
  "Document": {
    "sender": "did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya",
    "recipients": [
      "did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya"
    ],
    "UUID": "3469ca3e-8fe6-4d1f-9713-11418bb9a8f4",
    "metadata": {
      "contentURI": "3a90c9e8b249929e5e4e902728056acc31e2ea5f7fc92e9421aef494b3c4451f",
      "schema": {
        "URI": "foxsign.app/shareDocument/ppv",
        "version": "1.0.7"
      }
    },
    "contentURI": "3a90c9e8b249929e5e4e902728056acc31e2ea5f7fc92e9421aef494b3c4451f",
    "checksum": {
      "value": "00274a9d73d08d943959c9f5889cbe0e2195039e5d1ad76d7e5f92203367aa22",
      "algorithm": "sha-256"
    }
  }
}
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

Getting the document with ID `d83422c6-6e79-4a99-9767-fcae46dfa371`:

```
http://localhost:1317/commercionetwork/documents/document/d83422c6-6e79-4a99-9767-fcae46dfa371
```
