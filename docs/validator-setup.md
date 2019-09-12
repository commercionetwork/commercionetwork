# Run a Validator on the Commercio.network Mainnet

::: tip
Information on how to join the mainnet (`genenis.json` file and seeds) is held 
[in our `launch` repo](https://github.com/commercionetwork/launch).
:::

Before setting up your validator node, make sure you've already gone through the 
[Full Node Setup](./join-mainnet.md) guide.

## What is a Validator? 
[Validators](https://cosmos.network/docs/cosmos-hub/validators/overview.html) are responsible for committing new 
blocks to the blockchain through voting. 
A validator's stake is slashed if they become unavailable or sign blocks at the same height. 
Please read about [Sentry Node Architecture](https://cosmos.network/docs/cosmos-hub/validators/validator-faq.html#technical-requirements) 
to protect your node from DDOS attacks and to ensure high-availability.

::: danger 
If you want to become a validator for the Hub's `mainnet`, you should 
[research security](https://cosmos.network/docs/cosmos-hub/validators/security.html)
:::

You may want to skip the next section if you have already [set up a full-node](./join-mainnet.md). 

## Create your validator machine
In order to properly run a validator, some criteria must be satisfied to ensure that your machine won't stop working and 
all your stake won't be slashed for downtime. Three different hardware setup can be found inside the 
[validator hardware](./validator-hardware.md) page.


## Create your validator
Your `comnetvalconspub` can be used to create a new validator by staking tokens. 
You can find your validator pubkey by running:

```bash
cnd tendermint show-validator
```

To create your validator, just use the following command: 

::: warning
Don't use more `ucommercio` than you have
:::

 ```bash
cncli tx staking create-validator \
    --amount=50000000000ucommercio \
    --pubkey=$(cnd tendermint show-validator) \
    --moniker="choose a moniker" \
    --chain-id=<chain_id> \
    --commission-rate="0.10" \
    --commission-max-rate="0.20" \
    --commission-max-change-rate="0.01" \
    --min-self-delegation="1" \
    --gas="auto" \
    --gas-prices="0.025ucommercio" \
    --from=<key_name>
 ```
 
::: tip
When specifying the commission parameters, the `commision-max-change-rate` is used to measure the `% point` 
change over the `commission-rate`. E.g. 1% to 2% is a 100% rate increase but only 1 percentage point.
::: 

::: tip 
`Min-self-delegation` is a strictly positive integer that represents the minimum amount of self-delegated 
voting power your validator must always have. A `min-self-delegation` of 1 means your validator will never have a 
self-delegation lower than `1commercio`, or `1000000ucommercio`.
:::

You can confirm that you are in the validator set by using a third party explorer. 

## Edit the validator description
You can edit your validator's public description. This info is to identify your validator, and will be relied on by 
delegators to decide which validators to stake to.
Make sure to provide input for every flag below. If a flag is not included in the command the field will default to 
empty (`--moniker` defaults to the machine name) if the field has never been set or remain the same if it has been 
set in the past.

The `<key_name>` specifies which validator you are editing. If you choose to not include certain flags, remember that 
the `--from` flag must be included to identify the validator to update.

The `--identity` can be used as to verify identity with systems like Keybase or UPort. When using with 
Keybase `--identity` should be populated with a 16-digit string that is generated with a 
[keybase.io](https://keybase.io/) account. It's a cryptographically secure method of verifying your identity across 
multiple online networks. The Keybase API allows us to retrieve your Keybase avatar. 
This is how you can add a logo to your validator profile.

```bash
cncli tx staking edit-validator
  --moniker="choose a moniker" \
  --website="https://commercio.network" \
  --identity=6A0D65E29A4CBC8E \
  --details="The documents blockchain!" \
  --chain-id=<chain_id> \
  --gas="auto" \
  --gas-prices="0.025ucommercio" \
  --from=<key_name> \
  --commission-rate="0.10"
```

**Note**: The `commision-rate` value must adhere to the following inviariants: 

* Must be between 0 and the validator's `commission-max-rate`
* Must not exceed the validator's `commission-max-change-rate` which is maximum % point change rate per day.   
  In other words, a validator can only change its commission once per day and within `commission-max-change-rate` bounds.

## View the validator description
View the validator's information with this command:

```bash
cncli query staking validator <account_commercio_network>
```

## Track the Validator Signing Information
In order to keep track of a validator's signatures in the past you can do so by using the signing-info command:

```bash
cncli query slashing signing-info <validator-pubkey>\
  --chain-id=<chain_id>
```

## Unjail a validator
When a validator is "jailed" for downtime, you must submit an `Unjail` transaction from the operator account in order 
to be able to get block proposer rewards again (depends on the zone fee distribution).
 
```bash
cncli tx slashing unjail \
	--from=<key_name> \
	--chain-id=<chain_id>
```

## Confirm your validator is running
Your validator is active if the following command returns anything:

```bash
cncli query tendermint-validator-set | grep "$(cnd tendermint show-validator)"
```

You should now see your validator in one of the Cosmos Hub explorers. 
You are looking for the Bech32 encoded address in the `~/.cnd/config/priv_validator.json` file.

> NOTE. To be in the validator set, you need to have more total voting power than the 100th validator.
         
## Common problems
### Problem #1: My validator has `voting_power: 0`
Your validator has become jailed. Validators get jailed, i.e. get removed from the active validator set, 
if they do not vote on `500` of the last `10000` blocks, or if they double sign.

If you got jailed for downtime, you can get your voting power back to your validator. 
First, if `cnd` is not running, start it up again:

```bash
cnd start
```

Wait for your full node to catch up to the latest block. Then, you can [unjail your validator](#unjail-a-validator).

Lastly, check your validator again to see if your voting power is back.

```bash
cncli status
```

You may notice that your voting power is less than it used to be. That's because you got slashed for downtime!

### Problem #2: My `cnd` crashes because of `too many open files`
The default number of files Linux can open (per-process) is `1024`. 
`cnd` is known to open more than `1024` files. This causes the process to crash. 
A quick fix is to run `ulimit -n 4096` (increase the number of open files allowed) and then restart the process with 
`cnd start`. If you are using `systemd` or another process manager to launch `cnd` this may require some 
configuration at that level. A sample `systemd` file to fix this issue is below:

```
# /etc/systemd/system/cnd.service
[Unit]
Description=Commercio.network Node
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu
ExecStart=/home/ubuntu/go/bin/cnd start
Restart=on-failure
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
```