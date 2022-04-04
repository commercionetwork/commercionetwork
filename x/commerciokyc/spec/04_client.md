<!--
order: 4
-->

# Client

## CLI

A user can query and interact with the `commerciomint` module using the CLI.

### Query

The `query` commands allow users to query `commerciomint` state.

```bash
commercionetworkd query commerciokyc --help
```

#### invites

The `invites` command gets all invites:

```bash
commercionetworkd query commerciokyc invites [flags]
```

Example:

```bash
commercionetworkd query commerciokyc invites
```

Example Output:

```bash
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

#### trusted-service-providers

The `trusted-service-providers` command allows users to query all trusted service providers.

```bash
commercionetworkd query commerciokyc trusted-service-providers [flags]
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

#### pool-funds

The `pool-funds` command allows users to query a given ABR pool funds for the `commerciokyc` module.

```bash
commercionetworkd query commerciokyc pool-funds [flags]
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



## gRPC

A user can query the `commerciokyc` module using gRPC endpoints.

### Invite

The `Invite` endpoint allows users to query a given proposal.

```bash
commercionetwork/invite
```

Example:

```bash
grpcurl -plaintext \
    -d '{"invite":"1"}' \
    localhost:9090 \
    commercionetwork/invite
```

Example Output:

```bash
{
  "invite": {
  }
}
```



## REST

A user can query the `commerciokyc` module using REST endpoints.

### invites

The `invites` endpoint allows users to query a given proposal.

```bash
/commerciokyc/invites
```

Example:

```bash
curl localhost:1317/commerciokyc/invites
```

Example Output:

```bash
{
  "invites": {

  }
}
```

