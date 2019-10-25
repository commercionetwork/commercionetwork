# Sending a document reading receipt
Once you have received a document and you want to acknowledge the sender that you have properly read it, you can use 
the `MsgSendDocumentReceipt` message that allows you to do that. 

## Transaction message
In order to properly send a transaction to share a document, you will need to create and sign the
following message.

```json
{
  "type": "commercio/MsgSendDocumentReceipt",
  "value": {
    "uuid": "<Unique receipt identifier>",
    "sender": "<Document sender address>",
    "recipient": "<Document recipient address>",
    "tx_hash": "<Tx hash in which the document has been sent>",
    "document_uuid": "<Document UUID>",
    "proof": "<Optional reading proof>"
  }
}
```

### Fields requirements
| Field | Required | 
| :---: | :------: | 
| `uuid` | Yes |
| `sender` | Yes | 
| `recipient` | Yes | 
| `tx_hash` | Yes | 
| `document_uuid` | Yes |
| `proof` | No | 

## Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
sendDocumentReceipt
```