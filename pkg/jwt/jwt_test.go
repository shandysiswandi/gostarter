package jwt

import (
	"crypto/rsa"
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewJSONWebToken(t *testing.T) {
	type args struct {
		private string
		public  string
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONWebToken
		wantErr bool
	}{
		{
			name: "ErrorDecodePrivate",
			args: args{
				private: "abcd!@#",
				public:  "abcd!@#",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ErrorPemDecodePrivate",
			args: args{
				private: "none",
				public:  "abcd!@#",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ErrorParsePKCS8PrivateKey",
			args: args{
				private: "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1IUUNBUUVFSUVhNTZHRzJQVFVKeUl0NEZ5ZGFNTkl0WXNqTmo2WkliZDdqWHZEWTRFbGZvQWNHQlN1QkJBQUsKb1VRRFFnQUVKUURuOC92ZDhvUXBBL1ZFM2NoMGxNNlZBcHJPVGlWOVZMcDM4cndmT29nM3FVWWNUeHhYL3N4SgpsMU00SG5jcUVvcFlJS2trb3ZvRkZpNjJZcGg2bnc9PQotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0t",
				public:  "abcd!@#",
			},
			want:    nil,
			wantErr: true,
		},
		//
		{
			name: "ErrorDecodePublic",
			args: args{
				private: "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ4d2dnRTdBZ0VBQWtFQXg5d0R5UTdMUjd4Qk1MZ1BpREY4VkJGZkp3Rk9vRmYrZHhXbnZOZFJVcVBUb2NaWHlJTkExdVZUUjVhSStpSVMzM1FFdU9IL0E0QmFFNzdDU3pnL1BRSURBUUFCQWtBcFNsVjIzNmxVUWZxMjZjUDl0M21QOWNIU2tXVFE0RFVZbmI3NGx3UjhYYklNN29Vb3hWd3gyb052RDBPZ3RNRTkrcldYOU9tdmpZdWIvNENHNS9CSkFpRUE5czYwdEZpMUs3V2xFdXlRcW14OGxudmpYY0tuM0R1TktuM2pLS1Nldm1VQ0lRRFBUYW1WekpCV01mNkNrVWg1NTlkV0FMQzJmYlBUUWZST3k3ZkZPWEJqK1FJZ0RnN21JaU92WmliNW1TTmFkaXFweWhTU2RlUEJsZnphWktJNUR6YVpTRFVDSVFDd0g4dDArZGVuWTlKWUhCYjNlNEg0RDU0VGJiamFRNjdOUTBkZXlPNDBBUUloQUlGcVV6Mm5iSUNiZDNIWlhkRWlURzJDV3Q4eVVUZ2l4RWx1NHZtc2VoQjcKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==",
				public:  "abcd!@#",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ErrorPemDecodePrivate",
			args: args{
				private: "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ4d2dnRTdBZ0VBQWtFQXg5d0R5UTdMUjd4Qk1MZ1BpREY4VkJGZkp3Rk9vRmYrZHhXbnZOZFJVcVBUb2NaWHlJTkExdVZUUjVhSStpSVMzM1FFdU9IL0E0QmFFNzdDU3pnL1BRSURBUUFCQWtBcFNsVjIzNmxVUWZxMjZjUDl0M21QOWNIU2tXVFE0RFVZbmI3NGx3UjhYYklNN29Vb3hWd3gyb052RDBPZ3RNRTkrcldYOU9tdmpZdWIvNENHNS9CSkFpRUE5czYwdEZpMUs3V2xFdXlRcW14OGxudmpYY0tuM0R1TktuM2pLS1Nldm1VQ0lRRFBUYW1WekpCV01mNkNrVWg1NTlkV0FMQzJmYlBUUWZST3k3ZkZPWEJqK1FJZ0RnN21JaU92WmliNW1TTmFkaXFweWhTU2RlUEJsZnphWktJNUR6YVpTRFVDSVFDd0g4dDArZGVuWTlKWUhCYjNlNEg0RDU0VGJiamFRNjdOUTBkZXlPNDBBUUloQUlGcVV6Mm5iSUNiZDNIWlhkRWlURzJDV3Q4eVVUZ2l4RWx1NHZtc2VoQjcKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==",
				public:  "none",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ErrorParsePKCS8PrivateKey",
			args: args{
				private: "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ4d2dnRTdBZ0VBQWtFQXg5d0R5UTdMUjd4Qk1MZ1BpREY4VkJGZkp3Rk9vRmYrZHhXbnZOZFJVcVBUb2NaWHlJTkExdVZUUjVhSStpSVMzM1FFdU9IL0E0QmFFNzdDU3pnL1BRSURBUUFCQWtBcFNsVjIzNmxVUWZxMjZjUDl0M21QOWNIU2tXVFE0RFVZbmI3NGx3UjhYYklNN29Vb3hWd3gyb052RDBPZ3RNRTkrcldYOU9tdmpZdWIvNENHNS9CSkFpRUE5czYwdEZpMUs3V2xFdXlRcW14OGxudmpYY0tuM0R1TktuM2pLS1Nldm1VQ0lRRFBUYW1WekpCV01mNkNrVWg1NTlkV0FMQzJmYlBUUWZST3k3ZkZPWEJqK1FJZ0RnN21JaU92WmliNW1TTmFkaXFweWhTU2RlUEJsZnphWktJNUR6YVpTRFVDSVFDd0g4dDArZGVuWTlKWUhCYjNlNEg0RDU0VGJiamFRNjdOUTBkZXlPNDBBUUloQUlGcVV6Mm5iSUNiZDNIWlhkRWlURzJDV3Q4eVVUZ2l4RWx1NHZtc2VoQjcKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==",
				public:  "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUhRQ0FRRUVJRWE1NkdHMlBUVUp5SXQ0RnlkYU1OSXRZc2pOajZaSWJkN2pYdkRZNEVsZm9BY0dCU3VCQkFBSwpvVVFEUWdBRUpRRG44L3ZkOG9RcEEvVkUzY2gwbE02VkFwck9UaVY5VkxwMzhyd2ZPb2czcVVZY1R4eFgvc3hKCmwxTTRIbmNxRW9wWUlLa2tvdm9GRmk2MllwaDZudz09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==",
			},
			want:    nil,
			wantErr: true,
		},
		//
		{
			name: "Success",
			args: args{
				private: "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ4d2dnRTdBZ0VBQWtFQXg5d0R5UTdMUjd4Qk1MZ1BpREY4VkJGZkp3Rk9vRmYrZHhXbnZOZFJVcVBUb2NaWHlJTkExdVZUUjVhSStpSVMzM1FFdU9IL0E0QmFFNzdDU3pnL1BRSURBUUFCQWtBcFNsVjIzNmxVUWZxMjZjUDl0M21QOWNIU2tXVFE0RFVZbmI3NGx3UjhYYklNN29Vb3hWd3gyb052RDBPZ3RNRTkrcldYOU9tdmpZdWIvNENHNS9CSkFpRUE5czYwdEZpMUs3V2xFdXlRcW14OGxudmpYY0tuM0R1TktuM2pLS1Nldm1VQ0lRRFBUYW1WekpCV01mNkNrVWg1NTlkV0FMQzJmYlBUUWZST3k3ZkZPWEJqK1FJZ0RnN21JaU92WmliNW1TTmFkaXFweWhTU2RlUEJsZnphWktJNUR6YVpTRFVDSVFDd0g4dDArZGVuWTlKWUhCYjNlNEg0RDU0VGJiamFRNjdOUTBkZXlPNDBBUUloQUlGcVV6Mm5iSUNiZDNIWlhkRWlURzJDV3Q4eVVUZ2l4RWx1NHZtc2VoQjcKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==",
				public:  "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZzd0RRWUpLb1pJaHZjTkFRRUJCUUFEU2dBd1J3SkFiR1pycUhvaVRKV21kYzR2N1Fld09ETUR0UW9iK0NLcwpTY2RjTEZaZUdBWE9CMkpmOGFDeEk5MXc3WVBxQ2pHTVRNTDlRSmo1WDIvNkRMc203aUlBbndJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t",
			},
			want: &JSONWebToken{
				privateKey: &rsa.PrivateKey{},
				publicKey:  &rsa.PublicKey{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewJSONWebToken(tt.args.private, tt.args.public)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.want != nil {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestJSONWebToken_Generate(t *testing.T) {
	type args struct {
		c *Claim
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		mockFn  func(a args) *JSONWebToken
	}{
		{
			name: "Success",
			args: args{c: &Claim{
				AuthID: 101,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "test",
					Subject:   "test",
					Audience:  []string{"test"},
					ExpiresAt: jwt.NewNumericDate(time.Date(2034, time.December, 1, 0, 0, 0, 0, time.Local)),
					NotBefore: jwt.NewNumericDate(time.Date(2024, time.December, 1, 0, 0, 0, 0, time.Local)),
					IssuedAt:  jwt.NewNumericDate(time.Date(2024, time.December, 1, 0, 0, 0, 0, time.Local)),
				},
			}},
			want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX2lkIjoiMTAxIiwiaXNzIjoidGVzdCIsInN1YiI6InRlc3QiLCJhdWQiOlsidGVzdCJdLCJleHAiOjIwNDg1MTg4MDAsIm5iZiI6MTczMjk4NjAwMCwiaWF0IjoxNzMyOTg2MDAwfQ.Hj3JBUJFVIlquVpiR3ZPj0cw2gM7nLE2mzzZOnEihx2h6zcj7ZT8hVH9-0ZUsxS8UZf7xHuBfjxeAQkARragrg",
			wantErr: false,
			mockFn: func(a args) *JSONWebToken {
				j, _ := NewJSONWebToken("LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWQUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ0d2dnRTZBZ0VBQWtFQW8vUGJQL3BYSWJWRXZmUWoKeWJIMlVzY25iSmFNVjNqYTJZSTREN1Rtbk9nWVpTaHFSNDVYeHZZWEd2WWpsNjlkdm9KaXJOSUd4MEZBdzBOYgpDNEZWN1FJREFRQUJBa0VBanJRTGF5Vm52NlE2WUNmbkdvQjJ5VmdrL1lRUVJYYWc3bDlFa284L2h1T3FyVFJoClgwWVMxZll4UGFFZkdDRFM2ejZQSzY0Yk15aTdBMnZQOTZKUXRRSWhBT1V2ZE1uQlYzNWlXSW1JZ01CMmdKMzgKa2hXMDVDU2dZZW01RVJNS3BoZ25BaUVBdHlLTlFnVlZ1c05xbjhNYTZkOFplTk93d0QzN3M3Y0cxMW5lcnpubwp1Y3NDSUE5c0JCWFhkc1hBWkdqTTBLMGl6RURWVUJjNTF1aElDbzZwcjJaeW52NmRBaUJLcTlueEMzL1RNUTd1CnFYejEwelB0b2xNMWI1Q0x6SnNMZitkZWh6d3ZWUUlnSGhlcy9MVlQ1aFdLT2IyRXFISFQ0RlFlS1dpRFBlc3cKRkNickcrR1dtSmM9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0=", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBS1B6MnovNlZ5RzFSTDMwSThteDlsTEhKMnlXakZkNAoydG1DT0ErMDVwem9HR1VvYWtlT1Y4YjJGeHIySTVldlhiNkNZcXpTQnNkQlFNTkRXd3VCVmUwQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==")

				return j
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.mockFn(tt.args).Generate(tt.args.c)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJSONWebToken_Verify(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    *Claim
		wantErr bool
		mockFn  func(a args) *JSONWebToken
	}{
		{
			name:    "ErrorVerify",
			args:    args{token: ""},
			want:    nil,
			wantErr: true,
			mockFn: func(a args) *JSONWebToken {
				j, _ := NewJSONWebToken("LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ4d2dnRTdBZ0VBQWtFQXg5d0R5UTdMUjd4Qk1MZ1BpREY4VkJGZkp3Rk9vRmYrZHhXbnZOZFJVcVBUb2NaWHlJTkExdVZUUjVhSStpSVMzM1FFdU9IL0E0QmFFNzdDU3pnL1BRSURBUUFCQWtBcFNsVjIzNmxVUWZxMjZjUDl0M21QOWNIU2tXVFE0RFVZbmI3NGx3UjhYYklNN29Vb3hWd3gyb052RDBPZ3RNRTkrcldYOU9tdmpZdWIvNENHNS9CSkFpRUE5czYwdEZpMUs3V2xFdXlRcW14OGxudmpYY0tuM0R1TktuM2pLS1Nldm1VQ0lRRFBUYW1WekpCV01mNkNrVWg1NTlkV0FMQzJmYlBUUWZST3k3ZkZPWEJqK1FJZ0RnN21JaU92WmliNW1TTmFkaXFweWhTU2RlUEJsZnphWktJNUR6YVpTRFVDSVFDd0g4dDArZGVuWTlKWUhCYjNlNEg0RDU0VGJiamFRNjdOUTBkZXlPNDBBUUloQUlGcVV6Mm5iSUNiZDNIWlhkRWlURzJDV3Q4eVVUZ2l4RWx1NHZtc2VoQjcKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQ==", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZzd0RRWUpLb1pJaHZjTkFRRUJCUUFEU2dBd1J3SkFiR1pycUhvaVRKV21kYzR2N1Fld09ETUR0UW9iK0NLcwpTY2RjTEZaZUdBWE9CMkpmOGFDeEk5MXc3WVBxQ2pHTVRNTDlRSmo1WDIvNkRMc203aUlBbndJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t")

				return j
			},
		},
		{
			name:    "ErrorExpired",
			args:    args{token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImVtYWlsIiwiaXNzIjoiZ29zdGFydGVyIiwic3ViIjoiZW1haWwiLCJleHAiOjE3MzEyMjc5ODgsIm5iZiI6MTczMTIyNzk4OCwiaWF0IjoxNzMxMjI3OTg4fQ.mi8cqVth482eDxDrADwhVuVsXZiI6kLbpWq03xIJ1JL3g4hOsxt5LuVW-dLo7dG8fP4dAiZ00QzxiwRo3cobdA"},
			want:    nil,
			wantErr: true,
			mockFn: func(a args) *JSONWebToken {
				j, _ := NewJSONWebToken("LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWQUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ0d2dnRTZBZ0VBQWtFQW8vUGJQL3BYSWJWRXZmUWoKeWJIMlVzY25iSmFNVjNqYTJZSTREN1Rtbk9nWVpTaHFSNDVYeHZZWEd2WWpsNjlkdm9KaXJOSUd4MEZBdzBOYgpDNEZWN1FJREFRQUJBa0VBanJRTGF5Vm52NlE2WUNmbkdvQjJ5VmdrL1lRUVJYYWc3bDlFa284L2h1T3FyVFJoClgwWVMxZll4UGFFZkdDRFM2ejZQSzY0Yk15aTdBMnZQOTZKUXRRSWhBT1V2ZE1uQlYzNWlXSW1JZ01CMmdKMzgKa2hXMDVDU2dZZW01RVJNS3BoZ25BaUVBdHlLTlFnVlZ1c05xbjhNYTZkOFplTk93d0QzN3M3Y0cxMW5lcnpubwp1Y3NDSUE5c0JCWFhkc1hBWkdqTTBLMGl6RURWVUJjNTF1aElDbzZwcjJaeW52NmRBaUJLcTlueEMzL1RNUTd1CnFYejEwelB0b2xNMWI1Q0x6SnNMZitkZWh6d3ZWUUlnSGhlcy9MVlQ1aFdLT2IyRXFISFQ0RlFlS1dpRFBlc3cKRkNickcrR1dtSmM9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0=", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBS1B6MnovNlZ5RzFSTDMwSThteDlsTEhKMnlXakZkNAoydG1DT0ErMDVwem9HR1VvYWtlT1Y4YjJGeHIySTVldlhiNkNZcXpTQnNkQlFNTkRXd3VCVmUwQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==")

				return j
			},
		},
		{
			name: "Success",
			args: args{token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX2lkIjoiMTAxIiwiaXNzIjoidGVzdCIsInN1YiI6InRlc3QiLCJhdWQiOlsidGVzdCJdLCJleHAiOjIwNDg1MTg4MDAsIm5iZiI6MTczMjk4NjAwMCwiaWF0IjoxNzMyOTg2MDAwfQ.Hj3JBUJFVIlquVpiR3ZPj0cw2gM7nLE2mzzZOnEihx2h6zcj7ZT8hVH9-0ZUsxS8UZf7xHuBfjxeAQkARragrg"},
			want: &Claim{
				AuthID: 101,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "test",
					Subject:   "test",
					Audience:  []string{"test"},
					ExpiresAt: jwt.NewNumericDate(time.Date(2034, time.December, 1, 0, 0, 0, 0, time.Local)),
					NotBefore: jwt.NewNumericDate(time.Date(2024, time.December, 1, 0, 0, 0, 0, time.Local)),
					IssuedAt:  jwt.NewNumericDate(time.Date(2024, time.December, 1, 0, 0, 0, 0, time.Local)),
				},
			},
			wantErr: false,
			mockFn: func(a args) *JSONWebToken {
				j, _ := NewJSONWebToken("LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUJWQUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQVQ0d2dnRTZBZ0VBQWtFQW8vUGJQL3BYSWJWRXZmUWoKeWJIMlVzY25iSmFNVjNqYTJZSTREN1Rtbk9nWVpTaHFSNDVYeHZZWEd2WWpsNjlkdm9KaXJOSUd4MEZBdzBOYgpDNEZWN1FJREFRQUJBa0VBanJRTGF5Vm52NlE2WUNmbkdvQjJ5VmdrL1lRUVJYYWc3bDlFa284L2h1T3FyVFJoClgwWVMxZll4UGFFZkdDRFM2ejZQSzY0Yk15aTdBMnZQOTZKUXRRSWhBT1V2ZE1uQlYzNWlXSW1JZ01CMmdKMzgKa2hXMDVDU2dZZW01RVJNS3BoZ25BaUVBdHlLTlFnVlZ1c05xbjhNYTZkOFplTk93d0QzN3M3Y0cxMW5lcnpubwp1Y3NDSUE5c0JCWFhkc1hBWkdqTTBLMGl6RURWVUJjNTF1aElDbzZwcjJaeW52NmRBaUJLcTlueEMzL1RNUTd1CnFYejEwelB0b2xNMWI1Q0x6SnNMZitkZWh6d3ZWUUlnSGhlcy9MVlQ1aFdLT2IyRXFISFQ0RlFlS1dpRFBlc3cKRkNickcrR1dtSmM9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0=", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZ3d0RRWUpLb1pJaHZjTkFRRUJCUUFEU3dBd1NBSkJBS1B6MnovNlZ5RzFSTDMwSThteDlsTEhKMnlXakZkNAoydG1DT0ErMDVwem9HR1VvYWtlT1Y4YjJGeHIySTVldlhiNkNZcXpTQnNkQlFNTkRXd3VCVmUwQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==")

				return j
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tkn := tt.mockFn(tt.args)
			got, err := tkn.Verify(tt.args.token)
			log.Println(err)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
