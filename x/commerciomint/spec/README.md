<!--
order: 0
title: Commerciomint Overview
parent:
  title: "commerciomint"
-->

# CommercioMint 

## Abstract

This document specifies the commerciomint module of the Commercio Network.

The `commerciomint` module is the one that allows you to create Exchange Trade Position (*ETPs*) using your 
Commercio.network tokens (*ucommercio*) in order to get Commercio Cash Credits (*uccc*) in return.

A *Exchange Trade Position* (*ETP*) is a core component of the Commercio Network blockchain whose purpose is to
create Commercio Cash Credits (`uccc`) in exchange for Commercio Tokens (`ucommercio`) which it then holds in
escrow until the borrowed Commercio Cash Credits are returned.

In simple words, opening an ETP allows you to exchange any amount of `ucommercio` to get relative the amount of `uccc` with relative Conversion Rate value. 
For example, if you open an ETP lending `100 ucommercio` with 1.1 Conversion Rate value will result in you receiving `90 uccc` (approximation by default).  
Initial Conversion Rate value in Params is 1. 

## Contents

1. **[State](01_state.md)**
2. **[Messages](02_messages.md)**
   - [Mint Commercio Cash Credit (CCC)](02_messages.md#mint-commercio-cash-credit-(CCC))
   - [Burn Commercio Cash Credit (CCC)](02_messages.md#burn-commercio-cash-credit-(CCC))
   - [Set Parameters (Conversion Rate & Freeze Period)](02_messages.md#set-parameters-(conversion-rate-&-freeze-period))
3. **[Events](03_events.md)**
   - [Handlers](03_events.md#handlers)
4. **[Parameters](04_params.md)**
5. **[Client](05_client.md)**