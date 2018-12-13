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
$ nsd help
$ nscli help
```

## Running the live network and using the commands
```bash
# Initialize configuration files and genesis file
nsd init --chain-id testchain

# Copy the `Address` output here and save it for later use
nscli keys add jack

# Copy the `Address` output here and save it for later use
nscli keys add alice

# Add both accounts, with coins to the genesis file
nsd add-genesis-account $(nscli keys show jack --address) 1000mycoin,1000jackCoin
nsd add-genesis-account $(nscli keys show alice --address) 1000mycoin,1000aliceCoin

# Start the blockchain 
nsd start
```

Now, you can start using the commands to interact with the network.

```bash
# First check the accounts to ensure they have funds
nscli query account $(nscli keys show jack --address) \
    --indent --chain-id testchain --trust-node=true
    
nscli query account $(nscli keys show alice --address) \
    --indent --chain-id testchain --trust-node=true

# Buy your first name using your coins from the genesis file
nscli tx nameservice buy-name jack.id 5mycoin \
    --from     $(nscli keys show jack --address) \
    --chain-id testchain --trust-node=true

# Set the value for the name you just bought
nscli tx nameservice set-name jack.id 8.8.8.8 \
    --from     $(nscli keys show jack --address) \
    --chain-id testchain --trust-node=true

# Try out a resolve query against the name you registered
nscli query nameservice resolve jack.id --chain-id testchain --trust-node=true
# > 8.8.8.8

# Try out a whois query against the name you just registered
nscli query nameservice whois jack.id --chain-id testchain --trust-node=true
# > {"value":"8.8.8.8","owner":"cosmos1l7k5tdt2qam0zecxrx78yuw447ga54dsmtpk2s","price":[{"denom":"mycoin","amount":"5"}]}

# Alice buys name from jack
nscli tx nameservice buy-name jack.id 10mycoin \
    --from     $(nscli keys show alice --address) \
    --chain-id testchain --trust-node=true
```