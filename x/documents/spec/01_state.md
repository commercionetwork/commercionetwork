<!--
order: 1
-->

# State

The `x/documents` module keeps state of Documents and Document Receipt sharing.

## Documents

When a document gets shared, the module stores it in the store.

| Object | Prefix | Value |
| :-------: | :---------- | :---------- | 
| Document | `docs:document:` | `[documentID]` |

Also, the module updates the following lists with the ID of the shared document:

### Document w.r.t. addresses:

| List | Prefix | Value |
| :-------: | :---------- | :---------- | 
| _documentIDs_ | `docs:documents:sent:` | `[address]` |
| _documentIDs_ | `docs:documents:received:` | `[address]` |

## Document Receipts

When a document receipt gets shared, the module stores it in the store.

| List | Prefix | Value |
| :-------: | :---------- | :---------- | 
| Receipt | `docs:receipt:` | `[receiptID]` |

Also, the module updates the following lists with the ID of the shared document receipt:

### Document Receipts w.r.t. addresses:

| List | Prefix | Value |
| :-------: | :---------- | :---------- | 
| _receiptIDs_ | `docs:receipts:sent:` | `[address]` |
| _receiptIDs_ | `docs:receipts:received:` | `[address]` |

### Document Receipts w.r.t. documents:

| List | Prefix | Value |
| :-------: | :---------- | :---------- | 
| _receiptIDs_ | `docs:receipts:documents` | `[documentID]` |

## Type definitions

### `Document` definition

```
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

```
message DocumentChecksum {
  string value = 1;
  string algorithm = 2;
}
```

#### `DocumentEncryptionData` definition

```
message DocumentEncryptionData {
  repeated documents.DocumentEncryptionKey keys = 1;
  repeated string encryptedData = 2;
}
```

##### `DocumentEncryptionKey` definition

```
message DocumentEncryptionKey {
  string recipient = 1;
  string Value = 2;
}
```

#### `DocumentMetadata` definition

```
message DocumentMetadata {
	string contentURI = 1;
	//string schemaType = 2;
	DocumentMetadataSchema schema = 2;
}
```

##### `DocumentMetadataSchema` definition
```
message DocumentMetadataSchema {
  string URI = 1;
  string version = 2;
}
```

#### `DocumentDoSign` definition

```
message DocumentDoSign {
  string storageURI = 1;
  string signerInstance = 2;
  repeated string sdnData = 3; 
  string vcrID = 4;
  string certificateProfile = 5;
}
```

### `DocumentReceipt` definition

Please note that the former sender of a document becomes the recipient for a `DocumentReceipt`.
Conversely, one of the receivers (or it can be just one receiver) becomes the sender for a `DocumentReceipt`.

```
message DocumentReceipt {
    string UUID = 1; 
    string sender = 2; 
    string recipient = 3; 
    string txHash = 4; 
    string documentUUID = 5;
    string proof = 6;
}
```
