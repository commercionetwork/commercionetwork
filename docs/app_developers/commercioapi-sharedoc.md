# CommercioAPI eDelivery

The  CommercioAPI eDelivery allows you to operate with two type messages available in the Custom  `documents` module 

* [MsgShareDocument](/modules/documents/03_messages.html#share-a-document-with-msgsharedocument)
* [MsgSendDocumentReceipt](/modules/documents/03_messages.html#send-a-document-receipt-with-msgsenddocumentreceipt)


## MsgShareDocument

The message permit to
* Store a `hash of a document` you own in the blockchain at a specific time.
* Link the message to your identity account (your wallet address).
* Specify the recipient's identity account (receiver's wallet address).
* Provide basic metadata.
* Digitally sign (raw signature) the message using your account (which you exclusively control).


This features are knonw as `Notarization of a document in the blockchain`

### How a MsgShareDocument is structured and entites usage

The MsgShareDocument message is used to certify/notarize the existence of a file at a specific moment in time (the sharedoc's date) through its footprint (document's hash)


The complete structure of the message registered in the chain is documented [here](/modules/documents/03_messages.html#transaction-message) 

The following explains the use of some of the main entities of the message.

**sender** : The wallet address of the creator of the MsgShareDocument

**recipients** : The list of wallet addresses of the recipients of the MsgShareDocument

**uuid** : The unique identification code of the MsgShareDocument in the chain 

**checksum.value** : is the hash of the fisical document associated with the MsgShareDocument

**checksum.algorithm** : is the hashing method used to calculate `checksum.value`

**content_uri** : This is the URI where the document associated with MsgShareDocument is stored (usually an encrypted URI path and document).

**metadata** : The metadata section is designed and intended to provide additional functional information on how to 'handle' the file for which the MsgShareDocument is being performed. It is entirely at the discretion of the developer and creator of the message to determine what to indicate and how to use the entities that have been prepared for this purpose. There is no specific usage defined at the protocol level.

The aim is to make a file available, for example, hosted out of the chain, where you can indicate categorization information about the file for which you are creating a MsgShareDocument.

For example, a file (this is purely an example) could be a .json file containing a series of information.


```
{
  "creator": "Administration ACME ltd",  
  "user_code": "3452SFTa",
  "creation_date": "2023-07-12T13:27:17Z",
  "sender_email": "administration@acme.com",
  "receiver_email": "johndoe@user.com,
  "category": "orders,
  "file_format": "pdf/A",
    "product_category": [
    "trousers",
    "Tshirts",
  ],
  .....
  
  "delivery_date": "2023-07-30T13:27:17Z",
}
```


The **metadata** section allows you to indicate

``````
 "metadata": {
      "content_uri": "<Metadata content URI>",
      "schema": {
        "uri": "<Metadata schema definition URI>",
        "version": "<Metadata schema version>"
      },
    },

``````

**metadata.content_uri** : the URI where the metadata file is stored

**metadata.schema.uri** : the URI where an optional schema file that describes the content structure of the content file is stored

**metadata.schema.version** : the version of the schema file indicated  in **metadata.schema.uri**

In reference to the example, the URL to retrieve the file can be indicated in **metadata.content_uri**, preferably in an encrypted manner to avoid exposing the endpoint too much since the message is public.

Of course, the contents of the file could be change over time, and the **MsgShareDocument** only certifies the content of the hashed file, not the content of the .json file itself.

The important thing to understand is that whatever is indicated in the external metadata file has no certification, as it can still be externally manipulated.

Having said that, there are no rules or restrictions that prohibit using the entities in a different manner. It suggests that users have the freedom to use the entities as they see fit, without any specific limitations or guidelines. You could indicate custom information with total creativity and applying encryption as desired in the application.

For example if you have data or similar information that is smaller than 512 bytes and you wish to save it in a fixed manner within the MsgShareDocument, you can encrypt it and indicate it in the metadata/content_uri field.

In summary, the protocol allows flexibility in how you handle and encrypt the data, and you have full control over the content of metadata entities  and application logic required to work with the data within the metadata section.

### Questions and answers 

**IMPORTANT!!!**

* We are not actually sharing Documents on a blockchain. 
* We are Storing a Message through a transaction on a blockchain with a document footprint indicated (HASH) 

**Can I notarize only PDF file ?**

No, you can notarize any digital file for which you can generate a hash using the permitted algorithm. Therefore, you can notarize file extensions such as .pdf, .txt, .doc, .docx, .xls, .json, .xml, .zip, .ppt, etc.










## Check the hash of a document 

In order to check if a file you own correspond to the one notarized in the blockchain  you can perform a verification directly with <b>Almerico</b> throught a simple widget that permit  to  drag&drop a file  using a specific tool and verify if the calculated footprint of the file (`hash`) correspond to the one notarized in the message.

To perform the validation check:

**1) Open message**
 Open the spacific page of message for the specific transaction.

Example 

```

https://testnet.commercio.network/transactions/sharedoc/1348F1AB13E473A94D5656445D0F49FE1924CC3340B533C5C24EA8E2D7FACC43/uuid/6c509472-ead2-4f6f-89c3-f30206c7a737

```

<small>Nb: the url page is composed by tx has and uuid of the message</small>
<small>https://testnet.commercio.network/transactions/sharedoc/#TXHASH#/uuid/#MESSAGE_UUID#</small>

<br><br><br>

**2) Use the tool**
Drag&drop your file in the widget as asked  

