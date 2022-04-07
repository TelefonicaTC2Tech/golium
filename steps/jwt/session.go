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
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwe"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/tidwall/sjson"
)

// Session contains the information of a JWT session.
type Session struct {
	Token                      string
	Payload                    []byte
	ContentType                string
	SignedMessage              *jws.Message
	EncryptedMessage           *jwe.Message
	SignatureAlgorithm         jwa.SignatureAlgorithm
	KeyEncryptionAlgorithm     jwa.KeyEncryptionAlgorithm
	ContentEncryptionAlgorithm jwa.ContentEncryptionAlgorithm
	PublicKey                  interface{}
	PrivateKey                 interface{}
}

// ConfigureSignatureAlgorithm configures a signature algorithm for the JWT (JWS).
func (s *Session) ConfigureSignatureAlgorithm(ctx context.Context, alg string) error {
	s.SignatureAlgorithm = jwa.SignatureAlgorithm(alg)
	if err := s.SignatureAlgorithm.Accept(s.SignatureAlgorithm); err != nil {
		return fmt.Errorf("invalid sign algorithm '%s'", alg)
	}
	return nil
}

// ConfigureKeyEncryptionAlgorithm configures a key encryption algorithm for the JWT (JWE).
func (s *Session) ConfigureKeyEncryptionAlgorithm(ctx context.Context, alg string) error {
	s.KeyEncryptionAlgorithm = jwa.KeyEncryptionAlgorithm(alg)
	if err := s.KeyEncryptionAlgorithm.Accept(s.KeyEncryptionAlgorithm); err != nil {
		return fmt.Errorf("invalid key encryption algorithm '%s'", alg)
	}
	return nil
}

// ConfigureContentEncryptionAlgorithm configures a content encryption algorithm for the JWT (JWE).
func (s *Session) ConfigureContentEncryptionAlgorithm(ctx context.Context, alg string) error {
	s.ContentEncryptionAlgorithm = jwa.ContentEncryptionAlgorithm(alg)
	if err := s.ContentEncryptionAlgorithm.Accept(s.ContentEncryptionAlgorithm); err != nil {
		return fmt.Errorf("invalid content encryption algorithm '%s'", alg)
	}
	return nil
}

// ConfigureSymmetricKey configures the symmetric key. It sets this key as public and private key.
func (s *Session) ConfigureSymmetricKey(ctx context.Context, symmetricKey string) {
	s.PublicKey = []byte(symmetricKey)
	s.PrivateKey = s.PublicKey
}

// ConfigurePublicKey configures the public key to verify the signature
// of a JWT token or to encrypt a JWE token.
func (s *Session) ConfigurePublicKey(ctx context.Context, publicKeyPEM string) error {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return errors.New("public key is not a valid PEM")
	}
	var err error
	switch block.Type {
	case "PUBLIC KEY", "RSA PUBLIC KEY":
		s.PublicKey, err = x509.ParsePKIXPublicKey(block.Bytes)
	default:
		return fmt.Errorf("invalid PEM type '%s'", block.Type)
	}
	if err != nil {
		return fmt.Errorf("failed processing the public key: %w", err)
	}
	return nil
}

// ConfigurePrivateKey configures the private key to sign a JWT token or to decrypt a JWE token.
func (s *Session) ConfigurePrivateKey(ctx context.Context, privateKeyPEM string) error {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return errors.New("private key is not a valid PEM")
	}
	var err error
	switch block.Type {
	case "RSA PRIVATE KEY":
		s.PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return fmt.Errorf("invalid PEM type '%s'", block.Type)
	}
	if err != nil {
		return fmt.Errorf("failed processing the private key: %w", err)
	}
	return nil
}

// ConfigurePayloadWithContentType configures the payload and the content type (cty header).
func (s *Session) ConfigurePayloadWithContentType(
	ctx context.Context,
	payload, contentType string,
) {
	s.Payload = []byte(payload)
	s.ContentType = contentType
}

// ConfigureJSONPayload configures the JWT payload with a map of properties.
func (s *Session) ConfigureJSONPayload(ctx context.Context, props map[string]interface{}) error {
	var json string
	var err error
	for key, value := range props {
		if json, err = sjson.Set(json, key, value); err != nil {
			return fmt.Errorf("failed setting property '%s' with value '%s' in the request body: %w",
				key, value, err)
		}
	}
	s.ConfigurePayloadWithContentType(ctx, json, "JSON")
	return nil
}

// GenerateSignedJWTInContext builds a JWT with signed payload and stores it in the context.
func (s *Session) GenerateSignedJWTInContext(ctx context.Context, ctxtKey string) error {
	if s.Payload == nil {
		return errors.New("a payload is required")
	}
	if s.SignatureAlgorithm == "" {
		return errors.New("a signature algorithm is required")
	}
	if s.PrivateKey == nil {
		return errors.New("a private key is required")
	}
	token, err := sign(s.Payload, s.SignatureAlgorithm, s.PrivateKey, s.ContentType)
	if err != nil {
		return err
	}
	s.Token = token
	golium.GetContext(ctx).Put(ctxtKey, token)
	return nil
}

