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
    "@context": "https://www.w3.org/2019/did/v1",
    "id": "<Your Address>",
    "publicKey": [
      {
        "id": "<Your Address>#keys-1",
        "type": "Secp256k1VerificationKey2018",
        "controller": "<Your Address>",
        "publicKeyHex": "<Public key value, hex encoded>"
      },
      {
        "id": "<Your Address>#keys-2",
        "type": "RsaVerificationKey2018",
        "controller": "<Your Address>",
        "publicKeyHex": "<Public key value, hex encoded>"
      },
      {
        "id": "<Your Address>#keys-3",
        "type": "Secp256k1VerificationKey2018",
        "controller": "<Your Address>",
        "publicKeyHex": "<Public key value, hex encoded>"
      }
    ],
    "authentication": [
      "<Authentication key id>"
    ],
    "proof": {
      "type": "LinkedDataSignature2015",
      "created": "<Creation time, in ISO 8601 format>",
      "creator": "<Authentication key id>",
      "signatureValue": "<Signature value, Base64 encoded>"
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
| `@context` | Yes (Must be `https://www.w3.org/2019/did/v1`) |
| `id` | Yes |
| `publicKey` | Yes (Must be of length 3) |
| `authentication` | Yes |
| `proof` | Yes |
| `service` | No |

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setIdentity
```  