![Modal](./dragNdrop_hash_check.png)

**3) Check result** 
You could get the following results 

**Success**

![Modal](./verification_success.png)

It means that the  hash calculated with the method (sha-256,md5 ecc) indicated in the  MsgShareDocument on the file droped in the Drag&Drop area correspond to the one certified in the message


**Failure**

It means that the  hash calculated with the method (sha-256,md5 ecc) indicated in the  MsgShareDocument on the file droped in the Drag&Drop area <b>DO NOT</b> correspond to the one certified in the message

![Modal](./verification_failure.png)


<small>Same check could be obviuolsy done with other public hashing site or localy with internal hashing tools for example through the shell's functions</small>




## MsgSendDocumentReceipt
The message permit to 
* store in the blockchain a message in a precise time associated to a `MsgShareDocument` previously stored  in order to certify you have checked the corrispondance of the `hash of a document` in the  
`MsgShareDocument` message of the document yopu have revevied off the chain 
*  Sign (raw) the message with your account (You only control) as recevier of the `MsgShareDocument` to the Sender




For more detail refers to the [Document Module](/modules/documents/#documents)




## ShareDoc processes 


### Send a shareDoc Message
Permit to create a process to send a message in the block chain named `MsgShareDocument` or Sharedoc message throught the DOCS  Module


#### PATH

POST : `/sharedoc/process`


#### Step by step Example

Let's create a new process to send a Sharedocument message containig the hash (REMEMBER not the actual document, only the hash ) of a document associated with the given contentUri and having the given metadata and checksum. 

##### Step 1 - Define the first query 

Parameter value :
* Your account address (authenticated user): es `did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3`
* Account address of the recipient/s: es `did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3`
* Hash of the document with `sha-256` algorithm: `3cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824`
* Encripted content uri : `8cc590c1823ee24dae77eadfc3b6c62cac921f5e5d1526c99268ea3bc6f53fd9`

```
{
  "content_uri": "8cc590c1823ee24dae77eadfc3b6c62cac921f5e5d1526c99268ea3bc6f53fd9",
  "hash": "3cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
  "hash_algorithm": "sha-256",
  "recipients": [
    "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t"
  ],
  "sender": "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t",
  "type": "basic"
}
```






#### Step 2 - Use the API to Send the message 


**Use the tryout**

![Modal](./sharedoc_post.png)

**Corresponding Cli request**


```
curl -X 'POST' \
  'https://dev-api.commercio.app/v1/sharedoc/process' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer .....' \
  -H 'Content-Type: application/json' \
  -d '{
  "back_url": "http://example.com/callback",
  "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
  "hash": "3cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
  "hash_algorithm": "sha-256",
  "metadata": {
    "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
    "schema": {
      "uri": "http://example.com/schema.xml",
      "version": "1.0.0"
    }
  },
  "recipients": [
    "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
  ]
}'
```


**API : Body response**

```
{
  "process_id": "34669051-707f-4230-a960-e0ef8e517e43",
  "sender": "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t",
  "receivers": [
    "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t"
  ],
  "document_id": "b03c6c6e-90e4-49ae-a582-e6a3ff4726a3",
  "doc_hash": "3cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
  "doc_hash_alg": "sha-256",
  "doc_tx_hash": "",
  "doc_metadata": {
    "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
    "schema": {
      "uri": "http://example.com/schema.xml",
      "version": "1.0.0"
    }
  },
  "timestamp": "2021-05-20T08:27:56Z",
  "status": "queued"
}
``` 
Register the  process_id assigned `"process_id": "34669051-707f-4230-a960-e0ef8e517e43"` for future check 


##### Step 3 - Check the process status 

Use the API Get : /sharedoc/process with process_id = `34669051-707f-4230-a960-e0ef8e517e43`

see for more details below in the guide

**API : Body response**

``` 
{
  "process_id": "34669051-707f-4230-a960-e0ef8e517e43",
  "sender": "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t",
  "receivers": [
    "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t"
  ],
  "document_id": "b03c6c6e-90e4-49ae-a582-e6a3ff4726a3",
  "doc_hash": "3cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
  "doc_hash_alg": "sha-256",
  "doc_tx_hash": "78733941DE98F4D39424DD082F3516438E397A236BA28C0BBE2AC3CD3A66E94F",
  "doc_metadata": {
    "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
    "schema": {
      "uri": "http://example.com/schema.xml",
      "version": "1.0.0"
    }
  },
  "timestamp": "2021-05-20T08:27:56Z",
  "status": "processed"
}
```

Acquire the  "doc_tx_hash": "78733941DE98F4D39424DD082F3516438E397A236BA28C0BBE2AC3CD3A66E94F"

##### Step 4  - Check the transaction in the explorer 
 
Use the `doc_tx_hash`  in the explorer filter  

![Modal](./explorer_check_transaction.png)

Check the trasaction

![Modal](./explorer_transaction_doc_tx_hash.png)


##### Step 5  - Check the Message in the store (DB) of the chain 

Through the public API of the chain (LCD), it is also possible to directly check in the chain's database whether the Sharedoc message has been registered.

To perform this verification, use the following path:

https://lcd-testnet.commercio.network/#/Query/CommercionetworkCommercionetworkDocumentsDocument

You only need to know the UUID of the message.

Here's an example URL using the UUID "b03c6c6e-90e4-49ae-a582-e6a3ff4726a3":

https://lcd-testnet.commercio.network/commercionetwork/documents/document/b03c6c6e-90e4-49ae-a582-e6a3ff4726a3


```
{
  "Document": {
    "sender": "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t",
    "recipients": [
      "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t"
    ],
    "UUID": "b03c6c6e-90e4-49ae-a582-e6a3ff4726a3",
    "metadata": {
      "contentURI": "-",
      "schema": {
        "URI": "-",
        "version": "-"
      }
    },
    "contentURI": "8cc590c1823ee24dae77eadfc3b6c62cac921f5e5d1526c99268ea3bc6f53fd9",
    "checksum": {
      "value": "3cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
      "algorithm": "sha-256"
    },
    "encryptionData": null,
    "doSign": null
  }
}
```



#### Common error

The following are common error composing using a  POST Sharedocument message 


##### 1.Hashing Error

Message Example 

```
 {
    "error": "could not validate the ShareDocumentRequest: The hash field must have a length of 32, got instead 64"
}
```

It implies that the hash string indicated in entity `hash` has not a compliant format in respect of hashing algoritm indicated in field `hash_algorithm`

Example 

```
...
"hash": "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
"hash_algorithm": "md5",
... 
```


The hash `2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824` is a sha-256 format NOT Md5


##### 2. Recipients format not correct 

Message Example 

```
{
    "error": "could not build MsgShareDocument: could not derive account address from bech32 addr decoding bech32 failed: invalid bech32 string length 6: string"
}
```

It implies that the value indicated in the entity `recipients` has not a correct format  (Format example : did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr )

Entity  `recipients`  attend an array of wallet address (did) 


Example 

```
...
 "recipients": [
    "string"
  ]
... 
```

The value `string` is not a `did` format 


##### 3. Not enough CCC 

Message Example 

```
{
  "error": "account has only 0uccc, required more 10000uccc"
}
```

It implies that the wallet of the sender has not enough CCC to pay the chain fee for the transaction.   



### Sent Processes

Permit to get all the process of sharedoc sent by the authenticated user 


#### PATH

GET : `/sharedoc/process`


Moreover throught the following  parameters the API  permit to paginate and order the result.

* `limit` : Limit the max number of elements returned
* `next_cursor`:  Cursor that specifies an ID from starting to return following elements
* `order` :  Elements ordering by creation timestamp



#### Step by step Example

Let's create a query to get all messages sent by the sender with `did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr` thati is associated to the authenticated user 

##### Step 1 - Define the first query

Parameter value 

* limit = 30 (is the default value) 
* cursor = empty
* order : asc 

**Use the tryout**


![Modal](./sharedoc_processes.png)


**Corresponding Cli request**

```
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/sharedoc/process?limit=30&order=asc' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer ....'

```

**API : Body response**

Example Value

```
{
  "processes": [
    {
      "process_id": "38367565-ce60-4fb7-96ac-be591b5c65cb",
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "receivers": [
        "did:com:1ffjuspvy8sm8fw7wjyjtgvzg0wgv36pqxrah9n"
      ],
      "document_id": "be3a7460-4935-4434-b045-f0208d55c076",
      "doc_hash": "2f1ec16b9a177aabd5e1ff6bb5685a3df3a6b462dfa147e6b35585fa58954b6b",
      "doc_hash_alg": "sha-256",
      "doc_tx_hash": "390EF4F23974B3CF7663B5F3C8B263F9D0ED1A900167D02ED4760052003CC7F2",
      "tx_timestamp": "2021-06-30T09:49:32Z",
      "tx_type": "commercio/MsgShareDocument",
      "doc_metadata": {
        "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
        "schema": {
          "uri": "http://example.com/schema.xml",
           "version": "1.0.0"
        }
      },
      "chain_id": "commercio-testnet10k2",
      "timestamp": "2021-06-30T09:46:25Z",
      "status": "success"
    },
    {
      "process_id": "295c021f-b14b-4b26-859d-f310cc6a7a73",
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "receivers": [
        "did:com:1lustf0n3t6fr2sp2p07hrf5qzja47juzccz935"
      ],
      "document_id": "d5a84f6c-2540-47b6-8661-70ff99a7fbff",
      "doc_hash": "c2000af9444c2b4b949e86ab00c7521b8ecc8a5b6485dea84442f1e167b6a755",
      "doc_hash_alg": "sha-256",
      "doc_tx_hash": "390EF4F23974B3CF7663B5F3C8B263F9D0ED1A900167D02ED4760052003CC7F2",
      "tx_timestamp": "2021-06-30T09:49:32Z",
      "tx_type": "commercio/MsgShareDocument",
      "doc_metadata": {
        "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
        "schema": {
          "uri": "http://example.com/schema.xml",
           "version": "1.0.0"
        }
      },
      "chain_id": "commercio-testnet10k2",
      "timestamp": "2021-06-30T09:46:27Z",
      "status": "success"
    },
....
    {
      "process_id": "68b833ad-20e9-4887-bc6d-34431d4c2c03",
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "receivers": [
        "did:com:1aeugwtn2tdsqx5uznac5su4k7wscc4atmew04k"
      ],
      "document_id": "49c01045-17ab-4e75-a0bb-f683714d2f08",
      "doc_hash": "a0ed0e4c307bd0a91f5976bb17c444332343716c5ea48f453d623ca8c2d5f4ea",
      "doc_hash_alg": "sha-256",
      "doc_tx_hash": "FD2E1D5DD97E9589673A6BDB2F9A2468F4A856664F843619AF0FDC1D99F6560E",
      "tx_timestamp": "2021-06-30T10:25:06Z",
      "tx_type": "commercio/MsgShareDocument",
      "doc_metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "chain_id": "commercio-testnet10k2",
      "timestamp": "2021-06-30T10:19:44Z",
      "status": "success"
    }
  ],
  "paging": {
    "next_cursor": "MTYyNTA0ODM4NDg0Njk2NjAwMA==",
    "next_link": "https://dev.commercio.app/sharedoc/api/v1/sharedoc/process?limit=30&order=asc&cursor=MTYyNTA0ODM4NDg0Njk2NjAwMA==",
    "total_count": 418
  }
}
```

In order to get the following processes use the value of `next_cursor` ( that is `MTYyNTA0ODM4NDg0Njk2NjAwMA==` in the exmple ) in the parameter `next_cursor`



### Specific sent process details

#### PATH

GET :  `/sharedoc/process{process_id}`

Permit to check the status of a specific process knowing its process_id assigned by the system


#### Step by step Example
Let's create a query to get the details of a specific process 

##### Step 1 - Define the first query

Lets's check the process with `"process_id": "38367565-ce60-4fb7-96ac-be591b5c65cb"`


**Use the tryout**

![Modal](./sharedoc_process_by_process_id.png)


**Corresponding Cli request**

<pre style="color:#FFF;">
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/sharedoc/process/38367565-ce60-4fb7-96ac-be591b5c65cb' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOi.....'
</pre>

**API : Body response**

```
{
  "process_id": "38367565-ce60-4fb7-96ac-be591b5c65cb",
  "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
  "receivers": [
    "did:com:1ffjuspvy8sm8fw7wjyjtgvzg0wgv36pqxrah9n"
  ],
  "document_id": "be3a7460-4935-4434-b045-f0208d55c076",
  "doc_hash": "2f1ec16b9a177aabd5e1ff6bb5685a3df3a6b462dfa147e6b35585fa58954b6b",
  "doc_hash_alg": "sha-256",
  "doc_tx_hash": "390EF4F23974B3CF7663B5F3C8B263F9D0ED1A900167D02ED4760052003CC7F2",
  "tx_timestamp": "2021-06-30T09:49:32Z",
  "tx_type": "commercio/MsgShareDocument",
      "doc_metadata": {
        "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
        "schema": {
          "uri": "http://example.com/schema.xml",
           "version": "1.0.0"
        }
      },
  "chain_id": "commercio-testnet10k2",
  "timestamp": "2021-06-30T09:46:25Z",
  "status": "success"
}

```




## Sent Sharedoc

Permit to get all sharedocs messages sent by the did of the authenticated user.

Alse messages not sent throught an  APi process [Send Sharedoc process](commercioapi-sharedoc.html#send-a-sharedoc)


### PATH

GET :  `/sharedoc/sent`

Moreover throught the following  parameters the API  permit to paginate and order the result.

* `limit` : Limit the max number of elements returned
* `next_cursor`:  Cursor that specifies an ID from starting to return following elements
* `order` :  Elements ordering by creation timestamp




### Step by step Example

Let's create a query to get all messages sent by the sender with `did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr` thati is associated to the authenticated user 

#### Step 1 - Define the first query

Parameter value 

* limit = 3 
* cursor = empty
* order : asc 


#### Step 2 - Use the API

**Use the tryout**


![Modal](./sharedoc_sent.png)


**Corresponding Cli request**

```
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/sharedoc/process?limit=30&order=asc' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer ....'

```



**API : Body response**

```
{
  "documents": [
    {
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "recipients": [
        "did:com:1u35avnkvywzcxp2uty8u0y6xu3s22hycfgd2we"
      ],
      "uuid": "0093638d-841f-4782-8ddb-d9cb020338eb",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "fe64b6b9a51756a244893722917132e85fe0daea99f7cfffb353eab7e1996dcd",
      "checksum": {
        "value": "c3936c163751c60e428774b5b5d8f3bce430aa962c567d4be6f3a33b69e440aa",
        "algorithm": "sha-256"
      }
    },
    {
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "recipients": [
        "did:com:1lustf0n3t6fr2sp2p07hrf5qzja47juzccz935"
      ],
      "uuid": "00af0720-af3b-4140-b785-00d8ff92e460",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "08faa5bcc66f70f2eab0607c516d28ff1774757edf67d2805ab28d520b0c4300",
      "checksum": {
        "value": "d19eca1648d9440daaa5f9e3477c0f5d5fdae68a3935d17d91558c075dd0483a",
        "algorithm": "sha-256"
      }
    },
    {
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "recipients": [
        "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
      ],
      "uuid": "00bbcfc0-a688-44d4-81ed-547a256d40f5",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "ec8736178ed930786e889b5945c59e8126cfb95162263beedc0cc00c264409f8",
      "checksum": {
        "value": "002d545ce3050e75dc4c6cb93ef3e0c61df4c98f51caca644bc10659d9966229",
        "algorithm": "sha-256"
      }
    }
  ],
  "paging": {
    "next_cursor": "MDBiYmNmYzAtYTY4OC00NGQ0LTgxZWQtNTQ3YTI1NmQ0MGY1",
    "next_link": "https://dev.commercio.app/sharedoc/api/v1/sharedoc/sent?limit=3&order=asc&cursor=MDBiYmNmYzAtYTY4OC00NGQ0LTgxZWQtNTQ3YTI1NmQ0MGY1",
    "total_count": 418
  }
}
```

The response contains 
* the details of the first 3 Sharedoc messages 
* the `paging/next_cursor` entity that permit to extract the next page  messages 

#### Step 2 - Extract next page 

Use in the tryout the value 

* limit = 3 
* cursor = `MDBiYmNmYzAtYTY4OC00NGQ0LTgxZWQtNTQ3YTI1NmQ0MGY1`
* order : asc 
  


**API : Body response**

```

{
  "documents": [
    {
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "recipients": [
        "did:com:1e4wh3a2cp20edg7dtkmkrumt9mh4w3x0a4lvjs"
      ],
      "uuid": "00f929e4-44d8-4fd5-8328-2375b67f7357",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "f816e7ad9a72892cbbc70f2d4bd6dfb35cf9e4b0411842ee779f3ad51bdde030",
      "checksum": {
        "value": "10d110a6b1645482572f00af5a7f3bf396e13e37264c861bfa275f0ee7f8b85c",
        "algorithm": "sha-256"
      }
    },
    {
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "recipients": [
        "did:com:1ffjuspvy8sm8fw7wjyjtgvzg0wgv36pqxrah9n"
      ],
      "uuid": "0113ca97-c4e4-4690-b683-3515968600bb",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "b6c44223d34837a00c5947e70b0b84693e9fdbcf391693a7f850b8cc22afc1bf",
      "checksum": {
        "value": "6f1a002bc49f6c4b87878eb314956f04fc283a0300fc131668edc6ea10f10b8c",
        "algorithm": "sha-256"
      }
    },
    {
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "recipients": [
        "did:com:1t934ysywlz29lvjudwd6mr204wysfg34g7pwxs"
      ],
      "uuid": "02e421f6-9839-48c7-ad90-448bdb49d20a",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "923aa06b2cab8e66117c270ddc3c332e4c7fcd295b57b0471db41a40fe829a5c",
      "checksum": {
        "value": "9367b50b323992e3177bd52bcc1a3d6c105a25797aa047c131ba1e92780930cf",
        "algorithm": "sha-256"
      }
    }
  ],
  "paging": {
    "next_cursor": "MDJlNDIxZjYtOTgzOS00OGM3LWFkOTAtNDQ4YmRiNDlkMjBh",
    "next_link": "https://dev.commercio.app/sharedoc/api/v1/sharedoc/sent?limit=3&order=asc&cursor=MDJlNDIxZjYtOTgzOS00OGM3LWFkOTAtNDQ4YmRiNDlkMjBh",
    "total_count": 418
  }
}

``` 


## Received Sharedoc

Permit to get all sharedocs messages received (sent to the did of the authenticated user).

### PATH

GET :

 `/sharedoc/received`



Moreover throught the following  parameters the API  permit to paginate and order the result.

* `limit` : Limit the max number of elements returned
* `next_cursor`:  Cursor that specifies an ID from starting to return following elements
* `order` :  Elements ordering by creation timestamp

### Step by step Example

Let's create a query to get all messages received by the authenticatd user with the did   `did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr` 

#### Step 1 - Define the first query

Parameter value 

* limit = 3 
* cursor = empty
* order : asc 

**Use the tryout**

![Modal](./sharedoc_received.png)


**Corresponding Cli request**

```
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/sharedoc/received?limit=30&order=asc' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSU....'
```

**API : Body response**

```
{
  "documents": [
    {
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "recipients": [
        "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
      ],
      "uuid": "00bbcfc0-a688-44d4-81ed-547a256d40f5",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "ec8736178ed930786e889b5945c59e8126cfb95162263beedc0cc00c264409f8",
      "checksum": {
        "value": "002d545ce3050e75dc4c6cb93ef3e0c61df4c98f51caca644bc10659d9966229",
        "algorithm": "sha-256"
      }
    },
    {
      "sender": "did:com:1va3cl46zcmd9lga3mulvhyd7k5a23jg23fkypt",
      "recipients": [
        "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
      ],
      "uuid": "014203f1-cf7b-42c8-adbf-b898ad088b21",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "cb57fd1164a2bcd070a14e00c82ead963961c7028bc55af7ed492bcec6b92409",
      "checksum": {
        "value": "b56f2de7cde1285a49d0337869a1c0e52b917170df04c619df26a81b8e8d82d4",
        "algorithm": "sha-256"
      }
    },
    ....

    {
      "sender": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr",
      "recipients": [
        "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
      ],
      "uuid": "12ecae4c-bb8b-411c-b516-aba29e186a21",
      "metadata": {
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
        }
      },
      "content_uri": "a537ba5171d65c8c8aecd7971f9a65db93906e7dbb0a8618490ed2c5a5ac19b0",
      "checksum": {
        "value": "2d1203278986af1ac3a0d6e84b5d2cfb4d8cf2ce60dacddf91824b298189ff09",
        "algorithm": "sha-256"
      }
    }
  ],
  "paging": {
    "next_cursor": "MTJlY2FlNGMtYmI4Yi00MTFjLWI1MTYtYWJhMjllMTg2YTIx",
    "next_link": "https://dev.commercio.app/sharedoc/api/v1/sharedoc/received?limit=30&order=asc&cursor=MTJlY2FlNGMtYmI4Yi00MTFjLWI1MTYtYWJhMjllMTg2YTIx",
    "total_count": 397
  }
}

```


## Receipt

This API permit to manage the reading receipt message  `MsgSendDocumentReceipt` throught the DOCS  Module  used when the receivers wants  to acknowledge the sender that he has properly read a  specific `MsgShareDocument`

### Send a Receipt Message process

Permit to create a receipt message relative to a Sharedocument Message received (the did of the authenticated user is a receipient) 

#### PATH

POST : `/receipts/`

#### Step by step Example

Let's create a query to get all messages received by the authenticatd user with the did   `did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr` 

##### Step 1 - Define the first query


Let's create a new process to send a Receipt  relative to a `MsgShareDocument`
where your address is in the set of `recipients` 

In Review - Coming soon 



##### Step 2 - Define the first query 






### Sent Receipts processes


In Review - Coming soon 

### Sent Receipt Message specific process details
In Review - Coming soon 

### Received Receipt Message
In Review - Coming soon



## Use cases examples


### Timestamp a document

### Prove the existence of document

### Notarize a document

### Legally prove a document was shared with a third party



#