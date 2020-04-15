package customstaking

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestAppModule(t *testing.T) {
	module := &AppModule{}

	require.Equal(t, types.ModuleName, module.Name())
	require.Equal(t, types.RouterKey, module.Route())
	require.Equal(t, types.QuerierRoute, module.QuerierRoute())
}
