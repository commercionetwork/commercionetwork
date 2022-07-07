<!--
order: 4
-->

# Client

## Transactions

### Invite

Invite user to buy a membership

```bash
commercionetworkd tx commerciokyc invite \
  [subscriber]
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `subscriber`               | Address of the account you want to invite |


### Buy a membership

Tsp buy a membership for subscriber

```bash
commercionetworkd tx commerciokyc buy \
  [subscriber] \
  [membership-type] 
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `subscriber`      | Address of the account you want to buy a membership for |
| `membership-type` | Membership type to buy |

### Assign a membership

As government, assign membership to a user

```bash
commercionetworkd tx commerciokyc assign-membership \
  [subscriber] \
  [membership-type] 
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `subscriber`      | Address of the account you want to assign the membership |
| `membership-type` | Membership type to assign |


### Remove a membership

As government, remove membership of a user.

```bash
commercionetworkd tx commerciokyc remove-membership \
  [subscriber] \
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `subscriber`      | Address of the account from which you want to remove the membership |



### Add Trusted Service Provider

Government add a tsp.

```bash
commercionetworkd tx commerciokyc add-tsp \
  [tsp-address]
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `tsp-address`      | Address of the account you want to bacome a tsp |


### Remove Trusted Service Provider

Government remove a tsp.

```bash
commercionetworkd tx commerciokyc remove-tsp \
  [tsp-address]
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `tsp-address`      | Address of the account you want to remove from group of tsps |


### Deposit into pool (available soon)

Increments the membership rewards pool's liquidity by the given amount

```bash
commercionetworkd tx commerciokyc deposit \
  [amount]
```

**Parameters:**

| Parameter | Description |
| :------- | :---------- | 
| `amount`      | Amount of ucommercio tokens to deposit |



A user can query and interact with the `commerciokyc` module using the CLI.

## Queries

The `query` commands allow users to query `commerciokyc` state.

```bash
commercionetworkd query commerciokyc --help
```

### Invites


#### CLI


The `invites` command gets all invites:

```bash
commercionetworkd query commerciokyc invites 
```

Example:

```bash
commercionetworkd query commerciokyc invites
```

Example Output:

```yaml
invites:
- sender: did:com:1f06vm4x0ae978rtxvz5he82pg4mty3an6elt9x
  sender_membership: black
  status: "1"
  user: did:com:1xx88le4t8ateql77mzzyrg0damf43tt0qw2xms
- sender: did:com:1t5fz439f49zv39pmh73c2lvuhwfzqj0ze3kzj2
  sender_membership: black
  status: "1"
  user: did:com:1xz6ues73ahw5jdx9ukv8ruey5jqfg6qay0e6j8
```


#### REST

```
/commercionetwork/commerciokyc/invites
```

##### Example 

```
https://localhost:1317/commercionetwork/commerciokyc/invites
```

#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciokyc.Query/Invites
```

##### Example

```bash
grpcurl -plaintext \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciokyc.Query/Invites
```

##### Response
```json
{
  "invites": [
    {
      "sender": "did:com:1gdcxa02g5l3cm0mgqfsz3ju42jyur82z3cx45p",
      "senderMembership": "black",
      "user": "did:com:109fup66yms0e559l54tjaawz0rsj2gxvqzlr9z",
      "status": "1"
    },
    ...
  ],
  "pagination": {
    "total": "100"
  }
}
```

### Invite


#### CLI


The `invite` command gets user invite:

```bash
commercionetworkd query commerciokyc invites \
  [user]
```

Example:

```bash
commercionetworkd query commerciokyc invite \
  did:com:1xx88le4t8ateql77mzzyrg0damf43tt0qw2xms
```

Example Output:

```yaml
invite:
- sender: did:com:1f06vm4x0ae978rtxvz5he82pg4mty3an6elt9x
  sender_membership: black
  status: "1"
  user: did:com:1xx88le4t8ateql77mzzyrg0damf43tt0qw2xms
```

#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciokyc.Query/Invite
```

##### Example

```bash
grpcurl -plaintext \
    -d '{"address":"did:com:1gdcxa02g5l3cm0mgqfsz3ju42jyur82z3cx45p"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciokyc.Query/Invite
```

##### Response
```json
{
  "invite": {
    "sender": "did:com:1mj9h87yqjel0fsvkq55v345kxk0n09krtfvtyx",
    "senderMembership": "black",
    "user": "did:com:1gdcxa02g5l3cm0mgqfsz3ju42jyur82z3cx45p",
    "status": "1"
  }
}
```

#### REST

```
/commercionetwork/commerciokyc/{address}/invite
```

##### Example 

```
https://localhost:1317/commercionetwork/commerciokyc/did:com:1gdcxa02g5l3cm0mgqfsz3ju42jyur82z3cx45p/invite
```

### Memberships


#### CLI


The `memberships` command gets all memberships:

```bash
commercionetworkd query commerciokyc memberships 
```

Example:

```bash
commercionetworkd query commerciokyc memberships
```

Example Output:

