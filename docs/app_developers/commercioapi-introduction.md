# CommercioAPI introduction

<!-- npm run docs:serve  -->

<!-- https://lcd-testnet.commercio.network/docs/did:com:1ug9j7hgaxu6mvfu2kgfdt3hqxn4mrwuztxc7nu/received -->


CommercioAPI is a RESTful web service  that  allows anyone to create transactions with the set of permitted messages on the commercio.network  Blockchain and to query it. 

## Foreword

Querying and interacting with the blockchain  Commercio Network is free. 

In order to send messages to the chain however You need first to control a wallet 


### What is a wallet

A wallet (digital) is a software (Encryption) that provide a virtual equivalent of a wallet.

In the common sense a wallet is an instrumental where you can storage your coins. 

You don’t actually store any cryptocurrency in your wallet. You just manage the `keys` to access 
it on the blockchain.

Attention:
> The blockchain records the amount of coins associated with a key pair (your identity on the blockchain).
It calculates the amount of money the keys have access to based on all the transactions on the blockchain. 
Remember: the main function of a blockchain is to store all transactions in the correct order.

Thorugh this calculation of past transaction you are able to check your balance, receive, and 
send funds with another wallet registering transfer transaction in the blockchain

A keychain concept is similar to what a wallet does. 

To spend your money, you need the private key stored in your wallet. 

Is not needed to understand how public-key cryptography works in detail but 
is important to understand that <strong>if you don’t control your keys, you don’t control your funds.</strong>

In summary a wallet is a program that has three main functions:

* Generating, storing and handling your keys and addresses
* Showing you your balance
* Creating and signing message transactions for example  send funds to another wallet address

Mainly however it permit to generate, manage and store cryptographic keys - your public and private key

Another important things is that wallets generally don’t allow you to buy cryptocurrencies;
Exchanges perform this for you.

All exchanges propose you wallets where store the coins you buy,
but wallets usually DON'T provide you any exchange service.

Thus the only way to get coins in your wallet is to receive them from another one.  


### Manage the keys - types of wallet 
The way a user decide to manage the crypto keys of his wallet have two main scenarios 

#### 1. Non-custodial wallets
Those wallets provide an interface to check your funds or create transactions in your web browser, but you have to provide the keys with each login.

Registering with a central authority is not needed to create a wallet. This comes at the cost of you being responsible for the safety of your coins Nobody can help you recover your keys if you lose them. 
If anybody were able to recover it they would also be able to steal your funds. 

This would eliminate the trustless aspect of blockchains

There is a sort of recovery mechanism with many wallets called a mnemonic phrase or backup phrase. A mnemonic phrase usually consists of 12 or 24 words. With these words, you can recover your keys. You receive your mnemonic phrase when you install and set up your wallet. Be sure to write it down on a piece of paper and keep it in a safe place. You should have at least two versions of your backup phrase stored in different locations.


#### 2. Custodial wallets or hosted wallet
It’s called hosted because a third party keeps your 
crypto for you, similar to how a bank keeps your money in a checking or savings account.

With hosted web wallets, your keys are stored online by a trusted third party.

Is the most popular and easy-to-set-up crypto wallet

You may have heard about people searching for old hard drives because they have “lost their bitcoins”. More accurately, they lost the keys to access their bitcoin.
but with a hosted wallet you don’t have to worry about any of that.


<div style="color:red;">


<!-- 
https://academy.horizen.io/technology/advanced/types-of-wallets/ 

https://www.coinbase.com/it/learn/tips-and-tutorials/how-to-set-up-a-crypto-wallet
-->

</div>




## The Commercio.app 
The Commercio App is a platform that provide You through a web app a <strong>hosted wallet</strong>.

Moreover the Commercio.app provides you web functionality for 

### eID Module
Associate identity data to your wallett address through specific procedures

### Credits Module
Manage CCC associated to your wallet

### API
Gain access to a set of rest API some of which permit You to interact in the blockchain throught your hosted wallet account


## Getting Started

In order to operate with the CommercioAPI two environements are available 
* <strong>Develop & Test </strong>: to test the API and get familiar with the system in the Test-net
* <strong>Official </strong>: to operate with the real CommercioAPI in the Main-net

The following services are envolved with the CommercioAPI usage

