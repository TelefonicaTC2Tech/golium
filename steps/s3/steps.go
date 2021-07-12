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

package s3steps

import (
	"context"
	"fmt"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
)

// Steps to initialize common steps.
type Steps struct {
}

// InitializeSteps initializes all the steps to work with S3.
// It implements StepsInitializer interface.
func (cs Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the s3 session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)
	// Initialize the steps
	scenCtx.Step(`^I create a new S3 session$`, func() error {
		return session.newS3Session(ctx)
	})
	scenCtx.Step(`^I create a file in S3 bucket "([^"]+)" with key "([^"]+)" and the content$`, func(bucket, key string, message *godog.DocString) error {
		if session.Client == nil {
			return fmt.Errorf("failed creating S3 file: nil session: may forget step 'I create a new S3 session'")
		}
		return session.uploadS3FileWithContent(ctx, golium.ValueAsString(ctx, bucket), golium.ValueAsString(ctx, key), golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I create the S3 bucket "([^"]+)"$`, func(bucket string) error {
		if session.Client == nil {
			return fmt.Errorf("failed creating S3 file: nil session: may forget step 'I create a new S3 session'")
		}
		return session.createBucket(ctx, golium.ValueAsString(ctx, bucket))
	})
	scenCtx.Step(`^the file "([^"]+)" exists in S3 bucket "([^"]+)"$`, func(key, bucket string) error {
		if session.Client == nil {
			return fmt.Errorf("failed validating S3 file exists: nil session: may forget step 'I create a new S3 session'")
		}
		return session.validateS3FileExists(ctx, golium.ValueAsString(ctx, bucket), golium.ValueAsString(ctx, key))
	})
	scenCtx.Step(`^the file "([^"]+)" exists in S3 bucket "([^"]+)" with the content$`, func(key, bucket string, t *godog.DocString) error {
		if session.Client == nil {
			return fmt.Errorf("failed validating S3 file content: nil session: may forget step 'I create a new S3 session'")
		}
		return session.validateS3FileExistsWithContent(ctx, golium.ValueAsString(ctx, bucket), golium.ValueAsString(ctx, key), golium.ValueAsString(ctx, t.Content))
	})
	scenCtx.Step(`^I delete the file in S3 bucket "([^"]+)" with key "([^"]+)"$`, func(bucket, key string) error {
		if session.Client == nil {
			return fmt.Errorf("failed validating S3 file content: nil session: may forget step 'I create a new S3 session'")
		}
		return session.deleteS3File(ctx, golium.ValueAsString(ctx, bucket), golium.ValueAsString(ctx, key))
	})
	scenCtx.AfterScenario(func(sc *godog.Scenario, err error) {
		//	clean created documents
		cleanFiles := golium.Value(ctx, "[CONF:s3Autoclean]").(bool)
		if cleanFiles {
			session.CleanUp(ctx)
		}
	})
	return ctx
}
