package types

import (
	"fmt"
	"testing"
)

func TestService_isValid(t *testing.T) {

	ValidService := &Service{
		ID:              "https://bar.example.com",
		Type:            "agent",
		ServiceEndpoint: "https://commerc.io/agent/serviceEndpoint/",
	}

	tests := []struct {
		name    string
		service func() *Service
		wantErr bool
	}{
		{
			"valid",
			func() *Service {
				return ValidService
			},
			false,
		},
		{
			"not defined",
			func() *Service {
				return nil
			},
			true,
		},
		{
			"{ID} empty",
			func() *Service {
				ValidService.ID = ""
				return ValidService
			},
			true,
		},
		{
			"{ID} against the rules of RFC3986",
			func() *Service {
				ValidService.ID = fmt.Sprint("$", ValidService.ID)
				return ValidService
			},
			true,
		},
		{
			"{type} empty",
			func() *Service {
				ValidService.Type = ""
				return ValidService
			},
			true,
		},
		{
			"{serviceEndpoint} empty",
			func() *Service {
				ValidService.ServiceEndpoint = ""
				return ValidService
			},
			true,
		},
		{
			"{serviceEndpoint} against the rules of RFC3986",
			func() *Service {
				ValidService.ServiceEndpoint = fmt.Sprint("$", ValidService.ServiceEndpoint)
				return ValidService
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.service().isValid(); (err != nil) != tt.wantErr {
				t.Errorf("Service.isValid() for service %s error = %v, wantErr %v", tt.service(), err, tt.wantErr)
			}
		})
	}
}

func TestVerificationMethod_isValid(t *testing.T) {

	did := "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc"

	valid := &VerificationMethod{
		ID:                 did + "#key-1",
		Type:               "RsaVerificationKey2018",
		Controller:         did,
		PublicKeyMultibase: "m" + "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB",
	}

	tests := []struct {
		name               string
		verificationMethod func() *VerificationMethod
		wantErr            bool
	}{
		{
			"valid",
			func() *VerificationMethod {
				return valid
			},
			false,
		},
		{
			"not defined",
			func() *VerificationMethod {
				return nil
			},
			true,
		},
		{
			"{ID} empty",
			func() *VerificationMethod {
				valid.ID = ""
				return valid
			},
			true,
		},
		{
			"{ID} against the DID url specification",
			func() *VerificationMethod {
				valid.ID = "$" + valid.ID
				return valid
			},
			true,
		},
		{
			"{type} empty",
			func() *VerificationMethod {
				valid.ID = ""
				return valid
			},
			true,
		},
		{
			"{type} not supported",
			func() *VerificationMethod {
				valid.Type = "NotSupported2077"
				return valid
			},
			true,
		},
		{
			"{controller} empty",
			func() *VerificationMethod {
				valid.ID = ""
				return valid
			},
			true,
		},
		{
			"{controller} against the DID specification",
			func() *VerificationMethod {
				valid.ID = "$" + did
				return valid
			},
			true,
		},
		{
			"{publicKeyMultibase} empty",
			func() *VerificationMethod {
				valid.PublicKeyMultibase = ""
				return valid
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.verificationMethod().isValid(); (err != nil) != tt.wantErr {
				t.Errorf("VerificationMethod.isValid() for verificationMethod %s error = %v, wantErr %v", tt.verificationMethod(), err, tt.wantErr)
			}
		})
	}
}
