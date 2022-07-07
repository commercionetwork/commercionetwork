<!--
order: 5
-->

# Client

## Transactions

### Mint CCC

Mints a given amount of CCC
#### CLI



```bash
commercionetworkd tx commerciomint mint \
  [amount]
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `amount`               | Amount of required uCCC  |

### Burn CCC

Burns a given amount of tokens, associated with ETP id.
#### CLI

```bash
commercionetworkd tx commerciomint burn \
  [id] \
  [amount]
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `id`         | ETP id from wich get tokens  |
| `amount`         | Amount of required uCOM  |

### Set Params

Set the commerciomint params with conversion rate and freeze-period in seconds.
#### CLI

```bash
commercionetworkd tx commerciomint set-params \
  [conversion-rate] \
  [freeze-period]
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `conversion-rate`         | ETP conversion reate  |
| `freeze-period`         | ETP freeze period in seconds  |


## Queries

### Reading all Exchange Trade Position (ETP)

#### CLI

```bash
commercionetworkd query commerciomint get-all-etps
```

#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciomint.Query/Etps
```

##### Example

```bash
grpcurl -plaintext \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciomint.Query/Etps
```

##### Response
```json
{
  "Positions": [
    {
      "owner": "did:com:1rsyglnhpg7q6hvz3422wm63tehtkx5xa2uwp3j",
      "collateral": "400000000",
      "credits": {
        "denom": "uccc",
        "amount": "400000000"
      },
      "createdAt": "2022-07-04T10:19:15.704764997Z",
      "ID": "8f1a387b-dcbd-43ec-9376-026b45d1f5d2",
      "exchangeRate": "1000000000000000000"
    },

    ...

  ],
  "pagination": {
    "total": "10"
  }
}
```

#### REST

Endpoint:
   
```
/commercionetwork/commerciomint/etps
```

##### Example

Getting all users opened ETPs:

```
http://localhost:1317/commercionetwork/commerciomint/etps
```

##### Response
```json
{
  "Positions": [
    {
      "owner": "did:com:1zg4jreq2g57s4efrl7wnh2swtrz3jt9nfaumcm",
      "collateral": "7",
      "credits": {
        "denom": "uccc",
        "amount": "10"
      },
      "created_at": "2021-07-22T13:18:44.598560074Z",
      "ID": "090ca0c2-cf00-4119-8307-b51413a00cf4",
      "exchange_rate": "0.610000000000000000"
    },
    {
      "owner": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
      "collateral": "274500",
      "credits": {
        "denom": "uccc",
        "amount": "450000"
      },
      "created_at": "2022-02-15T09:02:46.475744007Z",
      "ID": "805a82db-a9e7-441a-a26b-d9dd9dc84a0b",
      "exchange_rate": "0.610000000000000000"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "2"
  }
}
```

### Reading all Exchange Trade Position (ETP) opened by a user

#### CLI

```bash
commercionetworkd query commerciomint get-etps [user-addr]
```
#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciomint.Query/EtpsByOwner
```

##### Example

```bash
grpcurl -plaintext \
    -d '{"Owner":"did:com:1rsyglnhpg7q6hvz3422wm63tehtkx5xa2uwp3j"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciomint.Query/EtpsByOwner
```

##### Response
```json
{
  "Positions": [
    {
      "owner": "did:com:1rsyglnhpg7q6hvz3422wm63tehtkx5xa2uwp3j",
      "collateral": "400000000",
      "credits": {
        "denom": "uccc",
        "amount": "400000000"
      },
      "createdAt": "2022-07-04T10:19:15.704764997Z",
      "ID": "8f1a387b-dcbd-43ec-9376-026b45d1f5d2",
      "exchangeRate": "1000000000000000000"
    },
    ...
  ],
  "pagination": {
    "total": "3"
  }
}
```

#### REST

Endpoint:
   
```
​/commercionetwork​/commerciomint/${Owner}/etpsOwner
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `owner` | Address of the user for which to read all the ETPs |

##### Example

Getting ETPs opened by `did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd`:

```
http://localhost:1317/commercionetwork/commerciomint/did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd/etpsOwner
```

##### Response
```json
{
  "Positions": [
    {
      "owner": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
      "collateral": "274500",
      "credits": {
        "denom": "uccc",
        "amount": "450000"
      },
      "created_at": "2022-02-15T09:02:46.475744007Z",
      "ID": "805a82db-a9e7-441a-a26b-d9dd9dc84a0b",
      "exchange_rate": "0.610000000000000000"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

### Reading all Exchange Trade Position (ETP) by ID

#### CLI

```bash
commercionetworkd query commerciomint get-etp [id]
```

#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciomint.Query/Etp
```

##### Example

```bash
grpcurl -plaintext \
    -d '{"ID":"8f1a387b-dcbd-43ec-9376-026b45d1f5d2"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciomint.Query/Etp
```

##### Response
```json
{
  "Position": {
    "owner": "did:com:1rsyglnhpg7q6hvz3422wm63tehtkx5xa2uwp3j",
    "collateral": "400000000",
    "credits": {
      "denom": "uccc",
      "amount": "400000000"
    },
    "createdAt": "2022-07-04T10:19:15.704764997Z",
    "ID": "8f1a387b-dcbd-43ec-9376-026b45d1f5d2",
    "exchangeRate": "1000000000000000000"
  }
}
```

#### REST

Endpoint:
   
```
​/commercionetwork​/commerciomint/${ID}/etp
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `id` | ID of the wanted etp |

##### Example

Getting ETPs with ID `805a82db-a9e7-441a-a26b-d9dd9dc84a0b`:

```
http://localhost:1317/commercionetwork/commerciomint/805a82db-a9e7-441a-a26b-d9dd9dc84a0b/etp
```

##### Response
```json
{
  "Position": {
    "owner": "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd",
    "collateral": "274500",
    "credits": {
      "denom": "uccc",
      "amount": "450000"
    },
    "created_at": "2022-02-15T09:02:46.475744007Z",
    "ID": "805a82db-a9e7-441a-a26b-d9dd9dc84a0b",
    "exchange_rate": "0.610000000000000000"
  }
}
```

### Reading the Params (conversion rate & freeze period)

#### CLI

```bash
commercionetworkd query commerciomint get-params
```
#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciomint.Query/Params
```

##### Example

```bash
grpcurl -plaintext \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciomint.Query/Params
```

##### Response
```json
{
  "params": {
    "conversionRate": "1000000000000000000",
    "freezePeriod": "1814400s"
  }
}
```

#### REST

Endpoint:
   
```
/commercionetwork​/commerciomint​/params
```

##### Example

Getting the parameters:

```
http://localhost:1317/commercionetwork/commerciomint/params
```

##### Response
```json
{
  "params": {
    "conversion_rate": "0.610000000000000000",
    "freeze_period": "1814400s"
  }
}
```
