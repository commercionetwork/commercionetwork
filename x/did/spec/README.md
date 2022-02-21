## x/did

The `did` module allows the management of _identitities_ by associating a 
DID document to a `did:com:` address.
The module is also responsible for the historicization of identities.

The `Commercio.network` blockchain is the Verifiable Data Registry that should be used to perform DID resolution for the `did:com:` method.
In fact, the query functionalities of the `did` module provide all the necessary information to perform DID resolution for a certain address, requesting:
- The latest DID document and the corresponding metadata.
- The list of updates to the DID document and corresponding metadata.


## Historicization 

The `did` module has been updated to support the historicization of DID documents.

Following the latest [W3C Decentralized Identifiers (DIDs) v1.0 specification](https://www.w3.org/TR/2021/PR-did-core-20210803/), a DID resolution with no additional options should result in the latest version of the DID document for a certain DID plus additional metadata.
A DID document can be updated and its previous versions should remain accessible.

In `commercionetwork`, an identity is represented as the history of DID document updates made by a certain address.

Type `Identity` is made of two fields: 
- `DidDocument` - the JSON-LD representation of a DID document
- `Metadata` - with the `Created` and `Updated` timestamps

Querying for an `IdentityHistory` means asking for the list of updates to an `Identity`, sorted in chronological order.

Querying for an `Identity` means asking for the most recent version of the `DidDocument`, along with the associated `Metadata`.

Updating an `Identity` means appending to the blockchain store a new version of the DID document.

<!---
This operation uses the block time (guaranteed to be deterministic and always increasing) to populate the `Updated` field of `Metadata`. This timestamp is also used to populate the `Created` field, but only for the first version of the `Identity`.
Cosmos SDK store considerations:
- The key for storing an `Identity` is parameterized with the `ID` field of `DidDocument` (a `did:com:` address) and the `Updated` field of `Metadata` (timestamp). 
- The resulting key will look like the following. `did:identities:[address]:[updated]:`
- Since the value used for the `Updated` field is a timestamp guaranteed to be always increasing, then a store iterator with prefix `did:identities:[address]:` will retrieve values in ascending update order.
- For the same reason, the last value obtained by the same iterator will be the last identity appended to the store. Cosmos SDK allows to obtain a `ReverseIterator` returning values in the opposite order and therefore its first value will be the last updated identity.
- For a certain address only one update per block will persist, as a consequence of using the block time in the key.
--->