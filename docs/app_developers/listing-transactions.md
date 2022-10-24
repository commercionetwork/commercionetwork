# Listing past transactions
In this page you will learn how to list all the transactions that have been made in the past that include a specific message type.

## REST API
The REST API endpoint that must be called is the following: 

```
/cosmos/tx/v1beta1/txs?events=message.action%3D%27{action}%27
``` 



### Examples
**Endpoint**   
`https://lcd-testnet.commercio.network/cosmos/tx/v1beta1/txs?events=message.action%3D%27send%27`


**Output**  
```json
{
  "txs": [
    {
      "body": {
        "messages": [
          {
            "@type": "/cosmos.bank.v1beta1.MsgSend",
            "from_address": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu",
            "to_address": "did:com:15gsp698s0mgmlvqcz29tayu6a632znqgmg0upq",
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "100000000"
              }
            ]
          }
        ],
        "memo": "",
        "timeout_height": "0",
        "extension_options": [
        ],
        "non_critical_extension_options": [
        ]
      },
      "auth_info": {
        "signer_infos": [
          {
            "public_key": {
              "@type": "/cosmos.crypto.secp256k1.PubKey",
              "key": "Ay+WFw3pst8XdjgX1VULdnYQNkFmWfigIfHxPScNx4C1"
            },
            "mode_info": {
              "single": {
                "mode": "SIGN_MODE_LEGACY_AMINO_JSON"
              }
            },
            "sequence": "0"
          }
        ],
        "fee": {
          "amount": [
            {
              "denom": "ucommercio",
              "amount": "10000"
            }
          ],
          "gas_limit": "200000",
          "payer": "",
          "granter": ""
        }
      },
      "signatures": [
        "mGrIiKJ7gYo/3NAkblYJ3L8hnIhAsNoLpKaA25Ob1/05AiwUJOBOrcOmOg5Hr36NzBmINvLVKGlwU5wUyDz8VA=="
      ]
    },
    {
      "body": {
        "messages": [
          {
            "@type": "/cosmos.bank.v1beta1.MsgSend",
            "from_address": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu",
            "to_address": "did:com:1slhug5z2z5jgdpc8sjkwkpflxmd3yekjydtsum",
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "100000000"
              }
            ]
          }
        ],
        "memo": "",
        "timeout_height": "0",
        "extension_options": [
        ],
        "non_critical_extension_options": [
        ]
      },
      "auth_info": {
        "signer_infos": [
          {
            "public_key": {
              "@type": "/cosmos.crypto.secp256k1.PubKey",
              "key": "Ay+WFw3pst8XdjgX1VULdnYQNkFmWfigIfHxPScNx4C1"
            },
            "mode_info": {
              "single": {
                "mode": "SIGN_MODE_LEGACY_AMINO_JSON"
              }
            },
            "sequence": "0"
          }
        ],
        "fee": {
          "amount": [
            {
              "denom": "ucommercio",
              "amount": "10000"
            }
          ],
          "gas_limit": "200000",
          "payer": "",
          "granter": ""
        }
      },
      "signatures": [
        "gMvJ1Bq7ik/lHrOSnQKgrvh0csqfExd97pQ0YRAudwcqihJ53jxnFWDRbm/psP9IfupT5BCveAWU2u+nHar7AA=="
      ]
    }
  ],
  "tx_responses": [
    {
      "height": "4950415",
      "txhash": "16E5914CC80D9627803DCD288264D9A155CF6C67EC32346840C4FC02B67BCB0A",
      "codespace": "",
      "code": 0,
      "data": "0A060A0473656E64",
      "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"did:com:15gsp698s0mgmlvqcz29tayu6a632znqgmg0upq\"},{\"key\":\"sender\",\"value\":\"did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu\"},{\"key\":\"amount\",\"value\":\"100000000ucommercio\"}]}]}]",
      "logs": [
        {
          "msg_index": 0,
          "log": "",
          "events": [
            {
              "type": "message",
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
              ]
            },
            {
              "type": "transfer",
              "attributes": [
                {
                  "key": "recipient",
                  "value": "did:com:15gsp698s0mgmlvqcz29tayu6a632znqgmg0upq"
                },
                {
                  "key": "sender",
                  "value": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu"
                },
                {
                  "key": "amount",
                  "value": "100000000ucommercio"
                }
              ]
            }
          ]
        }
      ],
      "info": "",
      "gas_wanted": "200000",
      "gas_used": "70267",
      "tx": {
        "@type": "/cosmos.tx.v1beta1.Tx",
        "body": {
          "messages": [
            {
              "@type": "/cosmos.bank.v1beta1.MsgSend",
              "from_address": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu",
              "to_address": "did:com:15gsp698s0mgmlvqcz29tayu6a632znqgmg0upq",
              "amount": [
                {
                  "denom": "ucommercio",
                  "amount": "100000000"
                }
              ]
            }
          ],
          "memo": "",
          "timeout_height": "0",
          "extension_options": [
          ],
          "non_critical_extension_options": [
          ]
        },
        "auth_info": {
          "signer_infos": [
            {
              "public_key": {
                "@type": "/cosmos.crypto.secp256k1.PubKey",
                "key": "Ay+WFw3pst8XdjgX1VULdnYQNkFmWfigIfHxPScNx4C1"
              },
              "mode_info": {
                "single": {
                  "mode": "SIGN_MODE_LEGACY_AMINO_JSON"
                }
              },
              "sequence": "0"
            }
          ],
          "fee": {
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "10000"
              }
            ],
            "gas_limit": "200000",
            "payer": "",
            "granter": ""
          }
        },
        "signatures": [
          "mGrIiKJ7gYo/3NAkblYJ3L8hnIhAsNoLpKaA25Ob1/05AiwUJOBOrcOmOg5Hr36NzBmINvLVKGlwU5wUyDz8VA=="
        ]
      },
      "timestamp": "2022-02-03T15:48:15Z",
      "events": [
        {
          "type": "tx",
          "attributes": [
            {
              "key": "YWNjX3NlcQ==",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHUvMA==",
              "index": true
            }
          ]
        },
        {
          "type": "tx",
          "attributes": [
            {
              "key": "c2lnbmF0dXJl",
              "value": "bUdySWlLSjdnWW8vM05Ba2JsWUozTDhobkloQXNOb0xwS2FBMjVPYjEvMDVBaXdVSk9CT3JjT21PZzVIcjM2TnpCbUlOdkxWS0dsd1U1d1V5RHo4VkE9PQ==",
              "index": true
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "cmVjaXBpZW50",
              "value": "ZGlkOmNvbToxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWw4amN0dHc=",
              "index": true
            },
            {
              "key": "c2VuZGVy",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHU=",
              "index": true
            },
            {
              "key": "YW1vdW50",
              "value": "MTAwMDB1Y29tbWVyY2lv",
              "index": true
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "c2VuZGVy",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHU=",
              "index": true
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "YWN0aW9u",
              "value": "c2VuZA==",
              "index": true
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "cmVjaXBpZW50",
              "value": "ZGlkOmNvbToxNWdzcDY5OHMwbWdtbHZxY3oyOXRheXU2YTYzMnpucWdtZzB1cHE=",
              "index": true
            },
            {
              "key": "c2VuZGVy",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHU=",
              "index": true
            },
            {
              "key": "YW1vdW50",
              "value": "MTAwMDAwMDAwdWNvbW1lcmNpbw==",
              "index": true
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "c2VuZGVy",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHU=",
              "index": true
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "bW9kdWxl",
              "value": "YmFuaw==",
              "index": true
            }
          ]
        }
      ]
    },
    {
      "height": "4950416",
      "txhash": "61D30459611D66E0D970A0B588FC730F391F15AF23DFD0B407C1C178497DB928",
      "codespace": "",
      "code": 0,
      "data": "0A060A0473656E64",
      "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"did:com:1slhug5z2z5jgdpc8sjkwkpflxmd3yekjydtsum\"},{\"key\":\"sender\",\"value\":\"did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu\"},{\"key\":\"amount\",\"value\":\"100000000ucommercio\"}]}]}]",
      "logs": [
        {
          "msg_index": 0,
          "log": "",
          "events": [
            {
              "type": "message",
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
              ]
            },
            {
              "type": "transfer",
              "attributes": [
                {
                  "key": "recipient",
                  "value": "did:com:1slhug5z2z5jgdpc8sjkwkpflxmd3yekjydtsum"
                },
                {
                  "key": "sender",
                  "value": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu"
                },
                {
                  "key": "amount",
                  "value": "100000000ucommercio"
                }
              ]
            }
          ]
        }
      ],
      "info": "",
      "gas_wanted": "200000",
      "gas_used": "70312",
      "tx": {
        "@type": "/cosmos.tx.v1beta1.Tx",
        "body": {
          "messages": [
            {
              "@type": "/cosmos.bank.v1beta1.MsgSend",
              "from_address": "did:com:1ejuvfc2ydcq7ym4ks052lu45kg5xk6us0srwdu",
              "to_address": "did:com:1slhug5z2z5jgdpc8sjkwkpflxmd3yekjydtsum",
              "amount": [
                {
                  "denom": "ucommercio",
                  "amount": "100000000"
                }
              ]
            }
          ],
          "memo": "",
          "timeout_height": "0",
          "extension_options": [
          ],
          "non_critical_extension_options": [
          ]
        },
        "auth_info": {
          "signer_infos": [
            {
              "public_key": {
                "@type": "/cosmos.crypto.secp256k1.PubKey",
                "key": "Ay+WFw3pst8XdjgX1VULdnYQNkFmWfigIfHxPScNx4C1"
              },
              "mode_info": {
                "single": {
                  "mode": "SIGN_MODE_LEGACY_AMINO_JSON"
                }
              },
              "sequence": "0"
            }
          ],
          "fee": {
            "amount": [
              {
                "denom": "ucommercio",
                "amount": "10000"
              }
            ],
            "gas_limit": "200000",
            "payer": "",
            "granter": ""
          }
        },
        "signatures": [
          "gMvJ1Bq7ik/lHrOSnQKgrvh0csqfExd97pQ0YRAudwcqihJ53jxnFWDRbm/psP9IfupT5BCveAWU2u+nHar7AA=="
        ]
      },
      "timestamp": "2022-02-03T15:48:21Z",
      "events": [
        {
          "type": "tx",
          "attributes": [
            {
              "key": "YWNjX3NlcQ==",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHUvMA==",
              "index": true
            }
          ]
        },
        {
          "type": "tx",
          "attributes": [
            {
              "key": "c2lnbmF0dXJl",
              "value": "Z012SjFCcTdpay9sSHJPU25RS2dydmgwY3NxZkV4ZDk3cFEwWVJBdWR3Y3FpaEo1M2p4bkZXRFJibS9wc1A5SWZ1cFQ1QkN2ZUFXVTJ1K25IYXI3QUE9PQ==",
              "index": true
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "cmVjaXBpZW50",
              "value": "ZGlkOmNvbToxN3hwZnZha20yYW1nOTYyeWxzNmY4NHoza2VsbDhjNWw4amN0dHc=",
              "index": true
            },
            {
              "key": "c2VuZGVy",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHU=",
              "index": true
            },
            {
              "key": "YW1vdW50",
              "value": "MTAwMDB1Y29tbWVyY2lv",
              "index": true
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "c2VuZGVy",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHU=",
              "index": true
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "YWN0aW9u",
              "value": "c2VuZA==",
              "index": true
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "cmVjaXBpZW50",
              "value": "ZGlkOmNvbToxc2xodWc1ejJ6NWpnZHBjOHNqa3drcGZseG1kM3lla2p5ZHRzdW0=",
              "index": true
            },
            {
              "key": "c2VuZGVy",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHU=",
              "index": true
            },
            {
              "key": "YW1vdW50",
              "value": "MTAwMDAwMDAwdWNvbW1lcmNpbw==",
              "index": true
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "c2VuZGVy",
              "value": "ZGlkOmNvbToxZWp1dmZjMnlkY3E3eW00a3MwNTJsdTQ1a2c1eGs2dXMwc3J3ZHU=",
              "index": true
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "bW9kdWxl",
              "value": "YmFuaw==",
              "index": true
            }
          ]
        }
      ]
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "452"
  }
}
```