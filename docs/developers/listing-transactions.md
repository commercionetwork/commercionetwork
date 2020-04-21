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
`https://lcd-testnet.commercio.network/txs?message.action=shareDocument`

**Output**  
```json
{
  "total_count": "5",
  "count": "5",
  "page_number": "1",
  "page_total": "1",
  "limit": "30",
  "txs": [
    {
      "height": "178021",
      "txhash": "6D7065766CF473401FC6565DC6B2B7A9ED5D852F19100D2B5744AAB8439DA3B9",
      "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
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
      "gas_wanted": "200000",
      "gas_used": "52647",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "shareDocument"
            }
          ]
        }
      ],
      "tx": {
        "type": "cosmos-sdk/StdTx",
        "value": {
          "msg": [
            {
              "type": "commercio/MsgShareDocument",
              "value": {
                "sender": "did:com:1zfhgwfgex8rc9t00pk6jm6xj6vx5cjr4ngy32v",
                "recipient": "did:com:1g5cxz9p7dqux80sw5tljwg2pwz6x7zlq84km56",
                "uuid": "6a881ef0-04da-4524-b7ca-6e5e3b7e61dc",
                "metadata": {
                  "content_uri": "https://www.vargroup.it/managed-security-services/",
                  "schema_type": "",
                  "schema": {
                    "uri": "https://www.vargroup.it/managed-security-services/metadata/schema",
                    "version": "1.0.0"
                  },
                  "proof": "yes"
                },
                "content_uri": "https://www.vargroup.it/managed-security-services/",
                "checksum": {
                  "value": "7815696ecbf1c96e6894b779456d330e",
                  "algorithm": "md5"
                }
              }
            }
          ],
          "fee": {
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "250"
              }
            ],
            "gas": "200000"
          },
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A2LI01hYyy76qNPCxq7+gYshtefvK30FeVL0TpK6VpDI"
              },
              "signature": "ZXQL8XZZKCj0k0N6AXkrs2NFiXYgxa5sciCZgNZ7N3ZJQwYSk4sXKtAACPJqEnXGZrQUeamSom8zGm6Jx5RwbA=="
            }
          ],
          "memo": ""
        }
      },
      "timestamp": "2019-09-25T22:40:41Z"
    },
    {
      "height": "178070",
      "txhash": "EEB043AF4834D5F3396FAAC600A582701C0A78E88F745A5A23F73339AD2519B1",
      "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
      "logs": [
        {
          "msg_index": 0,
          "success": true,
          "log": ""
        }
      ],
      "gas_wanted": "200000",
      "gas_used": "73371",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "shareDocument"
            }
          ]
        }
      ],
      "tx": {
        "type": "cosmos-sdk/StdTx",
        "value": {
          "msg": [
            {
              "type": "commercio/MsgShareDocument",
              "value": {
                "sender": "did:com:1zfhgwfgex8rc9t00pk6jm6xj6vx5cjr4ngy32v",
                "recipient": "did:com:1g5cxz9p7dqux80sw5tljwg2pwz6x7zlq84km56",
                "uuid": "6a881ef0-04da-4524-b7ca-6e5e3b7e61dc",
                "metadata": {
                  "content_uri": "https://www.vargroup.it/managed-security-services/",
                  "schema_type": "",
                  "schema": {
                    "uri": "https://www.vargroup.it/managed-security-services/metadata/schema",
                    "version": "1.1.0"
                  },
                  "proof": "yes"
                },
                "content_uri": "https://www.vargroup.it/managed-security-services/",
                "checksum": {
                  "value": "7815696ecbf1c96e6894b779456d330e",
                  "algorithm": "md5"
                }
              }
            }
          ],
          "fee": {
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "250"
              }
            ],
            "gas": "200000"
          },
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A2LI01hYyy76qNPCxq7+gYshtefvK30FeVL0TpK6VpDI"
              },
              "signature": "d+8qd5DyevoinYUd3OqCvdhwuiM+423/T7zUeDnzsUVNXhFu6a2NgAzahvjyuXHHItwn5PxwI1yXpWjTOvoz9Q=="
            }
          ],
          "memo": ""
        }
      },
      "timestamp": "2019-09-25T22:45:05Z"
    },
    {
      "height": "178094",
      "txhash": "4A1025518519A2BC6953437DED0D21452C4C11D719639668E3AB80F60DB33690",
      "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
      "logs": [
        {
          "msg_index": 0,
          "success": true,
          "log": ""
        }
      ],
      "gas_wanted": "200000",
      "gas_used": "75255",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "shareDocument"
            }
          ]
        }
      ],
      "tx": {
        "type": "cosmos-sdk/StdTx",
        "value": {
          "msg": [
            {
              "type": "commercio/MsgShareDocument",
              "value": {
                "sender": "did:com:1zfhgwfgex8rc9t00pk6jm6xj6vx5cjr4ngy32v",
                "recipient": "did:com:1g5cxz9p7dqux80sw5tljwg2pwz6x7zlq84km56",
                "uuid": "6a881ef0-04da-4524-b7ca-6e5e3b7e61dc",
                "metadata": {
                  "content_uri": "https://www.vargroup.it/managed-security-services/",
                  "schema_type": "",
                  "schema": {
                    "uri": "https://www.vargroup.it/managed-security-services/metadata/schema",
                    "version": "1.1.0"
                  },
                  "proof": "yes"
                },
                "content_uri": "https://www.vargroup.it/managed-security-services/",
                "checksum": {
                  "value": "7815696ecbf1c96e6894b779456d330e",
                  "algorithm": "md5"
                }
              }
            }
          ],
          "fee": {
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "250"
              }
            ],
            "gas": "200000"
          },
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A2LI01hYyy76qNPCxq7+gYshtefvK30FeVL0TpK6VpDI"
              },
              "signature": "jhmHY1xLCjAi5X8Mr96QVlf264VJ0FwxjFqCmWxpdgVYcw4ZrmGCBy2EJYk8O6vKjaeyzGjjNqxX9u+L5oKYZw=="
            }
          ],
          "memo": ""
        }
      },
      "timestamp": "2019-09-25T22:47:14Z"
    },
    {
      "height": "178108",
      "txhash": "A771A321B64CCD7E8FFB6916B60FC90DCD2F87313C9626A7028C23DB30FF14F1",
      "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
      "logs": [
        {
          "msg_index": 0,
          "success": true,
          "log": ""
        }
      ],
      "gas_wanted": "200000",
      "gas_used": "75255",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "shareDocument"
            }
          ]
        }
      ],
      "tx": {
        "type": "cosmos-sdk/StdTx",
        "value": {
          "msg": [
            {
              "type": "commercio/MsgShareDocument",
              "value": {
                "sender": "did:com:1zfhgwfgex8rc9t00pk6jm6xj6vx5cjr4ngy32v",
                "recipient": "did:com:1g5cxz9p7dqux80sw5tljwg2pwz6x7zlq84km56",
                "uuid": "6a881ef0-04da-4524-b7ca-6e5e3b7e61dc",
                "metadata": {
                  "content_uri": "https://www.vargroup.it/managed-security-services/",
                  "schema_type": "",
                  "schema": {
                    "uri": "https://www.vargroup.it/managed-security-services/metadata/schema",
                    "version": "1.1.0"
                  },
                  "proof": "yes"
                },
                "content_uri": "https://www.vargroup.it/managed-security-services/",
                "checksum": {
                  "value": "7815696ecbf1c96e6894b779456d330e",
                  "algorithm": "md5"
                }
              }
            }
          ],
          "fee": {
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "250"
              }
            ],
            "gas": "200000"
          },
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A2LI01hYyy76qNPCxq7+gYshtefvK30FeVL0TpK6VpDI"
              },
              "signature": "OAPtfLhv0kVrOPHmzoAdykWExE+BOy2/I0anruh93tUSN4nQXFxZwFzvICcMSCue32HtVEPioF6P/+sHOSiT3w=="
            }
          ],
          "memo": ""
        }
      },
      "timestamp": "2019-09-25T22:48:29Z"
    },
    {
      "height": "183114",
      "txhash": "955362C3DADF23C852992756D5B10C1A1C526E5A28F6320793B238E3C49C4353",
      "raw_log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
      "logs": [
        {
          "msg_index": 0,
          "success": true,
          "log": ""
        }
      ],
      "gas_wanted": "200000",
      "gas_used": "75255",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "shareDocument"
            }
          ]
        }
      ],
      "tx": {
        "type": "cosmos-sdk/StdTx",
        "value": {
          "msg": [
            {
              "type": "commercio/MsgShareDocument",
              "value": {
                "sender": "did:com:1zfhgwfgex8rc9t00pk6jm6xj6vx5cjr4ngy32v",
                "recipient": "did:com:1g5cxz9p7dqux80sw5tljwg2pwz6x7zlq84km56",
                "uuid": "6a881ef0-04da-4524-b7ca-6e5e3b7e61dc",
                "metadata": {
                  "content_uri": "https://www.vargroup.it/managed-security-services/",
                  "schema_type": "",
                  "schema": {
                    "uri": "https://www.vargroup.it/managed-security-services/metadata/schema",
                    "version": "1.1.0"
                  },
                  "proof": "yes"
                },
                "content_uri": "https://www.vargroup.it/managed-security-services/",
                "checksum": {
                  "value": "7815696ecbf1c96e6894b779456d330e",
                  "algorithm": "md5"
                }
              }
            }
          ],
          "fee": {
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "250"
              }
            ],
            "gas": "200000"
          },
          "signatures": [
            {
              "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A2LI01hYyy76qNPCxq7+gYshtefvK30FeVL0TpK6VpDI"
              },
              "signature": "XIZ19jS4krMG4bDAHqJp4uSEtffb7Nnr0SzJ6clRtoJJmFhdeHgCdF1gWei0gllqOiIjs907AMxdGnpnnqhoIw=="
            }
          ],
          "memo": ""
        }
      },
      "timestamp": "2019-09-26T06:17:34Z"
    }
  ]
}
```