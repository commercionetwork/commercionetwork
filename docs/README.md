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

Sdks are deprecated e no longer mantained. If you want to move forward the develop of the follow packages you can fork them and pull request upgrades.

We have released SDK in 4 main languages 

|  | Dart/Flutter | Kotlin/Java | C#/Dot.net | GoLang |
| ------ | ------ | ------ | ------ | ------ | ------ | ------ | ------ |
| **Sacco**  | [Repo](https://github.com/commercionetwork/sacco.dart) | [Repo](https://github.com/commercionetwork/sacco.kt) | [Repo](https://github.com/commercionetwork/sacco.cs) |  [Repo](https://github.com/commercionetwork/sacco.go) | 
| **CommercioSDK**  | [Repo](https://github.com/commercionetwork/commercio-sdk.dart) | [Repo](https://github.com/commercionetwork/commercio-sdk.kt) | [Repo](https://github.com/commercionetwork/commercio-sdk.cs) | Later |


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

##  Test-net Vs Main-Net

## The Test-Net

[testnet.commercio.network](https://testnet.commercio.network)

Testnet, as the name suggests, is an alternative network for the developers for testing purposes. You can view a testnet as a demo network for experimenting. It’s like the beta stage of a blockchain network. A testnet is a blockchain made available for developers. It allows anyone to conduct experiments without wasting real tokens. A testnet is like a demo network where tokens do not have any value. You can easily test out any app on a testnet because it provides you a sandbox environment separate from the main blockchain.


## The Main-net

[mainnet.commercio.network](https://mainnet.commercio.network)

Mainnet is the complete opposite of the testnet. Mainnet is the main blockchain of Commercio.network. If someone says Commercio mainnet, it means the real Commercio.network blockchain.

Unlike testnet which is an open network for testing purposes, mainnet is the real deal. Tokens on the Commercio.network mainnet have real economic value, be careful.






