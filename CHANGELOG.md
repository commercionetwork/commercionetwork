# Version 3.0.0

* Update cosmos sdk to v0.42.10
* Add protobuf
* Reimplemented docs module
* Reimplemented did module
* Reimplemented commerciokyc module
* Reimplemented commerciomin module
* Reimplemented government module
* Reimplemented vbr module
* Add epochs module
* Remove upgrade custom module
* Reimplemented partial tests module
* Add IBC module
* Add Wasm module
* Enable cosmos upgrade standard module
* Enable cosmos gov standard module
* Update npm dependencies

# Version 3.0.0-rc5

* Add protobuf
* Reimplemented docs module
* Reimplemented did module
* Reimplemented commerciokyc module
* Reimplemented commerciomin module
* Reimplemented government module
* Reimplemented vbr module
* Add epochs module
* Remove upgrade custom module
* Reimplemented partial tests module
* Add IBC module
* Add Wasm module
* Enable cosmos upgrade standard module
* Enable cosmos gov standard module
* Update npm dependencies
# Version 3.0.0-rc4
* Improve modules tests
* Fix documents store bug
* Fix minor bugs
* Fix query issues
* Add events
* Complete command line interface
* Remove hardcoded denom 

# Version 3.0.0-rc3
* Fix commerciomint migrate
* Fix commerciokyc migrate
* Fix documents migrate
* Add commerciomint tests
* Add commerciokyc tests
* Add documents tests
* Fix documents tests
# Version 3.0.0-rc2
* Fix commerciomint and commerciokyc migrate
* Add commerciomint and commerciokyc tests

# Version 3.0.0-rc1

* Fix cli command set identity
* Add hour identifier to vbr
# Version 3.0.0-beta.2

* Module Did
   * Add historization of Did Document
   * Improve tests and coverage
* Module CommercioKyc
   * Delete "remove memberships" trigger
   * Fix bug create of commerciomint position.
   * Improve tests and coverage

# Version 3.0.0-beta.1

* Complete commerciokyc module with epochs
* Use params in commerciokyc
* Add/Fix commerciomint tests
* Add/Fix commerciokyc tests
* Add/Fix vbr tests
# Version 3.0.0-alpha.3

* Fix messages bug Commerciomint
* Add Commerciomint tests
* Partial fix Epochs tests
* Improved Commerciomint  coverage
* Update Ante
* Add Ante tests
* Fix Did migrate
# Version 3.0.0-alpha.2

* Update commerciokyc module ABR for green membership
* Convert ABR reward to uccc
* Remove Vbr reward reate and and automatic-withdraw
* Add Vbr earn rate
* Update did to new W3C standards
* Remove proof form did
* Add epochs module and applied it to VBR module
* Add/fix tests*

# Version 3.0.0-alpha.1

New alpha release 3.0.0

# Version 2.2.0

# Version 2.2.0-pre.1
* Rewrite module membership to commercioKyc
* Rewrite module commercioMint
* Add migrate 2.2.0
* Add module update


# Version 2.1.2

* Security fix
* Added DDO pub keys duplicated check
* Added control to avoid multiple insertion of opencdp messages in the same block

# Version 2.1.1
## Bug fix

- Fixed tendermint/tmksm comunication with tendermint fork
- Fixed proof calculation in did document
- Fixed named `PublicKey`  to `PublicKeyPem` 
- Removed add-genesis-membership, assign black membership to government
- Deleted old cdp queries files

## Additions
- Add rest api for tumbler address
- Add documentations


# Version 2.1.0
## Bug fix

- Fixed p2p memory issue
- Removed aliasing from commerciomint
- Removed verify credentials
- Fixed price feed command flags
- Added invariants for Commerciomint and tests
- Fixed credit risk query pool 
- Fixed build issues
- Fixed price feed issues
- Removed module gov

## Additions
- Upgrade to cosmos sdk 0.38.3
- Added do_sign to shareDoc
- Added cdp auto liquidate
- Added cncli command module id:
   - getSetIdentityCommand
   - getCmdQueryPowerUpRequest
- Added cncli command module commerciomint :
  - setCollateralRateCmd
  - getCdpCollateralRate
- Added cncli command module creditrisk:
  - getPoolFundsCmd
- Added cncli command module government:
   - getCmdSetTumblerAddress
   - getCmdGetTumblerAddr
- Added cncli command module pricefeed:
  - getCmdBlacklistedDenoms
  - GetCmdBlacklistDenom
- Added rest query for Vbr membership:
  - /memebership/did:com:1address


# Version 1.5.0
## Bug fix

- Fixed incompatibility with cosmos sdk in the management of the rewards of the Vbr module
- Remove AppendIfMissing from Id module
- Added invariants for Commerciomint and tests
- Redesign Membership module and fixed tests 
- Modify Membership invariant

