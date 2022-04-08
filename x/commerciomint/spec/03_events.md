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
| transfer (ucommercio) | recipient     | {moduleAddress}   |
| transfer (ucommercio) | sender        | {depositorAddress} |
| transfer (ucommercio) | amount        | {ucomAmount}      |
| transfer (uccc) | recipient     | {depositorAddress}   |
| transfer (uccc) | sender        | {moduleAddress} |
| transfer (uccc) | amount        | {creditsCoins}      |
| message        | action           | mintCCC              |
| message        | sender           | {senderAddress}      |
| message        | sender           | {depositorAddress}      |
| message        | sender           | {moduleAddress}      |

### MsgBurnCCC (WIP)

| Type       | Attribute Key    | Attribute Value   |
| --------   | -------------    |----------------   |
| burned_ccc | position_id      | {position_id}     |
| burned_ccc | sender           | {userAddress}     |
| burned_ccc | amount           | {burnAmount}      |
| burned_ccc | position_deleted | {bool}            |
| message    | action           | burnCCC           |
| message    | sender           | {senderAddress}   |

### MsgSetParams

| Type                | Attribute Key       | Attribute Value |
| ------------------- | ------------------- | --------------- |
| new_params          | params              | {params}        |
| message             | action              | setParams      |
| message             | sender              | {senderAddress} |
