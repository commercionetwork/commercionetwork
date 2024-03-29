# The Commercio.network Blockchain Documentation

Commercio.network is an [Open Source Blockchain](https://github.com/commercionetwork) that allows people to create:

* eID  electronic identities
* eSignatures electronic signature
* eDelivery  certified delivery

Anyone can exchange electronic documents in a legally binding way thanks to our eIDAS Compliance.

## What is a blockchain ?

A blockchain is a big distributed database. Think of it as a huge spreadsheet runned simultaneously on millions of worldwide computers. It’s a peer to peer network of nodes where you can settle transactions without the need of any trusted third party.

A node is a computer that is running the Commercio.network node software 
it is connected to other computers on the same network and there are two kind of nodes:

* **Full node** Full nodes are nodes that stores the whole transaction history. They connect to the blockchain and each time a new block is finalized, they write it on their hard disk. This means that being a full node you will be able to read the whole chain transaction history, you will need to have a large hard disk space if you want to keep it running.

* **Validator node**  Validator nodes are  full nodes with the added ability of validating new transactions that should be added to the chain. In order to do so, they possess a private key with which they sign the transactions marking them as valid. In exchange of their work, they get  a reward that is given to them each time a new block is created.

## What is the the software that allow this blockchain to exist ?

`cn` is the name of the [Commercio.network](https://commercio.network) blockchain application.    
It is shipped with `commercionetworkd`: The Commercio.network software provided daemon to run a full-node of the `cn` application and the command-line interface, which enables interaction with a Commercio.network full-node.
 

`cn` is built on top of the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) using the following modules:

* `x/auth`: Accounts and signatures.
* `x/bank`: Tokens transfers.
* `x/staking`: Staking logic.
* `x/distribution`: Fee distribution logic.
* `x/slashing`: Slashing logic.
* `x/params`: Handles app-level parameters.
* `x/ibc`: (wip).
* `x/wasm`: (wip).

On top of that `cn` comes with the following custom modules: 

* [`x/ante`](modules/ante/README.md): Custom fees. 
* [`x/documents`](modules/documents/README.md): Documents sharing. 
* [`x/did`](modules/did/README.md): Self sovereign  identities creation.
* [`x/government`](modules/government/README.md): On-chain government. 
* [`x/commerciokyc`](modules/commerciokyc/README.md): Invite new members and get ABR rewards. 
* [`x/commerciomint`](modules/commerciomint/README.md): Mint CCCs. 
* [`x/vbr`](modules/vbr/README.md): Run Validator nodes and get VBR rewards 


### Can I run this `cn` Node software ?

Sure. Please follow the step by step instructions on Running Nodes

### How can I develop an app on commercio.Network ?

**API**

The EASY WAY is to signup to [commercio.app](https://commercio.app) and use the CommercioAPI. You can start developing blockchain solutions in minutes with the programming language you are most familiar.

The CommercioAPI removes the complexity and the security of managing your users' wallets.

**SDK**


If you want to move forward the develop of the follow packages you can fork them and pull request upgrades.

We have released SDK in 4 main languages 

|  | Dart/Flutter | Kotlin/Java | C#/Dot.net | GoLang |
| ------ | ------ | ------ | ------ | ------ | ------ | ------ | ------ |
| **Sacco**  | [Repo](https://github.com/commercionetwork/sacco.dart) | [Repo](https://github.com/commercionetwork/sacco.kt) | [Repo](https://github.com/commercionetwork/sacco.cs) |  [Repo](https://github.com/commercionetwork/sacco.go) | 
| **CommercioSDK**  | [Repo](https://github.com/commercionetwork/commercio-sdk.dart) | [Repo](https://github.com/commercionetwork/commercio-sdk.kt) | [Repo](https://github.com/commercionetwork/commercio-sdk.cs) | Later |

<span style="color:red">Actual Sdks available are deprecated e no longer mantained.</span> 

Basic procedure for sending autonomusly a message is decribed here [Create, sign and send a transaction](/docs2.2.0/developers/create-sign-broadcast-tx.html#_1-message-creation)



## What is eIDAS Compliance ?

The eIDAS directive (Electronic Identification, Authentication and Trust Services) is an EU regulation on electronic identification and trust services for electronic transactions in the European Single Market.

The eIDAS org oversees electronic identification and trust services for electronic transactions in the European Union's internal market. It regulates electronic signatures, electronic transactions, involved bodies, and their embedding processes to provide a safe way for users to conduct business online like electronic funds transfer or transactions with public services. Both the signatory and the recipient can have more convenience and security. Instead of relying on traditional methods, such as appearing in person to submit paper-based documents, they may now perform transactions across borders.


### The advantage of using a eIDAS Compliant blockchain

According to Article 25.1 of the eIDAS Regulation, a standard electronic signature may not be denied legal effect and admissibility as evidence in legal proceedings solely on the grounds that it is in an electronic form or that it does not meet the requirements for qualified electronic signatures. 


## The Commercio Token (COM)

Commercio.network is a sovereign network that has its own native crypto currency which serves for:

* Incentivize users to manage the nodes of this network
* Incentivize users to grow this network by inviting other users.
  
Through this Token all active participants can benefit from the growth of the network.


The Commercio Token (COM) in not inflationary since it has a 60 million limited supply

* The main purpose of the Token is to be a unit of value that can be placed on stake to secure the network of Commercio.network by the validators nodes It is a STAKING TOKEN (utility Token)
* The price is VARIABLE and is determined according supply and demand by the market 


## The Commercio Cash Credit (CCC)

Commercio.network  is an Enterprise-grade third generation blockchain that removes some complexity problems of second generation blockchain like BitCoin or Ethereum:

**Any transaction Cost on our chain is defined in EURO and it costs €0.01 which makes it maybe the first StableChain in history.**
  
Through this Token all active participants can benefit from the growth of the network.


The Commercio Cash Credit (CCC)  ha an unlimited supply and can be minted only by freezing the COM Token

* The main purpose of the Token is to be a unit of value that can be used to perform transactions on chain. It is a FEE TOKEN (utility Token)
* The price is FIXED and is 1 EURO/CCC + VAT
  
**NB**: transaction fees could be paid with COM, but the cost is fixed to 0.01 COM.

##  Main-net Vs Test-net  




### The Main-net

Mainnet is the main blockchain of Commercio.network. If someone says Commercio mainnet, it means the real Commercio.network blockchain.

Unlike testnet which is an open network for testing purposes, mainnet is the real deal. Tokens on the Commercio.network mainnet have real economic value, be careful.

A view of the status of mainnet is possible throught the Explorer, a Webapp  named <strong>Almerico</strong> (whose code is Opensource and available <a href="https://github.com/commercionetwork/almerico" target="_blank">here</a> ) that query the mainnet chain.  


Main-net Almerico:  [mainnet.commercio.network](https://mainnet.commercio.network)

#### Endpoint & resources

|Description| Endpoint |
| --- | --- | 
| Explorer | https://mainnet.commercio.network/  | 
| LCD |  https://lcd-mainnet.commercio.network/  | 
| RPC |  https://rpc-mainnet.commercio.network/  | 
| GRPC |  grpc-mainnet.commercio.network:9090  |
|Commercio wallet app IOS  |   https://apps.apple.com/it/app/commerc-io/id1397387586  |
|Commercio wallet app Android |  https://play.google.com/store/apps/details?id=io.commerc.preview.one   |
|Commercio  app |  https://commercio.app  |
|Commercio  app API |   https://api.commercio.app/v1/swagger/index.html   |



#### How can I get COM tokens?

COM Token are reserved to infrastructure operators (aka Validators) and power users 

If you are interested in becoming a validator contact info@commercio.network  



#### How can I get CCC tokens?

CCC tokens can be minted by owner of COM tokens creating a position with COM tokens through the [Commercio wallet app](https://github.com/commercionetwork/Commercio-Wallet-App) . The function is available in the CCC Menu and corresponds to the "Mint" button.

CCC tokens can also be purchased directly from commercio.network SPA by contacting info@commercio.network. Please note that the "Buy" function in the commercio.app is currently unavailable.


### The Test-net

Testnet, as the name suggests, is an alternative network for the developers for testing purposes. It's a playground. You can use testnet as a demo network for experimenting. It’s like the beta stage of a blockchain network. A testnet is a blockchain made available for developers. It allows anyone to conduct experiments without wasting real tokens. A testnet is like a demo network where tokens do not have any value. You can easily test out any app on a testnet because it provides you a sandbox environment separate from the main blockchain.


Test-net Almerico: [testnet.commercio.network](https://testnet.commercio.network)


#### Endpoint & resources
|Description| Endpoint |
| --- | --- | 
| Explorer |  https://testnet.commercio.network/  | 
| LCD |  https://lcd-testnet.commercio.network/  | 
| RPC |  https://rpc-testnet.commercio.network/  |
| GRPC |  grpc-testnet.commercio.network:9090  |
| Faucet |   https://faucet-testnet.commercio.network/  |
|Commercio wallet app IOS  |  https://apps.apple.com/it/app/commerc-io/id1397387586   |
|Commercio wallet app Android |  https://play.google.com/store/apps/details?id=io.commerc.preview.one   |
|Commercio  app |   https://dev.commercio.app/   |
|Commercio  app API |   https://dev-api.commercio.app/v1/swagger/index.html   |


####  How can I get COM tokens?
Getting COM token in testnet is quite easy and free. It is possible throught a function named 
<strong>Faucet</strong>. Is a tool that allows to recharge a wallet  (with COM token).

Getting COM tokens in the testnet is quite easy and free. It is possible through a function called the "<strong>Faucet</strong>." The Faucet is a tool that allows you to recharge a wallet with COM tokens.

A destination address (`addr`) and the amount to be recharged with  (`amount` expressed in ucommercio) must be provided to the faucet endpoint (https://faucet-testnet.commercio.network/give).

**Example** 
Suppose you want to recharge with `10 COM` your wallet address `did:com:1fjqvugs6dfwtax3k4zzh46pswmwryc8ff7f0mv`

This is the request you need to make: 

```
https://faucet-testnet.commercio.network/give?addr=did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr&amount=10000000
``` 

**Limits**

`amount` : There is a limit in the amount value of `100000000 ucommercio`

**Use on Discord**


Faucet is also available in [Discord](https://discord.com/channels/973149882032468029/984721374843121664)


Simply post a message like this 

$request #WALLET_ADDRESS#

Example 

$request did:com:17chk7ldgk99xdxqxszwsvm2ee64rut5dmtuawr



ATTENTION : A `faucet` for CCC is not available. The Buy function (`coming soon`) in the dev.commercio.app must be used  





####  How can I Get CCC  

A faucet for CCC is not available. The easiest way is to use a non hosted wallet 
in the [commercio wallet app](https://github.com/commercionetwork/Commercio-Wallet-App)  and mint CCC using COM obtained from the faucet. Then send CCC to any wallett you need 

**Procedure**

1. Create a wallet in testnet in the [commercio wallet app](https://github.com/commercionetwork/Commercio-Wallet-App) 

2. Check for your wallet address 

3. Send to it some COM throught the faucet  

4. Perform a Mint CCC  function in the [commercio wallet app](https://github.com/commercionetwork/Commercio-Wallet-App)  with the COM obtained 

5. Form the [commercio wallet app](https://github.com/commercionetwork/Commercio-Wallet-App) send the CCC minted to any the address you want in testnet


##  Tools  

### Explorer Almerico

Is an online Web app that enables you to search for real-time and historical information about Commercio blockchain, including data related to blocks, transactions, addresses, and more.


Main-net Almerico:  [mainnet.commercio.network](https://mainnet.commercio.network)

Test-net Almerico: [testnet.commercio.network](https://testnet.commercio.network)

### Keplr wallet extension 

Keplr is a `Chrome` browser extension wallet for the `Cosmos` interchain ecosystem.

Is possible to connect your Keplr wallet to Almerico explorer and perform some action such as delegate, claim, undelegate, redelegate, vote message.

Connecting keplr is easy throught a simple function button. 

<img src="/keplr_connection.png">

More details avaialble  <a href="https://www.keplr.app/#starters" target="_blank">here</a>


### Commercio Wallet App 
Is a Mobile app available in the store that provide a wallet to the user and 
permit to interact with the blockchain for some specific function such as send tokens, delegate, claim, undelegate, redelegate messages

The wallet app can connect both to main-net and test-net

* <a href="https://apps.apple.com/it/app/commerc-io/id1397387586" target="_blank">Apple store IOS Mobile</a> 
* <a href="https://play.google.com/store/apps/details?id=io.commerc.preview.one" target="_blank">Google play store Android Mobile</a> 


### LCD 
Light Client Daemon (LCD REST Server) is a piece of software that connects to a full node to interact with the blockchain

Thus you can also query the Commercio blockchain through LCD rest API  available at specific endpoint 

*  <a href="https://lcd-mainnet.commercio.network/ " target="_blank">Main-net LCD</a>   
*  <a href="https://lcd-testnet.commercio.network/ " target="_blank">Test-net LCD</a>


### Commercio app
Is a hosted wallet platform that permit to interact with the blockchain throught its API :
<a href="./app_developers/commercioapi-introduction.html#the-commercio-app " target="_blank">documentation available here</a>

##  Support
In order to suppor the community feel free to suggest enhancement or 
report bugs opening specific issues on the following Repository 

* [Present documentation](https://github.com/commercionetwork/commercionetwork/issues)
* [Explorer - Almerico](https://github.com/commercionetwork/almerico/issues) 
* [Commercio Wallet App](https://github.com/commercionetwork/Commercio-Wallet-App/issues) 
* [Commercio App](https://github.com/commercionetwork/Commercio-app/issues)

Support could be also asked to the community subscribing into our [Discord](https://discord.com/invite/N7DxaDj5sW) 

Keep in touch with us throught our socialmedia channels available on top menu  


##  Nominal Processing Capacity of commercio.network Blockchain


### Definitions  
  

Before proceeding, it is necessary to clarify some concepts and values explained below in a "simplified" manner related to Cøsmos blockchains.

* Block: A "container" of transactions processed by the chain.
* Transaction: A "container" of messages sent to the chain.
* Message: An "atomic unit" that is recorded once the block is processed in the chain's store (DB).

So, each block processed by the chain nodes can contain multiple transactions, which, in turn, can contain multiple messages.

When discussing on  the Nominal Processing Capacity of a blockchain, non-technical audiences typically refer to TPS (transactions per second). However, in the Cøsmos environment, we refer to MPS (messages per second) rather than TPS.

This distinction is crucial because messages come in various types, such as MsgShareDocument, MsgSendDocumentReceipt, MsgSend, and more, each with different weights. Consequently, the evaluation will vary depending on the specific message or messages chosen for assessment. 

### Premisis  

The estimation of the maximum processing capacity of the chain depends on several factors, with the main ones being:

* Maximum block size (Currently Max 21 Mb, the base value used by Cøsmos chains)
* Block processing time (approximately one every 5/6 seconds)
* Message type (MsgShareDocument, MsgSendDocumentReceipt, MsgSend, etc.)
* Processing capacity of the validator nodes



### 15K TPS 

By acting on the first parameter (increasing the block size), there is a theoretical possibility of achieving values of up to 15K TPS. For example, with a block size of 378 Mb, it would be possible to transmit 90K messages per block, resulting in 3.7 gigabytes per minute and 222 gigabytes per hour. Of course, this requires adjusting the "Processing capacity of the nodes," and these are values that would make sense only when the economic return justifies enhancing the servers on which the nodes reside.



### Analysis

It is not possible to define an exact benchmark, but it is possible to define a theoretical nominal value based on the current values of the relevant parameters (Actual Maximum block size,Block processing time,mean weight of messages, Processing capacity of the nodes) . Taking a single "atomic unit" as one MsgShareDocument, which is the most frequent message type in the network, it can weigh approximately 4 Kb. Within a block (21 Mb), there can be around 5000 messages. Considering that a block is processed every 5/6 seconds:

* In one minute, approximately 50,000 can be estimated.
* In one hour, 3,000,000.
* In one day, 72,000,000.
* In one year, 26,280,000,000.

The nominal estimate is based on a single message in a transaction (4K weight). Two MsgShareDocument messages in the same transaction weigh approximately 7K (3.5k per message), and more messages in a transaction will have a  even little less less per message weight. Other types of messages weigh roughly:

* send 5K 
* receipt 4K 

Estimating a volume of messages of this size at 4K, it should be noted that in a year, the chain's store would weigh approximately 105 terabytes. 

Economic evaluations:

The nominal value would correspond to a mere fee expense to the chain of 262,800,000 CCC (equivalent Euro value since 1 CCC = 1 Euro) . Also, considering any surcharges applied by the commercio app or similar software, assuming they were all executed using Gold Membership, it would result in a yield of 262 million  CCC per year.

* Green: 0.24 Surcharge per Year: 6,307M
* Bronze: 0.11 Surcharge per Year: 2,891M
* Silver: 0.05 Surcharge per Year: 1,314M