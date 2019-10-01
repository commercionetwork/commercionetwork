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