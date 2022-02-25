package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	distrAcc      = accountTypes.NewEmptyModuleAccount(types.ModuleName)
	TestFunder, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

	valAddrVal, _    = sdk.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")
	valDelAddr, _    = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
	PKs              = simapp.CreateTestPubKeys(10)
	TestValidator, _ = stakingTypes.NewValidator(valAddrVal, PKs[0], stakingTypes.Description{})
)
