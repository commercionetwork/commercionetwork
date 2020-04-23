# Listing past transactions
In this page you will learn how to list all the transactions that have been made in the past that 
include a specific message type.

## REST API
The REST API endpoint that must be called is the following: 

```
/txs?message.action={action}
``` 

### Supported actions
Please refer to the [supported messages section](./message-types.md) to know more about all the supported
transaction messages. Inside each one page there's the `Action type` section which tells you the action type identifier
to use for transactions that include that message. 

### Examples
**Endpoint**   
`https://lcd-testnet.commercio.network/txs?message.action=send`

**Output**  
```json
{
  "count": "3", 
  "limit": "3", 
  "page_number": "1", 
  "page_total": "26", 
  "total_count": "78", 
  "txs": [
    {
      "gas_used": "71461", 
      "gas_wanted": "200000", 
      "height": "28", 
      "logs": [
        {
          "events": [
            {
              "attributes": [
                {
                  "key": "action", 
                  "value": "send"
                }, 
                {
                  "key": "sender", 
                  "value": "did:com:1t5fz439f49zv39pmh73c2lvuhwfzqj0ze3kzj2"
                }, 
                {
                  "key": "module", 
                  "value": "bank"
                }
              ], 
              "type": "message"
            }, 
            {
              "attributes": [
                {
                  "key": "recipient", 
                  "value": "did:com:1j0ge8wgxcwx4l50lxkam5zkhqv28r7xyxt4zyp"
                }, 
                {
                  "key": "amount", 
                  "value": "51000000000ucommercio"
                }
              ], 
              "type": "transfer"
            }
          ], 
          "log": "", 
          "msg_index": 0
        }
      ], 
      "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"did:com:1t5fz439f49zv39pmh73c2lvuhwfzqj0ze3kzj2\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"did:com:1j0ge8wgxcwx4l50lxkam5zkhqv28r7xyxt4zyp\"},{\"key\":\"amount\",\"value\":\"51000000000ucommercio\"}]}]}]", 
      "timestamp": "2020-04-16T14:32:21Z", 
      "tx": {
        "type": "cosmos-sdk/StdTx", 
        "value": {
          "fee": {
            "amount": [
              {
                "amount": "10000", 
                "denom": "ucommercio"
              }
            ], 
            "gas": "200000"
          }, 
          "memo": "", 
          "msg": [
            {
              "type": "cosmos-sdk/MsgSend", 
              "value": {
                "amount": [
                  {
                    "amount": "51000000000", 
                    "denom": "ucommercio"
                  }
                ], 
                "from_address": "did:com:1t5fz439f49zv39pmh73c2lvuhwfzqj0ze3kzj2", 
                "to_address": "did:com:1j0ge8wgxcwx4l50lxkam5zkhqv28r7xyxt4zyp"
              }
            }
          ], 
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1", 
                "value": "AmkprhPNn/OGrnyhGSRhDzM4O97/m3LxIOnyHtBfgenr"
              }, 
              "signature": "iTSvKGWETSeQZeg4m15rKb7MO148ZMzVxgPOgrpueaUGhQGBgtHfyDGcLKYE9ogmcLk7EMLn6iGIyvd2RRpMlg=="
            }
          ]
        }
      }, 
      "txhash": "1188C71EE5185F9FAC4678E113D542BB02CFD604003141ACE4C98F99461B22A4"
    }, 
    {
      "gas_used": "71147", 
      "gas_wanted": "200000", 
      "height": "340", 
      "logs": [
        {
          "events": [
            {
              "attributes": [
                {
                  "key": "action", 
                  "value": "send"
                }, 
                {
                  "key": "sender", 
                  "value": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu"
                }, 
                {
                  "key": "module", 
                  "value": "bank"
                }
              ], 
              "type": "message"
            }, 
            {
              "attributes": [
                {
                  "key": "recipient", 
                  "value": "did:com:1u9lml9fnsxz03rlgavy8zkrv9arf8ynjc9r26a"
                }, 
                {
                  "key": "amount", 
                  "value": "100000000ucommercio"
                }
              ], 
              "type": "transfer"
            }
          ], 
          "log": "", 
          "msg_index": 0
        }
      ], 
      "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"did:com:1u9lml9fnsxz03rlgavy8zkrv9arf8ynjc9r26a\"},{\"key\":\"amount\",\"value\":\"100000000ucommercio\"}]}]}]", 
      "timestamp": "2020-04-16T15:00:14Z", 
      "tx": {
        "type": "cosmos-sdk/StdTx", 
        "value": {
          "fee": {
            "amount": [
              {
                "amount": "10000", 
                "denom": "ucommercio"
              }
            ], 
            "gas": "200000"
          }, 
          "memo": "", 
          "msg": [
            {
              "type": "cosmos-sdk/MsgSend", 
              "value": {
                "amount": [
                  {
                    "amount": "100000000", 
                    "denom": "ucommercio"
                  }
                ], 
                "from_address": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu", 
                "to_address": "did:com:1u9lml9fnsxz03rlgavy8zkrv9arf8ynjc9r26a"
              }
            }
          ], 
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1", 
                "value": "Ay+WFw3pst8XdjgX1VULdnYQNkFmWfigIfHxPScNx4C1"
              }, 
              "signature": "uFobPcoyTKjSGt4bfTrbYbV65GaXSUweUnJ6CX4oAHFOfRHe1hF770ccuEPezUhFqrjgfrGp4kzmGEl/xiKK6g=="
            }
          ]
        }
      }, 
      "txhash": "46EC499A14A70B6CC8A3E59F4D896801B72A6BA4C0AEE144CB2A81A89606657D"
    }, 
    {
      "gas_used": "71067", 
      "gas_wanted": "200000", 
      "height": "342", 
      "logs": [
        {
          "events": [
            {
              "attributes": [
                {
                  "key": "action", 
                  "value": "send"
                }, 
                {
                  "key": "sender", 
                  "value": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu"
                }, 
                {
                  "key": "module", 
                  "value": "bank"
                }
              ], 
              "type": "message"
            }, 
            {
              "attributes": [
                {
                  "key": "recipient", 
                  "value": "did:com:1mezjl60v8alvkfw9ax00sq2sjglhdznl67c25j"
                }, 
                {
                  "key": "amount", 
                  "value": "1000000ucommercio"
                }
              ], 
              "type": "transfer"
            }
          ], 
          "log": "", 
          "msg_index": 0
        }
      ], 
      "raw_log": "[{\"msg_index\":0,\"log\":\"\",\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"did:com:1mezjl60v8alvkfw9ax00sq2sjglhdznl67c25j\"},{\"key\":\"amount\",\"value\":\"1000000ucommercio\"}]}]}]", 
      "timestamp": "2020-04-16T15:00:24Z", 
      "tx": {
        "type": "cosmos-sdk/StdTx", 
        "value": {
          "fee": {
            "amount": [
              {
                "amount": "10000", 
                "denom": "ucommercio"
              }
            ], 
            "gas": "200000"
          }, 
          "memo": "", 
          "msg": [
            {
              "type": "cosmos-sdk/MsgSend", 
              "value": {
                "amount": [
                  {
                    "amount": "1000000", 
                    "denom": "ucommercio"
                  }
                ], 
                "from_address": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu", 
                "to_address": "did:com:1mezjl60v8alvkfw9ax00sq2sjglhdznl67c25j"
              }
            }
          ], 
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1", 
                "value": "Ay+WFw3pst8XdjgX1VULdnYQNkFmWfigIfHxPScNx4C1"
              }, 
              "signature": "iyXxLwMteZQVpKo5vpzM6oURveJNmq4e9kqCfVai5RA39xPlVv/RbEe8x56u5IptT6QnGf/jAjMW9SPXGrDOQQ=="
            }
          ]
        }
      }, 
      "txhash": "8A436FB66BFDE6863213579FE52D978113E0F19DDBEB0BB94D741EF64F3B7CA5"
    }
  ]
}
```