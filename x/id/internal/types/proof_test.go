package types_test

import (
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestProof_Equals(t *testing.T) {
	zone, _ := time.LoadLocation("UTC")
	date := time.Date(2019, 1, 1, 1, 1, 1, 1, zone)
	proof := types.NewProof("type-1", date, "creator-1", "")

	tests := []struct {
		name  string
		us    types.Proof
		them  types.Proof
		equal bool
	}{
		{
			"different type",
			proof,
			types.NewProof("type-2", proof.Created, proof.Creator, proof.SignatureValue),
			false,
		},
		{
			"different creator",
			proof,
			types.NewProof(proof.Type, proof.Created, "creator-2", proof.SignatureValue),
			false,
		},
		{
			"off-by-1 day date",
			proof,
			types.NewProof(proof.Type, proof.Created.AddDate(0, 0, 1), proof.Creator, proof.SignatureValue),
			false,
		},
		{
			"different signature value",
			proof,
			types.NewProof(proof.Type, proof.Created, proof.Creator, proof.SignatureValue+"1"),
			false,
		},
		{
			"two equal proofs",
			proof,
			proof,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equal, tt.us.Equals(tt.them))
		})
	}
}

func TestProof_Validate(t *testing.T) {
	zone, _ := time.LoadLocation("UTC")
	date := time.Date(2019, 1, 1, 1, 1, 1, 1, zone)

	tests := []struct {
		name    string
		p       types.Proof
		wantErr error
	}{
		{
			"no type",
			types.NewProof("", date, "creator", "signature"),
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid proof type, must be LinkedDataSignature2015"),
		},
		{
			"no creation date",
			types.NewProof("LinkedDataSignature2015", time.Time{}, "creator", "signature"),
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid proof creation time"),
		},
		{
			"valid proof",
			types.NewProof("LinkedDataSignature2015", date, "creator", "signatureValue"),
			nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != nil {
				require.EqualError(t, tt.p.Validate(), tt.wantErr.Error())
			} else {
				require.NoError(t, tt.p.Validate())
			}
		})
	}
}
