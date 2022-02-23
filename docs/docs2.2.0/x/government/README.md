# Government
In order to allow some specific operations to be performed only by a small set of individuals, 
inside Commercio.network we introduced the `government` module. 

This module allows for two simple operations: 

1. Set a government address that will later be used as a sort of on-chain authentication method. 
2. Read the government address that has been set. 

## Setting a government address 
The address identified as the `government` can be set **only during the genesis**.
This operation can be performed using the following command: 

```bash
cnd set-genesis-government-address <ADDRESS-TO-USE>
```

:::danger
**Note**: you can run this command only once.
Running it several times after the first value has been set will result in an error been thrown inside the console.
:::

## Retrieving the government address
### End user

The government address can be retrieved by using either `cncli` or by making a REST request:

 - via **CLI**, `cncli query government gov-address`
 - via **REST**, by making a GET request to the `/government/address` endpoint 


## Setting a tumbler address 
The address identified as the `tumbler` can be set during the genesis or at active chain.
This operation can be performed using the following command during the genesis: 

```bash
cnd set-genesis-tumbler-address <ADDRESS-TO-USE>
```

When the chain is active you can use the following command

```bash
cncli tx government set-tumbler-address <ADDRESS-TO-USE> --form <GOV-ADDRESS>
```

:::warning
**Note**: only the government can configure the tumbler
:::


## Retrieving the tumbler address
### End user

The tumbler address can be retrieved by using either `cncli` or by making a REST request:

 - via **CLI**, `cncli query government tumbler-address`
 - via **REST**, by making a GET request to the `/government/tumbler` endpoint 



### From the codebase
If you're developing a new module or implementing a new feature into an existing one and you need to access the current 
government address, you can use the `government.Keeper` object.  
This object allows you to use the `GetGovernmentAddress` method to read the current value.