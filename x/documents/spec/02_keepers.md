<!--
order: 2
-->

# Keepers

The keeper of the `documents` module enforces the following requirements.

## Requirements for Documents

Documents are identified with a unique ID in the UUID format. 
Therefore, IDs cannot be reused.
The keeper marks as invalid requests trying to reuse an ID for storing a Document.

## Requirements for Document Receipts

Document Receipts are identified with a unique ID in the UUID format. 
Therefore, IDs cannot be reused.
The keeper marks as invalid requests trying to reuse an ID for storing a Document Receipt.

Also, only one Document Receipt can be sent by a user for a certain Document.
The keeper marks as invalid requests trying to send a Document Receipt if a previous one has already been shared.