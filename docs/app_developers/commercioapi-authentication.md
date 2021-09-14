# Authentication

To use all CommercioAPI web services You need to authenticate. The Authenticatin  method used is a `Bearer authentication`. 

 You can obtain  the security token  through an interaction with the IDM using the credential of your registered user in the Web app. As per the <a href="https://openid.net/specs/openid-connect-core-1_0.html" target="_blank">OpenID Connect protocol</a> 

Check the <a href="/app_developers/commercioapi-introduction.html#prerequisite">Prerequistes</a> in order to  perform correctly the process.

## Get the ID token
The ID token resembles the concept of an identity card, in a standard JWT formatThe ID token statements, or claims, are packaged in a simple JSON object

The ID token header, claims JSON and signature are encoded into a base 64 URL-safe string, for easy passing arround, for example as URL parameter.

You can read more about the JWT data structure and its encoding in <a href="https://datatracker.ietf.org/doc/html/rfc7519" target="_blank">RFC 7519</a> 

The endpoint to interact with the IDM  has the following path

```
https://{{.commercio_login_url}}/auth/realms/commercio/protocol/openid-connect/token
```


The process can be performed via CLI for obtaining the ID Token from the IDM .

```bash
curl -s --request POST \
    'https://{{.commercio_login_url}}/auth/realms/commercio/protocol/openid-connect/token' \
    --header 'Content-Type: application/x-www-form-urlencoded'  \
    --data-urlencode 'client_id={{.openid_client_id}}'  \
    --data-urlencode 'grant_type=password'  \
    --data-urlencode 'scope=openid'  \
    --data-urlencode 'username=<EMAIL>'  \
    --data-urlencode 'password=<PASSWORD>' | jq -r '.id_token'

```

Where `<EMAIL>` and `<PASSWORD>` are those of the user you registered in Web app


The `id_token` obtaneined must be used to autheticate
using the API 


### Example 
Suppose to have the user 

* `<EMAIL>`: testuser001@commercio.app
* `<PASSWORD>`: Testuser001


**Acquire the ID_Token**

```
curl -s --request POST \
    'https://devlogin.commercio.app/auth/realms/commercio/protocol/openid-connect/token' \
    --header 'Content-Type: application/x-www-form-urlencoded'  \
    --data-urlencode 'client_id=dev.commercio.app'  \
    --data-urlencode 'grant_type=password'  \
    --data-urlencode 'scope=openid'  \
    --data-urlencode 'username=testuser001@commercio.app'  \
    --data-urlencode 'password=Testuser001' | jq -r '.id_token'
```

**Acquire `Bearer ID_Token`**

Simple way to compose `Bearer` `ID_token` string through curl


```

 echo "Bearer "$(curl -s --request POST  'https://devlogin.commercio.app/auth/realms/commercio/protocol/openid-connect/token'  --header 'Content-Type: application/x-www-form-urlencoded'   --header 'Cookie: KEYCLOAK_LOCALE=en'  --data-urlencode 'client_id=dev.commercio.app'   --data-urlencode 'grant_type=password'   --data-urlencode 'scope=openid'   --data-urlencode 'username=testuser001@commercio.app'   --data-urlencode 'password=Testuser001' | jq -r '.id_token')

```

**Identity Manager (IDM) reply**


```
eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjEwMDMwMjEsImlhdCI6MTYyMTAwMjcyMSwiYXV0aF90aW1lIjowLCJqdGkiOiJmNTA5YjQ0YS0xYzIxLTQ5NjktYjE5Ni03YWYxOGFmZDkyYTciLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiJhMmIzZGI5Yi03NzUwLTQzYTEtODExZC1iOGI3MjA2NmQzZDYiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6ImE5ZGNmMWFjLTdjMTctNDViYS1hY2JlLWZkMmY1MGNhZGEzMyIsImF0X2hhc2giOiJLZko4XzJfWGxCQmFFNjVBYVhOWWRnIiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMDk5NDk2MCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6Ik1hcmNvIEF1cm8iLCJwaG9uZV9udW1iZXIiOiIxMjM0NTY3ODkwMSIsInByZWZlcnJlZF91c2VybmFtZSI6Im1hcmNvLnJ1YXJvQGdtYWlsLmNvbSIsImdpdmVuX25hbWUiOiJNYXJjbyIsImZhbWlseV9uYW1lIjoiQXVybyIsImVtYWlsIjoibWFyY28ucnVhcm9AZ21haWwuY29tIiwidXNlcm5hbWUiOiJtYXJjby5ydWFyb0BnbWFpbC5jb20ifQ.hDParV3scvir8B9kkNN-e56IF5Jmqxuhkfd7B__s8Vn41VAaccJBTl1bwqLggcrNJ2Yjl3jAKOxfXX3PFf_RtsFloFyYSZDlOdt73qD1m-8TzdPGfMjNwgiCLc7IvKIFV3_8JYsgkm3fsqtMGqOdsqZSD_s9KrGK7oYcoMIWHqiBKqeymAX9urLFg4lbHlEY1rJJ6C0zpFhA1nrqSFqwu3MuYdfylmtkhvKVreOl9jR8kG326BvwEd7NnwaYtJI6Anoe2ojNHzWgRwFTzd3djhwhYLziJTt3Q8SE7ag_FKxQ4BhjaK3w4PlBz9HK15B4rp_shd_ZUohVaZtJsNrKwg
```

