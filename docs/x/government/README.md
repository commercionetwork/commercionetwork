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

Note that you can run this command only once.
Running it several times after the first value has been set will result in an error been thrown inside the console.

## Retrieving the government address
### End user

The government address can be retrieved by using either `cncli` or by making a REST request:

 - via CLI, `cncli query government address`
 - via REST, by making a GET request to the `/government/address` endpoint 

### From the codebase
If you're developing a new module or implementing a new feature into an existing one and you need to access the current 
government address, you can use the `government.Keeper` object.  
This object allows you to use the `GetGovernmentAddress` method to read the current value.