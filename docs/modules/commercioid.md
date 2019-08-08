# CommercioID
CommercioID is the module that allows you to associate to an existing account (also called *identity*) a 
Did Document Reference.  
 
## Identities
### Creating an identity
First of all, let's define what an **identity** is inside the Commercio Network blockchain.  

> An identity is the method used inside the Commercio Network blockchain in order to identify documents' sender.

In order to create an identity, you simply have to create a Commercio Network address, which will have the 
following form: 

```
did:cmc:<unique part>
```

In order to do so, you can use the CLI and execute the following command: 

```shell
cnd keys add <key-name>
``` 

You will be required to set a password in order to safely store the key on your computer.  

:::warning
Please note that password will be later asked you when signing the transactions so be sure you remember it.
:::  

After inserting the password, you will be shown the mnemonic that can be used in order to import your account 
(and identity) into a wallet. 

<summary>
<details>Identity creation example output</details>

- name: jack
  type: local
  address: did:comnet:1cqewhexuscu5q76qyna4rmnx48flp6dsh9kump
  pubkey: comnetpub1addwnpepqdvvek7t28g6nga954qhvrln379l22l95qag3n6wlpukzrlajvyakwe9v89
  mnemonic: ""
  threshold: 0
  pubkeys: []


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

level awesome panther copy lottery race flag label crunch illegal step atom breeze choose undo solar remove upset caught maze endorse cargo sunny armed
</summary>

#### Using an identity
Once you have created it, in order to start performing a transaction with your identity you firstly have to 
fund your identity. Each and every transaction on the blockchain has a cost, and to pay for it you have to have some 
tokens.  
If you want to receive some tokens, please tell us inside our [official Telegram group](https://t.me/commercionetwork) 
and we will send you some as soon as possible. 