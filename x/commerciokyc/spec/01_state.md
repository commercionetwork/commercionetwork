<!--
order: 1
-->

# State

The `commerciokyc` module keeps state of the following objects

## Store

### Memberships


| Key |  | Value |
| ------- | ---------- | ---------- | 
| `commerciokyc:storage:[owner]` | &rarr; | _Membership_ |

### Invites

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `commerciokycinvite:[user]` | &rarr; | _Invite_ |


### Invites

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `commerciokycinvite:[user]` | &rarr; | _Invite_ |


### Trusted service providers

| Key |  | Value |
| ------- | ---------- | ---------- | 
| `commerciokyc:signers` | &rarr; | _Tsps_ |



## Type definitions

#### `Membership` definition

```protobuf
message Membership {
  string owner = 1;
  string tsp_address = 2;
  string membership_type = 3;
  google.protobuf.Timestamp expiry_at = 4 [(gogoproto.stdtime) = true];
}
```



#### `Invite` definition

```protobuf
message Invite {
  string sender = 1;
  string sender_membership = 2;
  string user = 3;
  uint64 status = 4;
}
```

#### `Trusted service providers` definition


```protobuf
message TrustedServiceProviders {
  repeated string addresses = 1;
}
```
