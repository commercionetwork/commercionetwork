# Creating an identity
First of all, let's define what an **identity** is inside the Commercio Network blockchain.  

> An identity is the method used inside the Commercio Network blockchain in order to identify documents' sender.

In order to create an identity, you simply have to create a Commercio Network address, which will have the 
following form: 

```
did:com:<unique part>
```

In order to do so, you can use the CLI and execute the following command: 

```bash
cncli keys add <key-name>
``` 

You will be required to set a password in order to safely store the key on your computer.  

:::warning
Please note that password will be later asked you when signing the transactions so be sure you remember it.
:::  

After inserting the password, you will be shown the mnemonic that can be used in order to import your account 
(and identity) into a wallet. 

```
- name: jack
  type: local
  address: did:com:13jckgxmj3v8jpqdeq8zxwcyhv7gc3dzmrqqger
  pubkey: did:com:pub1addwnpepqfdl6s8hdwdya9zvn5wtx8ty3qsqqqd2ddvygc5zutnrryh5x9ju73jdfg8
  mnemonic: ""
  threshold: 0
  pubkeys: []


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

scorpion what indoor keen topic cricket uphold inch cactus six suffer coin popular honey vendor clown day twin during vague midnight emerge man inform
```

## Using an identity
Once you have created it, in order to start performing a transaction with your identity you firstly have to 
fund your identity. Each and every transaction on the blockchain has a cost, and to pay for it you have to have some 
tokens.  
If you want to receive some tokens, please tell us inside our [official Telegram group](https://t.me/commercionetwork) 
and we will send you some as soon as possible.