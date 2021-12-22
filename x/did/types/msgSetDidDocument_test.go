package types

import (
	"testing"
)

var validMsgSetDidDocument = &MsgSetDidDocument{
	Context: validContext,
	ID:      validDid,
	VerificationMethod: []*VerificationMethod{
		&validVerificationMethod,
		{
			ID:                 validDid + "#key-agreement-1",
			Type:               "RsaSignature2018",
			Controller:         validDid,
			PublicKeyMultibase: "H3C2AVvLMv6gmMNam3uVAjZpfkcJCwDwnZn6z3wXmqPV",
		},
	},
	Authentication: []string{
		validDid + "#key-1",
	},
	AssertionMethod: []string{
		validDid + "#key-1",
	},
	KeyAgreement: []string{
		validDid + "#key-agreement-1",
	},
	CapabilityInvocation: nil,
	CapabilityDelegation: nil,
	Service: []*Service{
		&validService,
		{
			ID:              "https://foo.example.com",
			Type:            "xdi",
			ServiceEndpoint: "https://commerc.io/xdi/serviceEndpoint/",
		},
	},
}

func TestMsgSetDidDocument_ValidateBasic(t *testing.T) {

	tests := []struct {
		name              string
		msgSetDidDocument func() *MsgSetDidDocument
		wantErr           bool
	}{
		{
			"valid",
			func() *MsgSetDidDocument {
				return validMsgSetDidDocument
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
				msgSetDidDocument := validMsgSetDidDocument
				msgSetDidDocument.Context = []string{}
				return msgSetDidDocument
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

var validContext = []string{
	ContextDidV1,
	"https://w3id.org/security/suites/ed25519-2018/v1",
	"https://w3id.org/security/suites/x25519-2019/v1",
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
