package accreditations

import (
	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute

	TrustedSignersStoreKey = types.TrustedSignersStoreKey
	LiquidityPoolKey       = types.LiquidityPoolKey
)

var (
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier

	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	// Tests
	GetTestInput   = keeper.GetTestInput
	TestUser       = keeper.TestUser
	TestAccrediter = keeper.TestAccrediter
	TestSigner     = keeper.TestSigner
)

type (
	Keeper = keeper.Keeper

	Accreditation = types.Accreditation

	MsgSetAccrediter            = types.MsgSetAccrediter
	MsgDistributeReward         = types.MsgDistributeReward
	MsgDepositIntoLiquidityPool = types.MsgDepositIntoLiquidityPool
	MsgAddTrustedSigner         = types.MsgAddTrustedSigner
)
