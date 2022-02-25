package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	distrAcc = accountTypes.NewEmptyModuleAccount(types.ModuleName)

	valAddrVal, _    = sdk.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")
	PKs              = simapp.CreateTestPubKeys(10)
	TestValidator, _ = stakingTypes.NewValidator(valAddrVal, PKs[0], stakingTypes.Description{})
)
