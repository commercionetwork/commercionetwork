# CommercioAPI introduction

<!-- npm run docs:serve  -->

<!-- https://lcd-testnet.commercio.network/docs/did:com:1ug9j7hgaxu6mvfu2kgfdt3hqxn4mrwuztxc7nu/received -->


CommercioAPI is a RESTful web service  that  allows anyone to create transactions with the set of permitted messages on the commercio.network  Blockchain 
and to query it. 

## Getting Started

In order to operate with the CommercioAPI two environements are available 
* Develop & Test: to test the API and get familiar with the system in the Test-net
* Official: to operate with the real CommercioAPI in the Main-net

The following services are envolved with the CommercioAPI usage

*  <strong>Web App</strong>: to obtain and manage your membership, generate and manage your hostedwallet, get your account address, see the accountability of your transactions
*  <strong>CommercioAPI base url</strong>: to have access to the documentation of the Web RESTful services and interact with the API in the proper subpath
*  <strong>IDM(OpenID)</strong>: The IDentity Management service to be used for proper auhentication in the APIs services 
*  <strong>Explorer</strong>: The web application "**Almerico**" that permits to get the transaction information from the Commercio.network blockchain ledger

### Directions 

These are the endpoints of the services in the **Develop & Test** environment and **Official** (Production) one.

| Develop & Test | Official  | Note |
| --- | --- | ---|
| <a href="https://dev.commercio.app" target="_blank">dev.commercio.app</a>| <a href="https://commercio.app" target="_blank">commercio.app</a>  | Web App   |
| [dev-api.commercio.app/v1/](https://dev-api.commercio.app/v1/) | [api.commercio.app/v1/](https://api.commercio.app/v1/)  | CommercioAPI base url  |
| [dev-api.commercio.app/v1/swagger/index.html](https://dev-api.commercio.app/v1/swagger/index.html) | [api.commercio.app/v1/swagger/index.html](https://api.commercio.app/v1/swagger/index.html)  | Swagger  |
| devlogin.commercio.app/auth/realms/commercio/protocol/openid-connect/token| login.commercio.app/auth/realms/commercio/protocol/openid-connect/token   | IDM(OpenID) authentication URL |
| <a href="https://testnet.commercio.network" target="_blank">testnet.commercio.network</a>  | <a href="https://mainnet.commercio.network" target="_blank">mainnet.commercio.network</a>   | Explorer |

### Develop & Test
Is a playground where everyone can test the characteristics of the commercio.network blockchain.

#### dev.commercio.app

Memberships can be bought once registered using a sandbox credit card provided by Stripe 

<a href="https://stripe.com/docs/testing#cards" target="_blank">Stripe Cards</a>


#### Faucet 
In the https://testnet.commercio.network chain, a tool that allows to recharge a wallet 
(with COM token) is available. 

A destination address (`addr`) and the amount to be recharged with  (`amount` expressed in ucommercio) must be provided to the faucet endpoint (https://faucet-testnet.commercio.network/give).

**Example** 
Suppose you want to recharge with `10 COM` your wallet address `did:com:1fjqvugs6dfwtax3k4zzh46pswmwryc8ff7f0mv`

This is the request you need to make: 

```
https://faucet-testnet.commercio.network/give?addr=did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr&amount=10000000
``` 

**Limits**

`amount` : There is a limit in the amount value of `100000000 ucommercio`


ATTENTION : A `faucet` for CCC is not available. The Buy function (`coming soon`) in the dev.commercio.app must be used  


### Official
Is the real-world environment where real credit card transactions are requested and real tokens (COM and CCC) are spent



## Prerequisites 

To use the APIs you need to: 

* Register and Login on the web App (dev.commercio.app or commercio.app)
* Own a valid membership (Bronze,Silver,Gold) and get your `ID account` (or your own wallet address, e.g. : `did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl` )
* Own enough CCCs (Commerce Cash Credits) to pay for the transaction fees.


## CommercioAPIs overview
Brief overview of the available functions. Refer to the specific APIs guide for more details.

For any support or questions regarding the APIs or the documentation, please open an <a href="https://github.com/commercionetwork/commercionetwork/issues" target="_blank">Issue </a>


### Swagger environement
In the **CommercioAPI base url**  in the path `/swagger/` the documentation of the set of released API is available, and a Tryout is possible through the Swagger interface.

Example 

https://dev-api.commercio.app/v1/swagger/
<br><br>

<img src="./img/swagger.png"> 

<br><br>
NB: Any Examples in the documentation refers to the **Develop & Test environement**


#### Hint : Basic Client for Major script Languages 

Downloading the `openapi.yaml` file from the Swagger interface page you can upload it in the  [https://editor.swagger.io/](https://editor.swagger.io/) 

Using the Generate Client menu you can obtain a basic stack software for the language you choose 



### Authentication process  
In order to gain proper access to the API an authetication process should be performed.

* <a href="/app_developers/commercioapi-authentication.html">AUTHORIZE</a> : Permit to authenticate through the IDM and get permission to API usage

### Available API

The following APIs are available

* <a href="/app_developers/commercioapi-wallet.html">Wallet</a> to manage the basic operations on your Wallet/s throught the <a href="/x/bank/#sending-tokens">Bank</a>  Module

* <a href="/app_developers/commercioapi-sharedoc.html">Sharedoc</a> to manage the <a href="/x/documents/#sending-a-document">MsgShareDocument</a> eDelivery Digital Time Stamping to certify document integrity using the <a href="/x/documents/#docs">DOCS MODULE</a> 


* ID (coming next): the eID. To Create and manage Self Sovereign Identities

* SIGN  (coming next) : the eSignature. to Electronically Sign any PDF e XML digital document
