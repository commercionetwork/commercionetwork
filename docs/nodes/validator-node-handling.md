# Handling a validator
Once you've properly set up a [validator node](validator-node-installation.md), it must be subject to certain rules and certain operations must be carried out to manage it.


## Downtime rules
The node can only stay offline for a certain amount of time.   
In the case of **Commercio Network** this period has been fixed at 10,000 blocks lost, approximately corresponding to 17/18 hours.    
If the node remains offline or fails to produce blocks for a period longer than the limit, it will incur slashing, i.e. an immediate loss of a certain amount of the staked tokens.    
**For "Commercio Network" the slashing percentage for downtime is set at `1%`.**     

## Double sign rules
A validator node must be unique on the chain, so only a node can sign with that private key.   
If there was another node with the same private key, this would result in a double signature, and therefore an immediate jail entry of the node with no exit possibility.   
**For "Commercio Network" the slashing percentage for double sing is set at `5%`.**       
If you run into double signatures all tokens must be unbond and you must create a new validator node with new private keys.    

:::warning
The unbond period is 21 days, so is necessary to wait this period to get back in possession of your tokens.     
:::

## Unjail procedure
In case a validator ended up jail for downtime, it is necessary that the wallet that created the validator performs a ujail transaction.   
The follow command must be performed 

```bash
cncli tx slashing \
  unjail \
  --from <your pub addr creator val> \
  --fees=10000ucommercio  \
  --chain-id="$CHAINID" \
  -y
```

If you are using the **Ledger device** you must first connect it to your computer, start the cosmos application and add `--ledger` flag to the command.


## Unbond procedure
Tokens can be delegated to any validator to increase its stake.      
In case you want to remove all or part of the delegated tokens, an `unbond transaction` must be performed.   
The undelegated period is **21 days**, so is necessary to wait this period to get back in possession of your tokens.     
To perform `unbond transaction` use the follow command

```bash
cncli tx staking \
  unbond \
  <validator-operator-addr> \
  50000000000ucommercio \
  --from <your pub addr delegator> \
  --fees=10000ucommercio  \
  --chain-id="$CHAINID" \
  -y
```

value of `<validator-operator-addr>` can be obtain from explorer:


[Commercio.network explorer Testnet](https://testnet.commercio.network/it/validators).       
[Commercio.network explorer Mainnet](https://mainnet.commercio.network/it/validators).       


If you see your validator in the list click on its name.     
The validator tab should have the value **Operator**. That value is your `<validator-operator-addr>`     
      
If you are using the **Ledger device** you must first connect it to your computer, start the cosmos application and add `--ledger` flag to the command.


## Redelegate procedure
It is possible to perform a procedure of moving tokens in stake from one validator to another through the `redelegate transaction`.     
To perform `redelegate transaction` use the follow command

```bash
cncli tx staking \
  redelegate \
  <source-validator-operator-addr> \
  <destination-validator-operator-addr> \
  50000000000ucommercio \
  --from <your pub addr delegator> \
  --fees=10000ucommercio \
  --chain-id="$CHAINID" \
  -y
```

value of `<source-validator-operator-addr>` and `<destination-validator-operator-addr>` can be obtains from explorer:


[Commercio.network explorer Testnet](https://testnet.commercio.network/it/validators).       
[Commercio.network explorer Mainnet](https://mainnet.commercio.network/it/validators).       


Search your source validator in the list, i.e. the validator you want to move the tokens from, and click on its name.     
The validator tab should have the value **Operator**. That value is your `<source-validator-operator-addr>`     
Return to the list of validators and search your destination validator, i.e. the validator you want to move the tokens to, and click on its name.   
The validator tab should have the value **Operator**. That value is your `<destination-validator-operator-addr>`     



If you are using the **Ledger device** you must first connect it to your computer, start the cosmos application and add `--ledger` flag to the command.
