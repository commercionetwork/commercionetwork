<!--
order: 6
-->

# Client

## Queries

### Reading all Exchange Trade Position (ETP)

#### CLI

```sh
commercionetworkd query commerciomint get-all-etps
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

#### Response
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

```sh
commercionetworkd query commerciomint get-etps [user-addr]
```

#### REST

Endpoint:
   
```
​/commercionetwork​/commerciomint/${Owner}/owner
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `owner` | Address of the user for which to read all the ETPs |

##### Example

Getting ETPs opened by `did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd`:

```
http://localhost:1317/commercionetwork/commerciomint/did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd/owner
```

#### Response
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

```sh
commercionetworkd query commerciomint get-etp [id]
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

#### Response
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

#### Response
```json
{
  "params": {
    "conversion_rate": "0.610000000000000000",
    "freeze_period": "1814400s"
  }
}
```