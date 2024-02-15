# Becoming a validator
Once you've properly set up a [full node](full-node-installation.md), if you wish you can become a validator node and
start in earning by  validating  the chain transactions. 


## Requirements
If you want to become a Commercio.network validator you need to:

1. Be a full node.  
   If you are not, please follow the [full node installation guide](full-node-installation.md).
   
2. Own enough tokens.  
   To become a validator you need two wallets: one with at least one token to create the validator and another with 50,000 tokens to delegate to the validator node.

:::tip  
If you have any problems with the procedure try to read the section **[Common errors](#_common-errors)**.   
:::


## 1. Add wallet key

:::tip
**`PAY ATTENTION`**     
From now on we will refer to the address of the wallet that will create and to which the validator node will be connected with the notation `<CREATOR ADDRESS>`.
We will refer to the wallet with all tokens (black card in mainnet) as `<DELEGATOR ADDRESS>` notation.   
:::


Inside the testnet or mainnet you can use the **Ledger** with wallet software `commercionetworkd`.    

If you wish to use **Ledger** with `commercionetworkd`, please add the `--ledger` flag to any command.   
You can also use the `commercionetworkd` software without the **Ledger**, and then using the system keyring or an encrypted file-based keyring.

:::warning  
Please remember to copy the 24 words seed phrase in a secure place.  
They are your mnemonic and if you loose them you lose all your tokens and the whole access to your validator.  
:::

:::warning  
Without the use of a ledger (or any other type of hsm that supports Cosmos SDK) or trade wallet app, private keys could be exposed and thus at risk.
:::


To transfer tokens from one wallet to another and to delegate you can use the commercio wallet app

<p style="text-align: center;"><a href="https://apps.apple.com/it/app/commerc-io/id1397387586"><img decoding="async" loading="lazy" class="size-full wp-image-433 alignright" src="https://commercio.network/wp-content/uploads/2016/02/appstore.png" alt="" width="154" height="52"></a>     <a href="https://play.google.com/store/apps/details?id=io.commerc.preview.one"><img decoding="async" loading="lazy" class="size-full wp-image-432 alignleft" src="https://commercio.network/wp-content/uploads/2016/02/play.png" alt="" width="154" height="52"></a></p>

Another possibility is the use of [Keplr](https://medium.com/chainapsis/how-to-use-keplr-wallet-40afc80907f6), an extension for Google Chrome

### Using `commercionetworkd` software

Create the `<CREATOR ADDRESS>` wallet with the following command
```bash
commercionetworkd keys add $NODENAME
# Enter a password that you can remember
```
The output of the command will provide the 24 words that are the mnemonic.    
:::warning  
Please remember to copy the 24 words seed phrase in a secure place.  
:::    


If you are using the **Ledger** device you must first connect it to your computer, start the `commercionetworkd` application and run the command 
```bash
commercionetworkd keys add $NODENAME --ledger
# Enter a password that you can remember
```
In this case the 24 words are not provided because they have already been configured in the **Ledger** initialization.

If you already have the 12 or 24 words use the `--recover` flag and you not use **Ledger**.
```bash
commercionetworkd keys add $NODENAME --recover
# Enter a password that you can remember
```

Copy your public address. It should have the format `did:com:139nvx0ugwxr2ql6ph0azjkkzf5lncq7jgglw8d`.

### Using the `commercio wallet app` (WIP)

Download the `commercio wallet app` and follow the instruction to setup the application.    
Choose "New wallet" and write down the mnemonic produced by the app.    

### Using `Keplr` (WIP)



## 2. Get the tokens

To get tokens inside our testnet you can use [Faucet](https://faucet-testnet.commercio.network). You can call api with `POST` or `GET` method.

To get your tokens inside our mainnet you are required to purchase them using or having received a black card.  
**The black card is the wallet `<DELEGATOR ADDRESS>`**

Read [Add wallet key](#1-add-wallet-key) to create your own `<CREATOR ADDRESS>`.

**Use the ledger or another hsm to make a recovery from 24 words for the `<DELEGATOR ADDRESS>` with the black card.**

If you not use the **Ledger** recover the wallet with your mnemonic and the `--recover` flag.

:::tip
**Note**: You can use `uccc` tokens instead `ucommercio` for the `fees` value
:::

### Sending tokens to `<CREATOR ADDRESS>` with `commercionetworkd`

Send one token to the `<CREATOR ADDRESS>` wallet using the following command

:::warning  
This transaction is expected to be done with an hsm as **Ledger** in mainnet. If you are using a **Ledger** add the `--ledger` flag.
:::

`$CHAINID` is the chain id of the chain. Use `commercio-3` if you work on the mainnet, or the `commercio-testnet11k` if you work on the testnet

```bash
commercionetworkd tx bank send \
<DELEGATOR ADDRESS> \
<CREATOR ADDRESS> \
1110000ucommercio \
--fees=10000ucommercio  \
--chain-id="$CHAINID" \
-y
```
Once you've been confirmed the successful transaction, please check using the following command:
```bash
commercionetworkd query bank balances <CREATOR ADDRESS>
```

The output should look like this **(WIP)**:
```
- denom: ucommercio
  amount: "1110000"
```



### Sending tokens to `<CREATOR ADDRESS>` with `commercio wallet app`


You can use the `Commercio Wallet App`. Please, if you find some bugs [open a issue](https://github.com/commercionetwork/Commercio-Wallet-App).

- Download the app
- Recover your wallet with mnemonic
- Use `Send` button
- Insert `<CREATOR ADDRESS>` address in the "Address To Send" field and 1.11 in the "Tokens To Send" field.
- Press `SEND TOKENS`  button and confirm with password app the transaction

### Sending tokens to `<CREATOR ADDRESS>` with `Keplr` (WIP)



## 3. Create a validator
Once you have the tokens, you can create a validator. If you want, while doing so you can also specify the following parameters
* `--moniker`: the name you want to assign to your validator. Spaces and special characters are accepted (**mandatory**)
* `--details`: a brief description about your node or your company
* `--identity`: your [Keybase](https://keybase.io) identity
* `--website`: a public website of your node or your company

The overall command to create a validator is the following:



### Testnet/Mainnet
If you use private key with file `~/.commercionetwork/config/priv_validator_key.json` use this command

```bash
export VALIDATOR_PUBKEY=$(commercionetworkd tendermint show-validator)
```

If you have a [**kms**](kms-installation.md) you got the value of the public address in the node from the keys registered in your **hsm**. If you have it put that value in the `pubkey` transaction parameter for creating the validator

```bash
export VALIDATOR_PUBKEY="did:com:valconspub1zcjduepq592mn5xucyqvfrvjegruhnx15rruffkrfq0rryu809fzkgwg684qmetxxs"
```



:::warning
:warning: The `<CREATOR ADDRESS>` **account can create one and only one validator** and **a validator can have one and only one creator address**
:::


```bash
commercionetworkd tx staking create-validator \
  --amount=1100000ucommercio \
  --pubkey=$VALIDATOR_PUBKEY \
  --moniker="$NODENAME" \
  --chain-id="$CHAINID" \
  --identity="" \
  --website="" \
  --details="" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=<CREATOR ADDRESS> \
  --fees=10000ucommercio \
  -y
```

**Note**: You can use `uccc` tokens instead `ucommercio` for the `fees` value


The output should look like this:
```
height: 0
txhash: C41B87615308550F867D42BB404B64343CB62D453A69F11302A68B02FAFB557C
codespace: ""
code: 0
data: ""
rawlog: '[]'
logs: []
info: ""
gaswanted: 0
gasused: 0
tx: null
timestamp: ""
```

The key thing is the transaction hash. In the example above it is `C41B87615308550F867D42BB404B64343CB62D453A69F11302A68B02FAFB557C`. Check on the 

- [Commercio.network explorer Testnet](https://testnet.commercio.network/validators/).       
- [Commercio.network explorer Mainnet](https://mainnet.commercio.network/validators/).  


## 4. Delegate tokens to the validator


### Confirm your validator is active
Please confirm that your validator is active by running the following command:

```bash
commercionetworkd query staking validators | fgrep -A 7 "moniker: $NODENAME" | fgrep operator_addres
```


The output should look like this:

```
  operatoraddress: did:com:valoper1zcjx15rruffkrfq0rryu809fzkgwg684qmetxxs
```

Copy the value of `operatoraddress`.


Also verify that the validator is active and that the `operator_addres` is correct by browsing the 


- [Commercio.network explorer Testnet validators](https://testnet.commercio.network/it/validators).       
- [Commercio.network explorer Mainnet validators](https://mainnet.commercio.network/it/validators).       


If you see your validator in the list click on its name.     


:::tip
Congratulations, you are now a Commercio.network validator 🎉
:::


The validator tab should have the value **Operator**. That value is your `operatoraddress`       


Register the value of `operatoraddress`.

```bash
export OPERATOR_ADDRESS="did:com:valoper1zcjx15rruffkrfq0rryu809fzkgwg684qmetxxs"
```


### Delegate tokens with `commercionetworkd` software

You should have already loaded the delegator wallet. If you haven't create the `<DELEGATOR ADDRESS>` on the **Ledger**, or create in keyring through the command

```bash
commercionetworkd keys add <name of second wallet> --recover
```
where `<name of second wallet>` is an arbitrary name.   
When requested, the 24 keywords must be entered


If you are using the **Ledger** device you must first connect it to your computer, start the commercionetworkd application and run the command 
```bash
commercionetworkd keys add <name of second wallet> --ledger
# Enter a password that you can remember
```
In this case the 24 words are not provided because they have already been configured in the **Ledger** initialization




Now you can delegate **50,000 tokens** to the validator node

:::warning  
This transaction is expected to be done with an hsm as a **Ledger** device . If you are using a **Ledger** add the `--ledger` flag.
:::



```bash
commercionetworkd tx staking delegate \
  $OPERATOR_ADDRESS \
  50000000000ucommercio \
  --from <DELEGATOR ADDRESS> \
  --chain-id="$CHAINID" \
  --fees=10000ucommercio \
  -y
```
**Note**: You can use `uccc` tokens instead `ucommercio` for the `fees` value


The output should look like this:
```
height: 0
txhash: 027B85834DA5486085BC56FFD2759443EFD3101BD1023FA9A681262E5C85A845
codespace: ""
code: 0
data: ""
rawlog: '[]'
logs: []
info: ""
gaswanted: 0
gasused: 0
tx: null
timestamp: ""
```

- **Testnet** You should now see your validator with 50001 staked tokens inside the [Commercio.network explorer testnet](https://testnet.commercio.network)
- **Mainnet** You should now see your validator with 50001 staked tokens inside the [Commercio.network explorer mainnet](https://mainnet.commercio.network)


## Note 

If you want to make transactions with the **Nano Ledger** from another machine a full node must be created locally or a full node must be configured to accept remote connections.   
Edit the `.commercionetwork/config/config.toml` file by changing from 

```
laddr = "tcp://127.0.0.1:26657"
```
to
```
laddr = "tcp://0.0.0.0:26657"
```

and restart your node
```
systemctl restart commercionetworkd
```

and use the transaction this way

```bash
commercionetworkd tx staking delegate \
  <validator-addr> \
  50000000000ucommercio \
  --from <DELEGATOR ADDRESS> \
  --node tcp://<ip of your fulle node>:26657 \
  --chain-id="$CHAINID" \
  --fees=10000ucommercio \
  --ledger \
  -y
```

### Delegate tokens with `commercio wallet app` (WIP)


### Delegate tokens with `Keplr` (WIP)


## Common errors

### Account does not exists

#### Problem
If I try to search for my address with the command 

```bash
commercionetworkd query account did:com:1sl4xupdgsgptr2nr7wdtygjp5cw2dr8ncmdsyp --chain-id $CHAINID
```

returns the message
```
ERROR: unknown address: account did:com:1sl4xupdgsgptr2nr7wdtygjp5cw2dr8ncmdsyp does not exist
```
#### Solution

Check if your node has completed the sync.
On https://mainnet.commercio.network or https://testnet.commercio.network you can view the height of the chain at the current state

Use the command 
```
journal -u commercionetworkd -f | egrep " commercionetworkd+.*Committed state"
```
to check the height that reach your node

### Failed validator creation

#### Problem

I executed the validator [creation transaction](#_3-create-a-validator) but I don't appear in validators explorer page.

#### Solution

It may be that by failing one or more transactions the tokens are not sufficient to execute the transaction.

Send more funds to your `<CREATOR ADDRESS>` and repeat the validator creation transaction

### DB errors

#### Problem

Trying to start the rest server or query the chain I get this error
```
panic: Error initializing DB: resource temporarily unavailable
```

#### Solution

Maybe `commercionetworkd` and/or `commercionetworkd` services have been left active.
Use the following commands 

```bash
systemctl stop commercionetworkd
pkill -9 commercionetworkd
```

and repeat the procedure