## Additions
- Added cncli command module membership:
   - getCmdDepositIntoPool
   - getCmdGovAssignMembership
   - getCmdInviteUser
   - getCmdBuyMembership
   - getCmdGetTrustedServiceProviders
   - getCmdGetInvitesForUser
   - getCmdGetInvites
- Added cncli command module vbr:
   - getRetrieveBlockRewardsPoolFunds
   - government:
   - getCmdGetGovernmentAddr
- Added cncli command module commerciomint :
  - openCDPCmd
  - closeCDPCmd
- Added cncli command module government:
   - getCmdGetGovernmentAddr
- Added rest  query for Vbr module

# Version 1.4.0

## Bug fix

- Duplication of receipts prevented
- Remove AppendIfMissing from Docs module
- Correct deserialized RawPrices and Cdps
- Use "require" instead of "assert" on every test source file
- Added invariants for docs and memberships
- Fix test execution for memberships
## Additions

- Added cdp query and tx command to cncli
- Removed the NFT module completely


# Version 1.3.4
## Bug fixes
- Renamed `mint` module to `commerciomint` to increase clearness
- Updated the old `context` type in `id` module with the new correct one https://www.w3.org/ns/did/v1
- Update `serialize-javascript` version to `2.1.1` to avoid old version's vulnerability
- Fixed a bug that let add document's receipts with the same UUID
- Added migrate for version `1.3.3` and version `1.3.4`

# Version 1.3.3
## Bug fixes
- update to btcd 0.20.1 to avoid go.sum problems.

# Version 1.3.2
## Bug fixes
- Fixed a bug while checking the validity of the authentication key while handling `MsgSetIdentity` messages
- Fixed how the UUID are validated (#63) 
- Removed double `docs` entry inside the `cncli query` command (#60)
- Fixed a bug in mint module's query

## Changes
**CommercioID**
- The `status` field is no longer required when using `MsgRequestDidDeposit` and `MsgRequestDidPowerUp` 

# Version 1.3.1
## Bug fixes
- Fixed a bug inside the migration command
- Fixed a bug during the serialization of the `bank` module genesis state

# Version 1.3.0
## Bug fixes
- Fixed the export command (#48)
- Fixed the TBR formula (#49)

## Changes 
**CommercioID**
- Implemented the pairwise Did power up system (#40)
- Changes the `MsgSetIdentity` so that it now requires a full Did Document inside its `value` field (#47)

**CommercioDOCS**
- Implemented the minimum fees when sending a `MsgShareDocument` (#38)

**CommercioMEMBERSHIP**
- Changed how the membership can be purchased.  
  It is now required to be invited and can be purchased on-chain using Commercio Cash Credits (#45) 

## Additions
- Implemented the `mint` module (#42)
- Implemented the possibility for the government to block specific accounts from sending tokens (#46)

## Migration
In order to migrate from v1.2.x to v1.3.0 you can use the following command:

```shell
cnd migrate v1.3.0 [genesis-file-path] --chain-id=<chain_id>
```

# Version 1.2.1
## Bug fixes
- Fixed a bug inside the migration command 

# Version 1.2.0
## Changes
**CommercioID (#30)**
- Changed the contents of `MsgSetIdentity`  
   - Now it is required to specify the hash of the Did Document's contents when setting a Did Document
     associated to an account.  

**CommercioDOCS (#31)**
- Changed the contents of `MsgSharedocument`
   - Added the possibility of sending the same document to multiple recipients
   - Added the possibility of specifying an encryption key for each recipient to be used when
     wanting to hide sensitive data inside the message itself
   - Removed the `metadata.proof` field
- Changed how the documents are stored inside the chain.  
   This should grant lower costs while sending a transaction containing a `MsgShareDocument` message.

**CommercioTBR (#34)**
- Fixed some bugs
- Added a genesis command to properly initialize the TBR pool

## Additions
- Implemented the `pricefeed` module (#33)

## Migration
In order to migrate from v1.1.0 to v1.2.0 you can use the following command:

```shell
cnd migrate v1.2.0 [genesis-file-path] --chain-id=<chain_id>
```

# Version 1.1.0
## Changes
**CommercioDOCS (#22)**
- Improved the contents of `MsgShareDocument`

## Additions
- Implemented the `memberships` module (#18)
- Implemented the `government` module (#22)
- Implemented the `tbr` module (#23)
- Implemented the `accreditations` module (#24)

## Migrate
In order to migrate from version 1.0.2 to 1.1.0, the chain needs to be **completely wiped**. 

# Version 1.0.2
- Updated the Cosmos SDK to v0.36.0-rc5

# Version 1.0.1
- Updated the Cosmos SDK to v0.35.0

# Version 1.0.0
- First release of the Commercio.network blockchain
