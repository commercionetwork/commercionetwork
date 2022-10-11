# Handling a validator
Once you've properly set up a [validator node](validator-node-installation.md), it must be subject to certain rules and certain operations must be carried out to manage it.



- [Handling a validator](#handling-a-validator)
  - [Downtime rules](#downtime-rules)
  - [Double sign rules](#double-sign-rules)
  - [Unjail procedure](#unjail-procedure)
  - [Unbond procedure](#unbond-procedure)
  - [Redelegate procedure](#redelegate-procedure)
  - [Move node to other server](#move-node-to-other-server)
    - [Move validator to another server with `priv_validator_key.json` file](#move-validator-to-another-server-with-priv_validator_keyjson-file)
    - [Move validator to another server with `kms`](#move-validator-to-another-server-with-kms)
    - [Resume validator after break down with `priv_validator_key.json` file](#resume-validator-after-break-down-with-priv_validator_keyjson-file)
    - [Resume validator after break down with `kms`](#resume-validator-after-break-down-with-kms)
  - [x% Loss of blocks](#x-loss-of-blocks)
  - [Add identity to your validator](#add-identity-to-your-validator)
    - [References keybase](#references-keybase)
    - [1. Registration](#1-registration)
    - [2. Adding one or more devices](#2-adding-one-or-more-devices)
      - [2.1. Program installation on computers](#21-program-installation-on-computers)
      - [2.2. App installation on iPhone and Android](#22-app-installation-on-iphone-and-android)
    - [3. Verify email](#3-verify-email)
    - [4. Add logo and key](#4-add-logo-and-key)
    - [5. Edit your validator adding identity](#5-edit-your-validator-adding-identity)

## Downtime rules
The node can only stay offline for a certain amount of time.   
In the case of **Commercio Network** this period has been fixed at 10,000 blocks lost, approximately corresponding to 17/18 hours. 
Validator must validate at least 5% of the 10,000 blocks.
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
commercionetworkd tx slashing \
  unjail \
  --from <your pub addr creator val> \
  --fees=10000ucommercio  \
  --chain-id="$CHAINID" \
  -y
```
**Note**: You can use `uccc` tokens instead `ucommercio` for the `fees` value

If you are using the **Ledger device** you must first connect it to your computer, start the cosmos application and add `--ledger` flag to the command.


## Unbond procedure
Tokens can be delegated to any validator to increase its stake.      
In case you want to remove all or part of the delegated tokens, an `unbond transaction` must be performed.   
The undelegated period is **21 days**, so is necessary to wait this period to get back in possession of your tokens.     
To perform `unbond transaction` use the follow command

```bash
commercionetworkd tx staking \
  unbond \
  <validator-operator-addr> \
  50000000000ucommercio \
  --from <your pub addr delegator> \
  --fees=10000ucommercio  \
  --chain-id="$CHAINID" \
  -y
```

**Note**: You can use `uccc` tokens instead `ucommercio` for the `fees` value


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
commercionetworkd tx staking \
  redelegate \
  <source-validator-operator-addr> \
  <destination-validator-operator-addr> \
  50000000000ucommercio \
  --from <your pub addr delegator> \
  --fees=10000ucommercio \
  --chain-id="$CHAINID" \
  -y
```

**Note**: You can use `uccc` tokens instead `ucommercio` for the `fees` value


value of `<source-validator-operator-addr>` and `<destination-validator-operator-addr>` can be obtains from explorer:


[Commercio.network explorer Testnet](https://testnet.commercio.network/it/validators).       
[Commercio.network explorer Mainnet](https://mainnet.commercio.network/it/validators).       


Search your source validator in the list, i.e. the validator you want to move the tokens from, and click on its name.     
The validator tab should have the value **Operator**. That value is your `<source-validator-operator-addr>`     
Return to the list of validators and search your destination validator, i.e. the validator you want to move the tokens to, and click on its name.   
The validator tab should have the value **Operator**. That value is your `<destination-validator-operator-addr>`     



If you are using the **Ledger device** you must first connect it to your computer, start the cosmos application and add `--ledger` flag to the command.


## Move node to other server

If you need to move your validator to another server, the only thing that you need to move is your private key.     
Your node structure should be something like below

```bash
.commercionetwork
├── config
│   └── app.toml
│   └── config.toml
│   └── genesis.json
│   └── node_id.json
│   └── priv_validator_key.json
├── data
│   └── priv_validator_state.json
└── cosmovisor
    └── current
    └── genesis
    └── bin
    │   └── commercionetword
    └── upgrades
    └── <name>
        └── bin
            └── commercionetword

```
If you don't use `kms` the private key of your validator is saved in `priv_validator_key.json` file.     

### Move validator to another server with `priv_validator_key.json` file

1. Install a full node in the new server
2. Stop the full node `commercionetword` service in the new server 
   ```bash
   systemctl stop commercionetword
   ```
3. Copy `data` folder to the new server
   ```bash
   rsync -av --delete ~/.commercionetwork/data/ <USER_NEW_SERVER>@<IP_NEW_SERVER>:.commercionetwork/data/
   ```
4. Stop and disable the `commercionetword` service in the current server of validator node
   ```bash
   systemctl stop commercionetword
   systemctl disable commercionetword
   ```
1. Sync again your new node `data` folder
   ```bash
   rsync -av --delete ~/.commercionetwork/data/ <USER_NEW_SERVER>@<IP_NEW_SERVER>:.commercionetwork/data/
   ```
2. Copy your `priv_validator_key.json` in the new server
   ```bash
   scp ~/.commercionetwork/config/priv_validator_key.json <USER_NEW_SERVER>@<IP_NEW_SERVER>:.commercionetwork/config/priv_validator_key.json
   ```
3. If you have some special setup in your `config.toml` and `app.toml` copy that in your new node.
4. Remove the `priv_validator_key.json` file from old server
   ```bash
   rm ~/.commercionetwork/config/priv_validator_key.json
   ```
1. Restart the node in your new server
   ```bash
   systemctl start commercionetword
   ```
1. Verify if your validator signs again
### Move validator to another server with `kms`

1. Install a full node in the new server
2. Stop the full node `commercionetword` service in the new server 
   ```bash
   systemctl stop commercionetword
   ```
3. Copy `data` folder to the new server
   ```bash
   rsync -av --delete ~/.commercionetwork/data/ <USER_NEW_SERVER>@<IP_NEW_SERVER>:.commercionetwork/data/
   ```
1. If you have some special setup in your `config.toml` and `app.toml` copy that in your new node. Especially setup your `priv_val_addr` in `config.toml` using the setup of your server. 
1. Enter in your `kms` server and stop the `tmkms` service
   ```bash
   systemctl stop tmkms
   ```
1. Modify `tmkms` config using `priv_val_addr` value of your validator 
   ```toml
   [[validator]]
   addr = "tcp://<VALIDATOR_IP>:26658"
   ``` 
1. Restart `tmkms` service
   ```bash
   systemctl start tmkms
   ```
1. Restart the node in your new server
   ```bash
   systemctl start commercionetword
   ```
1. Verify if your validator signs again
### Resume validator after break down with `priv_validator_key.json` file

To resume a validator after break down or some other terrible issue that destroy your server you need to have the `priv_validator_key.json` saved in a secure place.

1. Install a full node in the new server
2. After the node is synced with the chain stop it
   ```bash
   systemctl stop commercionetword
   ```
3. Copy your saved `priv_validator_key.json` file in the new server
4. Restart the node in your new server
   ```bash
   systemctl start commercionetword
   ```
1. Verify if your validator signs again

### Resume validator after break down with `kms`
1. Install a full node in the new server
2. After the node is synced with the chain stop it
   ```bash
   systemctl stop commercionetword
   ```
1. Setup your `priv_val_addr` in `config.toml` using the setup of your server. 
1. Enter in your `kms` server and stop the `tmkms` service
   ```bash
   systemctl stop tmkms
   ```
1. Modify `tmkms` config using `priv_val_addr` value of your validator 
   ```toml
   [[validator]]
   addr = "tcp://<VALIDATOR_IP>:26658"
   ``` 
1. Restart `tmkms` service
   ```bash
   systemctl start tmkms
   ```
1. Restart the node in your new server
   ```bash
   systemctl start commercionetword
   ```
1. Verify if your validator signs again
## x% Loss of blocks

 What does an x% loss of blocks mean for a validator? 

 In order to create blocks and extend the blockchain, active validators are selected in proportion to their stake by a pseudo-random mechanism that assigns at a time t a proposer (block proposer) and a fixed pool of validators (committee). The block proposer is assigned the task of proposing the block Bt (block at time t) while the committee is given the task of voting whether the block is hung on the chain or not.

 Losing a block for a validator means that this validator was inactive at the time of the committee's choice or that it did not vote towards the execution of the block. Thus, having a validator with an x% loss of blocks means that the validator is only active (100 - x)% of the time when the Bt block is assigned.


## Add identity to your validator

In order to publish the validator icon to any explorer, an identity must be created at [keybase](https://keybase.io).
Below are the steps to be performed 



### References keybase

- [Keybase book](https://book.keybase.io/account)
- [Keybase API](https://keybase.io/docs/api/1.0/intro)

### 1. Registration

In the [Keybase.io](https://keybase.io/) site click on `Login`; a popup opens where there is a `Join Keybase` button. After pressing this button fill in the form fields and submit the request.

### 2. Adding one or more devices


Install Keybase on a device, start the program and authenticate with the credentials entered during registration. The application displays a prompt to add the device in use.

#### 2.1. Program installation on computers

- [macOS](https://keybase.io/docs/the_app/install_macos)
- [Linux](https://keybase.io/docs/the_app/install_linux)
- [Windows](https://keybase.io/docs/the_app/install_windows)

#### 2.2. App installation on iPhone and Android

- [iPhone/iPad](https://apps.apple.com/us/app/keybase-crypto-for-everyone/id1044461770)
- [Android](https://play.google.com/store/apps/details?id=io.keybase.ossifrage&referrer=undefined)

### 3. Verify email

After authenticating, request verification of the email address.

### 4. Add logo and key

In the `people` section of the program, click on the `Edit Profile` button to complete the profile (data and image) and to access the proof list.

From the proof list, select `Add a PGP key`, then add a locally generated PGP public key or create a new one with the `Get a new PGP key` function provided by the application. Generation via the Keybase application is done automatically, the only actions required are:

- Fill out the form (`full name` and `email`);
- At the end of generation decide to save or not to Keybase the private key by encrypting it.

After adding the key you will be able to access your account image using the `https://keybase.io/_/api/1.0/user/lookup.json?key_suffix=[id]&fields=pictures` endpoint, where `[id]` consists of the first 16 hexadecimal digits of the PGP key.

### 5. Edit your validator adding identity

In order to add the identity to your validator you need to edit the validator profile with a transaction.      
**Warning: put the values of `moniker`, `website` and `details` according with your environment**

```bash
IDENTITY="[id]" #ID that you obtained in the previous step. Something like 21C53B7B20C1145D
commercionetworkd tx staking edit-validator \
  --moniker="$NODENAME" \
  --identity="$IDENTITY" \
  --website="your website (leave blank if you don't have it)" \
  --details="description of your node (leave blank if you don't have it)" \
  --chain-id="commercio-3" \
  --from=<CREATOR ADDRESS> \
  --fees=10000ucommercio \
  -y
```


*N.B.: after editing the profile the `Edit Profile` button will disappear, in its place will be the `Prove your identies` button until the identity is proven*

