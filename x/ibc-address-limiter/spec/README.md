<!--
order: 0
title: Ibc-address-limiter Overview
parent:
  title: "Ibc-address-limiter"
-->

# IBC Address Limiter

## Concept

The IBC Address Limiter module is responsible for adding a governance-configurable address limit to IBC transfers.
This is intended to handle a list of addresses that have permission to send tokens via IBC transfers.

The architecture of this package is a minimal go package which implements an [IBC Middleware](https://github.com/cosmos/ibc-go/blob/f57170b1d4dd202a3c6c1c61dcf302b6a9546405/docs/ibc/middleware/develop.md) that wraps the [ICS20 transfer](https://ibc.cosmos.network/main/apps/transfer/overview.html) app, and calls into a cosmwasm contract.
All the actual IBC address limiting logic is then implemented in the cosmwasm contract. 
The Cosmwasm code can be found in the [`contracts`](./contracts/) package, with bytecode findable in the [`bytecode`](./bytecode/) folder.

## Code structure

As mentioned at the beginning, the Go code is a relatively minimal ICS 20 wrapper, that dispatches relevant calls to a cosmwasm contract that implements the address limiting functionality.

## Contents

1. **[Go Middleware](01_go_middleware.md)**
2. **[Contract](02_contract.md)**
3. **[Contract state](03_contract_state.md)**
4. **[Contract messages](04_contract_messages.md)**
   - [Execute](04_contract_messages.md#ExecuteMsg)
   - [Sudo](04_contract_messages.md#SudoMsg)
   - [Query](04_contract_messages.md#QueryMsg)
5. **[Simulation](05_simulation.md)**