<!--
order: 1
-->

# State

The `did` module keeps state of Identities, represented as the evolution of the DID Document and Metadata for a certain DID.

## Store

The module appends in the store the updated Identity.

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `did:identities:[address]:[updated]` | &rarr; | _Identity_ |

This operation uses the block time (guaranteed to be deterministic and always increasing) to populate the `Updated` field of `Metadata`. 
This timestamp is also used to populate the `Created` field, but only for the first version of the `Identity`, that will be maintained in the newer versions.

## The `Identity` type

```protobuf
message Identity {
  DidDocument didDocument = 1;
  Metadata metadata = 2;
}
```

### `DidDocument` definition

```protobuf
message DidDocument {
  option (gogoproto.equal) = true;
  repeated string context = 1 [ (gogoproto.jsontag) = "@context,omitempty" ];
  string ID = 2;
  repeated VerificationMethod verificationMethod = 3;
  repeated string authentication = 4;
  repeated string assertionMethod = 5;
  repeated string keyAgreement = 6;
  repeated string capabilityInvocation = 7;
  repeated string capabilityDelegation = 8;
  repeated Service service = 9;
}
```

#### `VerificationMethod` definition

```protobuf
message VerificationMethod {
  option (gogoproto.equal) = true;
  string ID = 1;
  string Type = 2;
  string Controller = 3;
  string publicKeyMultibase = 4;
}
```

#### `Service` definition

```protobuf
message Service {
  option (gogoproto.equal) = true;
  string ID = 1;
  string type = 2;
  string serviceEndpoint = 3;
}
```

### `Metadata` definition

```protobuf
message Metadata {
  string created = 1;
  string updated = 2;
}
```
