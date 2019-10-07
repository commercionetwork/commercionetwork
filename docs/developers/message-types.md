# Supported messages types
When [creating, signing and broadcasting a transaction](create-sign-broadcast-tx.md) you need to 
know the `type` and `value` of the message(s) you are including inside the transaction itself. 

Inside this page you can find all the types and values for the different messages that are supported by Commercio. 

## `docs`
* Share a document: [MsgSendDocument](../x/docs/tx/send-document.md#transaction-message) 
* Send a document receipt: [MsgSendDocumentReceipt](../x/docs/tx/send-document-receipt.md#transaction-message)
* Add a metadata schema as officially supported: [MsgAddSupportedMetadataSchema](../x/docs/tx/add-supported-metadata-schema.md)
* Add a trusted schema proposer: [MsgAddTrustedMetadataSchemaProposer](../x/docs/tx/add-trusted-metadata-schema-proposer.md)

## `id`
* Associate a Did Document to your Did: [MsgSetIdentity](../x/id/tx/associate-a-did-document.md)
* Request a Did deposit: [MsgRequestDidDeposit](../x/id/tx/request-did-deposit.md)
* Invalidate a Did deposit request: [MsgInvalidateDidDepositRequest](../x/id/tx/invalidate-did-deposit-request.md)
* Request a Did power up: [MsgRequestDidPowerUp](../x/id/tx/request-did-power-up.md)
* Invalidate a Did power up request: [MsgInvalidateDidPowerUpRequest](../x/id/tx/invalidate-did-power-up-request.md))

## `mint`
* Open a CDP: [MsgOpenCdp](../x/mint/tx/open-cdp.md)
* Close a CDP: [MsgCloseCdp](../x/mint/tx/close-cdp.md)

## `pricefeed`
* Set an asset's price: [MsgSetPrice](../x/pricefeed/tx/set-raw-price.md)
