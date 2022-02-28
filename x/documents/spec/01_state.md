<!--
order: 1
-->

# State

The `x/documents` module keeps state of the following objects

- **Documents** `docs:document:[documentID]`
- _documentIDs_ sent by an user `docs:documents:sent:[address]`
- _documentIDs_ received by an user `docs:documents:received:[address]`
- **Receipts** `docs:receipt:[receiptID]`
- _receiptIDs_ sent by an user `docs:receipts:sent:[address]`
- _receiptIDs_ received by an user `docs:receipts:received:[address]`
- _receiptIDs_ associated to a certain document `docs:receipts:documents:[documentID]


## `Document`

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

### `DocumentChecksum`

```
message DocumentChecksum {
  string value = 1;
  string algorithm = 2;
}
```

### `DocumentEncryptionData`

```
message DocumentEncryptionData {
  repeated documents.DocumentEncryptionKey keys = 1;
  repeated string encryptedData = 2;
}
```

#### `DocumentEncryptionKey`

```
message DocumentEncryptionKey {
  string recipient = 1;
  string Value = 2;
}
```

### `DocumentMetadata`

```
message DocumentMetadata {
	string contentURI = 1;
	//string schemaType = 2;
	DocumentMetadataSchema schema = 2;
}
```

#### `DocumentMetadataSchema`
```
message DocumentMetadataSchema {
  string URI = 1;
  string version = 2;
}
```

### `DocumentDoSign`

```
message DocumentDoSign {
  string storageURI = 1;
  string signerInstance = 2;
  repeated string sdnData = 3; 
  string vcrID = 4;
  string certificateProfile = 5;
}
```

## `DocumentReceipt`

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

Please note that the former sender of a document becomes the recipient for a `DocumentReceipt`.
Conversely, one of the receivers (or it can be just one receiver) becomes the sender for a `DocumentReceipt`.