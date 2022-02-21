<!--
order: 1
-->

# State

## `Identity`

An identity consists of a DID document and metadata.

```
message Identity {
  DidDocument didDocument = 1;
  Metadata metadata = 2;
}
```

## `DidDocument`

A DID document consists of 

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

## `VerificationMethod`

```
message VerificationMethod {
  option (gogoproto.equal) = true;
  string ID = 1;
  string Type = 2;
  string Controller = 3;
  string publicKeyMultibase = 4;
}
```

## `Service`

```
message Service {
  option (gogoproto.equal) = true;
  string ID = 1;
  string type = 2;
  string serviceEndpoint = 3;
}
```

## `Metadata`

```
message Metadata {
  string created = 1;
  string updated = 2;
}
```