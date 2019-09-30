# Reading a user Did Document

## REST API
### Endpoint
```
/identities/${did}
```

### Parameters  
| Parameter | Description |
| :-------: | :---------- | 
| `did` | Address of the user for which to read the Did Document |

### Example 
#### Call
```
http://localhost:1317/identities/did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke
```

#### Response
```json
{
  "height": "0",
  "result": {
    "owner": "did:com:15erw8aqttln5semks0vnqjy9yzrygzmjwh7vke",
    "did_document": {
      "uri": "https://example.com/did-document",
      "content_hash": "9c5ef543dc05e7927da16e8d8a24372f0d064979e226a70cdea40a031d1daf51"
    }
  }
}
```