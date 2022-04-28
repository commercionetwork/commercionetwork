# Did

The `did` module allows the management of _identitities_ by associating a 
DID document to a `did:com:` address.
The module is also responsible for the historicization of identities.

The `Commercio.network` blockchain is the Verifiable Data Registry that should be used to perform DID resolution for the `did:com:` method.
In fact, the `did` module provides query functionalities providing all the necessary information to perform DID resolution for a certain address, allowing to request:
- The latest DID document and the corresponding metadata.
- The list of updates to the DID document and corresponding metadata.

## Creating an identity
First of all, let's define what an **identity** is inside the Commercio Network blockchain.  

> An identity is the method used inside the Commercio Network blockchain in order to identify documents' senders and recipients.

In order to create an identity, you should start by creating a Commercio Network address, which will have the following form: 

```
did:com:<unique part>
```

The address it itself a DID Decentralized Identifier.

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
Being your account address a DID, using the Commercio Network blockchain you can associate to it a DID document containing the information that are related to your public (or private) identity.  
In order to do so you will need to perform a transaction and so your account must have first received some tokens. 

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