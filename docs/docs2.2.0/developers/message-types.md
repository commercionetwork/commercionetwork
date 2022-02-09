# Supported messages types
When [creating, signing and broadcasting a transaction](create-sign-broadcast-tx.md) you need to 
know the `type` and `value` of the message(s) you are including inside the transaction itself. 

Inside this page you can find all the types and values for the different messages that are supported by Commercio. 


## `id`
* Associate a Did Document to your Did: [MsgSetIdentity](../x/id/#associating-a-did-document-to-your-identity)
* Request a Did power up: [MsgRequestDidPowerUp](../x/id/#did-power-up)
* Change Did power up status: [MsgChangePowerUpStatus](../x/id/#change-did-power-up-status-wip)

## `docs`
* Share a document: [MsgSendDocument](../x/documents/#sending-a-document) 
* Send a document receipt: [MsgSendDocumentReceipt](../x/documents/#sending-a-document-reading-receipt)

## `commerciomint`
* Open a ETP: [MsgMintCCC](../x/commerciomint/#mint-commercio-cash-credit-ccc)
* Close a ETP: [MsgBurnCCC](../x/commerciomint/#burn-commercio-cash-credit-ccc)
* Set ETP conversion rate: [MsgSetCCCConversionRate](../x/commerciomint/#set-ccc-conversion-rate)

## `commerciokyc`
* Sending an invite: [MsgInviteUser](../x/commerciokyc/#sending-an-invite)
* Buy a membership: [MsgBuyMembership](../x/commerciokyc/#buying-a-membership-2)
* Adding a TSP (Trust Service Provider): [MsgAddTsp](../x/commerciokyc/#adding-a-tsp)
* Removing a TSP (Trust Service Provider): [MsgRemoveTsp](../x/commerciokyc/#removing-a-tsp)
* Deposit into reward pool: [MsgDepositIntoLiquidityPool](../x/commerciokyc/#deposit-into-reward-pool)
* Set user membership: [MsgSetMembership](../x/commerciokyc/#set-user-membership)



