# Using the app
## Building the application
```bash
# Initialize dep and install dependencies
make get_tools && make get_vendor_deps

# Install the app into your $GOBIN
make install

# Now you should be able to run the following commands:
nsd help
nscli help
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