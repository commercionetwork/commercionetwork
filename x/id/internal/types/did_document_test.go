package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name        string
		theService  Service
		isValid     bool
		expectedErr error
	}{
		{
			"valid",
			Service{
				Id:              "did:example:123456789abcdefghi#vcr",
				Type:            "CredentialRepositoryService",
				ServiceEndpoint: "http://theUrl",
			},
			true,
			nil,
		},
		{
			"no id",
			Service{
				Type:            "CredentialRepositoryService",
				ServiceEndpoint: "http://theUrl",
			},
			false,
			sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "service field \"%s\" is required", "id"),
		},
		{
			"no type",
			Service{
				Id:              "did:example:123456789abcdefghi#vcr",
				ServiceEndpoint: "http://theUrl",
			},
			false,
			sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "service field \"%s\" is required", "type"),
		},
		{
			"no service point",
			Service{
				Id:   "did:example:123456789abcdefghi#vcr",
				Type: "CredentialRepositoryService",
			},
			false,
			sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "service field \"%s\" is required", "serviceEndpoint"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.theService.Validate()
			if tt.isValid {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}
