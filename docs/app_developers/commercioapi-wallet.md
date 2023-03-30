# CommercioAPI eID 


The  CommercioAPI eID permit to manage the basic operations on your in commercio chain  



### Wallet with the commercio app 
The commercio app provide You with a hosted wallet 

A wallet address is in the form :

`did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl`

It has a double function in commercio.network 

a) identifiy uniquely your identity as per an ID card code

b)  identifiy uniquely your account as per an IBAN

You could view the public information of your wallet 
the easiest way is to search for it in the explorer 

[https://testnet.commercio.network/account/did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl/](https://testnet.commercio.network/account/did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl/)


The Api `/wallet/` path permit you to interact with your wallet  and obtain some basic informations.



## Get your wallet address and balance
It permit to obtain your wallet address associated to your account in the commercio.app and its balance in terms of token


### Path

 `GET /wallet/address`


### Step by step example
Let's use the API

#### Step 1 - Send the message  

Use the API Get : /wallet/address


Example

```bash
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/wallet/address' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbG....jDr-g'
```


**API : Body response**
THe body response is a json containing  the following imortant entity 


* `address`:  The wallet id address associated to the authorized user the  example did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx

* `coins` : the "content" of the wallet each element defined by 
  * `denom` : the token type (mesure unit)
  * `amount` : the value of token owned

Example 

```json
{
  "account_number": "7659", 
  "address": "did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx", 
  "coins": [
    {
      "amount": "482380000", 
      "denom": "uccc"
    }
  ], 
  "public_key": {
    "type": "tendermint/PubKeySecp256k1", 
    "value": "A6In+vwf8iG3tr4T8wmeQXfTz2Yp1ztXSXyaka/wKd+M"
  }, 
  "sequence": "70"
}


```



## Get Wallet balance
The Api  permit you to obtain the tokens balance value associated with your wallet address 

### Api path 

GET ` /wallet/balance`


### Step by step example
Let's use the API  

#### Step 1 - Send the message  

Use the API GET : `/wallet/balance`


Example

```bash
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/wallet/balance' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbG....72i3g'

```

**API : Body response**

THe body response is a json containing information for each type of token owned  and its numeric amount. 
As per API `/wallet/address` entity response `coins`


  * `denom` : the token type (mesure unit)
  * `amount` : the value of token owned


Example 

```json
{
  "amount": [
    {
      "amount": "482380000",
      "denom": "uccc"
    }
  ]
}
```



## Send token from your wallet
The Api  basically permits you to send tokens to another wallet address from your wallet as per a wire transfer. 

In detail the Api permits to instatiate a process for sending tokens to another wallet in the commercio.app queue. The process request will be executed by the commercio.app in the blockchain. 

The process istantiated could have a positive or negative outcome so it should be checked.


### Api path 

POST  `/wallet/transfers`

### Step by step example
Let's use the API

We try to send
*  Amount = 1 ccc
*  To the following wallet `did:com:1u35avnkvywzcxp2uty8u0y6xu3s22hycfgd2we`


#### Step 1 - Send the message  

Use the API POST : `/wallet/transfers`


Example

```bash
curl -X 'POST' \
  'https://dev-api.commercio.app/v1/wallet/transfers' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSU.....AGguDwg' \
  -H 'Content-Type: application/json' \
  -d '{
  "amount": [
    {
      "amount": "1000000",
      "denom": "uccc"
    }
  ],
  "back_url": "http://example.com/callback",
  "recipient": "did:com:1u35avnkvywzcxp2uty8u0y6xu3s22hycfgd2we"
}'




```

**API : Body response**

The response contains data about the process  

`send_token_id` : is the unique indetifier of the process enqueued 

Is an important data that could be used later with the API Wallet transfers GET by send_token_id. See the next API description


```json
{
  "send_token_id": "af5b3a65-9c60-4241-9186-b655a1091dcc",
  "sender": "did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx",
  "receiver": "did:com:1u35avnkvywzcxp2uty8u0y6xu3s22hycfgd2we",
  "tx_type": "cosmos-sdk/MsgSend",
  "amount": [
    {
      "amount": "1000000",
      "denom": "uccc"
    }
  ],
  "created_at": "2021-06-22T10:20:28Z",
  "status": "enqueued",
  "back_url": "http://example.com/callback"
}
```

--- 

## Check sent token process 
The Api  permit you obtaini details on the process generated with a specific `send_token_id`


### Api path 

GET  `/wallet/transfers/{send_token_id}`

### Step by step example
Let's use the API


#### Step 1 - Send the message  

Use the API GET : `/wallet/transfers`

send_token_id = `af5b3a65-9c60-4241-9186-b655a1091dcc`


Example


```bash
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/wallet/transfers/af5b3a65-9c60-4241-9186-b655a1091dcc' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cC....GguDwg'
```


**API : Body response**


```json
{
  "send_token_id": "af5b3a65-9c60-4241-9186-b655a1091dcc",
  "sender": "did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx",
  "receiver": "did:com:1u35avnkvywzcxp2uty8u0y6xu3s22hycfgd2we",
  "tx_hash": "B6F551947DF246C1E35A8DBAC01DFD8371F0AEF286AEE1B37386035E6FB382C2",
  "tx_timestamp": "2021-06-22T10:20:49Z",
  "tx_type": "cosmos-sdk/MsgSend",
  "amount": [
    {
      "amount": "1000000",
      "denom": "uccc"
    }
  ],
  "chain_id": "commercio-testnet10k2",
  "created_at": "2021-06-22T10:20:28Z",
  "status": "success",
  "back_url": "http://example.com/callback"
}
```

Check in the explorer searching for `transaction tx_hash` value 

Example 

B6F551947DF246C1E35A8DBAC01DFD8371F0AEF286AEE1B37386035E6FB382C2


Direct Link 

[https://testnet.commercio.network/transactions/detail/B6F551947DF246C1E35A8DBAC01DFD8371F0AEF286AEE1B37386035E6FB382C2](https://testnet.commercio.network/transactions/detail/B6F551947DF246C1E35A8DBAC01DFD8371F0AEF286AEE1B37386035E6FB382C2)






## Sent token process list
The Api  permit you obtaini details on all the sending token processes istantiated by  your wallet (of the authenticated user).


### Api path 

GET  `/wallet/transfers/`

### Step by step example
Let's use the API 


#### Step 1 - Send the message  

Use the API GET : `/wallet/transfers`

Example

```bash
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/wallet/transfers' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR......ZDmig'

```


**API : Body response**


```json
[
  {
    "send_token_id": "af5b3a65-9c60-4241-9186-b655a1091dcc",
    "sender": "did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx",
    "receiver": "did:com:1u35avnkvywzcxp2uty8u0y6xu3s22hycfgd2we",
    "tx_hash": "B6F551947DF246C1E35A8DBAC01DFD8371F0AEF286AEE1B37386035E6FB382C2",
    "tx_timestamp": "2021-06-22T10:20:49Z",
    "tx_type": "cosmos-sdk/MsgSend",
    "amount": [
      {
        "amount": "1000000",
        "denom": "uccc"
      }
    ],
    "chain_id": "commercio-testnet10k2",
    "created_at": "2021-06-22T10:20:28Z",
    "status": "success",
    "back_url": "http://example.com/callback"
  },
  {
    "send_token_id": "8db5ad6d-3ee5-4d9b-bc10-d2c38b5187e8",
    "sender": "did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx",
    "receiver": "did:com:1j930xl8kr92wrxpmur0e5p8vlmy2ce6zg87w3t",
    "tx_hash": "1F9930A3FE5C7884ABAB05B88DF96982770698067F04020D59A403784444D702",
    "tx_timestamp": "2021-06-22T10:48:18Z",
    "tx_type": "cosmos-sdk/MsgSend",
    "amount": [
      {
        "amount": "5000000",
        "denom": "uccc"
      }
    ],
    "chain_id": "commercio-testnet10k2",
    "created_at": "2021-06-22T10:47:58Z",
    "status": "success",
    "back_url": "http://example.com/callback"
  },
  {
    "send_token_id": "e057944e-c163-4489-9edb-9154393613c2",
    "sender": "did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx",
    "receiver": "did:com:1wmlxvvglrmvpsdw3ju45apwjgha6al87h000e7",
    "tx_type": "cosmos-sdk/MsgSend",
    "amount": [
      {
        "amount": "120000000",
        "denom": "uccc"
      }
    ],
    "created_at": "2021-06-22T10:48:32Z",
    "status": "processing",
    "back_url": "http://example.com/callback"
  }
]

```

Obviously You can always look at all the transactions associated to your wallet directly throught the explorer 

Direct Link using your did address 

[https://testnet.commercio.network/account/did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx/](https://testnet.commercio.network/account/did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx)



## Manage your DDO

A DID Document is a digital document that describes a Decentralized Identifier (DID) and contains information about the DID's associated public keys, authentication mechanisms, and service endpoints.

A Decentralized Identifier (DID) is a unique identifier that is not dependent on any centralized authority or registry. It is designed to provide a decentralized, secure, and privacy-preserving way to identify and interact with people, organizations, and things on the internet.

A DID Document provides a standardized way for different systems to interact with each other using DIDs, making it easier to build decentralized applications and services that rely on the secure and private exchange of information. The information contained in a DID Document can be used to verify the authenticity and integrity of a DID and to establish secure communication channels between different parties.

For details on DID documents refers to [DID module](/modules/did/)

Through API is possible to interact with the DID Document of a user present in the chain knowing the eID ((Wallet address) or to manage your DID Document updating it  associated to your eID (Wallet address)

Mainly a DDO contains public keys of the user that can be used 
by applications for many scope

There are two type of keys that are directly managed by the hosted wallet 
* `RsaVerificationKey2018`
* `RsaSignature2018`

Public version are created in DDO at first generation

This type of kes are intended for encrypting method described [here](/modules/documents/03_messages.html#encryption-data-field-requirements)
not still implemented in commercio.app


### Get DDO of a user 
In Review - cooming soon

#### Api path 

GET ` /ddo/{wallet_address}`


#### Step by step example
Let's use the API  




### Get DDO history of a user 

In Review - cooming soon



### Update your DDO 

#### Api path 

POST ` /ddo/process`


##### Step by step example
Let's create a new process to create the first version of DDO   

**Step 1 - Define the first query**


**Step 2 - Use the API to istantiate the process**
Use the tryout


Corresponding Cli request


API : Body response


S




### Get your DDO Updating  process status


Step 3 - Check the process status
Use the API Get : /sharedoc/process with process_id = 34669051-707f-4230-a960-e0ef8e517e43

see for more details below in the guide

API : Body response


Acquire the "doc_tx_hash": "78733941DE98F4D39424DD082F3516438E397A236BA28C0BBE2AC3CD3A66E94F"

#Step 4 - Check the transaction in the explorer
Use the doc_tx_hash in the explorer filter

Modal

Check the trasaction

Modal

#Common error
The following are common error composing using a POST Sharedocument message

#