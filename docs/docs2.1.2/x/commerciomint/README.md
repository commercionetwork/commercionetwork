# CommercioMint

The `commerciomint` module is the one that allows you to create Collateralized Debt Positions (*CDPs*) using your 
Commercio.network tokens (*ucommercio*) in order to get Commercio Cash Credits (*uccc*) in return.

A *Collateralized Debt Position* (*CDP*) is a core component of the Commercio Network blockchain whose purpose is to
create Commercio Cash Credits (`uccc`) in exchange for Commercio Tokens (`ucommercio`) which it then holds in
escrow until the borrowed Commercio Cash Credits are returned.

In simple words, opening a CDP allows you to exchange any amount of `ucommercio` to get half the amount of `uccc`. 
For example, if you open a CDP lending `100 ucommercio` will result in you receiving `50 uccc`.    

## Transactions

### Open a CDP

#### Transaction message
To open a new CDP you need to create and sign the following message:
  
```json
{
  "type": "commercio/MsgOpenCdp",
  "value": {
    "deposited_amount": [
      {
        "amount": "<amount to be deposited>",
        "denom": "<token denom to be deposited>"
      }
    ],
    "depositor": "<user address>"
  }
}
```

A maximum of 1 `MsgOpenCdp` messages can be sent in each transaction.

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
openCdp
```  

### Close a CDP

#### Transaction message

To close a previously opened CDP you need to create and sign the following message:

```json
{
  "type": "commercio/MsgCloseCdp",
  "value": {
    "signer": "<user address>",
    "timestamp": "<block height at which the CDP is being inserted into the chain>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
closeCdp
```  

### Set CDP collateral rate

:::warning  
This transaction type is accessible only to the [government](../../government/README.md).  
Trying to perform this transaction without being the government will result in an error.  
:::

#### Transaction message

To set the CDP collateral rate you need to create and sign the following message:

```json
{
  "type": "commercio/MsgSetCdpCollateralRate",
  "value": {
    "signer": "<user address>",
    "cdp_collateral_rate": "<floating-point collateral rate>"
  }
}
```

#### Action type
If you want to [list past transactions](../../../developers/listing-transactions.md) including this kind of message,
you need to use the following `message.action` value: 

```
setCdpCollateralRate
```  

## Queries

### Reading all CDP opened by a user at a given timestamp

#### CLI

```bash
cncli query commerciomint get-cdp [user-addr] [block-height]
```

#### REST

Endpoint:  

```
/commerciomint/cdps/${address}/${timestamp}
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read the CDP |
| `timestamp`| Timestamp of when the CDP request was made |

##### Example 

Getting CDPs opened by `did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke` at timestamp `1570177686`:
```
http://localhost:1317/commerciomint/cdps/did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke/1570177686
```
#### Response
```json
{
  "height": "0",
  "result": {
    "deposited_amount": [
      {
        "amount": "10000000",
        "denom": "ucommercio"
      }
    ],
    "liquidity_amount": {
      "amount": "500000",
      "denom": "uccc"
    },
    "owner": "did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke/1570177686",
    "timestamp": "1570177686"
  }
}
```

### Reading all CDP opened by a user

#### CLI

```sh
$ cncli query commerciomint get-cdps [user-addr]
```

#### REST

Endpoint:
   
```
/commerciomint/cdps/${address}
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read all the CDPs |

##### Example

Getting CDPs opened by `did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke`:

```
http://localhost:1317/commerciomint/cdps/did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke
```

#### Response
```json
{
  "height": "0",
  "result": [
    {
      "deposited_amount": [
        {
          "amount": "10000000",
          "denom": "ucommercio"
        }
      ],
      "liquidity_amount": {
        "denom": "uccc",
        "amount": "500000"
      },
      "owner": "did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke/1570177686",
      "timestamp": "1570177686"
    }
  ]
}
```

### Reading the current CDP collateral rate

#### CLI

```bash
cncli query commerciomint collateral-rate
```

#### REST

Endpoint:
   
```
/commerciomint/collateral_rate
```

##### Example

Getting the current CDP collateral rate:

```
http://localhost:1317/commerciomint/collateral_rate
```

#### Response
```json
{
  "height": "0",
  "result": "2.000000000000000000"
}
```