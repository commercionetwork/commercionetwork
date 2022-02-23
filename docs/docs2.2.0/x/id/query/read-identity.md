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
    "owner": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
    "did_document": {
      "@context": "https://www.w3.org/2019/did/v1",
      "id": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
      "publicKey": [
        {
          "id": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-1",
          "type": "Secp256k1VerificationKey2018",
          "controller": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
          "publicKeyHex": "028b722575fe90167ccae99ab06a9f155fbb11f2045b35a452e635efc57182624a"
        },
        {
          "id": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-2",
          "type": "RsaVerificationKey2018",
          "controller": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
          "publicKeyHex": "3082010a0282010100a365a9d0ac3bc66e256268ef8126b1c9acbb977ecaa140f0f738f28e6645e038bf84ccce0e53052726c8cad0cd3eeacfd2036959917355765ecf43ebc487889dad4e388787b231c8351cafc5394572046942642f6062566a90dc309f4fe910707ed6bbb310e0fe879ee31d4ed3eb74ffff0eda3b8f0cfcfb70392ce936143c13cdcb6a11bd997d0405e7d1bcd043315e7851c30bacce8985d006f794bcb50b861b90c580fee6958a668983c0ba06bd70d6165b1b73b6666ecc0818cbd69bc09aab6d497fe5c58e46bb1b4a795bb99a40d5793fd23588d8d804e5473569bfd1454d1003c2bc74de8ef9db35a00911446df32e2071a964c7b606ffc665a5d879bd0203010001"
        },
        {
          "id": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-3",
          "type": "Secp256k1VerificationKey2018",
          "controller": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h",
          "publicKeyHex": "04239ef9669a71ee8ca95f25ab729714e5240e305fa66c0d3720997a456838f9bc9f2a52979bc79ef1062514430be2852f5e4702d605937cea67ecb9c897988486"
        }
      ],
      "authentication": [
        "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-1"
      ],
      "proof": {
        "type": "LinkedDataSignature2015",
        "created": "2019-11-11T13:44:48.829363Z",
        "creator": "did:com:16ymj373t0rz2r6a57gm8ffzm6xm4euwqecta7h#keys-1",
        "signatureValue": "3045022100d2070318a640077c202137c0ac4c64ea2a9274baf3b8b20f7d2d526d881b9a2602204ca3d0141fb1ebe4a87d9efdd4f03f89862aeadb192e87183671d1a2ec9dec11"
      },
      "service": null
    }
  }
}
```