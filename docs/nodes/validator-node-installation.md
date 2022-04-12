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
Inside the testnet you can use the **Ledger**, but you can also use the wallet software with the `commercionetworkd`.     
However, if you wish to use **Ledger**, please add the `--ledger` flat to any command.

:::warning  
Please remember to copy the 24 words seed phrase in a secure place.  
They are your mnemonic and if you loose them you lose all your tokens and the whole access to your validator.  
:::


Create the first wallet with the following command
```bash
commercionetworkd keys add $NODENAME
# Enter a password that you can remember
```
The output of the command will provide the 24 words that are the mnemonic.    
      

If you are using the **Ledger** device you must first connect it to your computer, start the commercionetworkd application and run the command 
```bash
commercionetworkd keys add $NODENAME --ledger
# Enter a password that you can remember
```
In this case the 24 words are not provided because they have already been configured in the **Ledger** initialization


Copy your public address. It should have the format `did:com:<data>`.


The second wallet must be requested through a message on the [Telegram group](https://t.me/commercionetworkvipsTelegram). With a private message will be sent the information of the second wallet.

 
**ATTENTION**: from now on we will refer to the value of your public address of the first wallet as `<CREATOR ADDRESS>` notation.
We will refer to the second wallet as `<DELEGATOR ADDRESS>` notation.   

## 2. Get the tokens

To get your tokens inside our mainnet or testnet, you are required to purchase them using an exchange or having received a black card.  
**The black card is the wallet `<DELEGATOR ADDRESS>`**

Read [Add wallet key](#1-add-wallet-key) to create your own `<CREATOR ADDRESS>`.

**Use the ledger or another hsm to make a recovery from 24 words for the second wallet `<DELEGATOR ADDRESS>` with the black card.**

Send one token to the `<CREATOR ADDRESS>` wallet using the following command

:::warning  
This transaction is expected to be done with an hsm as Ledger. If you are using a Ledger add the `--ledger` flag.
:::

`$CHAINID` is the chain id of the chain. Use `commercio-3` if you work on the mainnet, or the `commerico-testnet11k` if you work on the testnet

```bash
commercionetworkd tx bank send \
<DELEGATOR ADDRESS> \
<CREATOR ADDRESS> \
1110000ucommercio \
--fees=10000ucommercio  \
--chain-id="$CHAINID" \
-y
```
or you can use the [Commercio Wallet App](https://github.com/commercionetwork/Commercio-Wallet-App)

**Note**: You can use `uccc` tokens instead `ucommercio` for the `fees` value

Once you've been confirmed the successful transaction, please check using the following command:
```bash
commercionetworkd query bank balances <CREATOR ADDRESS>
```

The output should look like this **(WIP)**:
```
- denom: ucommercio
  amount: "1110000"
```


## 3. Create a validator
Once you have the tokens, you can create a validator. If you want, while doing so you can also specify the following parameters
* `--moniker`: the name you want to assign to your validator. Spaces and special characters are accepted (**mandatory**)
* `--details`: a brief description about your node or your company
* `--identity`: your [Keybase](https://keybase.io) identity
* `--website`: a public website of your node or your company

The overall command to create a validator is the following:



### Testnet
```bash
export VALIDATOR_PUBKEY=$(commercionetworkd tendermint show-validator)
```

### Mainnet
If you have a **kms** you got the value of the public address in the node from the keys registered in your **hsm**. If you have it put that value in the `pubkey` transaction parameter for creating the validator

```bash
export VALIDATOR_PUBKEY="did:com:valconspub1zcjduepq592mn5xucyqvfrvjegruhnx15rruffkrfq0rryu809fzkgwg684qmetxxs"
```

:::warning
Warning: a did address can create one and only one validator and a validator can have one and only one creation address
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


The output should look like this (WIP):
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


**ATTENTION**: Also verify that the validator is active and that the `operator_addres` is correct by browsing the 


[Commercio.network explorer Testnet](https://testnet.commercio.network/it/validators).       
[Commercio.network explorer Mainnet](https://mainnet.commercio.network/it/validators).       


If you see your validator in the list click on its name.     
The validator tab should have the value **Operator**. That value is your `operatoraddress`       


Register the value of `operatoraddress`.

```bash
export OPERATOR_ADDRESS="did:com:valoper1zcjx15rruffkrfq0rryu809fzkgwg684qmetxxs"
```


### Delegate tokens

Once received the second wallet must be loaded on the ledger or keyring through the command

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

**Testnet** You should now see your validator with 50001 staked tokens inside the [Commercio.network explorer testnet](https://testnet.commercio.network)
**Mainnet** You should now see your validator with 50001 staked tokens inside the [Commercio.network explorer mainnet](https://mainnet.commercio.network)

:::tip
Congratulations, you are now a Commercio.network validator 🎉
:::


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



