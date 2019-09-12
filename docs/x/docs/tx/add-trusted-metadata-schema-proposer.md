# Add a trusted metadata schema proposer
When [submitting a new officially recognized document metadata scheme](add-supported-metadata-schema.md) 
you must be a **trusted metadata scheme proposer**.  

Inside this page there are all the references for the *government* so that they can add you as such after you've
required them to do so. 

Please note that using the term *government* we're referring to the entity of 
[Commercio Consortium](https://commercioconsortium.org/) which is, inside our chain, the ultimate decision maker 
when adding new trusted accounts having particular permissions inside the chain itself. 

## Transaction message
In order to add a new trusted schema proposer, the government must create and 
sign a `commercio/MsgAddTrustedMetadataSchemaProposer` message:

```json
{
  "type": "commercio/MsgAddTrustedMetadataSchemaProposer",
  "value": {
    "proposer": "did:com:1fwr65ph34yfzejkgly2pj6druxyexn797gmmvp",
    "signer": "did:com:1ljw7ny7rx3jewq85hpkln82wzm9lxqwaruxjdg"
  }
}
```

Please note that the `signer` address must be the one of the government account that has been set 
during the genesis using the `set-genesis-government-address` command.  