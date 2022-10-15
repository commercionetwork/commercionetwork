# CommercioAPI DOCS

The  CommercioAPI DOCS allows you to share a document with another user, and retrieve the list of documents that you have received.


## What is an Electronic Cerified Delivery?

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


## ShareDoc trasaction processes 

See folowing guides for more technical details on  <a href="/x/documents/#sending-a-document">MsgShareDocument</a> using the <a href="/x/documents/#docs">DOCS MODULE</a>


### Send a shareDoc 

Use the API POST : `/sharedoc/process`

Permit to create a process to send a message in the block chain named `MsgShareDocument` or Sharedoc message throught the DOCS  Module



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
    "content_uri": "-",
    "schema": {
      "uri": "-",
      "version": "-"
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
    "content_uri": "-",
    "schema": {
      "uri": "-",
      "version": "-"
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
Use the API GET : `/sharedoc/process`

Permit to get all the process of sharedoc sent by the authenticated user 

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
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
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
        "content_uri": "-",
        "schema": {
          "uri": "-",
          "version": "-"
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



### Sent specific process details
Use the API GET : /sharedoc/process{process_id} 

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
    "content_uri": "-",
    "schema": {
      "uri": "-",
      "version": "-"
    }
  },
  "chain_id": "commercio-testnet10k2",
  "timestamp": "2021-06-30T09:46:25Z",
  "status": "success"
}

```




## Sent Sharedoc
Use the API GET : `/sharedoc/sent`

Permit to get all sharedocs messages sent by the did of the authenticated user. Alse messages not sent
throught an  APi process [Send Sharedoc process](commercioapi-sharedoc.html#send-a-sharedoc)


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
Use the API GET : `/sharedoc/received`


Permit to get all sharedocs messages received (sent to)  the did of the authenticated user.

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

Use the API POST : cooming soon

Permit to create a recipt message relative to a Sharedocument Message recived (the did of the authenticaded user is a receipient) 


### Sent Receipts processes
Coming soon 

### Sent Receipt Message specific process details
Coming soon 

### Received Receipt Message
Coming soon 



