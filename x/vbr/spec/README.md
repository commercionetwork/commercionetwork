<!--
order: 0
title: Vbr Overview
parent:
  title: "vbr"
-->

# `Vbr`

## Overview

The Validator Block Rewards (VBR) halving mechanism is a critical component of the [Commercio.network](https://docs.commercio.network/). It's designed to ensure a sustainable and balanced distribution of rewards within the network. This document introduces a new VBR halving mechanism, detailing the changes and their implications for network participants.

## Purpose

The revision of the VBR halving mechanism aims to align rewards more closely with the growth and stability of the network. By adjusting the rewards according to the pool size, it ensures a more equitable distribution and incentivizes long-term participation and investment in the network.

## The VBR Halving Mechanism

The revised halving mechanism adjusts staking rewards based on the total size of the staked pool. The rewards decrease incrementally as the pool grows, following a predetermined schedule. This approach is designed to manage the rate of new token generation, ensuring a sustainable ecosystem.

### Reward Schedule

The new VBR rewards are structured as follows:

- **From 12,500,000 to 6,250,000:** Rewards are set at 50% VBR.
- **From 6,250,000 to 3,125,000:** Rewards are reduced to 25% VBR.
- **From 3,125,000 to 1,562,500:** Rewards are further reduced to 12.5% VBR.
- **From 1,562,500 to 781,250:** Rewards are set at 6.25% VBR.
- **From 781,250 to 390,625:** Rewards are reduced to 3.125% VBR.
- **From 390,625 to 195,312:** Rewards are set at 1.5625% VBR.
- **From 195,312 to 97,656:** Rewards are reduced to 0.78125% VBR.
- **From 97,656 to 48,828:** Rewards are set at 0.390625% VBR.
- **From 48,828 to 24,414:** Rewards are reduced to 0.1953125% VBR.

## Impact on Stakeholders

### Validators
Validators will see a change in their reward structure, which now depends on the total pool size. This change encourages validators to focus on network security and performance, as rewards are more closely tied to the overall health and size of the network.

### Delegators
Delegators are incentivized to participate in staking, especially in the early stages of the network's growth. The new mechanism ensures that rewards are distributed more evenly over time, reducing the impact of large fluctuations in reward rates.

### Network Stability
By gradually reducing the reward rate as the pool size increases, the new mechanism aims to prevent rapid inflation and promote a stable economic environment within the network.

## VBR Pool Reconstitution

In a significant move to sustain the robustness of the reward system, Commercio.network S.p.A. has taken a proactive approach to refinance the VBR pool out of its pockets. Recognizing the importance of maintaining a healthy and incentivized network environment, we replenished the depleted VBR pool back to 3,125,000 from its own token reserves. This decisive action was necessary to maintain the  network's long-term viability and stakeholder value.

To align with the new VBR halving mechanism, the incentive rate has now been manually adjusted to 12.5% VBR. This adjustment ensures that the reward distribution remains consistent with the network's revised economic model. Furthermore, to facilitate a seamless and efficient operation of this mechanism, plans are underway to update the VBR smart contract. This update will automate the reward adjustment process, reflecting the new halving mechanism directly within the contract's functionality.


## Conclusion

The new VBR halving mechanism is a strategic update to the Commercio.network's economic model. It's designed to balance the distribution of rewards, promote network growth and stability, and encourage long-term participation from validators and delegators. This update is a step towards ensuring a sustainable and prosperous future for the Commercio.network.


## Contents

1. **[State](01_state.md)**
2. **[Messages](02_messages.md)**
   - [Increment Block Rewards Pool](02_messages.md#Increment-block-rewards-pool)
   - [Set Parameters](02_messages.md#Set-parameters)
3. **[Events](03_events.md)**
   - [Handlers](03_events.md#handlers)
4. **[Parameters](04_params.md)**
5. **[Client](05_client.md)**