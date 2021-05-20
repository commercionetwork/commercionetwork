# CommercioAPI introduction

CommercioAPI is a RESTfull web service  that  allows anyone to deliver a blockchain message through a transaction on the commercio.network 

## Getting Started

In order to operate with the CommercioAPI are available two environement 
* Develop & Test environement : to test the APi and get familiar with the system in the Test-net
* Official  : To operate with the real CommercioAPI in the Main-net

The following services are envolved  with the CommercioAPI usage

*  <strong>Web App</strong> : permit to obtain and mange your membership, generate and manage your hostedwallet, get your account address, see the accountability of your transactions
*  <strong>CommercioAPI base url</strong>  : have access to the documentation of the Web RESTful services 
*  <strong>IDM(OpenID)</strong>  : The Identity management service to be used for the correct auhentication in the API services 
*  <strong>Explorer</strong>  : The web application "**Almerico**" that permit to get the transaction informations from the Commercio blockchain ledger

### Directions 

This are the end point of the serivices in the Develop & Test environement and Offical (Production) one.

| Develop & Test | Official  | Note |
| --- | --- | ---|
| <a href="https://dev.commercio.app" target="_blank">dev.commercio.app</a>| Cooming soon  | Web App   |
| dev-api.commercio.app/v1/ | Cooming soon  | CommercioAPI base url  |
| <a href="https://devlogin.commercio.app" target="_blank">devlogin.commercio.app</a> | <a href="https://login.commercio.app" target="_blank">login.commercio.app</a>    | IDM(OpenID)  |
| <a href="https://testnet.commercio.network" target="_blank">testnet.commercio.network</a>  | <a href="https://mainnet.commercio.network" target="_blank">mainnet.commercio.network</a>   | Explorer |


## Prerequisite 

To use the API you need to

* Register and Login on web App (dev.commercio.app/commercio.app
* Own a valid membership (Bronze,Silver,Gold) and get your id account  (or own wallet address  example : `did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl` )
* Own enough CCC (Commerce Cash Credits) to pay the fees of the transactions you send


## CommercioAPI overview
Brief overview of the available functions. Refer tho the specific API guide for more details.

For any support or questions regarding the API or the documentation please open an <a href="https://github.com/commercionetwork/commercionetwork/issues" target="_blank">Issue </a>


### Swagger environement
In the **CommercioAPI base url**  is available the documentation of the set of released API  and a Tryout is possible
through the Swagger interface


NB: Any Examples in the documentation refers to the **Develop & Test environement**


### Authentication process  
In order to gain proper access to the API an authetication process should be performed.

* <a href="">AUTHORIZE</a> : Permit to authenticate through the IDM and get permission to API usage

### Available APi

The following api are available categorize by the Commercio.network Modules on which are based.


* <a href="/app_developers/commercioapi-sharedoc.html">Sharedoc</a> Permit to manage the  <a href="/x/documents/#sending-a-document">MsgShareDocument</a>  eDelivery Digital Time Stamping to certify document integrity  using the <a href="/x/documents/#docs">DOCS MODULE</a> 


* ID (coming next): the eID. To Create and manage Self Sovereign Identities

* SIGN  (coming next) : the eSignature. to Electronically Sign any PDF e XML digital document
