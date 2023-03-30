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
* He has to accept the Term of services
* He has to choose a password for future login in the app and use of the api services
* He will receive a Green Membership  with a 1CCC hfo freei his wallet for free



**Existing user**



``` 
{
  "wallet_address": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
}
``` 


#### Common error




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
In case you are the owner of the wallet address indicated the  Body response wil contain also the data fo the accreditation audit 


``` 

.....



``` 


#### Common error