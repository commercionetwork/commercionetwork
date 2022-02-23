<!--
order: 3
-->

# Events

The commerciomint module emits the following events:

## Handlers

### MsgMintCCC

| Type           | Attribute Key    | Attribute Value      |
| --------       | -------------    | ------------------   |
| new_position   | depositor        | {depositorAddress}   |
| new_position   | amount_deposited | {ucomAmount}         |
| new_position   | minted_coins     | {creditsCoins}       |
| new_position   | position_id      | {position_id}        |
| new_position   | timestamp        | {position_createdAt} |
| message        | module           | commerciomint        |
| message        | action           | mintCCC              |
| message        | sender           | {senderAddress}      |

### MsgBurnCCC

| Type       | Attribute Key    | Attribute Value   |
| --------   | -------------    |----------------   |
| burned_ccc | position_id      | {position_id}     |
| burned_ccc | sender           | {userAddress}     |
| burned_ccc | amount           | {burnAmount}      |
| burned_ccc | position_deleted | {bool}            |
| message    | module           | commerciomint     |
| message    | action           | burnCCC           |
| message    | sender           | {senderAddress}   |

### MsgSetParams

| Type                | Attribute Key       | Attribute Value |
| ------------------- | ------------------- | --------------- |
| new_freeze_period   | new_params          | {params}        |
| message             | module              | commerciomint   |
| message             | action              | set_params      |
| message             | sender              | {senderAddress} |
