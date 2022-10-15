<!--
order: 5
-->

# Client

## Transactions

### Setting an identity

#### CLI

```bash
tx did set-identity [did_document_proposal_path]
```

### Parameters  
| Parameter | Description |
| :-------: | :---------- | 
| `did_document_proposal_path` | The OS path to a `.json` file containing the DID document that must be associated with an identity. |

For example, the user controlling the DID `did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd` could be interested in proposing the following DID document, saved in a `.json` file:
```
{
    "@context": [
        "https://www.w3.org/ns/did/v1",
        "https://w3id.org/security/suites/ed25519-2018/v1"
    ],
    "id": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
    "verificationMethod": [
        {
            "type": "RsaVerificationKey2018",
            "id": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytdw#keys-1",
            "controller": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
            "publicKeyMultiBase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPicbLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
        },
        {
            "type": "RsaSignature2018",
            "id": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd#keys-2",
            "controller": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
            "publicKeyMultiBase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPicbLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
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
            "id": "A",
            "type": "agent",
            "serviceEndpoint": "https://commercio.network/agent/serviceEndpoint/"
        }
    ]
}
```


## Queries

### Showing an identity

#### CLI

```bash
query did show-identity [did]
```

### Parameters  
| Parameter | Description |
| :-------: | :---------- | 
| `did` | Address of the user for which to read the Did Document |

#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.did.Query/Identity
```

##### Example


```bash
grpcurl -plaintext \
    -d '{"ID":"did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.did.Query/Identity
