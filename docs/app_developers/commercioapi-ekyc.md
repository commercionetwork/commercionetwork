# CommercioAPI eKyC

<!-- npm run docs:serve  -->



## Invite user 

This API permit to invite a user or verify if a user is already present in commercio through a check of his email 

### PATH

POST : `/eidentity/invite`


### Step by step Example

Let's create a new process to check if a user already exist

#### Step 1 - Define the first query and payload

Parameter value :
* Email to chack : es `john.doe@email.com`




#### Step 2 - Use the API to Send the message 


**Use the tryout**

Fill the swagger tryout with the payload

```
{
  "email_address": "john.doe@email.com"
}
```

**Corresponding Cli request**


```
curl -X 'POST' \
  'https://dev-api.commercio.app/v1/eidentity/invite' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOi....knnxylF4CA' \
  -H 'Content-Type: application/json' \
  -d '{
  "email_address": "john.doe@email.com"
}'
```


**API : Body response**


**Not Existing user**


In case the user doesn't exist in the platform (unkown user) the paltfrom instatiate an invitation process

```
{
  "user_id": "6e10447f-a956-49b2-a22b-a94ba5d05276",
  "sent": true
}
``` 

An invitation email will be sent to the user containing a magick link 

![Modal](./invitation_email.png)


As the user click the link will be onboarded in the commercio.app 

**a) He has to accept the Term of services**

![Modal](./accept_tos.png)



**b) He has to choose a password for future login in the app and use of the api services**

![Modal](./choose_a_password.png)

**c) He has to wait for the account setup and membership assignement**

![Modal](./setup_account_and_assigning_membership_green.png)

**d) He will receive a Green Membership  with a 1CCC in his wallet for free**

![Modal](./account_and_membership_assigned.png)



**Existing user**

Simply will be returned the wallet address associated to the email user account of the platform  indicated as payload of the APi Request

``` 
{
  "wallet_address": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
}
``` 


#### Common error

The APi obviuosly return datas relative to  users account with a hosted wallet in the commercio.app platform. 



## Verify Credentials 

This API permit to check a user level of accreditation 
 

### PATH

POST : `/ekyc`


### Step by step Example
Let's create a new process to check the level of accreditation of a user


#### Step 1 - Define the first query and payload

Parameter value :
* Address to check : es `did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr`



#### Step 2 - Use the API to Send the message 


**Use the tryout**

Fill the swagger tryout with the payload

```
did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr`
```

**Corresponding Cli request**


```
curl -X 'GET' \
  'https://dev-api.commercio.app/v1/eKYC/did%3Acom%3A1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiI..... X1v0SyprA'

```


**API : Body response**


**Other user**


The common response is as follow 
 

```
{
  "credentials": [
    {
      "audit": "email",
      "started_at": "2022-02-02T10:27:46Z"
    },
    {
      "audit": "mobile_number",
      "started_at": "2022-03-03T10:27:46Z"
    },
    {
      "audit": "credit_card",
      "started_at": "2023-02-22T16:34:02Z"
    },
    {
      "audit": "public_digital_identity_system",
      "started_at": "2023-03-29T13:57:21Z"
    }
  ]
}
``` 



**Your wallet**
In case you are the owner of the wallet address indicated in the payload the  Body response will contain also the datas fo the accreditation audit in the optional entity `content`


Example : 

``` 
{
  "credentials": [
    {
      "audit": "email",
      "started_at": "2022-02-02T10:27:46Z"
    },
    {
      "audit": "mobile_number",
      "started_at": "2022-03-03T10:27:46Z"
    },
    {
      "audit": "credit_card",
      "started_at": "2023-02-22T16:34:02Z"
    },
    {
      "audit": "public_digital_identity_system",
      "started_at": "2023-03-29T13:57:21Z",
      "content": "ewogICJhZGRyZXNzIjogIlZpYSAgTWFyaWFuaSA0MCAzNjAxNSAgTWVsYm91cm5lIFZJIiwKICAiY29tcGFueV9maXNjYWxfbnVtYmVyIjogIiIsCiAgImNvbXBhbnlfbmFtZSI6ICItIiwKICAiY291bnRyeV9vZl9iaXJ0aCI6ICJWSSIsCiAgImRhdGVfb2ZfYmlydGgiOiAiMTk0My0xMi0wOCIsCiAgImRpZ2l0YWxfYWRkcmVzcyI6ICItIiwKICAiZG9taWNpbGVfbXVuaWNpcGFsaXR5IjogIiBNZWxib3VybmUiLAogICJkb21pY2lsZV9uYXRpb24iOiAiVVNBIiwKICAiZG9taWNpbGVfcG9zdGFsX2NvZGUiOiAiMzYwMTUiLAogICJkb21pY2lsZV9wcm92aW5jZSI6ICJWSSIsCiAgImRvbWljaWxlX3N0cmVldF9hZGRyZXNzIjogIlZpYSAgTWFyaWFuaSA0MCIsCiAgImVtYWlsIjogImpvaG4uZG9lQGVtYWlsbC5jb20iLAogICJmYW1pbHlfbmFtZSI6ICJNb3JyaXNvbiIsCiAgImZpc2NhbF9udW1iZXIiOiAiVElOSVQtSk1NTVJTNDNMMDhMODQwUiIsCiAgImdlbmRlciI6ICJNIiwKICAiaWRfY2FyZCI6ICJjYXJ0YUlkZW50aXRhIENBMjMxMTM0NUhQIGNvbXVuZSBNZWxib3VybmUgMjAzMS0wNy0yMyIsCiAgIml2YV9jb2RlIjogIi0iLAogICJtb2JpbGVfcGhvbmUiOiAiKzM5MzQ4MTExMTExMTExIiwKICAibmFtZSI6ICJKaW0iLAogICJyZWdpc3RlcmVkX29mZmljZSI6ICItIgp9"
    }
  ]
}
``` 

The `content` data is base 64 encoded. 

Thus decoding it will appear somthing like this 



``` 
{
  "address": "Via  Mariani 40 36015  Melbourne VI",
  "company_fiscal_number": "",
  "company_name": "-",
  "country_of_birth": "VI",
  "date_of_birth": "1943-12-08",
  "digital_address": "-",
  "domicile_municipality": " Melbourne",
  "domicile_nation": "USA",
  "domicile_postal_code": "36015",
  "domicile_province": "VI",
  "domicile_street_address": "Via  Mariani 40",
  "email": "john.doe@emaill.com",
  "family_name": "Morrison",
  "fiscal_number": "TINIT-JMMMRS43L08L840R",
  "gender": "M",
  "id_card": "cartaIdentita CA2311345HP comune Melbourne 2031-07-23",
  "iva_code": "-",
  "mobile_phone": "+39348111111111",
  "name": "Jim",
  "registered_office": "-"
}



``` 


#### Common error