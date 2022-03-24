# x/did

The `did` module allows the management of _identitities_ by associating a 
DID document to a `did:com:` address.
The module is also responsible for the historicization of identities.

The `Commercio.network` blockchain is the Verifiable Data Registry that should be used to perform DID resolution for the `did:com:` method.
In fact, the query functionalities of the `did` module provide all the necessary information to perform DID resolution for a certain address, requesting:
- The latest DID document and the corresponding metadata.
- The list of updates to the DID document and corresponding metadata.

## Creating an identity
First of all, let's define what an **identity** is inside the Commercio Network blockchain.  

> An identity is the method used inside the Commercio Network blockchain in order to identify documents' senders and recipients.

In order to create an identity, you simply have to create a Commercio Network address, which will have the following form: 

```
did:com:<unique part>
```

In order to do so, you can use the CLI and execute the following command: 

```bash
commercionetworkd keys add <key-name>
``` 

You will be required to set a password in order to safely store the key on your computer.  

:::warning
Please note that password will be later asked you when signing the transactions so be sure you remember it.
:::  

After inserting the password, you will be shown the mnemonic that can be used in order to import your account (and identity) into a wallet. 

```
- name: jack
  type: local
  address: did:com:13jckgxmj3v8jpqdeq8zxwcyhv7gc3dzmrqqger
  pubkey: did:com:pub1addwnpepqfdl6s8hdwdya9zvn5wtx8ty3qsqqqd2ddvygc5zutnrryh5x9ju73jdfg8
  mnemonic: ""
  threshold: 0
  pubkeys: []


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

scorpion what indoor keen topic cricket uphold inch cactus six suffer coin popular honey vendor clown day twin during vague midnight emerge man inform
```

### Using an identity
Once you have created it, in order to start performing a transaction with your identity you firstly have to fund your identity. 
Each and every transaction on the blockchain has a cost, and to pay for it you must have some tokens.  
If you want to receive some tokens in **Test-net**, please use faucet service or tell us inside our [official Telegram group](https://t.me/commercionetwork) 
and we will send you some as soon as possible.

### Associating a Did Document to your identity 
Being your account address a Did, using the Commercio Network blockchain you can associate to it a Did Document containing the information that are related to your public (or private) identity.  
In order to do so you will need to perform a transaction and so your account must have first received some tokens. 

### Updating an identity

Updating an `Identity` means appending to the blockchain store a new version of the DID document.
The transaction used to associate a DID document can be used to update the Identity.


<!-- TODO: check Msg format with Document -->

## Transaction message
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

The message must include a DID document conform to the rules of Decentralized Identitfiers (DIDs) v1.0 plus additional rules defined by commercionetwork. 
Please refer to [https://www.w3.org/TR/2021/PR-did-core-20210803/]() and to the following requirements.
A `commercio/MsgSetIdentity` transaction that **doesn't** meet these requirements will be discarded.

Fields that are NOT required can be omitted from the message.


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
- *<sup>2</sup> Supported values for the `type` field are: `Ed25519Signature2018`, `Ed25519VerificationKey2018`, `RsaSignature2018`, `RsaVerificationKey2018`, `EcdsaSecp256k1Signature2019`, `EcdsaSecp256k1VerificationKey2019`, `EcdsaSecp256k1RecoverySignature2020`, `EcdsaSecp256k1RecoveryMethod2020`, `JsonWebSignature2020`, `JwsVerificationKey2020`, `GpgSignature2020`, `GpgVerificationKey2020`, `JcsEd25519Signature2020`, `JcsEd25519Key2020`, `BbsBlsSignature2020`, `BbsBlsSignatureProof2020`, `Bls12381G1Key2020`, `Bls12381G2Key2020`.
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

The `id` and `serviceEndpoint` fields must conform to the rules of RFC3986 for URIs. 
Please refer to [https://datatracker.ietf.org/doc/html/rfc3986]().


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

