# Add a trusted metadata schema proposer
When [submitting a new officially recognized document metadata scheme](add-supported-metadata-schema.md) 
you must be a **trusted metadata scheme proposer**.  

Inside this page there are all the references for the *government* so that they can add you as such after you've
required them to do so. 

Please note that using the term *government* we're referring to the entity of 
[Commercio Consortium](https://commercioconsortium.org/) which is, inside our chain, the ultimate decision maker 
when adding new trusted accounts having particular permissions inside the chain itself. 

## Transaction message
In order to properly send a transaction to share a document, you will need to create and sign the
following message.

```json
{
  "type": "commercio/MsgAddTrustedMetadataSchemaProposer",
  "value": {
    "proposer": "did:com:1fwr65ph34yfzejkgly2pj6druxyexn797gmmvp",
    "signer": "did:com:1ljw7ny7rx3jewq85hpkln82wzm9lxqwaruxjdg"
  }
}
```

### Fields requirements
| Field | Required | 
| :---: | :------: |
| `proposer` | Yes |
| `signer` | Yes *<sup>1</sup> |  

- *<sup>1</sup> Please note that the `signer` address must be the one of the government account that has been set 
during the genesis using the `set-genesis-government-address` command.

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
addTrustedMetadataSchemaProposer
```  