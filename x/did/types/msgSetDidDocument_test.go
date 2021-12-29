package types

import (
	"testing"
)

func TestMsgSetDidDocument_ValidateBasic(t *testing.T) {

	tests := []struct {
		name              string
		msgSetDidDocument func() *MsgSetDidDocument
		wantErr           bool
	}{
		{
			"valid",
			func() *MsgSetDidDocument {
				return &ValidMsgSetDidDocument
			},
			false,
		},
		{
			"not defined",
			func() *MsgSetDidDocument {
				return nil
			},
			true,
		},
		{
			"{context} invalid",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.Context = []string{}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{ID} empty",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.ID = ""
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{ID} invalid DID com",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.ID = "$" + ValidMsgSetDidDocument.ID
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{verificationMethod} invalid",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.VerificationMethod = []*VerificationMethod{}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{service} invalid",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.Service = []*Service{{}}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{authentication} not a set",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.Authentication = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{authentication} no reference",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.Authentication = []string{
					"NoReference",
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{assertionMethod} not a set",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.AssertionMethod = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{assertionMethod} no reference",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.AssertionMethod = []string{
					"NoReference",
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{capabilityDelegation} not a set",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.CapabilityDelegation = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{capabilityDelegation} no reference",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.CapabilityDelegation = []string{
					"NoReference",
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{capabilityInvocation} not a set",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.CapabilityInvocation = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{capabilityInvocation} no reference",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.CapabilityInvocation = []string{
					"NoReference",
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{keyAgreement} not a set",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.KeyAgreement = []string{
					validVerificationMethodRsaSignature2018.ID,
					validVerificationMethodRsaSignature2018.ID,
				}
				return &msgSetDidDocument
			},
			true,
		},
		{
			"{keyAgreement} no reference",
			func() *MsgSetDidDocument {
				msgSetDidDocument := ValidMsgSetDidDocument
				msgSetDidDocument.KeyAgreement = []string{
					"NoReference",
				}
				return &msgSetDidDocument
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.msgSetDidDocument().ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("MsgSetDidDocument.ValidateBasic() for %s error = %v, wantErr %v", tt.msgSetDidDocument(), err, tt.wantErr)
			}
		})
	}
}

func Test_validateContext(t *testing.T) {

	tests := []struct {
		name    string
		context []string
		wantErr bool
	}{
		{
			"valid",
			validContext,
			false,
		},
		{
			"empty",
			[]string{},
			true,
		},
		{
			"not a set",
			[]string{
				ContextDidV1,
				ContextDidV1,
			},
			true,
		},
		{
			"no ContextDidV1",
			[]string{
				"https://w3id.org/security/suites/ed25519-2018/v1",
			},
			true,
		},

		{
			"not first ContextDidV1",
			[]string{
				"https://w3id.org/security/suites/ed25519-2018/v1",
				ContextDidV1,
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateContext(tt.context); (err != nil) != tt.wantErr {
				t.Errorf("validateContext() for %s error = %v, wantErr %v", tt.context, err, tt.wantErr)
			}
		})
	}
}

func Test_validateService(t *testing.T) {

	tests := []struct {
		name     string
		services []*Service
		wantErr  bool
	}{
		{
			"valid",
			validServices,
			false,
		},
		{
			"empty",
			[]*Service{},
			false,
		},
		{
			"not a set",
			[]*Service{
				&validServiceBar,
				&validServiceBar,
			},
			true,
		},
		{
			"contains invalid service",
			[]*Service{
				{},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateService(tt.services); (err != nil) != tt.wantErr {
				t.Errorf("validateContext() for %s error = %v, wantErr %v", tt.services, err, tt.wantErr)
			}
		})
	}
}

func Test_validateVerificationMethod(t *testing.T) {

	tests := []struct {
		name                string
		verificationMethods []*VerificationMethod
		wantErr             bool
	}{
		{
			"valid",
			validVerificationMethods,
			false,
		},
		{
			"empty",
			[]*VerificationMethod{},
			true,
		},
		{
			"not a set",
			[]*VerificationMethod{
				&validVerificationMethodRsaVerificationKey2018,
				&validVerificationMethodRsaVerificationKey2018,
			},
			true,
		},
		{
			"does not contain RsaSignature2018",
			[]*VerificationMethod{
				&validVerificationMethodRsaVerificationKey2018,
			},
			true,
		},
		{
			"does not contain RsaVerificationKey2018",
			[]*VerificationMethod{
				&validVerificationMethodRsaSignature2018,
			},
			true,
		},
		{
			"does not contain RsaSignature2018 and RsaVerificationKey2018",
			[]*VerificationMethod{},
			true,
		},
		{
			"contains invalid",
			[]*VerificationMethod{
				{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateVerificationMethod(tt.verificationMethods, didSubject); (err != nil) != tt.wantErr {
				t.Errorf("validateVerificationMethod() for %s error = %v, wantErr %v", tt.verificationMethods, err, tt.wantErr)
			}
		})
	}
}
