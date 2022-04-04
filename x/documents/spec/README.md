<!--
order: 0
title: Documents Overview
parent:
  title: "Documents"
-->

# Documents 

## Abstract

The `documents` module allows a user to share a document to other users.
Then, the receivers can send back to the sender a receipt proving that they have seen the document.

**IMPORTANT!!!**

* We are not actually sharing Documents on a blockchain. 
* We are sharing a transaction on a blockchain with a document footprint (HASH) 

An hash is the output of a hashing algorithm like SHA (Secure Hash Algorithm). These algorithms essentially aim to produce a unique, fixed-length string – the hash value, or “message digest” – for any given piece of data or “message”. 

As every electronic file is just data that can be represented in binary form, a hashing algorithm can take that data and run a complex calculation on it and output a fixed-length string as the result of the calculation. 

The result is the file’s hash value or message digest.

'Sharing a Document' on Commercio.network means sending a shareDoc transaction on a blockchain with your document hash.


### ShareDoc real world use cases

* Legally prove a document was shared with a third party
* Timestamp a document 
* Prove the existence of document 
* Notarize a document  


## Contents

1. **[State](01_state.md)**
2. **[Messages](03_messages.md)**
   - [MsgShareDocument](03_messages.md#share-a-document-with-msgsharedocument)
   - [MsgSendDocumentReceipt](03_messages.md#send-a-document-receipt-with-msgsenddocumentreceipt)
3. **[Events](04_events.md)**
   - [Handlers](04_events.md#handlers)
4. **[Client](05_client.md)**
   - [Query](05_client.md#query)
   - [gRPC](05_client.md#gRPC)
   - [Rest](05_client.md#rest)