// GenerateEncryptedJWTInContext builds a JWT with encrypted payload and stores it in the context.
func (s *Session) GenerateEncryptedJWTInContext(ctx context.Context, ctxtKey string) error {
	if s.Payload == nil {
		return errors.New("a payload is required")
	}
	if s.KeyEncryptionAlgorithm == "" {
		return errors.New("a key encryption algorithm is required")
	}
	if s.ContentEncryptionAlgorithm == "" {
		return errors.New("a content encryption algorithm is required")
	}
	if s.PublicKey == nil {
		return errors.New("a public key is required")
	}
	token, err := encrypt(
		s.Payload, s.KeyEncryptionAlgorithm, s.PublicKey, s.ContentEncryptionAlgorithm, s.ContentType)
	if err != nil {
		return err
	}
	s.Token = token
	golium.GetContext(ctx).Put(ctxtKey, token)
	return nil
}

// GenerateSignedEncryptedJWTInContext builds a JWT with signed encrypted payload
// and stores it in the context.
// The payload is signed first. Then the whole JWT is considered as payload for encryption phase.
// The content type header (cty) of the final token is set to JWT.
func (s *Session) GenerateSignedEncryptedJWTInContext(ctx context.Context, ctxtKey string) error {
	if err := s.GenerateSignedJWTInContext(ctx, ctxtKey); err != nil {
		return err
	}
	s.Payload = []byte(s.Token)
	s.ContentType = ContentTypeJWT
	return s.GenerateEncryptedJWTInContext(ctx, ctxtKey)
}

// ProcessSignedJWT reads a signed JWT and stores the data in the session.
// This method does not validate the token; use ValidateJWT for this purpose.
func (s *Session) ProcessSignedJWT(ctx context.Context, token string) error {
	s.Token = token
	var err error
	s.SignedMessage, s.Payload, err = parse(token)
	return err
}

// ProcessEncryptedJWT reads an encrypted JWT (JWE) and stores in the session
// the token, encrypted message and payload.
// There is no validation method for encrypted tokens.
func (s *Session) ProcessEncryptedJWT(ctx context.Context, token string) error {
	s.Token = token
	var err error
	s.EncryptedMessage, s.Payload, err = decrypt([]byte(token), s.KeyEncryptionAlgorithm, s.PrivateKey)
	return err
}

// ProcessSignedEncryptedJWT reads a signed encrypted JWT and stores in the session
// the embedded signed token, the encrypted message, the signed message and the signed payload.
// Note that this token expects that a signed JWT token is the payload of a JWE token.
func (s *Session) ProcessSignedEncryptedJWT(ctx context.Context, token string) error {
	if err := s.ProcessEncryptedJWT(ctx, token); err != nil {
		return err
	}
	if s.EncryptedMessage.ProtectedHeaders().ContentType() != ContentTypeJWT {
		return errors.New("content type of the encrypted token is not JWT")
	}
	return s.ProcessSignedJWT(ctx, string(s.Payload))
}

// ValidateJWT checks that the token is valid (the claims and the signature of the token).
// Note that JWE tokens are not validated.
func (s *Session) ValidateJWT(ctx context.Context) error {
	if err := s.ValidateJWTRequirements(); err != nil {
		return err
	}
	if err := verify(s.Token, s.SignatureAlgorithm, s.PublicKey); err != nil {
		return fmt.Errorf("token is invalid: %w", err)
	}
	return nil
}

// ValidateInvalidJWT checks that the token is invalid (the claims and the signature of the token).
// Note that JWE tokens are not validated.
func (s *Session) ValidateInvalidJWT(ctx context.Context, expectedError string) error {
	if err := s.ValidateJWTRequirements(); err != nil {
		return err
	}
	err := verify(s.Token, s.SignatureAlgorithm, s.PublicKey)
	if err == nil {
		return errors.New("token is valid")
	}
	if expectedError == "" || strings.Contains(err.Error(), expectedError) {
		return nil
	}
	return fmt.Errorf("token is invalid: %w", err)
}

func (s *Session) ValidateJWTRequirements() error {
	if s.Token == "" {
		return errors.New("no token has been processed")
	}
	if s.SignatureAlgorithm == "" {
		return errors.New("a signature algorithm is required")
	}
	if s.PublicKey == nil {
		return errors.New("a public key is required")
	}
	if s.SignedMessage == nil {
		// Encrypted messages cannot be verified
		return nil
	}
	return nil
}

// ValidatePayloadJSONProperties checks if the payload contains a map of expected properties.
func (s *Session) ValidatePayloadJSONProperties(
	ctx context.Context,
	expectedPayload map[string]interface{},
) error {
	if s.Payload == nil {
		return errors.New("no token has been processed")
	}
	m := golium.NewMapFromJSONBytes(s.Payload)
	for key, expectedValue := range expectedPayload {
		value := m.Get(key)
		if expectedValue != value {
			return fmt.Errorf(
				"mismatch payload property '%s': expected '%v', actual '%v'", key, expectedValue, value)
		}
	}
	return nil
}
