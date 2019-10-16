package custombank

import (
	"github.com/commercionetwork/commercionetwork/x/encapsulated/bank/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/encapsulated/bank/internal/types"
)

const (
	StoreKey = types.StoreKey
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier

	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	Keeper = keeper.Keeper

	MsgBlockAddressSend  = types.MsgBlockAccountSend
	MsgUnlockAddressSend = types.MsgUnlockAccountSend
)
