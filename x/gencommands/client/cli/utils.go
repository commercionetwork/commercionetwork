package cli

import (
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
)

func getAddressFromString(value string) (sdk.AccAddress, error) {
	minterAddr, err := sdk.AccAddressFromBech32(value)
	if err != nil {
		kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
		if err != nil {
			return nil, err
		}

		info, err := kb.Get(value)
		if err != nil {
			return nil, err
		}

		minterAddr = info.GetAddress()
	}

	return minterAddr, nil
}
