package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateSdnData(t *testing.T) {
	tests := []struct {
		name        string
		sdnData     SdnData
		isValid     bool
		expectedErr error
	}{
		{
			"invalid: invalid field empty",
			SdnData{
				"first_name",
				"",
			},
			false,
			fmt.Errorf("sdn_data value %s is not supported", ""),
		},
		{
			"invalid: strange value not included in supportedSdnData",
			SdnData{
				"first_name",
				"last_name",
				"age",
			},
			false,
			fmt.Errorf("sdn_data value %s is not supported", "age"),
		},
		{
			"valid: all supported fields",
			SdnData{
				"first_name",
				"last_name",
				"tin",
				"email",
				"organization",
				"country",
			},
			true,
			nil,
		},
		{
			"valid: subset of supported fields",
			SdnData{
				"first_name",
				"tin",
				"email",
				"country",
			},
			true,
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdnData.Validate()

			if tt.isValid {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
