<!--
order: 4
-->

# Events

The `documents` module emits the following events:

## Handlers

### MsgShareDocument

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| new_saved_document | sender     | {senderAddress}         |
| new_saved_document | doc_id     | {documentUUID}          |
| new_saved_document | receiver_0 | {firstReceiverAddress}  |
| new_saved_document | receiver_1 | {secondReceiverAddress} |
| new_saved_document | ...        | ...                     |
| message            | action     | shareDocument         |
| message            | sender     | {senderAddress}         |

### MsgSendDocumentReceipt
| Type     | Attribute Key | Attribute Value     |
| -------- | ------------- | ------------------  |
| new_saved_receipt | receipt_id  | {receiptUUID}      |
| new_saved_receipt | document_id | {documentUUID}     |
| new_saved_receipt | sender      | {senderAddress}    |
| new_saved_receipt | recipient   | {recipientAddress} |
| message            | action     | sendDocumentReceipt         |
| message  | sender        | {senderAddress}     |
  








