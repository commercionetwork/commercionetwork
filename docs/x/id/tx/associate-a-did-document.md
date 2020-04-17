# Associating a Did Document to your identity 
Being your account address a Did, using the Commercio Network blockchain you can associate to it a Did Document
containing the information that are related to your public (or private) identity.  
In order to do so you will need to perform a transaction and so your account must have first received some tokens. To
know how to get them, please take a look at the [*"Using an identity"* section](create-an-identity.md#using-an-identity). 

## Transaction message
In order to properly send a transaction to set a DID Document associating it to your identity, you will need
to create and sign the following message:

```javascript
{
  "type": "commercio/MsgSetIdentity",
  "value": {
    "@context": "https://www.w3.org/ns/did/v1",
    "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
    "publicKey": [
      {
        "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc#keys-1",
        "type": "RsaVerificationKey2018",
        "controller": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
        "publicKeyPem": "-----BEGIN PUBLIC KEY----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB-----END PUBLIC KEY-----\r\n"
      },
      {
        "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc#keys-2",
        "type": "RsaSignature2018",
        "controller": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
        "publicKeyPem": "-----BEGIN PUBLIC KEY----MIGfM3TvO3Ku3PJgZ9PO4qRw7+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCQVvTkCbc9A0GCSqGSIbqd4pNXtgbfbwJGviZ6kQIDAQAB-----END PUBLIC KEY-----\r\n"
      }
    ],
    "proof": {
      "type": "EcdsaSecp256k1VerificationKey2019",
      "created": "2019-02-08T16:02:20Z",
      "proofPurpose":"authentication",
      "controller": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
      "verificationMethod": "<did bech32 pubkey>",
      "signatureValue": "QNB13Y7Q91tzjn4w=="
    },
  }
}
```

### `value` fields requirements

| Field | Required | 
| :---: | :------: | 
| `@context` | Yes (Must be `https://www.w3.org/ns/did/v1`) |
| `id` | Yes |
| `publicKey` | Yes |
| `proof` | Yes |

The `id` field represents the DID you want to associate the provided identity to.

The `publicKey` field represents the public keys users can use to communicate safely with you.

Each key **must** have an `id` field defined by the concatenation of the content of the `id` field, along with a `#keys-NUMBER` suffix, where `NUMBER` must be an integer.

The `controller` key field must be equal to the `id` field content.

The commercio.network blockchain requires at least two keys, defined in the following way:

 - key with suffix `#keys-1` must be of type `RsaVerificationKey2018`, and must be a valid RSA PKIX public key;
 - key with suffix `#keys-2` must be of type `RsaSignature2018`, and must be a valid RSA PKIX public key.
 
A `commercio/MsgSetIdentity` transaction that **doesn't** meet these requirements will be discarded.

### Proof fields requirements

| Field | Required | Value | 
| :---: | :------: | :------: | 
| `type` | Yes | must always be `EcdsaSecp256k1VerificationKey2019` |
| `created` | Yes | creation date in UTC format |
| `proofPurpose` | Yes | must always be `authentication` |
| `controller` | Yes | same value specified in the `id` field |
| `verificationMethod` | Yes | bech32-encoded public key associated with the address specified in the `id` field |
| `signatureValue` | Yes | see explaination below |

### Creating the `signatureValue` value

In order to create `signatureValue`, the following steps must be followed

1. Create a `value` JSON as specified earlier, including only the `@context`, `id` and `publicKey` fields:
```javascript
{
 "@context": "https://www.w3.org/ns/did/v1",
 "id": "your DID",
 "publicKey": "your public keys",
}
```
2. alphabetically sort the `did_document_unsigned` and remove all the white spaces and line endings characters.
3. obtain hash of resulting string bytes using **SHA-256**. 
4. sign the result of the hashing process using your DID's public key, which you assigned to the `verificationMethod` `proof` JSON field
5. encode the result in **base64** obtaining `signatureValue`.

The signature commercio.network accepts is `EcdsaSecp256k1VerificationKey2019`, which is a type of elliptic-curve signature scheme.

The signature format produced in step 4, must be of the `r || s` kind, otherwise the identity creation **will** fail.

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setIdentity
```  