*You can decode the Id_token here <a href="jwt.io" target="_blank">jwt.io</a>*



For the `Tryout`  in the Swagger (available at the  CommercioAPI base url)  use  in the modal associated to the `Authorize` button composing the two element separated by a space   

* Method : `Beared` 
* id_token obtained.

Example : 

```
Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjEwMDMwMjEsImlhdCI6MTYyMTAwMjcyMSwiYXV0aF90aW1lIjowLCJqdGkiOiJmNTA5YjQ0YS0xYzIxLTQ5NjktYjE5Ni03YWYxOGFmZDkyYTciLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiJhMmIzZGI5Yi03NzUwLTQzYTEtODExZC1iOGI3MjA2NmQzZDYiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6ImE5ZGNmMWFjLTdjMTctNDViYS1hY2JlLWZkMmY1MGNhZGEzMyIsImF0X2hhc2giOiJLZko4XzJfWGxCQmFFNjVBYVhOWWRnIiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMDk5NDk2MCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6Ik1hcmNvIEF1cm8iLCJwaG9uZV9udW1iZXIiOiIxMjM0NTY3ODkwMSIsInByZWZlcnJlZF91c2VybmFtZSI6Im1hcmNvLnJ1YXJvQGdtYWlsLmNvbSIsImdpdmVuX25hbWUiOiJNYXJjbyIsImZhbWlseV9uYW1lIjoiQXVybyIsImVtYWlsIjoibWFyY28ucnVhcm9AZ21haWwuY29tIiwidXNlcm5hbWUiOiJtYXJjby5ydWFyb0BnbWFpbC5jb20ifQ.hDParV3scvir8B9kkNN-e56IF5Jmqxuhkfd7B__s8Vn41VAaccJBTl1bwqLggcrNJ2Yjl3jAKOxfXX3PFf_RtsFloFyYSZDlOdt73qD1m-8TzdPGfMjNwgiCLc7IvKIFV3_8JYsgkm3fsqtMGqOdsqZSD_s9KrGK7oYcoMIWHqiBKqeymAX9urLFg4lbHlEY1rJJ6C0zpFhA1nrqSFqwu3MuYdfylmtkhvKVreOl9jR8kG326BvwEd7NnwaYtJI6Anoe2ojNHzWgRwFTzd3djhwhYLziJTt3Q8SE7ag_FKxQ4BhjaK3w4PlBz9HK15B4rp_shd_ZUohVaZtJsNrKwg
```

<img src="./img/authorise_button.png">  


![Modal](./authorize_modal.png)



### Direct usage  in the api endpoint 

Example path /sharedoc/process


