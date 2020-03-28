package creditrisk

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/creditrisk/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	moduleAcc := keeper.getModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(data.Pool); err != nil {
			panic(err)
		}
		keeper.setModuleAccount(ctx, moduleAcc)
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data types.GenesisState) error {
	if !data.Pool.IsValid() {
		return fmt.Errorf("invalid pool amount: %s", data.Pool)
	}
	return nil
}
