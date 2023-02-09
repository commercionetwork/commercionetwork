# CommercioAPI Wallet 


The  CommercioAPI Wallet permit to manage the basic operations on your Wallet/s




### Wallet with the commercio app 
The commercio app provide You with a hosted wallet 

A wallet address is in the form :

`did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl`

It has a double function in commercio.network 

a) identifiy uniquely your account as per an IBAN

b) identifiy uniquely your identity as per an ID card code

You could view the public information of your wallet 
the easiest way is to search for it in the explorer 

[https://testnet.commercio.network/account/did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl/](https://testnet.commercio.network/account/did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl/)


The Api `/wallet/` path permit you to interact with your wallet  and obtain some basic informations.



## Get your address and balance
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
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjQzNTQxNjcsImlhdCI6MTYyNDM1MzI2NywiYXV0aF90aW1lIjowLCJqdGkiOiI2YTk4ZjIyZi02ZTNkLTQ4MzQtYmMwYy03MzhmZTI1ZWM1Y2MiLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiI0OWFhZjQ3OS1hMjE4LTRhZjItOWY3MS1kMTI2OThmNjk5YjkiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6ImVjYTg3ZWIwLWZmYWItNGMzMi05YzFlLWQ0MDE3ZTE4YmRhZSIsImF0X2hhc2giOiJieXFLRE5WLXowVUlfOGFHRVp6bkV3IiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMTUyNDY0NyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIEVudGVycHJpc2V1c2VyMDAzIiwicGhvbmVfbnVtYmVyIjoiMzQ4NTI0MTY0OSIsInByZWZlcnJlZF91c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIiwiZmFtaWx5X25hbWUiOiJFbnRlcnByaXNldXNlcjAwMyIsImVtYWlsIjoiZW50ZXJwcmlzZXVzZXIwMDNAem90c2VsbC5jb20iLCJ1c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIn0.VSPIg8fefKfP2rShWVTvap1IvGx-A64FPYUsV19lHF4KFl8oZq7AyP6ePXTbkI2G1ifiA7a8mVr7g3O8b8MofRbUHSxrzPLSh_eSSQYP618f4G1sTYPGOIZuRjzTX_liywryejvEXGBzt50E-KFpGwUA99CTyG2q8s1Z-gBpDbzTY5Wd7Kc_1GkbYsTKSx1hs1-4OiCCFJ8cTRkYgVq01JqdX-Ghf8KF9yrpvORIPvvKBo9ZjoqszVSJFOgm51Zp0NuxL3Vb9FsLIIuEjlR4ocdLNXJ6qeFKa2xUWKtxwdFL-sredJgiQyt-tixcGFtVKpivVNV7KoMuSlgikjDr-g'
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
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjQzNTQ0MzcsImlhdCI6MTYyNDM1MzUzNywiYXV0aF90aW1lIjowLCJqdGkiOiJhZDdmNzBhMi1iMDI2LTRhNmYtOTgzMC1iMjA4MmVlMzQwMDEiLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiI0OWFhZjQ3OS1hMjE4LTRhZjItOWY3MS1kMTI2OThmNjk5YjkiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6IjdmZGNiY2Y1LTRkMGYtNGUwYi1hOTRiLTc1ODQ5ZmMxNTk2YiIsImF0X2hhc2giOiJCVU9Ndjg3djNEUDVmZ0xHd1FuVjhBIiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMTUyNDY0NyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIEVudGVycHJpc2V1c2VyMDAzIiwicGhvbmVfbnVtYmVyIjoiMzQ4NTI0MTY0OSIsInByZWZlcnJlZF91c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIiwiZmFtaWx5X25hbWUiOiJFbnRlcnByaXNldXNlcjAwMyIsImVtYWlsIjoiZW50ZXJwcmlzZXVzZXIwMDNAem90c2VsbC5jb20iLCJ1c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIn0.i3NbUWvwB4NNLfPSmUrfziwA4lJOzbsT6J0Ngc8QfEHfZd7R2U_GlQSv_e94v-Hac-97bBGUHhdeqZCreieW2wc_6Gbwyi3CvbglBRcNWGNbbtX78aU0K5gOLBR0_KfJxMxZuZe4AcWKjdQ3urq85-A-_AGoq8OWvcGkzAzA1Pi8UX4q30imTaW3m-N2cvK9fAxSLCnf5c9XPKDMaHWF-ACi30_GM4Yrubzev8I7Dg6Jaf24jqZKBzKOL0MmOk2Iw7SuR2XoqaiUUkKk7iuI0fnrIhUDaGy88bXj9pwoQtrtw9_kPXIQSp3pXvsRjCfqoOGMVpks7sFVNh6oc72i3g'

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
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjQzNTgwMDEsImlhdCI6MTYyNDM1NzEwMSwiYXV0aF90aW1lIjowLCJqdGkiOiIxNTA5M2IwMS05MmYxLTRjMjctYWFlNi00ZGZlM2M3MGM0NDciLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiI0OWFhZjQ3OS1hMjE4LTRhZjItOWY3MS1kMTI2OThmNjk5YjkiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6IjkyMDkxMzdkLTYxYzYtNDkwZS1iMjRlLTIyOTU0OGZkNGI5OSIsImF0X2hhc2giOiJYY1hZLWxFYjNTSERRbmg0N1dmUDZBIiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMTUyNDY0NyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIEVudGVycHJpc2V1c2VyMDAzIiwicGhvbmVfbnVtYmVyIjoiMzQ4NTI0MTY0OSIsInByZWZlcnJlZF91c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIiwiZmFtaWx5X25hbWUiOiJFbnRlcnByaXNldXNlcjAwMyIsImVtYWlsIjoiZW50ZXJwcmlzZXVzZXIwMDNAem90c2VsbC5jb20iLCJ1c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIn0.GXcDC-HviuylBqryyeFQnR1g_sMZG70utlKr6OVUEkoJ4ysQGfLVMuptZUlkxEebHSspGLtB2vTPtuMaY6D7jN3AKLJbW0ceRTg1u1lfWbWJrqG7Ly2zKlklvcDK-VBcW38OLqi3JjJkQYgLJ6P_YuWqH6K9N0Jz6CHHNP1iGPM6T4Yx-AIihfVfy85xtbG4NnHHKm25FElh-PTUUCTXatsP8CTwWrA2CVPfKNoSttJJ3GYJSc7hq-Qf7pv8g7NTe2PWaZdeaOhGQwLKZSfrIN_Pxu1FeHjfTg0jyENSiBxPhJDdlWDnk1jZforMiXXu9WD294z-7E6vULDAGguDwg' \
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
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjQzNTgwMDEsImlhdCI6MTYyNDM1NzEwMSwiYXV0aF90aW1lIjowLCJqdGkiOiIxNTA5M2IwMS05MmYxLTRjMjctYWFlNi00ZGZlM2M3MGM0NDciLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiI0OWFhZjQ3OS1hMjE4LTRhZjItOWY3MS1kMTI2OThmNjk5YjkiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6IjkyMDkxMzdkLTYxYzYtNDkwZS1iMjRlLTIyOTU0OGZkNGI5OSIsImF0X2hhc2giOiJYY1hZLWxFYjNTSERRbmg0N1dmUDZBIiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMTUyNDY0NyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIEVudGVycHJpc2V1c2VyMDAzIiwicGhvbmVfbnVtYmVyIjoiMzQ4NTI0MTY0OSIsInByZWZlcnJlZF91c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIiwiZmFtaWx5X25hbWUiOiJFbnRlcnByaXNldXNlcjAwMyIsImVtYWlsIjoiZW50ZXJwcmlzZXVzZXIwMDNAem90c2VsbC5jb20iLCJ1c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIn0.GXcDC-HviuylBqryyeFQnR1g_sMZG70utlKr6OVUEkoJ4ysQGfLVMuptZUlkxEebHSspGLtB2vTPtuMaY6D7jN3AKLJbW0ceRTg1u1lfWbWJrqG7Ly2zKlklvcDK-VBcW38OLqi3JjJkQYgLJ6P_YuWqH6K9N0Jz6CHHNP1iGPM6T4Yx-AIihfVfy85xtbG4NnHHKm25FElh-PTUUCTXatsP8CTwWrA2CVPfKNoSttJJ3GYJSc7hq-Qf7pv8g7NTe2PWaZdeaOhGQwLKZSfrIN_Pxu1FeHjfTg0jyENSiBxPhJDdlWDnk1jZforMiXXu9WD294z-7E6vULDAGguDwg'
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

[https://testnet.commercio.network/tx/B6F551947DF246C1E35A8DBAC01DFD8371F0AEF286AEE1B37386035E6FB382C2](https://testnet.commercio.network/tx/B6F551947DF246C1E35A8DBAC01DFD8371F0AEF286AEE1B37386035E6FB382C2)






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
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjQzNTk2MzYsImlhdCI6MTYyNDM1ODczNiwiYXV0aF90aW1lIjowLCJqdGkiOiJjZjYxNDhkMi0yNWJkLTQyY2MtYWIyYy1iZGM4ZDFlN2EzY2IiLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiI0OWFhZjQ3OS1hMjE4LTRhZjItOWY3MS1kMTI2OThmNjk5YjkiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6ImM5NjVkOTc0LTZkMGYtNDU3Mi05NTNmLTJiZjM4NWFiNjBiYSIsImF0X2hhc2giOiI0eGVoS3RKTkZVMzlPSy1NRE9XMTF3IiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMTUyNDY0NyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIEVudGVycHJpc2V1c2VyMDAzIiwicGhvbmVfbnVtYmVyIjoiMzQ4NTI0MTY0OSIsInByZWZlcnJlZF91c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVudGVycHJpc2V1c2VyMDAzIiwiZmFtaWx5X25hbWUiOiJFbnRlcnByaXNldXNlcjAwMyIsImVtYWlsIjoiZW50ZXJwcmlzZXVzZXIwMDNAem90c2VsbC5jb20iLCJ1c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAzQHpvdHNlbGwuY29tIn0.B-gllJgvQhO2aEKv38eviYihsb0R2TCSTVQL5K2LTpKVnARye3pIBo7vJ-DeOjs48e0y_y0usD0_I-XocZxIGNupsRFcK46nKmiwJg289QSu-b4J0o2sDoDe3OpjUxBhtuZgO9zkfbTDcI1F3DSfAD9fhZV-LddNHlzOx7nbShmZ7mY0voR4d4xwMt-1QpE1Y_H43UXfFlvdFLeCr7mv8HI2yAFPGJ2B8BY6_eqXj-PIZzbgGzNgfe4JwjsYqYXuhgnbgp52TlNccIfWGHw1rCoB5doUzusdJ2K3gg-5EWPAapMZcwdfbsaRKx936mL9LxeGq5-iZkFJDyvm7ZDmig'

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

[https://testnet.commercio.network/validators/account/did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx](https://testnet.commercio.network/validators/account/did:com:1cjatcdv2uf20803mt2c5mwdrj87tjnuvk3rvsx)



