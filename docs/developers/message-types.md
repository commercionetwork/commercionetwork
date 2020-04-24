# Supported messages types
When [creating, signing and broadcasting a transaction](create-sign-broadcast-tx.md) you need to 
know the `type` and `value` of the message(s) you are including inside the transaction itself. 

Inside this page you can find all the types and values for the different messages that are supported by Commercio. 


## `id`
* Associate a Did Document to your Did: [MsgSetIdentity](../x/id/#associating-a-did-document-to-your-identity)
* Request a Did power up: [MsgRequestDidPowerUp](../x/id/#did-power-up)
* Change Did power up status: [MsgChangePowerUpStatus](../x/id/#change-did-power-up-status-wip)

## `docs`
* Share a document: [MsgSendDocument](../x/docs/#sending-a-document) 
* Send a document receipt: [MsgSendDocumentReceipt](../x/docs/#sending-a-document-reading-receipt)

## `pricefeed`

* Add an oracle: [MsgAddOracle](../x/pricefeed/#adding-an-oracle)
* Set an asset's price: [MsgSetPrice](../x/pricefeed/#set-a-price-for-an-asset)


## `commerciomint`
* Open a CDP: [MsgOpenCdp](../x/commerciomint/#open-a-cdp)
* Close a CDP: [MsgCloseCdp](../x/commerciomint/#close-a-cdp)
* Set CDP collateral rate: [MsgSetCdpCollateralRate](../x/commerciomint/#set-cdp-collateral-rate)

## `memberships`
* Sending an invite: [MsgInviteUser](../x/memberships/#sending-an-invite)
* Buy a membership: [MsgBuyMembership](../x/memberships/#buying-a-membership-2)
* Adding a TSP (Trust Service Provider): [MsgAddTsp](../x/memberships/#adding-a-tsp)
* Deposit into reward pool: [MsgDepositIntoLiquidityPool](../x/memberships/#deposit-into-reward-pool)
* Set user membership: [MsgSetMembership](../x/memberships/#set-user-membership)