```
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/sharedoc/process' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjE1MjE0MzYsImlhdCI6MTYyMTUyMTEzNiwiYXV0aF90aW1lIjowLCJqdGkiOiI1NjYzNjQ3Yi1jYjFhLTQxNDUtYWI5OS0xYjM0MGRhMmE4ZTciLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiI1MzAwOTI2Ny00YzNkLTQ3ZTktODg2MC1mZTgyNmEwM2Y4MGMiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6IjNiOWJkZDg3LWRjYmYtNDlhZC04ZTg5LWIyNDJhZDIxZjI1NiIsImF0X2hhc2giOiI1WHBPWXVaS3RLZTNPUi1KczdBa2lnIiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMTAwOTUwMyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6Ik1hcmNvIFJ1YXJvIiwicGhvbmVfbnVtYmVyIjoiMTIzNDU2Nzg5MDEiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJtYXJjby5ydWFyb0BnbWFpbC5jb20iLCJnaXZlbl9uYW1lIjoiTWFyY28iLCJmYW1pbHlfbmFtZSI6IlJ1YXJvIiwiZW1haWwiOiJtYXJjby5ydWFyb0BnbWFpbC5jb20iLCJ1c2VybmFtZSI6Im1hcmNvLnJ1YXJvQGdtYWlsLmNvbSJ9.p9cJYeRDCqPiQLWKV3JQEYoLTvWm7Phbsv_61umM5HbZN052ZDHa_WcF-HibhFkagphQRoXur7w2UK6UVpRzsRygViyOT8AeSQrJS0_H-ySluZxn-vfnwxsEVuew0mx7iQsYY7mXmVX4pGYTdjZ43cUjo8kMd2_-CjqJlvn3B2H_JJwmjjBOSE8jF5i92xmEX1oieeIpNc1rQkdggPwh9bpK43S4dKlm1okrxQCrADMNoLCDJSi8_AYoAYMUJhzkn8hgrbC2LMQmTavQHSwzajcmikrk16ricTFHSe3NPn_a7is3g0fYrH5gf3qq3bajLVHeTT8_8dbaKknQBC6P6A'
``` 


### Expiration time 

Pay attention that the `id_token` has an expiration time 

