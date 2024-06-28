#!/bin/bash

printf "\n\n🔧 Installing Hermes relayer...\n"

LOGS_DIR="/home/cheikh/Lavoro/commercionetwork/logs"
COMMERCIO_CHAIN_ID="commercionetwork"
COMMERCIO_KEYRING="--keyring-backend test"
COMMERCIO2_CHAIN_ID="commercio2"
COMMERCIO2_KEYRING="--keyring-backend test"
HERMES_CONFIG_DIR="/home/cheikh/.hermes"



# Validate Hermes configuration
echo -e "\n⏳ Waiting to validate Hermes configuration...\n"
hermes config validate

# Create Hermes wallet in both chains
echo -e "\n🔧 Creating Hermes wallet in both chains\n"

# Wallet mnemonic for Hermes
MNEMONIC_FOR_HERMES="trash travel uncle swarm side mobile project luggage ring reveal course rabbit window lonely either famous prize smooth charge sting pencil coil crime parrot"

# Transfer funds to Hermes wallet
commercionetworkd tx bank send \
  did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd \
  did:com:1xdjmpk8r0nzfuhgfvg6s2gxmtfxslqmn9newuv \
  1000000000ucommercio,1000000000uccc \
  --chain-id $COMMERCIO_CHAIN_ID $COMMERCIO_KEYRING \
  --fees 10000ucommercio \
  -o json -y 

commercionetworkd tx bank send \
  did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd \
  did:com:1xdjmpk8r0nzfuhgfvg6s2gxmtfxslqmn9newuv \
  1000000000ucommercio,1000000000uccc --chain-id $COMMERCIO2_CHAIN_ID $COMMERCIO2_KEYRING \
  --fees 10000ucommercio \
  --node tcp://localhost:26659 \
  -o json -y

echo $MNEMONIC_FOR_HERMES > $HERMES_CONFIG_DIR/.mnemonic

# Change permissions if on Linux
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sudo chmod -R 777 $HERMES_CONFIG_DIR
fi

# Create relayer keys
echo -e "\n🔧 Creating relayer keys\n"

hermes keys add \
    --key-name commercionetwork-relay \
    --chain $COMMERCIO_CHAIN_ID \
    --mnemonic-file $HERMES_CONFIG_DIR/.mnemonic

hermes keys add \
    --key-name commercio2-relay \
    --chain $COMMERCIO2_CHAIN_ID \
    --mnemonic-file $HERMES_CONFIG_DIR/.mnemonic

# Create relayer clients
echo -e "\n🔧 Creating relayer clients\n"

hermes create client \
    --host-chain $COMMERCIO_CHAIN_ID \
    --reference-chain $COMMERCIO2_CHAIN_ID

hermes create client \
    --host-chain $COMMERCIO2_CHAIN_ID \
    --reference-chain $COMMERCIO_CHAIN_ID

# Create relayer connections
echo -e "\n🔧 Creating relayer connections\n"

hermes create connection \
    --a-chain $COMMERCIO_CHAIN_ID \
    --a-client 07-tendermint-0 \
    --b-client 07-tendermint-0

echo -e "\n⏳ Waiting for connection handshake...\n"
sleep 30

# Create channel between chains
echo -e "\n🔧 Creating channel between chains\n"

hermes create channel --order unordered \
    --a-chain $COMMERCIO_CHAIN_ID \
    --a-connection connection-0 \
    --a-port transfer \
    --b-port transfer

# Query channels
echo -e "\n⏳ Waiting for query channels...\n"

hermes query channels --show-counterparty \
    --chain $COMMERCIO_CHAIN_ID

# Start Hermes
echo -e "\n🚀 Waiting for Hermes to start...\n"

hermes start > $LOGS_DIR/hermes_server_start.txt &
hermes_pid=$!

echo -e "\n💸 Transferring some tokens 💸\n"

