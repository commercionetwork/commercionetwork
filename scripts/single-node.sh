#!/bin/sh

#set -o errexit -o nounset
set -o errexit

CHAINID=$1
GENACCT=$2
HOMECOMMERCIO=$3
BINCOMMERCIO=$4

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input genesis account address..."
  exit 1
fi


if [ -z "$3" ]; then
  echo "Need to input home chain..."
  exit 1
fi

if [ -z "$4" ]; then
  BINCOMMERCIO="commercionetworkd"
fi




# Build genesis file incl account for passed address
coins="10000000000uccc,100000000000stake"
echo "[OK] inizializzo il genesis"
$BINCOMMERCIO init --chain-id $CHAINID $CHAINID --home $HOMECOMMERCIO
echo "[OK] aggiungo account di test al keyring"
$BINCOMMERCIO keys add validator --keyring-backend="test"
echo "[OK] aggiungo account di test alla chain"
$BINCOMMERCIO add-genesis-account $($BINCOMMERCIO keys show validator -a --keyring-backend="test") $coins --home $HOMECOMMERCIO
echo "[OK] aggiungo account di in input"
$BINCOMMERCIO add-genesis-account $GENACCT $coins --home $HOMECOMMERCIO
echo "[OK] aggiungo government"
$BINCOMMERCIO set-genesis-government-address $GENACCT --home $HOMECOMMERCIO
echo "[OK] aggiungo vbr pool"
$BINCOMMERCIO set-genesis-vbr-pool-amount 1000000000stake --home $HOMECOMMERCIO
echo "[OK] aggiungo vbr pool"
$BINCOMMERCIO set-genesis-vbr-reward-rate 0.01 --home $HOMECOMMERCIO
echo "[OK] aggiungo validatore"
$BINCOMMERCIO gentx validator 5000000000stake --keyring-backend="test" --keyring-dir="~/.commercionetwork" --chain-id $CHAINID --home $HOMECOMMERCIO
$BINCOMMERCIO collect-gentxs --home $HOMECOMMERCIO

# Set proper defaults and change ports
sed 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' \
   $HOMECOMMERCIO/config/config.toml > \
   $HOMECOMMERCIO/config/config.toml.tmp
sed 's/timeout_commit = "5s"/timeout_commit = "1s"/g' \
   $HOMECOMMERCIO/config/config.toml.tmp > \
   $HOMECOMMERCIO/config/config.toml.tmp2
sed 's/timeout_propose = "3s"/timeout_propose = "1s"/g' \
   $HOMECOMMERCIO/config/config.toml.tmp2 > \
   $HOMECOMMERCIO/config/config.toml.tmp
sed 's/index_all_keys = false/index_all_keys = true/g' \
   $HOMECOMMERCIO/config/config.toml.tmp > \
   $HOMECOMMERCIO/config/config.toml
# Start the gaia
$BINCOMMERCIO start --pruning=nothing --home $HOMECOMMERCIO