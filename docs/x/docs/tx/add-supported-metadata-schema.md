# Adding a metadata schema as officially supported
When [sharing a document](send-document.md), you have the option to specify an officially recognized metadata schema
using the `schema_type` field. 

In this page we describe how to add a new schema specification as officially recognized. 

## Requirements
First of all, you need to be an **trusted metadata schema proposer**. If you wish to become one, please refer 
to the [proper page](../trusted-metadata-schema-proposers.md). 

If you have more than one account that is a trusted schema proposer, you can use whichever you want. 

## Transaction message
In order to add a metadata schema as officially recognized, you need to use the 
`commercio/MsgAddSupportedMetadataSchema` message:

```json
{
  "type": "commercio/MsgAddSupportedMetadataSchema",
  "value": {
    "signer": "<Proposal signer>",
    "schema": {
      "type": "<Unique metadata schema type>",
      "schema_uri": "<Uri linking to the schema definition>",
      "version": "<Version of the schema>"
    }
  }
}
```

### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `signer` | Yes *<sup>1</sup> |
| `schema` | Yes |

- *<sup>1</sup> The `signer` value should be the address of the account used to sign the transaction, 
which should also be a **trusted metadata schema proposer**. Read the [requirements](#requirements) for more details.

#### `metadata`
| Field | Required | 
| :---: | :------: |
| `type` | Yes |
| `schema_uri` | Yes | 
| `version` | Yes |

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
addSupportedMetadataSchema
```  