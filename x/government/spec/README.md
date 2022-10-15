<!--
order: 0
title: Government Overview
parent:
  title: "covernment"
-->

# Government 

## Abstract

In order to allow some specific operations to be performed only by a small set of individuals, 
inside Commercio.network we introduced the `government` module. 

This module allows for two simple operations: 

1. Set a government address that will later be used as a sort of on-chain authentication method. 
2. Read the government address that has been set. 

### Setting a government address 
The address identified as the `government` can be set **only during the genesis**.
This operation can be performed using the following command: 

```bash
commercionetworkd set-genesis-government-address <ADDRESS-TO-USE>
```

:::danger
**Note**: you can run this command only once.
Running it several times after the first value has been set will result in an error been thrown inside the console.
:::

### Query government address 

The government address can be get:

 - via **CLI**, `commercionetworkd query government gov-address`
 - via **REST**, by making a GET request to the `/commercionetwork/government/governmentAddress` endpoint 
 - via **gRPC**, by making a Query to the `commercionetwork.commercionetwork.government.Query` method

#### gRPC

Endpoint:

```
commercionetwork.commercionetwork.government.Query/GovernmentAddr
```

##### Example

```bash
grpcurl -plaintext \
    localhost:9090 \
    commercionetwork.commercionetwork.government.Query/GovernmentAddr
```

##### Response
```json
{
  "governmentAddress": "did:com:1mj9h87yqjel0fsvkq55v345kxk0n09krtfvtyx"
}
```

## Contents

1. **[State](01_state.md)**
