# Becoming a validator (**WIP**)
Once you've properly set up a [full node](full-node-installation.md), if you wish you can become a validator node and
start in earning by  validating  the chain transactions. 

Before you start, we recommend that you run the command 

```bash
cncli config chain-id $CHAINID
```

In this way you can omit the flag `--chain-id="$CHAINID"` in every command of the **cncli**


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
Inside the testnet you can use the ledger, but you can also use the wallet software with the `cncli`.     
However, if you wish to use Ledger, please add the `--ledger` flat to any command.

:::warning  
Please remember to copy the 24 words seed phrase in a secure place.  
They are your mnemonic and if you loose them you lose all your tokens and the whole access to your validator.  
:::


Create the first wallet with the following command
```bash
cncli keys add $NODENAME
# Enter a password that you can remember
```

Copy your public address. It should have the format `did:com:<data>`.


The second wallet must be requested through a message on the [Telegram group](https://t.me/commercionetworkvipsTelegram). With a private message will be sent the information of the second wallet.

    
From now on we will refer to the value of your public address of the first wallet as `<your pub addr creator val>` notation.
We will refer to the second wallet as `<your pub addr delegator>` notation.   

## 2. Get the tokens
- [Testnet](#testnet)
- [Mainnet](#mainnet)

### Testnet
In order to get  `<your pub addr creator val>` you can use the following command: 
```bash
cncli keys show $NODENAME --address
```

From the command line
```bash
curl "https://faucet-testnet.commercio.network/invite?addr=<your pub addr creator val>"
```

Or on a browser copy and paste the following address
```
https://faucet-testnet.commercio.network/invite?addr=<your pub addr creator val>
```

The call should return something like

```json
{"tx_hash":"4AB05DF5BEB7321059A6724BF18A7B95631AB55773BBD55DFC448351101BE972"}
```

Now you can request the token from faucet's service

```bash
cncli keys show $NODENAME --address
```

From the command line
```bash
curl "https://faucet-testnet.commercio.network/give?addr=<your pub addr creator val>&amount=1100000"
```

Or on a browser copy and paste the following address
```
https://faucet-testnet.commercio.network/give?addr=<your pub addr creator val>&amount=1100000
```

The call should return something like

```json
{"tx_hash":"BB733FDB2665265D3B3A32576F23B10B10EA8F56EEBAD08C1BF39D5E2FAC601C"}
```



Once you've been confirmed the successful transaction, please check using the following command:
```bash
cncli query account <your pub addr creator val> --chain-id $CHAINID
```

The output should look like this:
```
- denom: ucommercio
  amount: "1100000"
```

When you have received the second wallet `<your pub addr delegator>` via telegram message, check if the tokens are actually present
```bash
cncli query account <your pub addr delegator> --chain-id $CHAINID
```
The output should look like this:
```
- denom: ucommercio
  amount: "5100000000"
```



### Mainnet
To get your tokens inside our mainnet, you are required to purchase them using an exchange or having received a black card.  
**The black card is the wallet `<your pub addr delegator>`**


Create the first wallet `<your pub addr creator val>` with the command 
```bash
cncli keys add $NODENAME
# Enter a password that you can remember
```

**Use the ledger or another hsm to make a recovery from 24 words for the second wallet with the black card.**

Send one token to the first wallet using the following command

:::warning  
This transaction is expected to be done with an hsm as Ledger. If you are using a Ledger add the `--ledger` flag.
:::

```bash
cncli tx send \
<your pub addr delegator> \
<your pub addr creator val> \
1000000ucommercio \
--fees=10000ucommercio  \
--chain-id="$CHAINID" \
-y
```

Once you've been confirmed the successful transaction, please check using the following command:
```bash
cncli query account <your pub addr creator val> --chain-id $CHAINID
```

The output should look like this:
```
- denom: ucommercio
  amount: "1000000"
```


## 3. Create a validator
Once you have the tokens, you can create a validator. If you want, while doing so you can also specify the following parameters
* `--details`: a brief description about your node or your company
* `--identity`: your [Keybase](https://keybase.io) identity
* `--website`: a public website of your node or your company

The overall command to create a validator is the following:



### Testnet
```bash
export VALIDATOR_PUBKEY=$(cnd tendermint show-validator)
```

### Mainnet
If you have a **kms** you got the value of the public address in the node from the keys registered in your **hsm**. If you have it put that value in the `pubkey`

```bash
export VALIDATOR_PUBKEY="did:com:valconspub1zcjduepq592mn5xucyqvfrvjegruhnx15rruffkrfq0rryu809fzkgwg684qmetxxs"
```



```bash
cncli tx staking create-validator \
  --amount=1000000ucommercio \
  --pubkey=$VALIDATOR_PUBKEY \
  --moniker="$NODENAME" \
  --chain-id="$CHAINID" \
  --identity="" --website="" --details="" \
  --commission-rate="0.10" --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" --min-self-delegation="1" \
  --from=<your pub addr creator val> \
  --fees=10000ucommercio \
  -y
##Twice input password required
```


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

## 4. Delegate tokens to the validator


### Confirm your validator is active
Please confirm that your validator is active by running the following command:

```bash
cncli query staking validators --chain-id $CHAINID | fgrep -B 1 $VALIDATOR_PUBKEY
```
Something like this

```
  operatoraddress: did:com:valoper1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  conspubkey: did:com:valconspub1zcjduepq592mn5xucyqvfrvjegruhnx15rruffkrfq0rryu809fzkgwg684qmetxxs
```

Copy the value of `operatoraddress`. Below we will refer to this value with `<validator-addr>`

### Delegate tokens

Once received the second wallet must be loaded on the ledger or keyring through the command

```bash
cncli keys add <name of second wallet> --recover
```
where `<name of second wallet>` is an arbitrary name.   
When requested, the 24 keywords must be entered


Now you can delegate the tokens to the validator node

:::warning  
This transaction is expected to be done with an hsm as a Ledger device . If you are using a Ledger add the `--ledger` flag.
:::



```bash
cncli tx staking delegate \
 <validator-addr> \
 50000000000ucommercio \
 --from <your pub addr delegator> \
 --chain-id="$CHAINID" \
 --fees=10000ucommercio \
 -y
```

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

**Testnet** You should now see your validator inside the [Commercio.network explorer testnet](https://testnet.commercio.network)
**Mainnet** You should now see your validator inside the [Commercio.network explorer mainnet](https://mainnet.commercio.network)

:::tip
Congratulations, you are now a Commercio.network validator 🎉
:::


## Note 

If you want to make transactions with the nano ledger from another machine a full node must be created locally or a full node must be configured to accept remote connections.   
Edit the `.cnd/config/config.toml` file by changing from 

```
laddr = "tcp://127.0.0.1:26657"
```
to
```
laddr = "tcp://0.0.0.0:26657"
```

and restart your node
```
systemctl cnd restart
```

and use the transaction this way

```bash
cncli tx staking delegate \
 <validator-addr> \
 50000000000ucommercio \
 --from <your pub addr delegator> \
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
cncli query account did:com:1sl4xupdgsgptr2nr7wdtygjp5cw2dr8ncmdsyp --chain-id $CHAINID
```

returns the message
```
ERROR: unknown address: account did:com:1sl4xupdgsgptr2nr7wdtygjp5cw2dr8ncmdsyp does not exist
```
#### Solution

Check if your node has completed the sync.
On https://testnet.commercio.network you can view the height of the chain at the current state

Use the command 
```
tail -1f /var/log/syslog | egrep " cnd+.*Committed state"
```
to check the height that reach your node

### Failed validator creation

#### Problem

I executed the validator [creation transaction](#_3-create-a-validator) but I don't appear on https://testnet.commercio.network/it/validators.

#### Solution

It may be that by failing one or more transactions the tokens are not sufficient to execute the transaction.

Request more funds from the faucet with the command

```bash
curl "https://faucet-testnet.commercio.network/give?addr=<your pub addr creator val>&amount=1100000"
```
and repeat the validator creation transaction

### DB errors

#### Problem

Trying to start the rest server or query the chain I get this error
```
panic: Error initializing DB: resource temporarily unavailable
```

#### Solution

Maybe `cnd` and/or `cncli` services have been left active.
Use the following commands 

```bash
systemctl stop cnd
systemctl stop cncli
pkill -9 cnd
pkill -9 cncli
```

and repeat the procedure