# hermes tx ft-transfer \
#     --timeout-seconds 1000 \
#     --dst-chain $COMMERCIO2_CHAIN_ID \
#     --src-chain $COMMERCIO_CHAIN_ID \
#     --src-port transfer \
#     --src-channel channel-0 \
#     --amount 1000000 \
#     --denom ucommercio

# sleep 10

# hermes tx ft-transfer \
#     --timeout-seconds 1000 \
#     --dst-chain $COMMERCIO_CHAIN_ID \
#     --src-chain $COMMERCIO2_CHAIN_ID \
#     --src-port transfer \
#     --src-channel channel-0 \
#     --amount 500000 \
#     --denom ibc/EF505ADD9D17C91A013DC23BD975CEB7AA3CC118052940A20A3D623A75658B4C 

# sleep 10

# hermes tx ft-transfer \
#     --timeout-seconds 1000 \
#     --dst-chain $COMMERCIO_CHAIN_ID \
#     --src-chain $COMMERCIO2_CHAIN_ID \
#     --src-port transfer \
#     --src-channel channel-0 \
#     --amount 1000000 \
#     --denom ucommercio 

# sleep 10

# hermes tx ft-transfer \
#     --timeout-seconds 1000 \
#     --dst-chain $COMMERCIO2_CHAIN_ID \
#     --src-chain $COMMERCIO_CHAIN_ID \
#     --src-port transfer \
#     --src-channel channel-0 \
#     --amount 500000 \
#     --denom ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518 

# Wait for Hermes to finish
wait $hermes_pid

########################################################################################################################
# Set variables
WALLET_CREATOR="did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"
ADDRESS_LIMITER="/home/cheikh/Lavoro/commercionetwork/x/ibc-address-limiter/target/wasm32-unknown-unknown/release/address_limiter.wasm"
LIMITER_CODE_ID=1
INIT='{"gov_module": "did:com:10d07y265gmmuvt4z0w9aw880jnsr700jdfgwkq","ibc_module": "did:com:1yl6hdjhmkf37639730gffanpzndzdpmhe5dzhs","addrs_whitelist": ["did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"]}'
LIMITER_ADDRESS="did:com:14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sh7yll8"
PROPOSAL_PATH="/home/cheikh/Lavoro/proposal.json"

# Store the address limiter contract
commercionetworkd tx wasm store $ADDRESS_LIMITER --from did:com:1zg4jreq2g57s4efrl7wnh2swtrz3jt9nfaumcm --fees 100000000ucommercio -o text --gas 50000000 -y

sleep 5

# Instantiate the address limiter contract
commercionetworkd tx wasm instantiate $LIMITER_CODE_ID "$INIT" --label "IBC ADDRESS LIMITER" --admin $WALLET_CREATOR --from $WALLET_CREATOR --fees 10000ucommercio -o text --gas 50000000 -y

sleep 5

# Submit a legacy governance proposal for parameter change
commercionetworkd tx gov submit-legacy-proposal param-change $PROPOSAL_PATH --from did:com:1zg4jreq2g57s4efrl7wnh2swtrz3jt9nfaumcm --fees 100000000ucommercio -o text --gas 50000000 -y

sleep 5

# Vote on the proposal
commercionetworkd tx gov vote 1 yes --from did:com:1zg4jreq2g57s4efrl7wnh2swtrz3jt9nfaumcm -y --fees 10000ucommercio
commercionetworkd tx gov vote 1 yes --from did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd -y --fees 10000ucommercio

echo -e "\n⏳ Waiting for 2 minutes...\n"

sleep 130

# Query the votes
commercionetworkd query gov votes 1

# Query the address limiter contract parameter
commercionetworkd query params subspace address-limited-ibc contract

# # Perform an IBC transfer
# commercionetworkd tx ibc-transfer transfer transfer channel-0 did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd 200000ucommercio --from did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd --fees 100000000ucommercio -o text --gas 50000000 -y

echo -e "\n\n🏁 🏆 ${BOLD}Well done, you're ready! 🎇 🎆${NC}\n"