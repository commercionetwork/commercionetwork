package types

import "testing"

func TestDidDocument_Validate(t *testing.T) {

	tests := []struct {
		name    string
		ddo     func() *DidDocument
		wantErr bool
	}{
		{
			"valid",
			func() *DidDocument {
				return &validDidDocument
			},
			false,
		},
		{
			"not defined",
			func() *DidDocument {
				return nil
			},
			true,
		},
		{
			"{context} invalid",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.Context = []string{}
				return &ddo
			},
			true,
		},
		{
			"{ID} empty",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.ID = ""
				return &ddo
			},
			true,
		},
		{
			"{ID} invalid DID com",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.ID = "$" + ddo.ID
				return &ddo
			},
			true,
		},
		{
			"{verificationMethod} invalid",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.VerificationMethod = []*VerificationMethod{}
				return &ddo
			},
			true,
		},
		{
			"{service} invalid",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.Service = []*Service{{}}
				return &ddo
			},
			true,
		},
		{
			"{authentication} not a set",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.Authentication = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &ddo
			},
			true,
		},
		{
			"{authentication} no reference",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.Authentication = []string{
					"NoReference",
				}
				return &ddo
			},
			true,
		},
		{
			"{assertionMethod} not a set",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.AssertionMethod = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &ddo
			},
			true,
		},
		{
			"{assertionMethod} no reference",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.AssertionMethod = []string{
					"NoReference",
				}
				return &ddo
			},
			true,
		},
		{
			"{capabilityDelegation} not a set",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.CapabilityDelegation = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &ddo
			},
			true,
		},
		{
			"{capabilityDelegation} no reference",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.CapabilityDelegation = []string{
					"NoReference",
				}
				return &ddo
			},
			true,
		},
		{
			"{capabilityInvocation} not a set",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.CapabilityInvocation = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &ddo
			},
			true,
		},
		{
			"{capabilityInvocation} no reference",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.CapabilityInvocation = []string{
					"NoReference",
				}
				return &ddo
			},
			true,
		},
		{
			"{keyAgreement} not a set",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.KeyAgreement = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &ddo
			},
			true,
		},
		{
			"{keyAgreement} no reference",
			func() *DidDocument {
				ddo := validDidDocument
				ddo.KeyAgreement = []string{
					"NoReference",
				}
				return &ddo
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ddo().Validate(); (err != nil) != tt.wantErr {
				t.Errorf("MsgSetDidDocument.ValidateBasic() for %s error = %v, wantErr %v", tt.ddo(), err, tt.wantErr)
			}
		})
	}
}
