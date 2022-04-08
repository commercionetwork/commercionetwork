<!--
order: 3
-->

# Messages

## Update an identity with `MsgSetIdentity`

Updating an `Identity` means appending to the blockchain store a new version of the DID document.


### Protobuf message

```protobuf
message MsgSetIdentity { DidDocument didDocument = 1; }
```

### Transaction message
In order to properly send a transaction to set a DID Document associating it to your identity, you will need to create and sign the following message:

```javascript
{
  "type": "commercio/MsgSetIdentity",
  "value": {
    "didDocument":
        {
            "@context": [
                "https://www.w3.org/ns/did/v1",
                "string"
            ],
            "id": "string",
            "verificationMethod": [
                {
                    "type": "string",
                    "id": "string",
                    "controller": "string",
                    "publicKeyMultiBase": "string"
                },
                {
                    "type": "string",
                    "id": "string",
                    "controller": "string",
                    "publicKeyMultiBase": "string"
                }
            ],
            "authentication": [
                "string"
            ],
            "assertionMethod": [
                "string"
            ],
            "keyAgreement": [
                "string"
            ],
            "capabilityInvocation": [
                "string"
            ],
            "capabilityDelegation": [
                "string"
            ],
            "service": [
                {
                    "id": "string",
                    "type": "string",
                    "serviceEndpoint": "string"
                }
            ]
        }
  }
}
```

The message must include a DID document conforming to the rules of [Decentralized Identifiers (DIDs) v1.0](https://www.w3.org/TR/2021/PR-did-core-20210803/) plus additional rules defined by commercionetwork.
A `commercio/MsgSetIdentity` transaction that **doesn't** meet these requirements will be discarded.

**Fields that are NOT required can be omitted from the message.**

### `didDocument` field requirements

| Field                  | Required | Value |
| ---                  | ------ | --- |
| `@context`             | Yes      | `["https://www.w3.org/ns/did/v1","https://w3id.org/security/suites/ed25519-2018/v1"]` |
| `id`                   | Yes      | `"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"` |
| `verificationMethod`   | Yes      | Consider the values in the description below | 
| `authentication`       | No       | `"#keys-authentication"` |
| `assertionMethod`      | No       | `"#keys-assertionMethod"` |
| `keyAgreement`         | No       | `"#keys-keyAgreement"` |
| `capabilityInvocation` | No       | `"#keys-capabilityInvocation"` |
| `capabilityDelegation` | No       | `"#keys-capabilityDelegation"` |
| `service`              | No       | Consider the values in the description below |

### `verificationMethod` field requirements
| Field                  | Required | Value |
| ---                  | ------ | --- |
| `id`                   | Yes *<sup>1</sup> | `"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytdw#keys-1"` (absolute) or `"#keys-1"` (relative) | 
| `type`                 | Yes *<sup>2</sup> | `"RsaVerificationKey2018"` | 
| `controller`           | Yes *<sup>3</sup> | `"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"` | 
| `publicKeyMultiBase`   | Yes *<sup>4</sup> | `"mMIIBIjANBgkqh...3awGwIDAQAB"` | 

- *<sup>1</sup> The `id` field supports both absolute (e.g. `"did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytdw#keys-1"`) and relative (`"#keys-1"`) identifiers.
- *<sup>2</sup> Supported values for the `type` field are: 
    - `Ed25519Signature2018`|
    - `Ed25519VerificationKey2018` 
    - `RsaSignature2018`
    - `RsaVerificationKey2018`
    - `EcdsaSecp256k1Signature2019`
    - `EcdsaSecp256k1VerificationKey2019`
    - `EcdsaSecp256k1RecoverySignature2020`
    - `EcdsaSecp256k1RecoveryMethod2020`
    - `JsonWebSignature2020`
    - `JwsVerificationKey2020`
    - `GpgSignature2020`
    - `GpgVerificationKey2020`
    - `JcsEd25519Signature2020`
    - `JcsEd25519Key2020`
    - `BbsBlsSignature2020`
    - `BbsBlsSignatureProof2020`
    - `Bls12381G1Key2020`
    - `Bls12381G2Key2020`
- *<sup>3</sup> `controller` must be equal to the DID document `id` field.
- *<sup>4</sup> For more information about this field format, please refer to [The Multibase Data Format](https://tools.ietf.org/id/draft-multiformats-multibase-00.html). The example value `"mMIIBIjANBgkqh...3awGwIDAQAB"` start with `m` and therefore the rest of the string is in base64 [RFC 4648](https://datatracker.ietf.org/doc/html/rfc4648) no padding.

Additional requirements:
- a verification method of type `RsaVerificationKey2018` must have the suffix `#keys-1` in the `id` field, and must be a valid _RSA PKIX_ public key;
- a verification method of type `RsaSignatureKey2018`, must have the suffix `#keys-2` in the `id` field, and must be a valid _RSA PKIX_ public key.



### `service` field requirements
| Field                  | Required | Value |
| ---                  | ----- | --- |
| `id`                   | Yes      | `"Service001"` | 
| `type`                 | Yes      | `"agent"` | 
| `serviceEndpoint`      | Yes      | `"https://commercio.network/agent/serviceEndpoint/"`      | 

The `id` and `serviceEndpoint` fields must conform to the [rules of RFC3986 for URIs](https://datatracker.ietf.org/doc/html/rfc3986).