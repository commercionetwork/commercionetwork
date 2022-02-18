<!--
order: 6
-->

# Client

## Queries

### Getting the total Pool Funds

#### CLI

```sh
commercionetworkd query vbr pool-funds
```

#### REST

Endpoint:
   
```
/commercionetwork/vbr/funds
```

##### Example

Getting all the block rewards pool Funds:

```
http://localhost:1317/commercionetwork/vbr/funds
```

#### Response
```json
{
  "funds": [
    {
      "denom": "ucommercio",
      "amount": "10000000.000000000000000000"
    }
  ]
}
```

### Reading the Params

#### CLI

```bash
commercionetworkd query vbr get-params
```

#### REST

Endpoint:
   
```
/commercionetwork/vbr/params
```

##### Example

Getting the parameters:

```
http://localhost:1317/commercionetwork/vbr/params
```

#### Response
```json
{
  "params": {
    "distr_epoch_identifier": "day",
    "earn_rate": "0.500000000000000000"
  }
}
```