# CommercioAPI introduction

CommercioAPI is a RESTfull web service  that  allows anyone to create transactions  with the set of permitted messages on the commercio.network  Blockchain 
and to query it. 

## Getting Started

In order to operate with the CommercioAPI are available two environements 
* Develop & Test : to test the APi and get familiar with the system in the Test-net
* Official  : to operate with the real CommercioAPI in the Main-net

The following services are envolved  with the CommercioAPI usage

*  <strong>Web App</strong> : permit to obtain and mange your membership, generate and manage your hostedwallet, get your account address, see the accountability of your transactions
*  <strong>CommercioAPI base url</strong>  : have access to the documentation of the Web RESTful services and interact with the API in the proper subpath
*  <strong>IDM(OpenID)</strong>  : The IDentity Management service to be used for proper auhentication in the APIs services 
*  <strong>Explorer</strong>  : The web application "**Almerico**" that permit to get the transaction informations from the Commercio.network blockchain ledger

### Directions 

These are the end points of the services in the **Develop & Test** environement and **Offical** (Production) one.

| Develop & Test | Official  | Note |
| --- | --- | ---|
| <a href="https://dev.commercio.app" target="_blank">dev.commercio.app</a>| Cooming soon  | Web App   |
| dev-api.commercio.app/v1/ | Cooming soon  | CommercioAPI base url  |
| <a href="https://devlogin.commercio.app" target="_blank">devlogin.commercio.app</a> | <a href="https://login.commercio.app" target="_blank">login.commercio.app</a>    | IDM(OpenID)  |
| <a href="https://testnet.commercio.network" target="_blank">testnet.commercio.network</a>  | <a href="https://mainnet.commercio.network" target="_blank">mainnet.commercio.network</a>   | Explorer |


## Prerequisite 

To use the APIs you need to : 

* Register and Login on web App (dev.commercio.app or commercio.app)
* Own a valid membership (Bronze,Silver,Gold) and get your `ID account`  (or own wallet address  example : `did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl` )
* Own enough CCCs (Commerce Cash Credits) to pay for the sending of the transactions.


## CommercioAPIs overview
Brief overview of the available functions. Refer tho the specific APIs guide for more details.

For any support or questions regarding the API or the documentation please open an <a href="https://github.com/commercionetwork/commercionetwork/issues" target="_blank">Issue </a>


### Swagger environement
In the **CommercioAPI base url**  in the path `/swagger/` is available the documentation of the set of released API  and a Tryout is possible through the Swagger interface

Example 

https://dev-api.commercio.app/v1/swagger/


NB: Any Examples in the documentation refers to the **Develop & Test environement**


#### Hint : Basic Client for Major scritp Languages 

Downloading the `openapi.yaml` file from the Swagger interface page you can upload it in the  [https://editor.swagger.io/](https://editor.swagger.io/) 

Using the Generate Client menu you can obtain a basic stack software for the language you choose 



### Authentication process  
In order to gain proper access to the API an authetication process should be performed.

* <a href="/app_developers/commercioapi-authentication.html">AUTHORIZE</a> : Permit to authenticate through the IDM and get permission to API usage

### Available APi

The following api are available categorize by the Commercio.network Modules on which are based.


* <a href="/app_developers/commercioapi-sharedoc.html">Sharedoc</a> Permit to manage the  <a href="/x/documents/#sending-a-document">MsgShareDocument</a>  eDelivery Digital Time Stamping to certify document integrity  using the <a href="/x/documents/#docs">DOCS MODULE</a> 


* ID (coming next): the eID. To Create and manage Self Sovereign Identities

* SIGN  (coming next) : the eSignature. to Electronically Sign any PDF e XML digital document
