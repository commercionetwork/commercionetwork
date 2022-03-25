<!--
order: 1
-->

# State

The `x/did` module keeps state of Identities, represented as the evolution of the DID Document and Metadata for a certain DID.

## Store

The module appends in the store the updated Identity.

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `did:identities:[address]:[updated]` | &rarr; | _Identity_ |

This operation uses the block time (guaranteed to be deterministic and always increasing) to populate the `Updated` field of `Metadata`. 
This timestamp is also used to populate the `Created` field, but only for the first version of the `Identity`.

## The `Identity` type

```
message Identity {
  DidDocument didDocument = 1;
  Metadata metadata = 2;
}
```

### `DidDocument` definition

```
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

```
message VerificationMethod {
  option (gogoproto.equal) = true;
  string ID = 1;
  string Type = 2;
  string Controller = 3;
  string publicKeyMultibase = 4;
}
```

#### `Service` definition

```
message Service {
  option (gogoproto.equal) = true;
  string ID = 1;
  string type = 2;
  string serviceEndpoint = 3;
}
```

### `Metadata` definition

```
message Metadata {
  string created = 1;
  string updated = 2;
}
```

## DID Resolution

In `commercionetwork`, an identity is represented as the history of DID document updates made by a certain address.

Following the latest [W3C Decentralized Identifiers (DIDs) v1.0 specification](https://www.w3.org/TR/2021/PR-did-core-20210803/), a DID resolution with no additional options should result in the latest version of the DID document for a certain DID plus additional metadata.

Querying for an `Identity` means asking for the most recent version of the `DidDocument`, along with the associated `Metadata`.
The result will be an `Identity` made of two fields: 
- `DidDocument` - the stored DID document JSON-LD representation
- `Metadata` - including the `Created` and `Updated` timestamps

### Historicization

The `did` module has been updated to support the historicization of DID documents.
A DID document can be updated and its previous versions should remain accessible.

Querying for an `IdentityHistory` means asking for the list of updates to an `Identity`, sorted in chronological order.

