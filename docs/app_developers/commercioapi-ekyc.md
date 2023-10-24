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


**A) Not Existing user**


In case the user doesn't exist in the platform (unkown user) the paltfrom instatiate an invitation process

```
{
  "user_id": "6e10447f-a956-49b2-a22b-a94ba5d05276",
  "sent": true
}
``` 


**B) Existing user**

Simply will be returned the wallet address associated to the email user account of the platform  indicated as payload of the APi Request

``` 
{
  "wallet_address": "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
}
``` 


### Not Existing user workflow 

When a A) Not Existing user reply is received by the api an oboarding workflow start 

1) An invitation email will be sent to the user containing a magick link 

![Modal](./invitation_email.png)


2) As the user click the link will be onboarded in the commercio.app 

**2.a) He has to accept the Term of services**

![Modal](./accept_tos.png)



**2.b) He has to choose a password for future login in the app and use of the api services**

![Modal](./choose_a_password.png)

**2.c) He has to wait for the account setup and membership assignement**

![Modal](./setup_account_and_assigning_membership_green.png)

**2.d) He will receive a Green Membership  with a 1CCC in his wallet for free**

![Modal](./account_and_membership_assigned.png)

### send_email parameter 

Is a parameter permitted only for invites from Gold membership members that allows deactivation of the email sending invitation . The invite URL will be displayed as a response to the API.

### password parameter 
Is a parameter permitted only for invites from Gold membership members that allows to directly set the password for the invited user. The password to the user must be comunicated out of the platform or the user could use the forgot password procedire in the platform login page  

### workflow_completed_webhook_callback

At the end of the workflow  process istantiated through the API parameters , a POST call is made to the URL specified in the workflow_completed_webhook_callback parameter.

The body of the POST request is as follows:

```json
{
  "message": "The user has successfully completed the onboarding workflow.",
  "wallet_address": "did:com:1ffsmvvt29r....a5gwdre46tkz22n7vdj",
  "user": "john.doe@userdomain.com",
  "success": true
}
```

Where user is the email provided as a parameter to the invite endpoint, and thus the email of the new user.
This Post Call is invoked at the end of all steps of the workflow defined. Differently from other parameter 
`workflow_wallet_created_callback` invoked just after wallet creation and membership assignement

###  requires_spid_identification

If this parameter is set = true the user after onboarding will be redirected to SPID authentication process 


Example Body Payload

```json
{
  "email_address": "john.doe@yourdomain.com",
  "requires_spid_identification": true,
  "password": "JhonDOe8",
  "workflow_completed_redirect_uri": "https://www.yourdomain.com"
}

```

With these parameters in the payload, the user, after clicking the invitation link and accepting the TOS, will be redirected to the SPID authentication process. After completing it, they will be directed to the URL 'https://www.yourdomain.com.'"


###  workflow_wallet_created_callback
At the end of the onboarding process, which occurs after the user has accepted the invitation throught a magic link and completed the onboarding procedure (Tos acceptance,assignment of membership) , a POST call is made to the URL specified in the workflow_wallet_created_callback parameter.


The body of the POST request is as follows:

```json
{
....

  "user": "john.doe@userdomain.com",
  "success": true
}
```

Where user is the email provided as a parameter to the invite endpoint, and thus the email of the new user.

Pay attention this POST is invoked just after wallet creation and membership assignement even if the workflow 
defined by other parameters is not completed


#### Common Question

<strong>Which are the users recognized by the APIs ? </strong>

The API obviuosly return datas relative only to existing users account with a hosted wallet in the commercio.app platform. 


<strong>Does the invitation has a validity time ? </strong>

Yes, the invitation lasts for 24 hours. When it expires and is clicked, the application will show this page.

![Modal](./expired_invite.png)


Regardless, the user is created and can retrieve the password through the 'Forgot Password' link.











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

Thus decoding it will appear somthing like this depending on the data received form de IDP provider  



``` 
{
  "name": "Jim",
  "family_name": "Morrison",
  "fiscal_number": "TINIT-JMMMRS43L08L840R",
  "date_of_birth": "1943-12-08",
  "country_of_birth": "IT",
  "gender": "M",
  "company_name": "-",
  "company_fiscal_number": "",
  "iva_code": "-",
.... ....

  "address": "Via  Mariani 40 36015  Melbourne VI",
  "digital_address": "-",
  "domicile_municipality": " Melbourne",
  "domicile_nation": "USA",
  "domicile_postal_code": "36015",
  "domicile_province": "VI",
  "domicile_street_address": "Via  Mariani 40",
  "email": "john.doe@emaill.com",
  "id_card": "cartaIdentita CA2311345HP comune Melbourne 2031-07-23",
  "mobile_phone": "+39348111111111",
  "registered_office": "-"
}



``` 


#### Common error




## Request SPID session   

This API permit to obtain a direct link for a spid session for the user logged in the api 

### PATH

POST : `/ekyc/spid`


### Step by step Example
Let's create a request to obtain the spid session url 


#### Step 1 - Define the first query and payload


Parameter value :
* success_url : es `"https://www.yourlandingurl.com`

Is the url where the user are redirected at the end of Spid recognition process (Spid authentication)


#### Step 2 - Use the API to Send the message 


**Use the tryout**

Fill the swagger tryout with the Body payload

```
{
  "success_url": "https://www.yourlandingurl.com"
}
```

**Corresponding Cli request**


```
curl -X 'GET' \
curl -X 'POST' \
  'https://dev-api.commercio.app/v1/eKYC/spid' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1....54Eg' \
  -H 'Content-Type: application/json' \
  -d '{
  "success_url": "https://www.yourlandingurl.com"
}'

```


**API : Body response**

``````

{
  "authentication_url": "https://spid.commercio.app/?subject_id=19fcd4f6-80cc-42ce-8eac-3b276d9794eb&success_redirect_uri=https://www.yourlandingurl.com&attributes_uri=https://commercio.app"
}

```

The URL can be used by the user to perform SPID authentication and save data on the platform for recognition. After the SPID authentication, the user will be redirected to 'https://www.yourlandingurl.com.'



#### Common error