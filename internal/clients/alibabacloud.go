/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/crossplane-contrib/provider-upjet-alibabacloud/internal/version"
	"github.com/crossplane-contrib/provider-upjet-alibabacloud/internal/xpprovider"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane-contrib/provider-upjet-alibabacloud/apis/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal alicloud credentials as JSON"
)

// SetupConfig holds configuration for the Terraform provider setup.
type SetupConfig struct {
	TerraformProvider *schema.Provider
	Logger            logging.Logger
}

// SelectTerraformSetup returns a Terraform setup function configured for upjet v2.
func SelectTerraformSetup(config *SetupConfig) terraform.SetupFn {
	return func(ctx context.Context, c client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(c, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		// Extract credentials (existing logic)
		creds, err := extractAndUnmarshalCredentials(ctx, c, configRef)
		if err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// Get region (existing logic)
		region, err := getRegion(mg, creds)
		if err != nil {
			return ps, errors.Wrap(err, "cannot get region")
		}

		// Configure the embedded provider
		if config.TerraformProvider == nil {
			return ps, errors.New("terraform provider cannot be nil")
		}

		return ps, errors.Wrap(configureNoForkAlibabaCloudClient(ctx, &ps, config, region, creds),
			"could not configure the no-fork Alibaba Cloud client")
	}
}

// configureNoForkAlibabaCloudClient configures the embedded Terraform provider with credentials.
func configureNoForkAlibabaCloudClient(ctx context.Context, ps *terraform.Setup, config *SetupConfig, region string, creds map[string]string) error {
	// Configure provider with credentials
	tfConfig := xpprovider.AlibabaCloudConfig{
		AccessKey:           creds["access_key"],
		SecretKey:           creds["secret_key"],
		SecurityToken:       creds["security_token"],
		Region:              region,
		ConfigurationSource: getUserAgent(),
	}

	// Get the provider client/meta
	tfClient, diags := tfConfig.GetClient(ctx, config.TerraformProvider)
	if diags.HasError() {
		return errors.Errorf("cannot construct TF Alibaba Cloud Client from TF Config, %v", diags)
	}

	ps.Meta = tfClient
	// Note: FrameworkProvider is not needed as alicloud only uses SDK
	// ps.FrameworkProvider = nil

	return nil
}

// TerraformSetupBuilder is deprecated. Use SelectTerraformSetup instead.
// Kept for backward compatibility during migration.
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, c client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(c, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		creds, err := extractAndUnmarshalCredentials(ctx, c, configRef)
		if err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		region, err := getRegion(mg, creds)
		if err != nil {
			return ps, errors.Wrap(err, "cannot get region")
		}

		// Set credentials in Terraform provider configuration.
		ps.Configuration = map[string]any{
			"region": region,
		}
		if v, ok := creds["access_key"]; ok {
			ps.Configuration["access_key"] = v
		}
		if v, ok := creds["secret_key"]; ok {
			ps.Configuration["secret_key"] = v
		}
		if v, ok := creds["security_token"]; ok {
			ps.Configuration["security_token"] = v
		}
		ps.Configuration["configuration_source"] = getUserAgent()
		return ps, nil
	}
}

func getRegion(obj runtime.Object, creds map[string]string) (string, error) {
	fromMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return "", errors.Wrap(err, "cannot convert to unstructured")
	}
	credsRegion := creds["region"]
	if credsRegion == "" {
		// region_id is used as a fallback for old version
		credsRegion = creds["region_id"]
	}
	r, err := fieldpath.Pave(fromMap).GetString("spec.forProvider.region")
	if fieldpath.IsNotFound(err) {
		// Region is not required for all resources, e.g. resource in "ram" group.
		return credsRegion, nil
	}
	return r, err
}
func extractAndUnmarshalCredentials(ctx context.Context, c client.Client, configRef *v1.Reference) (map[string]string, error) {
	pc := &v1beta1.ProviderConfig{}
	creds := map[string]string{}
	if err := c.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
		return creds, errors.Wrap(err, errGetProviderConfig)
	}

	data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, c, pc.Spec.Credentials.CommonCredentialSelectors)
	if err != nil {
		return creds, errors.Wrap(err, errExtractCredentials)
	}
	if err = json.Unmarshal(data, &creds); err != nil {
		return creds, errors.Wrap(err, errUnmarshalCredentials)
	}
	return creds, nil
}
func getUserAgent() string {
	// user agent formats as "crossplane/<CROSSPLANE_VERSION> <PROJECT_NAME>/<PROJECT_VERSION>"
	return fmt.Sprintf("crossplane/%s provider-upjet-alibabacloud/%s", version.CrossplaneVersion, version.ProviderVersion)
}
