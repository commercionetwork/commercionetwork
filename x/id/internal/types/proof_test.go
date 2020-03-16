package types_test

import (
	"testing"
	"time"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestProof_Equals(t *testing.T) {
	zone, _ := time.LoadLocation("UTC")
	date := time.Date(2019, 1, 1, 1, 1, 1, 1, zone)
	proof := types.NewProof("proof", date, "purpose", "controller", "verificationMethod", "sigvalue")

	tests := []struct {
		name  string
		us    types.Proof
		them  types.Proof
		equal bool
	}{
		{
			"different type",
			proof,
			types.NewProof("prooff", date, "purpose", "controller", "verificationMethod", "sigvalue"),
			false,
		},
		{
			"off-by-1 day date",
			proof,
			types.NewProof("proof", date.AddDate(0, 0, 1), "purpose", "controller", "verificationMethod", "sigvalue"),
			false,
		},
		{
			"different signature value",
			proof,
			types.NewProof("proof", date, "purpose", "controller", "verificationMethod", "sigvaluee"),
			false,
		},
		{
			"different purpose",
			proof,
			types.NewProof("proof", date, "purposee", "controller", "verificationMethod", "sigvalue"),
			false,
		},
		{
			"different controller",
			proof,
			types.NewProof("proof", date, "purpose", "controllerr", "verificationMethod", "sigvalue"),
			false,
		},
		{
			"different verificationMethod",
			proof,
			types.NewProof("proof", date, "purpose", "controller", "verificationMethodd", "sigvalue"),
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
	types.ConfigTestPrefixes()
	var testZone, _ = time.LoadLocation("UTC")
	var testTime = time.Date(2016, 2, 8, 16, 2, 20, 0, testZone)
	var testOwnerAddress, _ = sdk.AccAddressFromBech32("did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	var testAnotherAddress, _ = sdk.AccAddressFromBech32("cosmos1gdpsu89prllyw49eehskv6t8800p6chefyuuwe")

	validProof := types.Proof{
		Type:               "EcdsaSecp256k1VerificationKey2019",
		Created:            testTime,
		ProofPurpose:       "authentication",
		Controller:         testOwnerAddress.String(),
		SignatureValue:     "4T2jhs4C0k7p649tdzQAOLqJ0GJsiFDP/NnsSkFpoXAxcgn6h/EgvOpHxW7FMNQ9RDgQbcE6FWP6I2UsNv1qXQ==",
		VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
	}

	tests := []struct {
		name    string
		p       types.Proof
		wantErr error
	}{
		{
			"invalid type",
			types.Proof{
				Type: "wrongType",
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid proof type, must be EcdsaSecp256k1VerificationKey2019"),
		},
		{
			"invalid created",
			types.Proof{
				Type:    validProof.Type,
				Created: time.Time{},
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, ("Invalid proof creation time")),
		},
		{
			"invalid proof purpose",
			types.Proof{
				Type:         validProof.Type,
				Created:      validProof.Created,
				ProofPurpose: "notvalid",
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "proof purpose must be \"authentication\""),
		},
		{
			"invalid controller",
			types.Proof{
				Type:         validProof.Type,
				Created:      validProof.Created,
				ProofPurpose: validProof.ProofPurpose,
				Controller:   "not bech32",
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "invalid controller, must be a valid bech32-encoded address"),
		},
		{
			"valid proof",
			validProof,
			nil,
		},
		{
			"valid proof but controller is not associated to the public key",
			types.Proof{
				Type:               validProof.Type,
				Created:            validProof.Created,
				ProofPurpose:       validProof.ProofPurpose,
				Controller:         testAnotherAddress.String(),
				VerificationMethod: validProof.VerificationMethod,
				SignatureValue:     validProof.SignatureValue,
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "verification method-derived AccAddress differs from controller"),
		},
		{
			"valid proof but signature does not match",
			types.Proof{
				Type:               validProof.Type,
				Created:            validProof.Created,
				ProofPurpose:       validProof.ProofPurpose,
				Controller:         validProof.Controller,
				VerificationMethod: validProof.VerificationMethod,
				SignatureValue:     validProof.SignatureValue + "oh no!",
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "signature value must be base64 encoded"),
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
