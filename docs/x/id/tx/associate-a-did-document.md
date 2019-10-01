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
  "type": "commercio/SetIdentity",
  "value": {
    "owner": "<Your Did>",
    "did_document": {
      "uri": "<Uri of the Did Document content>",
      "content_hash": "<Sha256 hash of the Did document content, hex encoded>"
    }
  }
}
```

### Transaction example
```json
{
  "chain_id": "test-chain-GrNuU0",
  "account_number": "0",
  "sequence": "1",
  "fee": {
    "amount": [],
    "gas": "200000"
  },
  "msgs": [
    {
      "type": "commercio/SetIdentity",
      "value" : {
        "owner": "did:com:1flzcn7yy9p04qwhh67hu8r38ar7ylxde2k47pr",
        "did_document": {
          "uri": "https://example.com/my-did",
          "content_hash": "9c5ef543dc05e7927da16e8d8a24372f0d064979e226a70cdea40a031d1daf51"
        }
      }
    }
  ],
  "memo": ""
}
```