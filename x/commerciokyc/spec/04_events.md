<!--
order: 4
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
| message  | sender        | {senderAddress}    |

### MsgBuyMembership

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| assign_membership | owner     | {ownerAddress} |
| assign_membership | membership_type        | {membershipMessage}           |
| assign_membership | tsp_address        | {tspAddress}           |
| assign_membership | expiry_at        | block.Time()           |
| message  | module        | commerciokyc               |
| message  | action        | multisend          |
| message  | sender        | {senderAddress}    |


### MsgRemoveMembership

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |




### MsgSetMembership


| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |



### MsgAddTsp

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |

### MsgRemoveTsp

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |


### MsgDepositIntoLiquidityPool

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |






