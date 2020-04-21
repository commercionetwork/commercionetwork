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
* Request a Did power up: [MsgRequestDidPowerUp](../x/id/tx/request-did-power-up.md)
* Change Did power up status: [MsgChangePowerUpStatus](../x/id/tx/invalidate-did-power-up-request.md))


## `commerciomint`
* Open a CDP: [MsgOpenCdp](../x/commerciomint/tx/open-cdp.md)
* Close a CDP: [MsgCloseCdp](../x/commerciomint/tx/close-cdp.md)
* Set collateral rate: [MsgSetCdpCollateralRate](../x/commerciomint/#)

## `memberships`
* Sending an invite: [MsgInviteUser](../x/memberships/#sending-an-invite)
* Buy a membership: [MsgBuyMembership](../x/memberships/#buying-a-membership-2)
* Adding a TSP (Trust Service Provider): [MsgAddTsp](../x/memberships/#adding-a-tsp)
* Deposit into reward pool: [MsgDepositIntoLiquidityPool](../x/memberships/#deposit-into-reward-pool)
* Set user membership: [MsgSetMembership](../x/memberships/#set-user-membership)

## `pricefeed`

* Add an oracle: [MsgAddOracle](../x/pricefeed/#adding-an-oracle)
* Set an asset's price: [MsgSetPrice](../x/pricefeed/#set-a-price-for-an-asset)

