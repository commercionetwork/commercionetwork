#Open a CDP (Collateralized Debt Position)
:::warning  
Before doing this type of transaction be sure to understand what a CDP is, and how it's works.  
You could loose your token if you don't be careful.  
:::
##Transaction message
To open a new CDP you need to create and sign the following message.  
```json
{
  "type": "commercio/MsgOpenCDP",
  "value": {
    "cdp_request": {
      "deposited_amount": "<Token to be deposited as a collateral (supports only integers)>",
      "signer": "<User address>",
      "timestamp": "<Timestamp of when the CDP request was made>"
    }
  }
}
```

##About the Timestamp


