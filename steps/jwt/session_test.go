// Copyright 2021 Telefonica Cybersecurity & Cloud Tech SL
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwt

import (
	"context"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
)

const (
	rs256         = "RS256"
	rsa15         = "RSA1_5"
	a128cbchs256  = "A128CBC-HS256"
	signPublicKey = `
-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoNNaEB/t0c4kZNhoz9G5
t5fq+1OQVl4szeB+MwicYdzf2Ho23w7/0TFBa8azF++o0CYJYAgqRh/MCKXD7gsF
iXjm/TSj2M0GgJXGtPn/ZS7ULOLKm/Fp+mB6209qtcBPON9L4c8++tUCn6wwIqaO
x6OzdjUMMC9jhSkOJChBxRJtkJjwf4usRi9nYCCxsVfVeJGOxwEOghMxjdA5vJCx
XcXzFnhGQigT6EHUoxLg7JRvKgAdvN9+lpAvn8lDnCJsrCjDsFrz0BAewiKBcM8C
tR5AcCFf3pG+oQ7Uq62idvjmKsCB96jkTfVidr19Rj03E1Lg+2RCuS+4Qtylnl4b
RQIDAQAB
-----END RSA PUBLIC KEY-----
`
	encryptPublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu+5LoHNwgUG8+bn5Mn51
uTQ6EfP6XlblPv23i1LmTHpoGsgGirnIwYkrjHUqxNOP1XOq32p+SJwZdRZOLyoQ
Gigf6gSoYHuCG2AfsPdkVjd7JBAiejbNwdODn9coqGbgVJFi3iR4qfI9GSMmE5Qj
pvLvv/XKSVTkAkobgxmeKs6RzdWepyWgOXUgWdyYJXj5B7yCotMeYDrhbwtfmX0j
yoyMN0hyLJRDG7UPbuvl+PrHzYiC4TVo+cQm89qOJnnvAsGTg7QYcK1854pU8evh
CDRvfCc1KNcyPbJj8ZjrvamK16KbxBxwsMXBuRznKQufNx60+Ej63vwRBxcH+y7T
MQIDAQAB
-----END PUBLIC KEY-----
`
	wrongTypeKey = `
-----BEGIN EC PRIVATE KEY-----
MHQCAQEEIFL3sLnioGcDvHWM/BPlNw96BOx1KKco2qsq4UwhQUosoAcGBSuBBAAK
oUQDQgAEXs1Fmq4QdPAbn3NycdEU+HOjc3kW9efbso2kI/vdDTWcSCMk310s53G3
tRClDBPPuuJAsKghbPfaTaUpmXFCNA==
-----END EC PRIVATE KEY-----
`
	signPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAoNNaEB/t0c4kZNhoz9G5t5fq+1OQVl4szeB+MwicYdzf2Ho2
3w7/0TFBa8azF++o0CYJYAgqRh/MCKXD7gsFiXjm/TSj2M0GgJXGtPn/ZS7ULOLK
m/Fp+mB6209qtcBPON9L4c8++tUCn6wwIqaOx6OzdjUMMC9jhSkOJChBxRJtkJjw
f4usRi9nYCCxsVfVeJGOxwEOghMxjdA5vJCxXcXzFnhGQigT6EHUoxLg7JRvKgAd
vN9+lpAvn8lDnCJsrCjDsFrz0BAewiKBcM8CtR5AcCFf3pG+oQ7Uq62idvjmKsCB
96jkTfVidr19Rj03E1Lg+2RCuS+4Qtylnl4bRQIDAQABAoIBABKWu0c32Y4xjIVX
ei3jKNsupQttqjZBZl5Zf1y4txKcrAbigWsg2bK9RbmGWvb+TX3Zl6XQ68n1LOkm
99GQ1pAAOHq11eZeNE0ygqgyaTGxyvZxNEf4DG7TLgAhWs0tMDr7nFK6WKY3brkz
9tBafhBXPIwCL6l2IMOobikBujBj9Pe5mpTwCzbqE8TEzfnVB5DzqroFR2W/O8/6
af6f7Co5cY+kta38G9hCq78wW4iU5qesynVw7K50mzUMCHRVF5SsKZwHmgNxjv7O
+B7jt1mzIKjPhBhN8ZjEIsFsKgZgOnZFSn8CKhxNGU9FYAo/r1Ih5TEKsPufAVbd
ChWkSoECgYEA0G8xd8lVM0WylukoI+A9u7uvpiiRkiksPmNuwhEtFSxvSLFKhqeM
gG/iDSe6DJGKSCgyCnLv6QVAwLcIhKmzJKUYrJb9OzxsBlaMpXucisWpKqRCKB7/
RQhK+rkxrKtExWGDvk/0+nbTvqXVPnd5ifp0QtfSHtp6OT5GX0xttWECgYEAxYbW
roeN1B2b+/+WFTafguA5e72g4WUDa6OQNRdePmfUdZHyzGOP7CMF/BJ59VeLRsn+
ZfHI2cD8YCb+8jxtmDK8M1h0dr3mUcrbExU09hRYhBKiL76xMsV8QU7zerHXcyDp
Mjk0D743lhUmllm02sLJFEGEfGc+Idk2+TW1DGUCgYEAyAGa303zkrKTr1nmKZ7Y
vhdYckHFhhI6IVe6hUCEGSg9VOzDDbkjCm/R4zu2vK6/mYPwmLQ34EspGoPICbzp
aQV/SsXMExZktiRA695UlZkcPg3Gacdsvio6AKLKttzVre1nxKvm8JwrjWqF2F4+
4xbQjv+X4gFVfS5zyqiFMaECgYBt8+YTFw/jEGxg9WAVBOf8EVbOQ7uHXBRwWYcP
lqd2c5O3snuGPLHDz6coLxzGbmnwCMbc9p9IX33dBDgMnYigHTXYGxgRdRn9U79p
OvfVN3QiaMDxdOPskDPfotQz60U0KBDHTUJmtQr6N2HYda0PzTfjV6kpGstiSiio
xrW2ZQKBgHlbw5kRWNKe/UCmyqL/EgPH8mE3Q/s+zFaMfu9bV3D/n5iffAS9POHv
pFUZU0CeuJcHF9D8VjnzjeisJW5uK/EEOVKaerTL01qjYE9nzOl18xvevKEgfD2l
jXW8Z4nFjhcqsv1DlylqVga/B47YHfg0qTcfLJ78gMaiGMNaLQhF
-----END RSA PRIVATE KEY-----
`
	encryptPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAu+5LoHNwgUG8+bn5Mn51uTQ6EfP6XlblPv23i1LmTHpoGsgG
irnIwYkrjHUqxNOP1XOq32p+SJwZdRZOLyoQGigf6gSoYHuCG2AfsPdkVjd7JBAi
ejbNwdODn9coqGbgVJFi3iR4qfI9GSMmE5QjpvLvv/XKSVTkAkobgxmeKs6RzdWe
pyWgOXUgWdyYJXj5B7yCotMeYDrhbwtfmX0jyoyMN0hyLJRDG7UPbuvl+PrHzYiC
4TVo+cQm89qOJnnvAsGTg7QYcK1854pU8evhCDRvfCc1KNcyPbJj8ZjrvamK16Kb
xBxwsMXBuRznKQufNx60+Ej63vwRBxcH+y7TMQIDAQABAoIBAQCon6SUD4C/OfEK
UehbekS/LTF9smDQDUAdSSJLjNK/hIWsohXcm96aaS3+FZOOkBXa/LIxTSiKliXx
fVYh06gnECGypQM/rxKK5bEJ5LDO+3EuZpvI9Suh4tuTrEb683QN7XW8xRTPyF8y
EuuzXZSv1ANzRmN/cQA1XbFZ7L6SwMKNiYn3FVkORZvABnCd80Nc4tqqMz22nRdl
kWFw/gI+h6igIv1DeSe2gJbHe0HrbqatW2IbAgHiBptr7u4BxhxB2ppDaXoxzksU
vyoc0QmVRPFy01EXt6kxiCG2wF+87VtY/cCiOqUxy3xV3SHriVQx8ewOR2W2GHgG
md7b6FlRAoGBAMLrAnD2jqUWVvt7b/51z+dwcbLiAm0VV75szRhnOuY52SIfB+3h
JdTW6ujnsPvHaWWVxv8tX+QaxzS4UFsakrgFcaBvKBjaicyxxoJ+/Nu0wok8ZFmy
qzfLBDV2919TfaDMwML+N4pHIpTAxFFxKlqahf6cIoJ86RWcL3kCkV83AoGBAPbS
wGiZERX5NWZoPqj01Beq8VdrW4P9v9kBixg4lRTHog57Uhk9zaKrE67uGhVwEEBU
nnZ+D0+90AwUu1Cy+HITGhaUBFIUJLY9RyT0nih944MzUwsybcxmlucR6bNLRlRp
v2mjawGPDoFDq/IPhttZOxZPHBd4CnUgZm3OvQTXAoGAK2hw471U6Rj/iAPmXhHY
mh8lgwPoLGjbYJIUXsHmkQ0C+SFV/7jrVuoB6JpohLnVFAV2CrANMdxwzqHZa2CQ
miDEPElk8ZwBoi9ZGQi0wS0RQcTMSFmM3eD9b/atgnIygRP4PbSlo8rRvbTsQ4Lj
Psg43QnieZLdya09uUJEI6MCgYEAkN2XYozcU1pGNkne5Ql1ZkLFjbqMJwcKz9Ix
ElE7ZsvY2MkWoYv9oojob5Z+JrD0SN2heAh68iGE92I/opi4azO87x2G/6mk9nU2
yYDtRvTEUOAR0JOTkBFyZkLEOKBoseizGMx6ZJrTN5lBVTw5uYpAvNJHuZqSALa4
h6B8nlcCgYEAkcS2h37xeCNfdLZh2GDEAnEcy8w/op6F7NzKWSnpFskOmIEUsCvM
5vEiKbAHyDTDPyoy0Zx8waMV6eAhZGxD65uqifsH8dj6ltMAGvOfHR9Cs1qD3TQt
HippkcLN9e1ETEx1zNxmWAXFHOX2Ia2kSccmuuCKEzglN47DrJE5/Gs=
-----END RSA PRIVATE KEY-----
`
	token = "eyJhbGciOiJSU0ExXzUiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwidHlwIj" +
		"oiSldFIn0.UxTx-Yh41FVC1V9FJdjKDbBcoX2-0As-S1IBr1KZibxwTDT_n1wlw0x3MzckR78cRWMBajU80a" +
		"4rsSRCAY0Gzi0aehuQgbrYufg1EURCE7X_IQAtDH-9nvgcVzfWxonfUCVvWlH95Zk9FrWPcQhhxBbR1aDuyr" +
		"Ozl4iHt2p0Grrt37e8EPyqJFE7edC2LVI-H3CfYaHlrILJFukJR1-swpEzl9r-uXGifAgEtSrT9DL56wLLCm" +
		"OxIwusMQvRvMl53uakmVYfQtWF7-Ibr4aMUOBYr9H31UhZP5ZbaVAzONRngFpnTPWbvhJg12kycT1snd8-8m" +
		"f8uRnfqcPkQWmrMw.2wNQ3hktIGHz9dpt3nOZvw.Ki9j5uT7JeSNGpO-JMWvLMpJRHgBHqSasV3dBoDH4pHI" +
		"aTT7n2A19_vLiLo4df0xtQGqaHFSgS49vcuV2N9yOuW49fZv6nJuGvXkk8HJcRDHtrSS3_AhNSX_zBrJw2do" +
		"FwHjKixXxyS1nboX3Q-p7AaTRIx9l6mRetOa_xwXogEusM9GEMqKP6GkxNE669j9MwR3DIDFO83S3Ntj4GiK" +
		"XXEhJgRQSQgsgL4qvuyfqK4VSY22m_Z6ndCtL5hsvC1chhF0PvB5M-6U36ynRc9_tx8iwv-Zwy_Ja5q0gPHY" +
		"R7fcxJ-u9eHz7ZZMLB0x6qyjDit3OsVF30ehiJvC6WfFq8v1MuDPdtHuLScTyCdG8jE7UtP0Djt9JxvvjvHY" +
		"dt9-RwgOLgVMAlQCBt4PA1NvUcjzKkfQaWrKR65Uc_MVtXAuYVqcR3lcZH5DbjS84UyMFAYf-5QqW7xctRg8" +
		"0VW_H7rvtI_NsPGtI6Qah3XKu5dkKmCokLW9MhSReQpiJW5oebpAL7vJ6fh_S9xPNF-msw.cUoLbK2mGgvaO" +
		"Od3K_Jh9Q"
)

func TestConfigureSignatureAlgorithm(t *testing.T) {
	tests := []struct {
		name    string
		alg     string
		wantErr bool
	}{
		{
			name:    "Valid signature algorithm",
			alg:     "HS512",
			wantErr: false,
		},
		{
			name:    "Invalid signature algorithm",
			alg:     "invalid_signature",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			if err := s.ConfigureSignatureAlgorithm(
				context.Background(), tt.alg); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigureSignatureAlgorithm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigureKeyEncryptionAlgorithm(t *testing.T) {
	tests := []struct {
		name    string
		alg     string
		wantErr bool
	}{
		{
			name:    "Valid encryption algorithm",
			alg:     rsa15,
			wantErr: false,
		},
		{
			name:    "Invalid encryption algorithm",
			alg:     "invalid_encryption",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			if err := s.ConfigureKeyEncryptionAlgorithm(
				context.Background(), tt.alg); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigureKeyEncryptionAlgorithm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigureContentEncryptionAlgorithm(t *testing.T) {
	tests := []struct {
		name    string
		alg     string
		wantErr bool
	}{
		{
			name: "Valid content encryption algorithm",
			alg:  a128cbchs256,

			wantErr: false,
		},
		{
			name:    "Invalid content encryption algorithm",
			alg:     "invalid_content_encryption",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			if err := s.ConfigureContentEncryptionAlgorithm(
				context.Background(), tt.alg); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigureContentEncryptionAlgorithm() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestConfigureSymmetricKey(t *testing.T) {
	type args struct {
		ctx          context.Context
		symmetricKey string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Configure Symmetric Key",
			args: args{
				ctx: context.Background(),
				symmetricKey: "sign_symmetric_key_that_is_long_enough_for_algorithm_" +
					"HS512_(with_more_than 256 bits!)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.ConfigureSymmetricKey(tt.args.ctx, tt.args.symmetricKey)
		})
	}
}

func TestConfigurePublicKey(t *testing.T) {
	type args struct {
		ctx          context.Context
		publicKeyPEM string
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name:    "Not valid pem",
			wantErr: true,
			args: args{
				ctx:          context.Background(),
				publicKeyPEM: "WRONG",
			},
		},
		{
			name:    "Valid Public Key",
			wantErr: false,
			args: args{
				ctx:          context.Background(),
				publicKeyPEM: encryptPublicKey,
			},
		},
		{
			name:    "Valid RSA Public Key",
			wantErr: false,
			args: args{
				ctx:          context.Background(),
				publicKeyPEM: signPublicKey,
			},
		},
		// {
		// 	name:    "Invalid Key",
		// 	wantErr: true,
		// 	args: args{
		// 		ctx:          context.Background(),
		// 		publicKeyPEM: invalidPublicKey,
		// 	},
		// },
		{
			name:    "Default",
			wantErr: true,
			args: args{
				ctx:          context.Background(),
				publicKeyPEM: wrongTypeKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			if err := s.ConfigurePublicKey(tt.args.ctx, tt.args.publicKeyPEM); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigurePublicKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigurePrivateKey(t *testing.T) {
	type args struct {
		ctx           context.Context
		privateKeyPEM string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not valid pem",
			wantErr: true,
			args: args{
				ctx:           context.Background(),
				privateKeyPEM: "WRONG",
			},
		},
		{
			name:    "Valid Private Key",
			wantErr: false,
			args: args{
				ctx:           context.Background(),
				privateKeyPEM: signPrivateKey,
			},
		},
		// {
		// 	name:    "Invalid Key",
		// 	wantErr: true,
		// 	args: args{
		// 		ctx:          context.Background(),
		// 		publicKeyPEM: invalidPublicKey,
		// 	},
		// },
		{
			name:    "Default",
			wantErr: true,
			args: args{
				ctx:           context.Background(),
				privateKeyPEM: wrongTypeKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			if err := s.ConfigurePrivateKey(tt.args.ctx, tt.args.privateKeyPEM); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigurePrivateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigureJSONPayload(t *testing.T) {
	paramsInput := make(map[string]interface{})
	paramsInput["title"] = "foo1"
	type args struct {
		ctx   context.Context
		props map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid payload",
			args: args{
				ctx:   context.Background(),
				props: paramsInput,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			if err := s.ConfigureJSONPayload(tt.args.ctx, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigureJSONPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateSignedJWTInContext(t *testing.T) {
	type args struct {
		ctx     context.Context
		ctxtKey string
	}
	tests := []struct {
		name               string
		args               args
		payload            []byte
		signatureAlgorithm string
		privateKey         interface{}
		wantErr            bool
	}{
		{
			name:    "Nil payload",
			payload: nil,
			wantErr: true,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name:               "Empty signature",
			payload:            []byte("payload"),
			signatureAlgorithm: "",
			wantErr:            true,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name:               "Nil private key",
			payload:            []byte("payload"),
			signatureAlgorithm: rs256,
			privateKey:         nil,
			wantErr:            true,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name:               "Valid token",
			payload:            []byte("payload"),
			signatureAlgorithm: rs256,
			privateKey:         "valid",
			wantErr:            false,
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxGolium := golium.InitializeContext(tt.args.ctx)
			ctx := InitializeContext(ctxGolium)
			s := &Session{}
			s.Payload = tt.payload
			s.SignatureAlgorithm = jwa.SignatureAlgorithm(tt.signatureAlgorithm)
			if tt.privateKey == nil {
				s.PrivateKey = tt.privateKey
			} else {
				s.ConfigurePrivateKey(ctx, signPrivateKey)
			}
			if err := s.GenerateSignedJWTInContext(
				ctx, tt.args.ctxtKey); (err != nil) != tt.wantErr {
				t.Errorf("Session.GenerateSignedJWTInContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateEncryptedJWTInContext(t *testing.T) {
	type args struct {
		ctx     context.Context
		ctxtKey string
	}
	tests := []struct {
		name                       string
		args                       args
		payload                    []byte
		keyEncryptionAlgorithm     string
		contentEncryptionAlgorithm string
		publicKey                  interface{}
		wantErr                    bool
	}{
		{
			name:    "Nil payload",
			payload: nil,
			wantErr: true,
			args: args{
				ctx:     context.Background(),
				ctxtKey: "jwt.jwse",
			},
		},
		{
			name:                   "Empty Encryption Algorithm",
			payload:                []byte("payload"),
			keyEncryptionAlgorithm: "",
			wantErr:                true,
			args: args{
				ctx:     context.Background(),
				ctxtKey: "jwt.jwse",
			},
		},
		{
			name:                       "Empty Content Encryption Algorithn",
			payload:                    []byte("payload"),
			keyEncryptionAlgorithm:     rsa15,
			contentEncryptionAlgorithm: "",
			wantErr:                    true,
			args: args{
				ctx:     context.Background(),
				ctxtKey: "jwt.jwse",
			},
		},
		{
			name:                       "Nil Public Key",
			payload:                    []byte("payload"),
			keyEncryptionAlgorithm:     rsa15,
			contentEncryptionAlgorithm: a128cbchs256,
			publicKey:                  nil,
			wantErr:                    true,
			args: args{
				ctx:     context.Background(),
				ctxtKey: "jwt.jwse",
			},
		},
		{
			name:                       "Valid token",
			payload:                    []byte("payload"),
			keyEncryptionAlgorithm:     rsa15,
			contentEncryptionAlgorithm: a128cbchs256,
			publicKey:                  "valid",
			wantErr:                    false,
			args: args{
				ctx:     context.Background(),
				ctxtKey: "jwt.jwse",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxGolium := golium.InitializeContext(tt.args.ctx)
			ctx := InitializeContext(ctxGolium)
			s := &Session{}
			s.Payload = tt.payload
			s.KeyEncryptionAlgorithm = jwa.KeyEncryptionAlgorithm(tt.keyEncryptionAlgorithm)
			s.ContentEncryptionAlgorithm = jwa.ContentEncryptionAlgorithm(tt.contentEncryptionAlgorithm)
			if tt.publicKey == nil {
				s.PublicKey = tt.publicKey
			} else {
				s.ConfigurePublicKey(ctx, encryptPublicKey)
			}
			if err := s.GenerateEncryptedJWTInContext(ctx, tt.args.ctxtKey); (err != nil) != tt.wantErr {
				t.Errorf("Session.GenerateEncryptedJWTInContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateSignedEncryptedJWTInContext(t *testing.T) {
	type args struct {
		ctx     context.Context
		ctxtKey string
	}
	tests := []struct {
		name        string
		args        args
		signedError bool
		wantErr     bool
	}{
		{
			name:        "Error generating signed JWT in Context",
			signedError: true,
			args: args{
				ctx:     context.Background(),
				ctxtKey: "jwt.jwse",
			},
			wantErr: true,
		},
		{
			name:        "Valid generation",
			signedError: false,
			args: args{
				ctx:     context.Background(),
				ctxtKey: "jwt.jwse",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxGolium := golium.InitializeContext(tt.args.ctx)
			ctx := InitializeContext(ctxGolium)
			s := &Session{}

			s.Payload = []byte("payload")
			s.SignatureAlgorithm = rs256
			if !tt.signedError {
				s.ConfigurePrivateKey(ctx, signPrivateKey)
			}
			s.KeyEncryptionAlgorithm = rsa15
			s.ContentEncryptionAlgorithm = a128cbchs256
			s.ConfigurePublicKey(ctx, encryptPublicKey)

			if err := s.GenerateSignedEncryptedJWTInContext(
				ctx, tt.args.ctxtKey); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.GenerateSignedEncryptedJWTInContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessSignedEncryptedJWT(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		tokenError bool
		args       args
		wantErr    bool
	}{
		{
			name: "Valid token",
			args: args{
				ctx: context.Background(),
			},
			tokenError: false,
			wantErr:    false,
		},
		{
			name: "Wrong token",
			args: args{
				ctx: context.Background(),
			},
			tokenError: true,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxGolium := golium.InitializeContext(tt.args.ctx)
			ctx := InitializeContext(ctxGolium)
			s := &Session{}
			if tt.tokenError {
				s.Token = "fakeToken"
			} else {
				s.Token = token
			}
			s.ConfigurePrivateKey(ctx, encryptPrivateKey)
			s.KeyEncryptionAlgorithm = rsa15
			if err := s.ProcessSignedEncryptedJWT(ctx, s.Token); (err != nil) != tt.wantErr {
				t.Errorf("Session.ProcessSignedEncryptedJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWTRequirements(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name               string
		token              string
		signatureAlgorithm jwa.SignatureAlgorithm
		publicKey          interface{}
		signedMessage      *jws.Message
		args               args
		wantErr            bool
	}{
		{
			name: "Empty token",
			args: args{
				ctx: context.Background(),
			},
			token:   "",
			wantErr: true,
		},
		{
			name: "Empty signature",
			args: args{
				ctx: context.Background(),
			},
			token:              token,
			signatureAlgorithm: "",
			wantErr:            true,
		},
		{
			name: "Empty Public Key",
			args: args{
				ctx: context.Background(),
			},
			token:              token,
			signatureAlgorithm: rs256,
			publicKey:          nil,
			wantErr:            true,
		},
		{
			name: "Nil Signed Message",
			args: args{
				ctx: context.Background(),
			},
			token:              token,
			signatureAlgorithm: rs256,
			publicKey:          "valid",
			signedMessage:      nil,
			wantErr:            false,
		},
		{
			name: "Signed Message",
			args: args{
				ctx: context.Background(),
			},
			token:              token,
			signatureAlgorithm: rs256,
			publicKey:          "valid",
			signedMessage:      &jws.Message{},
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxGolium := golium.InitializeContext(tt.args.ctx)
			ctx := InitializeContext(ctxGolium)
			s := &Session{}
			s.Token = tt.token
			s.SignatureAlgorithm = tt.signatureAlgorithm
			if tt.publicKey == nil {
				s.PublicKey = tt.publicKey
			} else {
				s.ConfigurePublicKey(ctx, encryptPublicKey)
			}

			s.SignedMessage = tt.signedMessage

			if err := s.ValidateJWTRequirements(); (err != nil) != tt.wantErr {
				t.Errorf("Session.validateJWTRequirements() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSession_ValidatePayloadJSONProperties(t *testing.T) {
	testPayload := make(map[string]interface{})
	testPayload["title"] = "foo1"
	tests := []struct {
		name            string
		payload         []byte
		expectedPayload map[string]interface{}
		wantErr         bool
	}{
		{
			name:    "Nil payload",
			payload: nil,
			wantErr: true,
		},
		{
			name:            "Valid payload",
			wantErr:         false,
			payload:         []byte("nonEmptyPayload"),
			expectedPayload: testPayload,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxGolium := golium.InitializeContext(context.Background())
			ctx := InitializeContext(ctxGolium)
			s := &Session{}
			if tt.payload == nil {
				s.Payload = tt.payload
			} else {
				s.ConfigureJSONPayload(ctx, testPayload)
			}
			if err := s.ValidatePayloadJSONProperties(ctx, tt.expectedPayload); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidatePayloadJSONProperties() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
