# CommercioAPI Authentication

To access all CommercioAPI web services


## Get the APIkey

### Prerequisites 

* Create a membership on [commercio.app](https://commercio.app)
* Purchase a membership (Bronze,Silver,Gold) and get your own wallet
* Have some CCC (Commerce Cash Credits) to pay the transaction fees

After you have created your account you can use the API with our Authentication method

To use the API you need to get the APIkey generated for your registered user in commercio.app.



### Ask the Identity Manager for your APIkey

To get the APIkey you need to interact with our Identity Manager (IDM) 

You need: 
 
<EMAIL> and <PASSWORD> are those of the user you registered in commerce.app


<EMAIL>: testuser001@commercio.app
<PASSWORD>: Testuser001


Below is an example via command line of how to obtain the bearer ($ID_TOKEN) to be used later in the interactions with the API (For the tryout in the Swagger use the Authorize function key to insert the Token obtained):

```bash
curl -s --request POST \
    'https://devlogin.commercio.app/auth/realms/commercio/protocol/openid-connect/token' \
    --header 'Content-Type: application/x-www-form-urlencoded'  \
    --data-urlencode 'client_id=dev.commercio.app'  \
    --data-urlencode 'grant_type=password'  \
    --data-urlencode 'scope=openid'  \
    --data-urlencode 'username='  \
    --data-urlencode 'password=' | jq -r '.id_token'
```


### Identity Manager (IDM) reply:

```
eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjEwMDMwMjEsImlhdCI6MTYyMTAwMjcyMSwiYXV0aF90aW1lIjowLCJqdGkiOiJmNTA5YjQ0YS0xYzIxLTQ5NjktYjE5Ni03YWYxOGFmZDkyYTciLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiJhMmIzZGI5Yi03NzUwLTQzYTEtODExZC1iOGI3MjA2NmQzZDYiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6ImE5ZGNmMWFjLTdjMTctNDViYS1hY2JlLWZkMmY1MGNhZGEzMyIsImF0X2hhc2giOiJLZko4XzJfWGxCQmFFNjVBYVhOWWRnIiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMDk5NDk2MCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6Ik1hcmNvIEF1cm8iLCJwaG9uZV9udW1iZXIiOiIxMjM0NTY3ODkwMSIsInByZWZlcnJlZF91c2VybmFtZSI6Im1hcmNvLnJ1YXJvQGdtYWlsLmNvbSIsImdpdmVuX25hbWUiOiJNYXJjbyIsImZhbWlseV9uYW1lIjoiQXVybyIsImVtYWlsIjoibWFyY28ucnVhcm9AZ21haWwuY29tIiwidXNlcm5hbWUiOiJtYXJjby5ydWFyb0BnbWFpbC5jb20ifQ.hDParV3scvir8B9kkNN-e56IF5Jmqxuhkfd7B__s8Vn41VAaccJBTl1bwqLggcrNJ2Yjl3jAKOxfXX3PFf_RtsFloFyYSZDlOdt73qD1m-8TzdPGfMjNwgiCLc7IvKIFV3_8JYsgkm3fsqtMGqOdsqZSD_s9KrGK7oYcoMIWHqiBKqeymAX9urLFg4lbHlEY1rJJ6C0zpFhA1nrqSFqwu3MuYdfylmtkhvKVreOl9jR8kG326BvwEd7NnwaYtJI6Anoe2ojNHzWgRwFTzd3djhwhYLziJTt3Q8SE7ag_FKxQ4BhjaK3w4PlBz9HK15B4rp_shd_ZUohVaZtJsNrKwg
```

This APIKey is needed to authenticate and interact with the Swagger endpoint (CommercioAPI)


## Interacting with CommercioAPI


### End points

Each API has its on Endpoint. For example the Share API endpoint is available at this address  

* https://dev-api.commercio.app/v1/sharedoc (Testnet)
* https://api.commercio.app/v1/sharedoc (mainNet)

Testnet is for testing purpose.

Main net is the real Blockchain. 



## Securing your App

Coming next

