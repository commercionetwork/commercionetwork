
# CommercioAPI surcharge

<!-- npm run docs:serve  -->

<!-- https://lcd-testnet.commercio.network/docs/did:com:1ug9j7hgaxu6mvfu2kgfdt3hqxn4mrwuztxc7nu/received -->


In the POST methd API is available the possibility to indicate a surcharge.

The surcharge is intended to allow those implementing external applications that rely on APIs to send an arbitrary surcharge value to their own wallet, in addition to the standard platform fee for the process execution.

The recipient of the surcharge amount must be equipped with a gold membership.

This option is available in all POST APIs.

* /sharedoc/process
* /receipts
* /wallet/transfer
* /eidentity/invite/
* /ddo/process/
* /sign/process/

Below is an example interaction schema for an external application.


COMING SOON


## Example 

Taking as an example  a sharedoc process executed 

* From a User with Green Membership (Sender) :  did:com:1kschysacm4zag3d9j7rf0pfjpxmx4waa0sc43d
* To user with Bronze membership (Recipient) :  did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr
* Surcharge towards Gold membership (Beneficiary): did:com:18kx6zp6crnagcq98hz008x7ze5w78h2ch6h3gz 
* Amount of surcharge: 2100000 uccc equivalent to 2.1 CCC 


### POST example

```
curl -X POST 'https://dev-api.commercio.app/v1/sharedoc/process' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer $ACCESS_TOKEN' \
  -d '{
    "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
    "hash": "f4e454f802b88d2f64168ff1742e8cf413fd677d38b87cbefb45821f8981b912",
    "hash_algorithm": "sha-256",
    "metadata": {
        "content_uri": "55fa8b74d91bc8443f46b9dc7a179bd3f709bb803f9dccda467310f0fb656a7f",
        "schema": {
            "uri": "http://example.com/schema.xml",
            "version": "1.0.0"
        }
    },
    "recipients": [
        "did:com:1tq5mvp7j4vtew08htaswsyjugzewe4jyph20qr"
    ],
    "surcharge": {
        "beneficiary": "did:com:18kx6zp6crnagcq98hz008x7ze5w78h2ch6h3gz",
        "amount": "2100000"
    }
  }'
``` 

Example of transaction 

https://testnet.commercio.network/transactions/detail/8718836DB8545935D3B3B603C02C5189CB4FE6B31E9A3890A298969C08472153

## Total cost of the process execution   

The total cost of each process (POST) for a Green membership user will therefore be the sum of

* Chain fee for the sharedoc Message 
* Platform cost for the sharedoc Message as 
   * Net Platform cost
   * Chain fee  for sending Platform cost to the platfrom
* Surcharge as 
   * Net surcharge amount
   * Chain fee  for sending Surcharge to the Beneficiary


Practically this permit to avoid to pay the Platform cost for sending the Surcharge to the Beneficiary . Only chain fee will be paid for sending the surcharge. Otherwise with a simple send message in case of Green membership a 0.23 Platform fee will be paid


| Costs components | Cost in CCC| Note |
| --- | --- | --- | 
| Sharedoc message  | 0.01 | Chain fee  | 
| Platform fee   | 0.23 |  cost fro green membership | 
| Platform fee send  | 0.01 |  Chain fee Transfer message of Platform fee (0.23) to Platform TSP| 
| Surcharge amount  | 2.1 |  Amount of surcharge to Gold | 
| Surcharge fee send  | 0.01 |  Chain fee Send message of Surcharge (2.1) to Gold membership| 



**TOTAL = 2.45 CCC**  

For other cases of platform costs based on membership, see [Platfrom Costs for each message section](https://docs.commercio.network/app_developers/commercioapi-introduction.html#platform-costs) .




