<!--
order: 3
-->

# Events

The commerciokyc module emits the following events:

## Handlers

### MsgInviteUser


| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| invite | recipient     | {recipientAddress} |
| invite | sender        | {senderAddress}           |
| invite  | sender_membership_type        | {senderMembership}               |
| message  | module        | commerciokyc    |
| message  | action        | invite          |
| message  | sender        | {senderAddress}    |

### MsgBuyMembership

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| assign_membership | owner     | {ownerAddress} |
| assign_membership | membership_type | {membershipMessage} |
| assign_membership | tsp_address | {tspAddress} |
| assign_membership | expiry_at | block.Time() |
| distribute_reward * | invite_sender | {inviteeAdress} |
| distribute_reward * | reward_coins | {rewardedCoins} |
| distribute_reward * | sender_membership_type | {inviteeMembership} |
| distribute_reward * | recipient_membership_type | {invitedMembership} |
| distribute_reward * | invite_recipient |  {invitedAddress} |
| message  | module        | commerciokyc |
| message  | action        | buy_membership |
| message  | sender        | {senderAddress} |

\* Could be not present

### MsgRemoveMembership

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| remove_membership | subscriber | {ownerAddress} |
| message | module | commerciokyc |
| message | action | remove_membership |
| message | sender | {govAddress} |



### MsgSetMembership


| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| assign_membership | owner     | {ownerAddress} |
| assign_membership | membership_type        | {membershipMessage}           |
| assign_membership | tsp_address        | {tspAddress}           |
| assign_membership | expiry_at        | block.Time()           |
| distribute_reward * | invite_sender | {inviteeAdress} |
| distribute_reward * | reward_coins | {rewardedCoins} |
| distribute_reward * | sender_membership_type | {inviteeMembership} |
| distribute_reward * | recipient_membership_type | {invitedMembership} |
| distribute_reward * | invite_recipient |  {invitedAddress} |
| message | module | commerciokyc |
| message | action | assign_membership |
| message | sender | {govAddress} |

\* Could be not present
### MsgAddTsp

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| add_tsp | tsp     | {tspAddress} |
| message | module | commerciokyc |
| message | action | add_tsp |
| message | sender | {govAddress} |




### MsgRemoveTsp

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| remove_tsp | tsp | {tspAddress} |
| message | module | commerciokyc |
| message | action | remove_tsp |
| message | sender | {govAddress} |

### MsgDepositIntoLiquidityPool

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| deposit_into_pool | depositor | {ownerAddress} |
| deposit_into_pool | amount | {amount} |
| message | module | commerciokyc |
| message | action | deposit_into_liquidity_pool |
| message | sender | {ownerAddress} |












