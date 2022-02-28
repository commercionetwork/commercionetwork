<!--
order: 5
-->

# Client (WiP)

<!-- 
## CLI (WiP)

## gRPC (WiP)

## REST (WiP) 
-->

## Queries


### List sent documents

#### CLI

```bash
cncli query docs sent-documents [address]
```


#### REST

```
/docs/{address}/sent
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current sent documents |

##### Example 

Getting sent docs from `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/docs/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/sent
```

### List received documents

#### CLI

```bash
cncli query docs received-documents [address]
```

#### REST

```
/docs/{address}/received
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current received documents |


##### Example 

Getting docs for `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/docs/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/received
```



### List sent receipts

#### CLI

```bash
cncli query docs sent-receipts [address]
```

#### REST

```
/receipts/{address}/sent
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current sent receipts |

##### Example 

Getting sent receipts from `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/receipts/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/sent
```

### List received receipts

#### CLI

```bash
cncli query docs received-receipts [address]
```
   

#### REST

```
/receipts/{address}/received
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of the user for which to read current received receipts |


##### Example 

Getting receipts for `did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf`:

```
http://localhost:1317/receipts/did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf/received
```

### List receipts associated to a certain document

#### CLI

```bash
cncli query docs documents-receipts [documentUUID]
```

#### REST

```
documents/document/{UUID}/receipts
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `UUID` | Document ID of the document for which to read current received receipts |

##### Example 

Getting receipts associated to the document with ID `d83422c6-6e79-4a99-9767-fcae46dfa371`:


```
http://localhost:1317/document/d83422c6-6e79-4a99-9767-fcae46dfa371/receipts
```

