# Becoming a validator
Once you've properly set up a [full node](full-node-installation.md), if you wish you can become a validator node and
start earning to validate the chain transactions. 

## Requirements
If you want to become a Commercio.network validator you need to:

1. Be a full node.  
   If you are not, please follow the [full node installation guide](full-node-installation.md).
   
2. Have enough tokens.  
   Currently in order to be a validator you are required to have 50.000 commercio tokens. 
   To get some please read the [*Getting tokens* paragraph](#getting-tokens)

## 1. Add wallet key
Inside the testnet we don't use the Ledger. 
However, if you wish to do so, please add the `--ledger` flat to any command.

:::warning  
Please remember to copy the 24 words seed phrase in a secure place.  
They are your mnemonic and if you loose them you lose all your tokens and the whole access to your validator.  
:::

```bash
cncli keys add $NODENAME
# Enter a password that you can remember
```

Copy your public address. It should have the format `did:com:<data>`.

```bash
cncli keys show $NODENAME --address
```
    
From now on we will refer to the value of your public address using the `<your pub addr>` notation.

## 2. Get the tokens
- [Testnet](#testnet)
- [Mainnet](#mainnet)

### Testnet
In order to receive your tokens on our testnet, please send your public inside our 
[Telegram group](https://t.me/commercionetworkvipsTelegram) requesting them. 
We will make sure to send them to you as soon as possible.

In order to get `<your pub address>` you can use the following command: 
```
cncli keys show $NODENAME --address
```

Once you've been confirmed the successful transaction, please check using the following command:
```bash
cncli query account <your pub addr> --chain-id $CHAINID
```

The output should look like this:
```
- denom: ucommercio
  amount: "52000000000"
```

### Mainnet
To get your tokens inside our mainnet, you are required to purchase them using an exchange.

## 3. Create a validator
Once you have the tokens, you can create a validator. If you want, while doing so you can also specify the following parameters
* `--details`: a brief description about your node or your company
* `--identity`: your [Keybase](https://keybase.io) identity
* `--website`: a public website of your node or your company

The overall command to create a validator is the following:

```bash
cncli tx staking create-validator \
  --amount=50000000000ucommercio \
  --pubkey=$(cnd tendermint show-validator) \
  --moniker="$NODENAME" \
  --chain-id="$CHAINID" \
  --identity="" --website="" --details="" \
  --commission-rate="0.10" --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" --min-self-delegation="1" \
  --from=<your pub addr> \
  -y
```

The output should look like this:
```
rawlog: '[{"msg_index":0,"success":true,"log":""}]'
```

### Confirm your validator is active
Please confirm that your validator is active by running the following command:

```bash
cncli query staking validators --chain-id $CHAINID | fgrep $(cnd tendermint show-validator)
```

You should now see your validator inside the [Commercio.network explorer](https://test.explorer.commercio.network)

:::tip
Congratulations, you are now a Commercio.network validator 🎉
:::