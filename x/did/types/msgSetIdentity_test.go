package types

import (
	"testing"
)

func Test_validateContextSlice(t *testing.T) {

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
			if err := validateContextSlice(tt.context); (err != nil) != tt.wantErr {
				t.Errorf("validateContext() for %s error = %v, wantErr %v", tt.context, err, tt.wantErr)
			}
		})
	}
}

func Test_validateServiceSlice(t *testing.T) {

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
			if err := validateServiceSlice(tt.services); (err != nil) != tt.wantErr {
				t.Errorf("validateContext() for %s error = %v, wantErr %v", tt.services, err, tt.wantErr)
			}
		})
	}
}

func Test_validateVerificationMethodSlice(t *testing.T) {

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
			if err := validateVerificationMethodSlice(tt.verificationMethods, validDidSubject); (err != nil) != tt.wantErr {
				t.Errorf("validateVerificationMethod() for %s error = %v, wantErr %v", tt.verificationMethods, err, tt.wantErr)
			}
		})
	}
}
