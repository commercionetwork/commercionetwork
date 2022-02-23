<!--
order: 3
-->

# Events

The vbr module emits the following events:

## Handlers

### MsgIncrementBlockRewardsPool

| Type           | Attribute Key    | Attribute Value      |
| --------       | -------------    | ------------------   |
| increment_block_rewards_pool   | funder   | {funderAddress}   |
| increment_block_rewards_pool   | amount   | {ucomAmount}      |
| message        | module           | vbr                       |
| message        | action           | incrementBlockRewardsPool |
| message        | sender           | {senderAddress}           |


### MsgSetParams

| Type           | Attribute Key    | Attribute Value      |
| --------       | -------------    | ------------------   |
| new_params     | government       | {governmentAddress}  |
| new_params     | distr_epoch_identifier   | {distrEpochIdentifier}  |
| new_params     | earn_rate        | {earnRate}           |
| message        | module           | vbr                  |
| message        | action           | setParams            |
| message        | sender           | {senderAddress}      |