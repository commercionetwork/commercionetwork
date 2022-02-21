<!--
order: 5
-->

# Client

## Transactions

### Setting an identity

#### CLI

```sh
tx did set-identity [did_document_proposal_path]
```


## Queries

### Showing an identity

#### CLI

```sh
query did show-identity [did]
```

#### REST

Endpoint:
   
```
/commercionetwork/did/identities/{did}
```

##### Example

Getting the latest identity of `did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd`:

```
http://localhost:1317/commercionetwork/did/identities/did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd
```

#### Response
```json
{
  "identity": {
    "didDocument": {
      "@context": [
        "https://www.w3.org/ns/did/v1",
        "https://w3id.org/security/suites/ed25519-2018/v1",
        "https://w3id.org/security/suites/x25519-2019/v1"
      ],
      "ID": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
      "verificationMethod": [
        {
          "ID": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytdw#keys-1",
          "Type": "RsaVerificationKey2018",
          "Controller": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
          "publicKeyMultibase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPicbLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
        },
        {
          "ID": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd#keys-2",
          "Type": "RsaSignature2018",
          "Controller": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
          "publicKeyMultibase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPicbLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
        }
      ],
      "authentication": [
        "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytdw#keys-1"
      ],
      "keyAgreement": [
        "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd#keys-2"
      ],
      "service": [
        {
          "ID": "A",
          "type": "agent",
          "serviceEndpoint": "https://commerc.io/agent/serviceEndpoint/"
        }
      ]
    },
    "metadata": {
      "created": "2022-02-21T09:20:59Z",
      "updated": "2022-02-21T09:20:59Z"
    }
  }
}
```
