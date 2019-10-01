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
There is currently no way to retrieve the government address using a CLI interface or a REST API endpoint.  
The only way that a user can see the government address is by looking inside the `genesis.json` file, 
specifically inside the `app_state.government.government_address` field value.  

### From the codebase
If you're developing a new module or implementing a new feature into an existing one and you need to access the current 
government address, you can use the `government.Keeper` object.  
This object allows you to use the `GetGovernmentAddress` method to read the current value.