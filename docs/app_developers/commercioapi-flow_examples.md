
# CommercioAPI Flows example

<!-- npm run docs:serve  -->

<!-- https://lcd-testnet.commercio.network/docs/did:com:1ug9j7hgaxu6mvfu2kgfdt3hqxn4mrwuztxc7nu/received -->

These flows are provided solely as examples and may be used as starting points for further development. 

It is assumed that an external software is in place to handle the logic and sequence of events, as well as to control access to the API. 

In the examples, this software stack is referred to as the 'Client App', which could be either a mobile or web application supported by a backend


## Basic document notarization 
`In Review  - cooming soon`

## Edelivery process with sharedoc and Receipt
`In Review  - cooming soon`


## User invitation with document notarization
`In Review  - cooming soon`

Suppose the target is to provide a function within the `Client App` (An app developed externally) that enables users NOT accredited in commercio.app to notarize a document to an existing eID, (suppose an administration's wallet address), while being authenticated in the `Client App` .

Thus the premisis are : 
* An administration account (receiver), referred to as the `administration account`, already exists in commercio.app. 

* The user (Sender) ,  referred to as the `user`, is authenticated within the 'Client App' but not accredited in commercio.app.

### Sequence of the workflow

* The `Client App`  present to the `user` an accreditation page asking to choose between two options

    *  a) **New user process**: Start commercio.app onboarding
    *  b) **Registered user process**:  Enter commercio.app password of the user as registered user in commercio.app

The `Client App` through the backend  uses the [Invite user api](/app_developers/commercioapi-ekyc.html#invite-user) loggin in with   `administration account` and using as input parameter the `user` email

a) If the user doesn't exist in commercio.app the Api will reply with an invitaition process Body response  and the user will receive an invitation email to commercio.app starting **New user workflow**
b) If the user exist in commercio.app  the Api will reply with a wallet address and the `Client App` should notify the user to enter his commercio.app credentials ad start the **Registered user workflow**


Our workflow regard the case  a) 

#### Invite  a) New user
The user click the link received by email perform through an external browser the step expected by the invitation process (Not Existing user) described [here](/app_developers/commercioapi-ekyc.html#invite-user) 

A `workflow_completed_redirect_uri` could be set in the [Invite user api](/app_developers/commercioapi-ekyc.html#invite-user) in order to redirect the user to the `Client App` if a web app 


#### Present the sharedoc form 
At this point the user is registered in commercio.app  


..... 


## User invitation with Spid recognition and request for document notarization 

Suppose the target is to provide a function within the `Client App` that enables users to notarize a document to an existing eID, (suppose an administration's wallet address), while being authenticated through their email account in the `Client App` .

 In addition, the user must be authenticated via a SPID process before they can notarize the document. This is necessary to retrieve the user's authentication private data from SPID within the `Client App` and store it in the backend for user identification purposes.

Thus the premisis are : An administration account, referred to as the `administration account`, already exists in commercio.app. Additionally, the user is authenticated within the 'Client App' using the email user001@email.com.


### Sequence of the workflow

* The `Client App`  present to the user an accreditation page asking to choose between two options

    *  a) **New user process**: Start commercio.app onboarding
    *  b) **Registered user process**:  Enter commercio.app password of the user as registered user in commercio.app

The `Client App` uses the [Invite user api](/app_developers/commercioapi-ekyc.html#invite-user) loggin in with  administrator account already present in commercio.app

If the user doesn't exist the Api will reply with an invitaition process Body response  and the user will receive an invitation email to commercio.app starting **New user workflow** othervise the Api will reply with a wallet address and the `Client App` will notify the user to enter his credentials ad start the **Registered user workflow**


#### Workflow a) New user
* 1. The user perform through an external browser the step expected by the invitation process (Not Existing user)
* 2. The  `Client App` check the [API KYC Verify Credentials](/app_developers/commercioapi-ekyc.html#verify-credentials) loggin in with  administrator account already present in commercio.app until `"audit": "public_digital_identity_system"` alternatively in the `Client App` present to the user an incomplete status message
* 3. Once the `Client App` veryfy that the `public_digital_identity_system` it propose  to enter the commercio account password in order to be used for log in the API service with the user credentials
* 4. Using with the user account the [API KYC Verify Credentials](/app_developers/commercioapi-ekyc.html#verify-credentials) the body response will reply with the  `public_digital_identity_system` element `content` filled. It must be Base 64 decoded and then sent to the backend of the `Client App` thus registering al the identification data of the user

Appropriete warning for exchanging personal data with the `Client App` app must be presented to the user and must be accepted. Is appropriate to register an audit of this acceptance and than make a sharedoc of the file containing the audit


* 5. At this point could be presented to the user the interface to istantiate a sharedoc message process [Send a shareDOc Message](/app_developers/commercioapi-sharedoc.html#send-a-sharedoc-message) from user wallet address to the one of the administrator notarizing the hash of the  document the app want to notarize in the Blockchain


#### Workflow b) Registered user

As registered user the process is similar to the previous one  except from the fact that after asking the user to enter his password of the commercio.app this workflow start from point 3. 


## Use sharedoc as a signature process 


## Use sharedoc and receipt as acceptance signature process 


## User request for pades like signatureand notarization 

In Review  - cooming soon