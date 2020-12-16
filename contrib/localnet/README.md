# Commercio.network - Using localnet full stack

It is possible to test the blockchain and its functionality through the use of a complete docker stack including 4 nodes + endpoints for lcd + rpc with websocket.    
From root project folder execute the below commands

```bash
make build-docker-cndode localnet-start
```

To send in backgroud the stack use 

```bash
make build-docker-cndode localnet-start >/dev/null 2>&1 &
```


To stop stack use `ctrl + c` or

```bash
make localnet-stop
```


If there are no problems, the nodes listen on the follow ports

```
cndnode0   0.0.0.0:26656-26657->26656-26657/tcp                 
cndnode1   0.0.0.0:26659->26656/tcp, 0.0.0.0:26660->26657/tcp
cndnode2   0.0.0.0:26661->26656/tcp, 0.0.0.0:26662->26657/tcp   
cndnode3   0.0.0.0:26663->26656/tcp, 0.0.0.0:26664->26657/tcp
```

Lcd and Rpc + websocket
```
proxy-nginx   0.0.0.0:7123->7123/tcp, 0.0.0.0:7124->7124/tcp 
```

## Lcd local access

```
http://localhost:7123
```

## Rpc local access
```
http://localhost:7124
```

## Websocket local access
```
ws://localhost:7124/websocket
```

## Node data

Every node account mnemonics are in 

```
/build/node<N>/cncli/key_seed.json
```

where `<N>` is the id of node

Every node configs are under 


```
/build/node<N>/cnd/config
```

Logs

```
/build/node<N>/cnd/cnd.log
```

