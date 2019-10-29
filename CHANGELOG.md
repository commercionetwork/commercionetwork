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