*  <strong>Web App</strong>: to obtain and manage your membership, generate and manage your hostedwallet, get your account address, see the accountability of your transactions
*  <strong>CommercioAPI base url</strong>: to have access to the documentation of the Web RESTful services and interact with the API in the proper subpath
*  <strong>IDM(OpenID)</strong>: The IDentity Management service to be used for  auhentication in the APIs services 
*  <strong>Explorer</strong>: The web application "**Almerico**" that permits to get the transaction information from the Commercio.network blockchain ledger

### Develop & Test
Is a playground where everyone can test the characteristics of the commercio.network blockchain.


| Develop & Test | Official  | Note |
| --- | --- | ---|
| <a href="https://dev.commercio.app" target="_blank">dev.commercio.app</a> | Web App   |
| [dev-api.commercio.app/v1/](https://dev-api.commercio.app/v1/) | CommercioAPI base url  |
| [dev-api.commercio.app/v1/swagger/index.html](https://dev-api.commercio.app/v1/swagger/index.html) | Swagger  |
| devlogin.commercio.app/auth/realms/commercio/protocol/openid-connect/token | IDM(OpenID) authentication URL |
| <a href="https://testnet.commercio.network" target="_blank">testnet.commercio.network</a> | Explorer |

Memberships on this environment can be bought for free  once registered using a dummy  credit card provided by Stripe (example 4242 4242 4242 4242)

Reed <a href="https://stripe.com/docs/testing#cards" target="_blank">Stripe Cards</a> for details


### Official
Is the real-world environment where real credit card transactions are requested and real tokens (COM and CCC) are spent


| Endpoint | Official  | Note |
| --- | --- | ---|
| <a href="https://commercio.app" target="_blank">commercio.app</a>  | Web App   |
| [api.commercio.app/v1/](https://api.commercio.app/v1/)  | CommercioAPI base url  |
| [api.commercio.app/v1/swagger/index.html](https://api.commercio.app/v1/swagger/index.html)  | Swagger  |
| login.commercio.app/auth/realms/commercio/protocol/openid-connect/token   | IDM(OpenID) authentication URL |
| <a href="https://mainnet.commercio.network" target="_blank">mainnet.commercio.network</a>   | Explorer |


Memberships on this environment MUST be bought with real credit card 






## Prerequisites 

To use the APIs you need to: 

* Register and Login on the web App (dev.commercio.app or commercio.app)
* Own a valid membership (Bronze,Silver,Gold) and get your `ID account` (or your own wallet address, e.g. : `did:com:1r0sk6stfm6d5jtfcne2jxd7s7n2whp35tjm7zl` )
* Own enough CCCs (Commerce Cash Credits) to pay for the transaction fees.


### Message Chain fee vs Commercio.app costs

#### Chain costs  

Sending message in the blockchain has a "protocol" 's  cost that is due in order to  support Validators expenses.

The cost in Commercio blockchain for each message coul be paid with [CCC](/#the-commercio-cash-credit-ccc) eihter with [COM](/#the-commercio-token-com) 
and it correspont to : 
* 0.01 COM per message sent
* 0.01 CCC per message sent

Anyway due to its nature it is better to bay with CCC due to the fact that CCC has a fix 
value of 1 Euro intead of COM that has a VARIABLE Euro value and is determined according supply and demand by the market

You can easily make your own wallet management and send messages to the chain 
at protocol's costs using the [following procedure](/docs2.2.0/developers/create-sign-broadcast-tx.html#_1-message-creation) but obviusly this need to be done by developer experts   

Anotehr way is using Restfull commercio API 


#### Platform costs  
Moreover sending messages autonoumously could also be done with Resfull API of the commercio.app

Using the platform comes with a cost for the usage that depends on the type of [Membrship you subscribe](modules/commerciokyc/#membership-types) 

Costs of membership expressed in Euro when bought from commercio.app IVA not included


|Membersip| Platform cost| Chain fee | Total cost per msg |
| --- | --- | --- | --- |
| Green | 0.24 | 0.01  | 0.25 |
| Bronze | 0.11 | 0.01 | 0.12 | 
| Silver | 0.05 | 0.01 | 0.06 | 
| Gold | 0.02 | 0.01 | 0.03 | 

* <sub>All costs in CCC.</sub> 
* <sub>Platform cost is  comprehensive of a Chain fee (0.01) for sending Platform costs to a wallet of Commercio Platform</sub> 

NB: ACTUAL COST OF PLATFORM is 0 CCC thus actually you are paying only Chain fee. Platform Surcharge will be actuated at the end of 1Q 2023 





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


