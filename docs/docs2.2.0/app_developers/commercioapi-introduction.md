# CommercioAPI introduction

<!-- npm run docs:serve  -->

<!-- https://lcd-testnet.commercio.network/docs/did:com:1ug9j7hgaxu6mvfu2kgfdt3hqxn4mrwuztxc7nu/received -->


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
| <a href="https://dev.commercio.app" target="_blank">dev.commercio.app</a>| <a href="https://commercio.app" target="_blank">commercio.app</a>  | Web App   |
| [dev-api.commercio.app/v1/](https://dev-api.commercio.app/v1/) | [api.commercio.app/v1/](https://api.commercio.app/v1/)  | CommercioAPI base url  |
| [dev-api.commercio.app/v1/swagger/index.html](https://dev-api.commercio.app/v1/swagger/index.html) | [api.commercio.app/v1/swagger/index.html](https://api.commercio.app/v1/swagger/index.html)  | Swagger  |
| devlogin.commercio.app/auth/realms/commercio/protocol/openid-connect/token| login.commercio.app/auth/realms/commercio/protocol/openid-connect/token   | IDM(OpenID) authentication URL |
| <a href="https://testnet.commercio.network" target="_blank">testnet.commercio.network</a>  | <a href="https://mainnet.commercio.network" target="_blank">mainnet.commercio.network</a>   | Explorer |

### Develop & Test
Is a playground where everyone can test the charcteristics 

#### dev.commercio.app

Memberships can be bought once registered using a sandbox credit card provided by stripe 

<a href="https://stripe.com/docs/testing#cards" target="_blank">Stripe Cards</a>


#### Faucet 
In the https://testnet.commercio.network chain is available a tool that permits to recharge a wallet 
with COM token . 

In order to use the tool is needed 

* `Endpoint` :  https://faucet-testnet.commercio.network/give
* `addr` : destination address
* `amount` : expressed in `ucommercio` 



**Example** 
Suppose you want to recharge with `10 COM`  your wallet address `did:com:1fjqvugs6dfwtax3k4zzh46pswmwryc8ff7f0mv`

This is the request to be used 

```
https://faucet-testnet.commercio.network/give?addr=did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr&amount=10000000
``` 

**Limits**

`amount` : There is a limit in the amount value of `100000000 ucommercio`


ATTENTION : Is not available a `faucet` for CCC . The Buy function (`coming soon`) in the dev.commercio.app must be used  


### Official
Is the realworld environment where real credit card transacions are requested and real token (COM and CCC) are spent



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
<br><br>

<img src="./img/swagger.png"> 

<br><br>
NB: Any Examples in the documentation refers to the **Develop & Test environement**


#### Hint : Basic Client for Major scritp Languages 

Downloading the `openapi.yaml` file from the Swagger interface page you can upload it in the  [https://editor.swagger.io/](https://editor.swagger.io/) 

Using the Generate Client menu you can obtain a basic stack software for the language you choose 



### Authentication process  
In order to gain proper access to the API an authetication process should be performed.

* <a href="/app_developers/commercioapi-authentication.html">AUTHORIZE</a> : Permit to authenticate through the IDM and get permission to API usage

### Available APi

The following api are available categorize by the Commercio.network Modules on which are based.

* <a href="/app_developers/commercioapi-wallet.html">Wallet</a> Permit to manage the the basic operations on your Wallet/s  throught the <a href="/x/bank/#sending-tokens ">Bank</a>  Module

* <a href="/app_developers/commercioapi-sharedoc.html">Sharedoc</a> Permit to manage the  <a href="/x/documents/#sending-a-document">MsgShareDocument</a>  eDelivery Digital Time Stamping to certify document integrity  using the <a href="/x/documents/#docs">DOCS MODULE</a> 


* ID (coming next): the eID. To Create and manage Self Sovereign Identities

* SIGN  (coming next) : the eSignature. to Electronically Sign any PDF e XML digital document
