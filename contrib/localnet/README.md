# Commercio.network - Using localnet full stack

It is possible to test the blockchain and its functionality through the use of a complete docker stack including 4 nodes + endpoints for lcd + rpc with websocket.    
From root project folder execute the below commands.

**NB**: if you use linux system you need work as root, or add `sudo` to all command.

If it is the first time that you build the docker images:
```bash
make build-image-to-download-libraries localnet-start
```

If you have already built the project once, use:
```bash
make build-image-libraries-cached localnet-start
```

In particular, build-image-to-donwload-libraries, build two docker images, one containeining the libraries for the project and one with the project built inside, the second one start building from the first permitting to avoid the donwloading of all the dependencies each time.


To send in backgroud the stack use 

```bash
make build-image-libraries-cached localnet-start >/dev/null 2>&1 &
```

or 

```bash
make build-image-libraries-cached localnet-start-daemon
```

You can see logs with

```bash
docker compose logs
```


To stop stack use `ctrl + c` or

```bash
make localnet-stop
```

You can reset chain without delete genesis or accounts with 

```bash
make localnet-reset
```

You can delete all data of the local chain and all binaries with

```bash
make clean
```


If there are no problems, the nodes listen on the follow ports

```
commercionetworknode0   0.0.0.0:26656-26657->26656-26657/tcp, 0.0.0.0:9090->9090/tcp, 0.0.0.0:1317->1317/tcp              
commercionetworknode1   0.0.0.0:26659->26656/tcp, 0.0.0.0:26660->26657/tcp, 0.0.0.0:9091->9090/tcp
commercionetworknode2   0.0.0.0:26661->26656/tcp, 0.0.0.0:26662->26657/tcp, 0.0.0.0:9092->9090/tcp   
commercionetworknode3   0.0.0.0:26663->26656/tcp, 0.0.0.0:26664->26657/tcp, 0.0.0.0:9093->9090/tcp
```

Lcd and Rpc + websocket + Grpc
```
proxy-nginx   0.0.0.0:7123->7123/tcp, 0.0.0.0:7124->7124/tcp, 0.0.0.0:7125->7125/tcp  
```

## Lcd local access

```
http://localhost:7123
```

## Rpc local access
```
http://localhost:7124
```

## GRpc local access
```
http://localhost:7125
```


## Websocket local access
```
ws://localhost:7124/websocket
```

## Node data

Every node account mnemonics are in 

```
/build/node<N>/key_seed.json
```

where `<N>` is the id of node.  
Government wallet is in node0.    

Every node configs are under 


```
/build/node<N>/commercionetwork/config
```

Logs

```
/build/node<N>/commercionetwork/commercionetwork.log
```

## Add Node

If you want add a new node you can start a new container of `commercionetworknode` with new configuration

### Compile binary


```bash
make build
```


### Create new configuration



```bash
./build/commercionetworkd init node4 --home ./build/node4/commercionetwork
```

### Copy default genesis in config file

```bash
cp ./build/base_config/genesis.json ./build/node4/commercionetwork/config/
```

### Setup persistent

```bash
PERSISTENT=$(cat ./build/base_config/persistent.txt)
sed -i -e "s/persistent_peers = \".*\"/persistent_peers = \"$PERSISTENT\"/g" ./build/node4/commercionetwork/config/config.toml
```

### Start docker node



```bash
docker run \
   -v $(pwd)/build:/commercionetwork:Z \
   -e ID=4 \
   -p 26691-26692:26656-26657 \
   -p 9191:9090 \
   --ip 192.168.10.10 \
   --name node4 \
   --network commercionetwork_localnet \
   -d \
   commercionetwork/commercionetworknode
```



You can see logs with

```bash
docker logs node4 -f
```


### Make the validator node

You can discover the account with a lot of tokens in the first node using

```bash
./build/commercionetworkd keys list \
  --keyring-backend test \
  --home ./build/node0/commercionetwork/
```

The output should be something like below

```
- name: node0
  type: local
  address: did:com:1fm4ktq7t2282kmgcsptgm3j7f4k58r4zswseqw
  pubkey: did:com:pub1addwnpepqgc5jkcd0yfky2zssfl3096u6aj6xruj6wjqjnw56s4zgvtkw358sw3qmn9
  mnemonic: ""
  threshold: 0
  pubkeys: []
```

Create a new wallet with

```
./build/commercionetworkd keys add wc_node4 \
  --keyring-backend test \
  --home ./build/node4/commercionetwork/
```

The output should be something like below

```
- name: wc_node4
  type: local
  address: did:com:1xnju336hjcjkgv7mk96z2sckh6y6axeglznrpl
  pubkey: did:com:pub1addwnpepqgsgf2cq29p5k2s7dckh8yuq92wdju4smkjhjjva5wjqf4s3ansty30h8k9
  mnemonic: ""
  threshold: 0
  pubkeys: []


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.
similar copy session lens elbow bind custom leopard pyramid inspire situate feature large intact boat ensure penalty envelope sign action sentence boost rebel tower
```
**Save the mnemonic**

Transfer a minumun amount of token to the new wallet from the first one to create the new account on chain

```bash
# did:com:1fm4ktq7t2282kmgcsptgm3j7f4k58r4zswseqw is the wallet in the first node with a lot of tokens
# did:com:1xnju336hjcjkgv7mk96z2sckh6y6axeglznrpl is the wallet that you created before

./build/commercionetworkd tx bank send \
  did:com:1fm4ktq7t2282kmgcsptgm3j7f4k58r4zswseqw \
  did:com:1xnju336hjcjkgv7mk96z2sckh6y6axeglznrpl \
  20000000ucommercio \
  --keyring-backend test \
  --home ./build/node0/commercionetwork/ \
  --chain-id $(jq -r '.chain_id' ./build/base_config/genesis.json) \
  --fees 10000ucommercio \
  -y
```

Check the balances of your wallet

```bash
./build/commercionetworkd \
  query bank balances \
  did:com:1xnju336hjcjkgv7mk96z2sckh6y6axeglznrpl
```

Create validator

```bash
NODENAME=node4
CHAINID=$(jq -r '.chain_id' ./build/base_config/genesis.json)
VALIDATOR_PUBKEY=$(./build/commercionetworkd \
  tendermint show-validator \
  --home ./build/node4/commercionetwork/)
WALLET_CREATOR="did:com:1xnju336hjcjkgv7mk96z2sckh6y6axeglznrpl"

./build/commercionetworkd tx staking create-validator \
  --amount=1000000ucommercio \
  --pubkey=$VALIDATOR_PUBKEY \
  --moniker="$NODENAME" \
  --chain-id="$CHAINID" \
  --identity="" \
  --website="" \
  --details="" \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --from=$WALLET_CREATOR \
  --keyring-backend test \
  --home ./build/node4/commercionetwork/ \
  --fees=10000ucommercio \
  -y
```

Check if your validator is present

```bash
./build/commercionetworkd query staking validators
```