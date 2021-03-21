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
	"fmt"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

const (
	// TypeJWT for signed tokens in typ header.
	TypeJWT = "JWT"
	// TypeJWE for encrypted tokens in typ header.
	TypeJWE = "JWE"
	// ContentTypeJWT for signed encrypted tokens where the encrypted token will include a cty header with this value.
	ContentTypeJWT = "JWT"
)

////////////////////////////////////////////////////////////////////////////////////////
// JWS (signed tokens)
////////////////////////////////////////////////////////////////////////////////////////

// sign applies JWS to sign the payload with an algorithm, a private key.
// It also sets a JWT header with the content type.
func sign(payload []byte, signatureAlgorithm jwa.SignatureAlgorithm, privateKey interface{}, contentType string) (string, error) {
	headers := jws.NewHeaders()
	headers.Set(jws.TypeKey, TypeJWT)
	if contentType != "" {
		headers.Set(jws.ContentTypeKey, contentType)
	}
	withHeaders := jws.WithHeaders(headers)
	signed, err := jws.Sign(payload, signatureAlgorithm, privateKey, withHeaders)
	if err != nil {
		return "", fmt.Errorf("failed signing JWT: %w", err)
	}
	return string(signed), nil
}

// parse processes a signed token.
// It returns the message (including payload, headers and signature), payload, and error.
func parse(token string) (*jws.Message, []byte, error) {
	msg, err := jws.ParseString(token)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing the JWT: %w", err)
	}
	return msg, msg.Payload(), nil
}

// verify validates the claims and the signature of the token.
func verify(token string, signatureAlgorithm jwa.SignatureAlgorithm, publicKey interface{}) error {
	if err := verifyClaims(token); err != nil {
		return err
	}
	return verifySignature(token, signatureAlgorithm, publicKey)
}

// verifyClaims validates claims of the token.
func verifyClaims(token string) error {
	_, err := jwt.ParseString(token, jwt.WithValidate(true))
	return err
}

// verify validates a JWS token checking the signature.
func verifySignature(token string, signatureAlgorithm jwa.SignatureAlgorithm, publicKey interface{}) error {
	_, err := jws.Verify([]byte(token), signatureAlgorithm, publicKey)
	return err
}

////////////////////////////////////////////////////////////////////////////////////////
// JWE (encrypted tokens)
////////////////////////////////////////////////////////////////////////////////////////

// encrypt a payload.
// It returns the token and error.
func encrypt(payload []byte, keyEncryptionAlgorithm jwa.KeyEncryptionAlgorithm, publicKey interface{}, contentEncryptionAlgorithm jwa.ContentEncryptionAlgorithm, contentType string) (string, error) {
	headers := jwe.NewHeaders()
	headers.Set(jws.TypeKey, TypeJWE)
	if contentType != "" {
		headers.Set(jws.ContentTypeKey, contentType)
	}
	withHeaders := jwe.WithProtectedHeaders(headers)
	encrypted, err := jwe.Encrypt(payload, keyEncryptionAlgorithm, publicKey, contentEncryptionAlgorithm, jwa.NoCompress, withHeaders)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt token: %w", err)
	}
	return string(encrypted), nil
}

// decrypt processes an encrypted token.
// It returns the message (including headers), decrypted payload, and error.
func decrypt(encrypted []byte, keyEncryptionAlgorithm jwa.KeyEncryptionAlgorithm, privateKey interface{}) (*jwe.Message, []byte, error) {
	msg, err := jwe.Parse(encrypted)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing the encrypted token: %w", err)
	}
	payload, err := msg.Decrypt(keyEncryptionAlgorithm, privateKey)
	if err != nil {
		return msg, nil, fmt.Errorf("failed decrypting the encrypted token: %w", err)
	}
	return msg, payload, nil
}
