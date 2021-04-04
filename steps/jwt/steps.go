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
	"fmt"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
)

// Steps type is responsible to initialize the JWT steps in godog framework.
type Steps struct {
}

// InitializeSteps adds JWT steps to the scenario context.
// It implements StepsInitializer interface.
// It returns a new context (context is immutable) with the JWT Context.
func (s Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the HTTP session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)
	// Initialize the steps
	scenCtx.Step(`^the JWT signature algorithm "([^"]*)"$`, func(alg string) error {
		return session.ConfigureSignatureAlgorithm(ctx, alg)
	})
	scenCtx.Step(`^the JWT key encryption algorithm "([^"]*)"$`, func(alg string) error {
		return session.ConfigureKeyEncryptionAlgorithm(ctx, alg)
	})
	scenCtx.Step(`^the JWT content encryption algorithm "([^"]*)"$`, func(alg string) error {
		return session.ConfigureContentEncryptionAlgorithm(ctx, alg)
	})
	scenCtx.Step(`^the JWT payload with the JSON properties$`, func(t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the payload: %w", err)
		}
		return session.ConfigureJSONPayload(ctx, props)
	})
	scenCtx.Step(`^the JWT symmetric key$`, func(message *godog.DocString) error {
		return session.ConfigureSymmetricKey(ctx, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^the JWT public key$`, func(message *godog.DocString) error {
		return session.ConfigurePublicKey(ctx, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^the JWT private key$`, func(message *godog.DocString) error {
		return session.ConfigurePrivateKey(ctx, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I generate a signed JWT and store it in context "([^"]*)"$`, func(ctxtKey string) error {
		return session.GenerateSignedJWTInContext(ctx, golium.ValueAsString(ctx, ctxtKey))
	})
	scenCtx.Step(`^I generate an encrypted JWT and store it in context "([^"]*)"$`, func(ctxtKey string) error {
		return session.GenerateEncryptedJWTInContext(ctx, golium.ValueAsString(ctx, ctxtKey))
	})
	scenCtx.Step(`^I generate a signed encrypted JWT and store it in context "([^"]*)"$`, func(ctxtKey string) error {
		return session.GenerateSignedEncryptedJWTInContext(ctx, golium.ValueAsString(ctx, ctxtKey))
	})
	scenCtx.Step(`^I process the signed JWT$`, func(message *godog.DocString) error {
		return session.ProcessSignedJWT(ctx, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I process the encrypted JWT$`, func(message *godog.DocString) error {
		return session.ProcessEncryptedJWT(ctx, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I process the signed encrypted JWT$`, func(message *godog.DocString) error {
		return session.ProcessSignedEncryptedJWT(ctx, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^the JWT must be valid$`, func() error {
		return session.ValidateJWT(ctx)
	})
	scenCtx.Step(`^the JWT must be invalid by "([^"]*)"$`, func(msg string) error {
		return session.ValidateInvalidJWT(ctx, golium.ValueAsString(ctx, msg))
	})
	scenCtx.Step(`^the JWT payload must have the JSON properties$`, func(t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the payload: %w", err)
		}
		return session.ValidatePayloadJSONProperties(ctx, props)
	})
	return ctx
}