```yaml
memberships:
- expiry_at: "2022-03-22T00:00:00Z"
  membership_type: black
  owner: did:com:1q8mkesv6kcyr8ft69mvtmy6lxzfvn5y6ywhgh9
  tsp_address: did:com:1mj9h87yqjel0fsvkq55v345kxk0n09krtfvtyx
- expiry_at: "2022-06-25T19:12:45.276830498Z"
  membership_type: bronze
  owner: did:com:1py237er2h2jdgdpzggeqmat556u65fv6ql22ya
  tsp_address: did:com:1x4hpem28uhrlh2sdvf3a2f5rw56jtvsgmjz5yp
```
#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciokyc.Query/Memberships
```

##### Example

```bash
grpcurl -plaintext \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciokyc.Query/Memberships
```

##### Response
```json
{
  "memberships": [
    {
      "owner": "did:com:1q8mkesv6kcyr8ft69mvtmy6lxzfvn5y6ywhgh9",
      "tspAddress": "did:com:1mj9h87yqjel0fsvkq55v345kxk0n09krtfvtyx",
      "membershipType": "black",
      "expiryAt": "2022-03-22T00:00:00Z"
    },
    ...
  ],
  "pagination": {
    "total": "92"
  }
}
```

#### REST


```
/commercionetwork/commerciokyc/memberships
```

##### Example 

```
https://localhost:1317/commercionetwork/commerciokyc/memberships
```






### Membership


#### CLI


The `membership` command gets user membership:

```bash
commercionetworkd query commerciokyc memberships \
  [user] \
  
```

Example:

```bash
commercionetworkd query commerciokyc membership \
  did:com:1q8mkesv6kcyr8ft69mvtmy6lxzfvn5y6ywhgh9
```

Example Output:

```yaml
membership:
  expiry_at: "2022-03-22T00:00:00Z"
  membership_type: black
  owner: did:com:1q8mkesv6kcyr8ft69mvtmy6lxzfvn5y6ywhgh9
  tsp_address: did:com:1mj9h87yqjel0fsvkq55v345kxk0n09krtfvtyx
```

#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciokyc.Query/Membership
```

##### Example

```bash
grpcurl -plaintext \
    -d '{"address":"did:com:1q8mkesv6kcyr8ft69mvtmy6lxzfvn5y6ywhgh9"}' \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciokyc.Query/Membership
```

##### Response
```json
{
  "membership": {
    "owner": "did:com:1q8mkesv6kcyr8ft69mvtmy6lxzfvn5y6ywhgh9",
    "tspAddress": "did:com:1mj9h87yqjel0fsvkq55v345kxk0n09krtfvtyx",
    "membershipType": "black",
    "expiryAt": "2022-03-22T00:00:00Z"
  }
}
```

#### REST


```
/commercionetwork/commerciokyc/memberships/{address}
```

Parameters:

| Parameter | Description |
| :-------: | :---------- | 
| `address` | Address of membership user |


##### Example 

```
https://localhost:1317/commercionetwork/commerciokyc/memberships/did:com:1q8mkesv6kcyr8ft69mvtmy6lxzfvn5y6ywhgh9
```


### Trusted Service Providers

#### CLI


The `trusted-service-providers` command allows users to query all trusted service providers.

```bash
commercionetworkd query commerciokyc trusted-service-providers 
```


Example:

```bash
commercionetworkd query commerciokyc trusted-service-providers
```

Example Output:

```bash
tsps:
- did:com:1t5fz439f49zv39pmh73c2lvuhwfzqj0ze3kzj2
- did:com:1cc65t29yuwuc32ep2h9uqhnwrregfq230lf2rj
- did:com:14rcpqu0y8jgjrc823ejylgjnsh2jkkeg8kchl3
```
#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciokyc.Query/Tsps
```

##### Example

```bash
grpcurl -plaintext \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciokyc.Query/Tsps
```

##### Response
```json
{
  "tsps": [
    "did:com:1mj9h87yqjel0fsvkq55v345kxk0n09krtfvtyx",
    "did:com:1ft3ggfazm9yakmhl79r0qukgufesadkw3xpsmx",
    "did:com:1x4hpem28uhrlh2sdvf3a2f5rw56jtvsgmjz5yp"
  ]
}
```

#### REST

```
/commercionetwork/commerciokyc/tsps
```

##### Example 

```
https://localhost:1317/commercionetwork/commerciokyc/tsps
```





### Pool Funds


#### CLI


The `pool-funds` command allows users to query a given ABR pool funds for the `commerciokyc` module.

```bash
commercionetworkd query commerciokyc pool-funds 
```

Example:

```bash
commercionetworkd query commerciokyc pool-funds
```

Example Output:

```bash
funds:
- amount: "974677500000"
  denom: ucommercio
```

#### gRPC
Endpoint:

```
commercionetwork.commercionetwork.commerciokyc.Query/Funds
```

##### Example

```bash
grpcurl -plaintext \
    localhost:9090 \
    commercionetwork.commercionetwork.commerciokyc.Query/Funds
```

##### Response
```json
{
  "funds": [
    {
      "denom": "ucommercio",
      "amount": "10674787750000"
    }
  ]
}
```

#### REST

```
/commercionetwork/commerciokyc/funds
```

##### Example 

```
https://localhost:1317/commercionetwork/commerciokyc/funds
```
