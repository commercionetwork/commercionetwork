# Sending a document reading receipt
Once you have received a document and you want to acknowledge the sender that you have properly read it, you can use 
the `SendDocumentReceipt` transaction that allows you to do that. 

## Transaction message
In order to properly create and send a transaction representing a document receipt you need to compose the 
`commercio/MsgSendDocumentReceipt` message:

```json
{
  "type": "commercio/MsgSendDocumentReceipt",
  "value": {
    "sender": "<Document sender address>",
    "recipient": "<Document recipient address>",
    "tx_hash": "<Tx hash in which the document has been sent>",
    "document_uuid": "<Document UUID>",
    "proof": "<Optional reading proof>"
  }
}
```

## Using the CLI
In order to send such a transaction using the CLI, you can execute the following command:

```bash
cncli tx commerciodocs send-document-receipt \
  [document-sender] \
  [document-recipient] \ 
  [tx-hash] \
  [document-uuid] \
  [proof]
```

### Parameters 
| Parameter | Type | Required | Description |  
| :-------- | :---: | :-----: | :---------- |
| `document-sender` | Address | Yes | The address of the original document sender | 
| `document-recipient` | Address | Yes | The address of the original document recipient |
| `tx-hash` | String | Yes | Hash of the transaction inside which is contained the sent document |
| `document-uuid` | Uuid | Yes | UUID of the document to which this receipt is related to |
| `proof` | String | No | Optional proof that the recipient has read the document | 

### Example usage 
```bash
cncli tx commerciodocs send-document \
  [document-sender] \
  [document-recipient] \ 
  [tx-hash] \
  [document-uuid] \
  [proof]
```