```
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjE0OTU1NjcsImlhdCI6MTYyMTQ5NTI2NywianRpIjoiODE4NzdlZjMtNDQ3Yy00N2JlLTlhYjYtN2E2ZjBkZDhlMjEwIiwiaXNzIjoiaHR0cHM6Ly9kZXZsb2dpbi5jb21tZXJjaW8uYXBwL2F1dGgvcmVhbG1zL2NvbW1lcmNpbyIsImF1ZCI6WyJkZXYuY29tbWVyY2lvLmFwcCIsImFjY291bnQiXSwic3ViIjoiNDhhMDRiMjctZDY3OS00NDJiLWJlNmItMDJhNmQ4MmYxMGIzIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzZXNzaW9uX3N0YXRlIjoiNTQ0ZWZiMDctNGE2ZC00N2RkLWEwYjYtYTg0ZWJkZDg1ZGU2IiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwczovL2Rldi5jb21tZXJjaW8uYXBwIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMTQzMjUzOCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6IkVudGVycHJpc2V1c2VyMDAyIEVudGVycHJpc2V1c2VyMDAyIiwicGhvbmVfbnVtYmVyIjoiMzQ4NTI0MTY0OSIsInByZWZlcnJlZF91c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAyQHpvdHNlbGwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVudGVycHJpc2V1c2VyMDAyIiwiZmFtaWx5X25hbWUiOiJFbnRlcnByaXNldXNlcjAwMiIsImVtYWlsIjoiZW50ZXJwcmlzZXVzZXIwMDJAem90c2VsbC5jb20iLCJ1c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAyQHpvdHNlbGwuY29tIn0.Wkho3o1Ef535BdGnyeQOiVQsOYKDvlNQGYCY_cIJHg23Ep2kr9gSM7KhSWvXz9o4Z7hMgCdX2jT7T2JL_6KLrv1sGeu8XYYbG1AbQDtoL5ZsKGtbSKl8yiL8QZZ5my-6lSPHmbg-xF8zePJYe3xhYR1evNaoO8WnDmTlyVyGsBoIu7Y2cBKVQtSyivD20XJx6V1ijp-Nr88wJTQFZYq4MSQS4IVdOeXTUbcVq3Ebc53tmOcfKg10OdZLKC2JoZQzh8Igomup-PaVB8MZUIv54Yxwg8nC45VNrgq6gVF32hJhlWVGf-LlhnP1Vmqjv7gSali6FsHg8sOBbxHbz99cew", 
  "expires_in": 300, 
  "id_token": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJwSnpWTkVBa1JieGJvazJGajZPenlmR3RNR25IRVhYNjA4bEVDOXJyNTlRIn0.eyJleHAiOjE2MjE0OTU1NjcsImlhdCI6MTYyMTQ5NTI2NywiYXV0aF90aW1lIjowLCJqdGkiOiI3NGMzNTE1OS1lOTQyLTQ4OWUtOWNlNS0wMDY2NmQ0NzQyNjIiLCJpc3MiOiJodHRwczovL2RldmxvZ2luLmNvbW1lcmNpby5hcHAvYXV0aC9yZWFsbXMvY29tbWVyY2lvIiwiYXVkIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzdWIiOiI0OGEwNGIyNy1kNjc5LTQ0MmItYmU2Yi0wMmE2ZDgyZjEwYjMiLCJ0eXAiOiJJRCIsImF6cCI6ImRldi5jb21tZXJjaW8uYXBwIiwic2Vzc2lvbl9zdGF0ZSI6IjU0NGVmYjA3LTRhNmQtNDdkZC1hMGI2LWE4NGViZGQ4NWRlNiIsImF0X2hhc2giOiI1SkEyYXRjYWIyaUt5WnE3bkg3ZjVRIiwiYWNyIjoiMSIsInRlcm1zX2FuZF9jb25kaXRpb25zIjoiMTYyMTQzMjUzOCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJhZGRyZXNzIjp7fSwibmFtZSI6IkVudGVycHJpc2V1c2VyMDAyIEVudGVycHJpc2V1c2VyMDAyIiwicGhvbmVfbnVtYmVyIjoiMzQ4NTI0MTY0OSIsInByZWZlcnJlZF91c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAyQHpvdHNlbGwuY29tIiwiZ2l2ZW5fbmFtZSI6IkVudGVycHJpc2V1c2VyMDAyIiwiZmFtaWx5X25hbWUiOiJFbnRlcnByaXNldXNlcjAwMiIsImVtYWlsIjoiZW50ZXJwcmlzZXVzZXIwMDJAem90c2VsbC5jb20iLCJ1c2VybmFtZSI6ImVudGVycHJpc2V1c2VyMDAyQHpvdHNlbGwuY29tIn0.T680TaiH_MfZqAnRNTlZa5rdC7k_dLRpmyAHRwIBVFNOj4wv6KHrZ-COSkm37l2hQvPutdehpJSa4B_9CGjL7h2ywmkZh3lAJnb1xEJsRtvBrfpiCUCHiWaGcZZ7KMJCWhCXYuaidGud4GbEb2wLRVPr_l1IVvbeKQ3d1OI2ure_4sNChVUuGpFFiNwHR4G-Q2sEWdmWd3GsLn10ninuafS2uJUH2VubbtGBVqbFLPGuYxNXYZjzId-p4Z12z7QXVvtc0zKHEhhpmq1ay2Mr60OjKefqshttfKg6zGB_E91zuUkEzsARLa3i8m6hulz2NyyAzGoz6kaN9fazWnMptg", 
  "not-before-policy": 0, 
  "refresh_expires_in": 1800, 
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI1YzcxNGJmZC03NThmLTRkZDItYTIyYy1jODJmMTk3NGIwZmMifQ.eyJleHAiOjE2MjE0OTcwNjcsImlhdCI6MTYyMTQ5NTI2NywianRpIjoiZmQ0OTdjZTItYjc0Yi00ZDJmLWJhYjMtZDQzYWUyYzJkMjZhIiwiaXNzIjoiaHR0cHM6Ly9kZXZsb2dpbi5jb21tZXJjaW8uYXBwL2F1dGgvcmVhbG1zL2NvbW1lcmNpbyIsImF1ZCI6Imh0dHBzOi8vZGV2bG9naW4uY29tbWVyY2lvLmFwcC9hdXRoL3JlYWxtcy9jb21tZXJjaW8iLCJzdWIiOiI0OGEwNGIyNy1kNjc5LTQ0MmItYmU2Yi0wMmE2ZDgyZjEwYjMiLCJ0eXAiOiJSZWZyZXNoIiwiYXpwIjoiZGV2LmNvbW1lcmNpby5hcHAiLCJzZXNzaW9uX3N0YXRlIjoiNTQ0ZWZiMDctNGE2ZC00N2RkLWEwYjYtYTg0ZWJkZDg1ZGU2Iiwic2NvcGUiOiJvcGVuaWQgZW1haWwgcHJvZmlsZSJ9.8c_9L0RSA52E3CgmvA-ahLyB5dH6Bc7eBNfGD6IcLx4", 
  "scope": "openid email profile", 
  "session_state": "544efb07-4a6d-47dd-a0b6-a84ebdd85de6", 
  "token_type": "Bearer"
}
```

* Rif : [Bearer Authentication](https://swagger.io/docs/specification/authentication/bearer-authentication/)
* Usefull guide for common Client available can be found here  <a href="https://openid.net/developers/certified/">Certified OpenID Connect Implementations </a>

--- 


## Securing your App 



Coming next
