package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrator_Migrate1to2(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeperWithV300Government(t, governmentTestAddress)

			m := NewMigrator(*k)

			if err := m.Migrate1to2(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Migrator.Migrate1to2() error = %v, wantErr %v", err, tt.wantErr)
			}

			actualGovernment := k.GetGovernmentAddress(ctx)

			require.Equal(t, governmentTestAddress, actualGovernment)
		})
	}
}
