<!--
order: 1
-->

# State

The `x/did` module keeps state of Identities, represented as DID Document and Metadata.

## Store


<!---
This operation uses the block time (guaranteed to be deterministic and always increasing) to populate the `Updated` field of `Metadata`. This timestamp is also used to populate the `Created` field, but only for the first version of the `Identity`.
Cosmos SDK store considerations:
- The key for storing an `Identity` is parameterized with the `ID` field of `DidDocument` (a `did:com:` address) and the `Updated` field of `Metadata` (timestamp). 
- The resulting key will look like the following. `did:identities:[address]:[updated]:`
- Since the value used for the `Updated` field is a timestamp guaranteed to be always increasing, then a store iterator with prefix `did:identities:[address]:` will retrieve values in ascending update order.
- For the same reason, the last value obtained by the same iterator will be the last identity appended to the store. Cosmos SDK allows to obtain a `ReverseIterator` returning values in the opposite order and therefore its first value will be the last updated identity.
- For a certain address only one update per block will persist, as a consequence of using the block time in the key.
--->


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
