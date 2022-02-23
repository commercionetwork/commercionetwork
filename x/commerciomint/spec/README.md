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
Initial Conversion Rate value in Params is 0.61. 


## Module Accounts

The supply functionality introduces a new type of `auth.Account` which can be used by
modules to allocate tokens and in special cases mint or burn tokens. At a base
level these module accounts are capable of sending/receiving tokens to and from
`auth.Account`s and other module accounts. This design replaces previous
alternative designs where, to hold tokens, modules would burn the incoming
tokens from the sender account, and then track those tokens internally. Later,
in order to send tokens, the module would need to effectively mint tokens
within a destination account. The new design removes duplicate logic between
modules to perform this accounting.

The `ModuleAccount` interface is defined as follows:

```go
type ModuleAccount interface {
  auth.Account               // same methods as the Account interface

  GetName() string           // name of the module; used to obtain the address
  GetPermissions() []string  // permissions of module account
  HasPermission(string) bool
}
```

> **WARNING!**
> Any module or message handler that allows either direct or indirect sending of funds must explicitly guarantee those funds cannot be sent to module accounts (unless allowed).

The supply `Keeper` also introduces new wrapper functions for the auth `Keeper`
and the comerciokyc `Keeper` that are related to `ModuleAccount`s in order to be able
to:

- Get and set `ModuleAccount`s by providing the `Name`.
- Send coins from and to other `ModuleAccount`s or standard `Account`s
  (`BaseAccount` or `VestingAccount`) by passing only the `Name`.
- `Mint` or `Burn` coins for a `ModuleAccount` (restricted to its permissions).

### Permissions

Each `ModuleAccount` has a different set of permissions that provide different
object capabilities to perform certain actions. Permissions need to be
registered upon the creation of the supply `Keeper` so that every time a
`ModuleAccount` calls the allowed functions, the `Keeper` can lookup the
permissions to that specific account and perform or not the action.

The available permissions are:

- `Minter`: allows for a module to mint a specific amount of coins.
- `Burner`: allows for a module to burn a specific amount of coins.

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