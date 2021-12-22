package types

import (
	"testing"
)

var validService = Service{
	ID:              "https://bar.example.com",
	Type:            "agent",
	ServiceEndpoint: "https://commerc.io/agent/serviceEndpoint/",
}

func TestService_isValid(t *testing.T) {

	tests := []struct {
		name    string
		service func() *Service
		wantErr bool
	}{
		{
			"valid",
			func() *Service {
				return &validService
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
				service := validService
				service.ID = ""
				return &service
			},
			true,
		},
		{
			"{ID} against the rules of RFC3986",
			func() *Service {
				service := validService
				service.ID = "$" + validService.ID
				return &service
			},
			true,
		},
		{
			"{type} empty",
			func() *Service {
				service := validService
				service.Type = ""
				return &service
			},
			true,
		},
		{
			"{serviceEndpoint} empty",
			func() *Service {
				service := validService
				service.ServiceEndpoint = ""
				return &service
			},
			true,
		},
		{
			"{serviceEndpoint} against the rules of RFC3986",
			func() *Service {
				service := validService
				service.ServiceEndpoint = "$" + validService.ServiceEndpoint
				return &service
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

var validVerificationMethod = VerificationMethod{
	ID:                 validDid + "#key-1",
	Type:               RsaVerificationKey2018,
	Controller:         validDid,
	PublicKeyMultibase: "m" + "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDMr3V+Auyc+zvt2qX+jpwk3wM+m2DbfLjimByzQDIfrzSHMTQ8erL0kg69YsXHYXVX9mIZKRzk6VNwOBOQJSsIDf2jGbuEgI8EB4c3q1XykakCTvO3Ku3PJgZ9PO4qRw7QVvTkCbc91rT93/pD3/Ar8wqd4pNXtgbfbwJGviZ6kQIDAQAB",
}

func TestVerificationMethod_isValid(t *testing.T) {

	tests := []struct {
		name               string
		verificationMethod func() *VerificationMethod
		wantErr            bool
	}{
		{
			"valid",
			func() *VerificationMethod {
				return &validVerificationMethod
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
				verificationMethod := validVerificationMethod
				verificationMethod.ID = ""
				return &verificationMethod
			},
			true,
		},
		{
			"{ID} against the DID url specification",
			func() *VerificationMethod {
				verificationMethod := validVerificationMethod
				verificationMethod.ID = "$" + validVerificationMethod.ID
				return &verificationMethod
			},
			true,
		},
		{
			"{type} empty",
			func() *VerificationMethod {
				verificationMethod := validVerificationMethod
				verificationMethod.Type = ""
				return &verificationMethod
			},
			true,
		},
		{
			"{type} not supported",
			func() *VerificationMethod {
				verificationMethod := validVerificationMethod
				verificationMethod.Type = "NotSupported2077"
				return &verificationMethod
			},
			true,
		},
		{
			"{controller} empty",
			func() *VerificationMethod {
				verificationMethod := validVerificationMethod
				verificationMethod.Controller = ""
				return &verificationMethod
			},
			true,
		},
		{
			"{controller} against the DID specification",
			func() *VerificationMethod {
				verificationMethod := validVerificationMethod
				verificationMethod.Controller = "$" + validDid
				return &verificationMethod
			},
			true,
		},
		{
			"{publicKeyMultibase} empty",
			func() *VerificationMethod {
				verificationMethod := validVerificationMethod
				verificationMethod.PublicKeyMultibase = ""
				return &verificationMethod
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

const validDid = "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc"

func Test_isValidDidCom(t *testing.T) {

	tests := []struct {
		name    string
		did     string
		wantErr bool
	}{
		{"valid", validDid, false},
		{"not valid", "$" + validDid, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isValidDidCom(tt.did); (err != nil) != tt.wantErr {
				t.Errorf("isValidDidCom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
