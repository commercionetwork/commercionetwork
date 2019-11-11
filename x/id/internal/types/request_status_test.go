package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestRequestStatus_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rq      types.RequestStatus
		wantErr sdk.Error
	}{
		{
			"invalid status type",
			types.NewRequestStatus("invalid", "message"),
			sdk.ErrUnknownRequest("Invalid status type: invalid"),
		},
		{
			"\"rejected\" type",
			types.NewRequestStatus("rejected", ""),
			nil,
		},
		{
			"\"canceled\" type",
			types.NewRequestStatus("canceled", ""),
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != nil {
				assert.EqualError(t, tt.rq.Validate(), tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.rq.Validate())
			}
		})
	}
}