```

##### Response
```json
{
  "identity": {
    "didDocument": {
      "context": [
        "https://www.w3.org/ns/did/v1"
      ],
      "ID": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20",
      "verificationMethod": [
        {
          "ID": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20#keys-1",
          "Type": "RsaVerificationKey2018",
          "Controller": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20",
          "publicKeyMultibase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq1/cwtXIoJxvp6jJpqYqxNlXGhdW9Elwpc8OW//hgSMknHqDzKOOfDUiJbiPq0/BMD00im0aBbVCezkSfQxmu1aB4RMCh39CYe5qbkSDXnKfwCI72e26Om/lJGfSUe11xVBHsjT9r5SugGFVJ3ljwG5MFY50sdeAeX16nUZbHjRIh/E+n2oZlNV0uhsxsV4HtzA0Km9D+D8bBCEt4PUXHCE6kOmPORfPS7UvG/q5IZUlPeSGPXJAl3kYOT/s58ptxe4Id8Np0h6g/iEfdxw7SG6NS7GVA52blg8mFbGv0ynQCor/P+Vfh7fRnZC9wzCSXqis4PGNVq1s1MGQ13bM5wIDAQAB"
        },
        {
          "ID": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20#keys-2",
          "Type": "RsaSignature2018",
          "Controller": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20",
          "publicKeyMultibase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxk51GIYeuiZ24rlmKXQ/LoiYOuMNFM1/HhK4pEh/G6tu59U13zGx5eiIZWtdbpesEGt3Uw67oQUV+3YSGjalUMl+HUtVmk1yL0zQz85b2np+SZxxE4DbaQunaFfFy1hfNEmsrh6mpXuPeU4TCEm7vri4wfWheqn939UuIgsOoIc6vRbVeQWHU13O/EvJsquItuunCWNMOD3FEAuSDTf4X/2LgE2q07Gwf643lipS7AzXB/sNJL0hg/QH8dNz9uaW2TcIqbmr8TtOejPTwBsTcOM/mG+WGfhAmR1HYmcdiZPn47bnjamMv5OwaoSrEjMxB8tvn8W+AGm8DuUHkk5TjwIDAQAB"
        }
      ],
      "authentication": [
        "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20#keys-1"
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
      "created": "2022-02-03T13:53:28Z",
      "updated": "2022-02-03T13:53:28Z"
    }
  }
}
```






#### REST

Endpoint:
   
```
/commercionetwork/did/{ID}/identities
```

##### Example

Getting the latest identity of `did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd`:

```
http://localhost:1317/commercionetwork/did/did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd/identities
```

##### Response
```json
{
  "identity": {
    "didDocument": {
      "@context": [
        "https://www.w3.org/ns/did/v1",
        "https://w3id.org/security/suites/ed25519-2018/v1"
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
          "serviceEndpoint": "https://commercio.network/agent/serviceEndpoint/"
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

Please note that in the metadata the fields `created` and `updated` are equal, meaning that this was the first DID document association made for this DID. 

### Showing the history of updates to an identity

#### CLI

```bash
query did show-history [did]
```

### Parameters  
| Parameter | Description |
| :-------: | :---------- | 
| `did` | Address of the user for which to read the history of Did Document updates|

#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.did.Query/IdentityHistory
```

##### Example

```bash
grpcurl -plaintext \
    -d '{"ID":"did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.did.Query/IdentityHistory
```

##### Response
```json
{
  "identities": [
    {
      "didDocument": {
        "context": [
          "https://www.w3.org/ns/did/v1"
        ],
        "ID": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20",
        "verificationMethod": [
          {
            "ID": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20#keys-1",
            "Type": "RsaVerificationKey2018",
            "Controller": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20",
            "publicKeyMultibase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAq1/cwtXIoJxvp6jJpqYqxNlXGhdW9Elwpc8OW//hgSMknHqDzKOOfDUiJbiPq0/BMD00im0aBbVCezkSfQxmu1aB4RMCh39CYe5qbkSDXnKfwCI72e26Om/lJGfSUe11xVBHsjT9r5SugGFVJ3ljwG5MFY50sdeAeX16nUZbHjRIh/E+n2oZlNV0uhsxsV4HtzA0Km9D+D8bBCEt4PUXHCE6kOmPORfPS7UvG/q5IZUlPeSGPXJAl3kYOT/s58ptxe4Id8Np0h6g/iEfdxw7SG6NS7GVA52blg8mFbGv0ynQCor/P+Vfh7fRnZC9wzCSXqis4PGNVq1s1MGQ13bM5wIDAQAB"
          },
          {
            "ID": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20#keys-2",
            "Type": "RsaSignature2018",
            "Controller": "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20",
            "publicKeyMultibase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxk51GIYeuiZ24rlmKXQ/LoiYOuMNFM1/HhK4pEh/G6tu59U13zGx5eiIZWtdbpesEGt3Uw67oQUV+3YSGjalUMl+HUtVmk1yL0zQz85b2np+SZxxE4DbaQunaFfFy1hfNEmsrh6mpXuPeU4TCEm7vri4wfWheqn939UuIgsOoIc6vRbVeQWHU13O/EvJsquItuunCWNMOD3FEAuSDTf4X/2LgE2q07Gwf643lipS7AzXB/sNJL0hg/QH8dNz9uaW2TcIqbmr8TtOejPTwBsTcOM/mG+WGfhAmR1HYmcdiZPn47bnjamMv5OwaoSrEjMxB8tvn8W+AGm8DuUHkk5TjwIDAQAB"
          }
        ],
        "authentication": [
          "did:com:133nf49v83dts7pq30csnnl7ul5unsvzz55qt20#keys-1"
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
        "created": "2022-02-03T13:53:28Z",
        "updated": "2022-02-03T13:53:28Z"
      }
    }
  ]
}
```


#### REST

Endpoint:
   
```
/commercionetwork/did/{ID}/identities/history
```

##### Example

Getting the history of identity updates of `did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd`:

```
http://localhost:1317/commercionetwork/did/did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd/identities/history
```

You can notice that 

```
{
  "identities": [
    {
      "didDocument": {
        "context": [
          "https://www.w3.org/ns/did/v1",
          "https://w3id.org/security/suites/ed25519-2018/v1"
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
        "assertionMethod": [],
        "keyAgreement": [
          "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd#keys-2"
        ],
        "capabilityInvocation": [],
        "capabilityDelegation": [],
        "service": [
          {
            "ID": "A",
            "type": "agent",
            "serviceEndpoint": "https://commercio.network/agent/serviceEndpoint/"
          }
        ]
      },
      "metadata": {
        "created": "2022-02-21T09:20:59Z",
        "updated": "2022-02-21T09:20:59Z"
      }
    },
    {
      "didDocument": {
        "context": [
          "https://www.w3.org/ns/did/v1"
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
        "assertionMethod": [],
        "keyAgreement": [
          "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd#keys-2"
        ],
        "capabilityInvocation": [],
        "capabilityDelegation": [],
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
        "updated": "2022-02-21T17:19:10Z"
      }
    }
  ]
}
```

Please note that the field `context` of the second DID document contains only `"https://www.w3.org/ns/did/v1"` and that the `updated` is greater than `created`.

A DID resolution would consider this last updated DID document as result.