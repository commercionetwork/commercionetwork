# CommercioAPI ShareDoc

The  CommercioAPI ShareDoc allows you to send a document to another user, and retrieve the list of documents that you have received.


## What is an Electronic Cerified Delivery?

**IMPORTANT!!!**

* We are not sending Documents on a blockchain. 
* We are sending a transaction on a blockchain with a document footprint (HASH) 

An hash is the output of a hashing algorithm like SHA (Secure Hash Algorithm). These algorithms essentially aim to produce a unique, fixed-length string – the hash value, or “message digest” – for any given piece of data or “message”. 

As every electronic file is just data that can be represented in binary form, a hashing algorithm can take that data and run a complex calculation on it and output a fixed-length string as the result of the calculation. 

The result is the file’s hash value or message digest.

'Sharing a Document' on Commercio.network  means sending a shareDoc transaction on a blockchain with your document hash.


## shareDoc 

### shareDoc real world use cases

* Legally prove a document was shared with a third party
* Timestamp a document 
* Prove the existence of document 
* Notarize a document  


### Step by step Example

Let's create a new transaction to share the document hash (REMEMBER not the actual document, only the hash ) associated with the given contentUri and having the given metadata and checksum. 

Step 1


Step 2


Step 3




### API Code Examples


PHP Python C# Java Go


## sendReceipt

## documentList sent 

## documentList sent received

## receiptList sent 

## receiptList sent received