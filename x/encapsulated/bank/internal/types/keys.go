package types

import (
	"github.com/cosmos/cosmos-sdk/x/bank"
)

const (
	ModuleName = bank.ModuleName
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	BlockedAddressesStoreKey = ModuleName + "blockedAddresses"

	QueryBlockedAccounts = "blockedAccounts"

	MsgTypeBlockAccountSend = "blockAccountSend"
)
