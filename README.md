# Commercio.network Cosmos Blockchain application
## Downloading the app
### Requirements
In order to be sure everything works properly, due to [Go](https://golang.org/) restrictions, the following requirements
must be matched.

1. You must have Go 11.2+ installed. The download is available [here](https://golang.org/dl/)
2. You must have a valid `GOPATH` environment variable set
3. You must have a valid `GOBIN` environment variable set

### Installation
The installation process is composed of the following steps
1. Creating a `src` folder inside the `GOPATH` folder.
2. Creating a `commercio-network` folder inside the `src` folder. 
3. Cloning the project inside the `commercio-network` folder. 
4. Installing all the tools. 
5. Updating the dependencies and installing the app into `GOBIN`

The following commands must be issued:
```bash
# Create a src folder inside GOPATH
mkdir $GOPATH/src && cd $GOPATH/src

# Create the commercio-network folder
mkdir commercio-network && cd commercio-network

# Clone the project
git clone https://scw-gitlab.zotsell.com/Commercio.network/Cosmos-application .

# Install all the tools necessary
make get_tools && dep init -v

# Update the dependencies and install the app 
dep ensure -update -v && make install
```

After all of this, you should be able to run the following commands successfully 
```bash
cnd help
cncli help
```

## Running the live network
```bash
# Initialize configuration files and genesis file
cnd init testchain

# Copy the `Address` output here and save it for later use
cncli keys add jack

# Copy the `Address` output here and save it for later use
cncli keys add alice

# Add both accounts, with coins to the genesis file
cnd add-genesis-account $(cncli keys show jack --address) 100000000stake,1000jackcoin
cnd add-genesis-account $(cncli keys show alice --address) 100000000stake,1000alicecoin

# Create the genesis transaction signing it with the jack private key
cnd gentx --name jack

# Collect all the genesis transactions
cnd collect-gentxs
```

Once you've done this, you are ready to start using the Command Line Interface (CLI) to interact with the blockchain.  
However, to make things easier, you should follow the following steps in order to set some useful default 
configurations.

```bash
# 1. Set the proper CLI output format and indentation
cncli config output json
cncli config indent true

# 2. Tell the CLI we are running on a trusted node
cncli config trust-node true

# 3. Set the default chain id
# First of all, read the generated chain id value
cat ~/.cnd/config/genesis.json

# Copy the value of the `chain_id` field, then execute the following
cnd config chain-id <CHAIN-ID>
# Example: cnd config chain-id test-chain-RKFXWR
``` 

Now, you can start using the commands to interact with the network.

## Using the Command Line Interface
Before you can start using the CLI, you need to effectively start the blockchain.  
To do so, just run
```bash
cnd start
``` 

### CommercioAUTH
#### Creating an account
Create an account specifying the address, the public key type and its value. 
Note: 
* the address value must be an hex value, **without** the `0x` prefix
* the key type value must be either `Secp256k1` or `Ed25519`
* the key value must be an hex value, **without** the `0x` prefix
```bash
cncli tx commercioauth register \
    0b13c55fc6c3496796258d8637330ff7e269cac8 \
    Ed25519 \
    917325912fbdbe78a5d6c92678d1d89f09a89c0ca35a856688a13f6e7ef5e951 \
    --from $(cncli keys show jack --address)
```

#### Reading the accounts' details
##### Read all the accounts
List all the registered accounts
```bash
cncli query commercioauth account list
```

##### Read the details of a specific account
Retrieve the details of an account based on his address expressed as an hex value **without** the `0x` prefix
```bash
cncli query commercioauth account 0b13c55fc6c3496796258d8637330ff7e269cac8
```

### CommercioID
#### Managing an identity
##### Creating an identity
Create an identity specifying the DID and the DDO reference.  
The first parameter is the DID, the second is the DDO reference.
```bash
cncli tx commercioid upsert-identity \
    0xa971c43e6c26c01e744a57db57cf9982b2e195ba \
    QmWCnEEqSaBcKtKLUurzi2Zs9LAPxJkpzE8as21fvmeXrj \
    --from $(cncli keys show jack --address)
```

##### Resolving an identity
Read the identity details by its DID.
```bash
cncli query commercioid resolve \
    0xa971c43e6c26c01e744a57db57cf9982b2e195ba
``` 

#### Managing a connection
##### Creating a connection
Create a connection specifying the two users that needs to connect to each other.\
The first parameter is the first DID, the second parameter is the second DID.
```bash 
cncli tx commercioid create-connection \
    0xa971c43e6c26c01e744a57db57cf9982b2e195ba \
    0x9f2ae6af2545076e7a55816dd4f8e45b650b07f0 \
    --from $(cncli keys show jack --address)
```

##### Reading all the connections of a user
List all the connections that a user has established.
The required parameter is the DID for which to retrieve the connections
```bash
cncli query commercioid connections \
    0xa971c43e6c26c01e744a57db57cf9982b2e195ba
```

### CommercioDOCS
#### Managing a document
##### Storing a document
Store a document with an associated metadata inside the blockchain. The params are: 
1. The identity that will be set as the creator of the document.  
   This identity must be owned by the user that signs the transaction, or must be free. 
2. The reference to the document
3. The reference to the document metadata

```bash
cncli tx commerciodocs store \
    0xa971c43e6c26c01e744a57db57cf9982b2e195ba \
    QmPMcMY6VAkLfDVZwrGYE48bFRbEovAHkRiJ7t7Lp7qD3n \
    QmRCxjZrUQ29aYNcrmdtnzDeB1GAf56tpdwpymgZhf2ifp \
    --from $(cncli keys show jack --address)
```

##### Read the metadata of a document
Retrieve the metadata of a document by its reference.
```bash
cncli query commerciodocs metadata \
    QmPMcMY6VAkLfDVZwrGYE48bFRbEovAHkRiJ7t7Lp7qD3n
```

#### Sharing a document
##### Share a document
Share a document with a given user specifying his DID. The params are:
1. The reference to the document. 
2. The sender identity (as a DID).  
   This identity must be owned by the user that signs the transaction, or must be free. 
3. The receiver identity (as a DID).

```bash
cncli tx commerciodocs share \
    QmPMcMY6VAkLfDVZwrGYE48bFRbEovAHkRiJ7t7Lp7qD3n \
    0xa971c43e6c26c01e744a57db57cf9982b2e195ba \
    0x9f2ae6af2545076e7a55816dd4f8e45b650b07f0 \
    --from $(cncli keys show jack --address)
```

##### Read the authorized users
List all the users that are authorized to read a document. These include the creator of the document and all the users 
with whom the creator has shared the document.  
The only parameter required is the reference of the document.

```bash
cncli query commerciodocs readers \
    QmPMcMY6VAkLfDVZwrGYE48bFRbEovAHkRiJ7t7Lp7qD3n
```


## Using the REST APIs
In order to use the REST APIs, you firstly need to run `cnd start` you did not run this previously.  
After that, you need to gather the address of the registered user (in this case, Jack). To do so, run
```bash
cncli keys show jack --address
```

Now, start the REST server by running
```bash
cncli rest-server --laddr=tcp://0.0.0.0:1317
```

Now, with the previously output address, run the following in other terminal shell:
```bash
curl -s https://localhost:1317/auth/accounts/${address}
```

Example:
```bash
$ curl -s -k https://localhost:1317/auth/accounts/cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0
{
  "type": "auth/Account",
  "value": {
    "address": "cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0",
    "coins": [
      {
        "denom": "jackCoin",
        "amount": "1000"
      },
      {
        "denom": "mycoin",
        "amount": "1000"
      }
    ],
    "public_key": {
      "type": "tendermint/PubKeySecp256k1",
      "value": "A8dJWr6t9Yh31YYvXkb0N/HtkC5J+KAP75dqg8pr3uws"
    },
    "account_number": "0",
    "sequence": "5"
  }
}
```

We will need those data later, specially the `address`, `account_nummber` and `sequence` values.

### Requests
All the requests examples can be found inside the [documentation](https://documenter.getpostman.com/view/3509480/Rzfnk6c7). 
