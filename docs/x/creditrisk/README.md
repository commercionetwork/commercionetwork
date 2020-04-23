# Creditrisk

The Creditrisk module keeps track and stores of all the Commercio tokens that were used to buy memberships and/or 
are the results of the CDP autoliquidation feature, which are concentrated on the **credit risk pool**.

## Queries

### Getting current credit risk pool amount

#### CLI

```bash
cncli query creditrisk pool
```

#### REST

Endpoints:
     
```
/creditrisk/pool
```

##### Example 

Getting current credit risk pool amount:

```
http://localhost:1317/creditrisk/pool
```

##### Response
```json
{
  "height": "0",
  "result": [
    "<amount of coins to deposit into the pool>"
  ]
}
```

where components of `result` are object of the following type:

```json
{
    "denom": "<token denom>",
    "amount": "<integer amount of tokens of denom>"
}