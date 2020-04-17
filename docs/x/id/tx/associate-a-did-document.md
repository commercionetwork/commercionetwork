# Associating a Did Document to your identity 
Being your account address a Did, using the Commercio Network blockchain you can associate to it a Did Document
containing the information that are related to your public (or private) identity.  
In order to do so you will need to perform a transaction and so your account must have first received some tokens. To
know how to get them, please take a look at the [*"Using an identity"* section](create-an-identity.md#using-an-identity). 

## Transaction message
In order to properly send a transaction to set a Did Document associating it to your identity, you will need
to create and sign the following message:

```json
{
  "type": "commercio/MsgSetIdentity",
  "value": {
    "@context": "https://www.w3.org/ns/did/v1",
    "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc",
    "publicKey": [
      {
        "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc#keys-1",
        "type": "RsaVerificationKey2018",
        "controller": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc ",
        "publicKeyPem": "-----BEGIN PUBLIC KEY----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB-----END PUBLIC KEY-----\r\n"
      },
      {
        "id": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc#keys-2",
        "type": "RsaSignature2018",
        "controller": "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc ",
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
    "service": [
      {
        "id": "<Service id reference>",
        "type": "<Service type>",
        "serviceEndpoint": "<Service endpoint>"
      }
    ]
  }
}
```

### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `@context` | Yes (Must be `https://www.w3.org/ns/did/v1`) |
| `id` | Yes |
| `publicKey` | Yes (Must be of length 2) |
| `proof` | Yes |
| `service` | No |

### Proof fields requirements
| Field | Required | Value | 
| :---: | :------: | :------: | 
| `type` | Yes | "EcdsaSecp256k1VerificationKey2019" |
| `created` | Yes | Creation date in UTC format |
| `proofPurpose` | Yes | "authentication" |
| `controller` | Yes | User did |
| `verificationMethod` | Yes | Public key associated to user did hex encoded |
| `signatureValue` | Yes | Read the explanation below |

### Creating the `signatureValue` value

In order to create `signatureValue`, the following steps must be followed

1. Create the `did_document_unsigned` json formed as follow.
```json
{
 "@context": "https://www.w3.org/ns/did/v1",
 "id": "<User Did bech32 format>",
 "publicKey": "<json contains public kyes>",
}
```
2. Alphabetically sort the `did_document_unsigned` and remove all the white spaces and line endings characters.
3. Obtain hash of resulting string bytes using **Sha3-256**. 
4. Sign the resulting hash using `VerificationMethod` value with **Secp256k1Sign** algorithm.
5. Encode the result using **Base64** obtaining `signatureValue`.

### Service 

`Service` contains a list of Trusted Service End Point or Service End Point for a specific purposes
 

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setIdentity
```  