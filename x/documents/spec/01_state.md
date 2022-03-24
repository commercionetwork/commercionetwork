<!--
order: 1
-->

# State

The `x/documents` module keeps state of Documents and Document Receipts sharing.

## Store

### Documents tracking

The module appends in the store the shared document.

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `docs:document:[documentID]` | &rarr; | _Document_ |

Also, the module appends to store the ID of the document, in the following lists:

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `docs:documents:sent:[senderAddress]:[documentID]` | &rarr; | _documentID_ |
| `docs:documents:received:[receiverAddress]:[documentID]` | &rarr; | _documentID_ |

### Document Receipts tracking

The module appends in the store the shared document receipt

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `docs:receipt:[receiptID]` | &rarr; | _Receipt_ |

Also, the module appends to store the ID of the document receipt, in the following lists:

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `docs:receipts:sent:[senderAddress]:[receiptID]` | &rarr; | _receiptID_ |
| `docs:receipts:received:[receiverAddress]:[receiptID]` | &rarr; | _receiptID_ |
| `docs:receipts:documents:[documentID]:[receiptID]` | &rarr; | _receiptID_ |

## Type definitions

### The `Document` type

#### `Document` definition
```proto
message Document {
  string sender = 1; 
  repeated string recipients = 2; 
  string UUID = 3; 
  documents.DocumentMetadata metadata = 4; 
  string contentURI = 5; 
  documents.DocumentChecksum checksum = 6; 
  documents.DocumentEncryptionData encryptionData = 7 [
    (gogoproto.moretags)   = "yaml:\"encryption_data\"",
    (gogoproto.customtype) = "DocumentEncryptionData",
    (gogoproto.nullable)   = true
  ]; 
  documents.DocumentDoSign doSign = 8 [
    (gogoproto.moretags)   = "yaml:\"do_sign\"",
    (gogoproto.customtype) = "DocumentDoSign",
    (gogoproto.nullable)   = true
  ];
}
```

#### `DocumentChecksum` definition

```proto
message DocumentChecksum {
  string value = 1;
  string algorithm = 2;
}
```

#### `DocumentEncryptionData` definition

```proto
message DocumentEncryptionData {
  repeated documents.DocumentEncryptionKey keys = 1;
  repeated string encryptedData = 2;
}
```

##### `DocumentEncryptionKey` definition

```proto
message DocumentEncryptionKey {
  string recipient = 1;
  string Value = 2;
}
```

#### `DocumentMetadata` definition

```proto
message DocumentMetadata {
	string contentURI = 1;
	//string schemaType = 2;
	DocumentMetadataSchema schema = 2;
}
```

##### `DocumentMetadataSchema` definition

```proto
message DocumentMetadataSchema {
  string URI = 1;
  string version = 2;
}
```

#### `DocumentDoSign` definition

```proto
message DocumentDoSign {
  string storageURI = 1;
  string signerInstance = 2;
  repeated string sdnData = 3; 
  string vcrID = 4;
  string certificateProfile = 5;
}
```

### The`DocumentReceipt` type

Please note that the former sender of a document becomes the recipient for a `DocumentReceipt`.
Conversely, one of the receivers (or it can be just one receiver) becomes the sender for a `DocumentReceipt`.

```proto
message DocumentReceipt {
    string UUID = 1; 
    string sender = 2; 
    string recipient = 3; 
    string txHash = 4; 
    string documentUUID = 5;
    string proof = 6;
}
```
