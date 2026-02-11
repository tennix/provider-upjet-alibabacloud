/*
Copyright 2021 Upbound Inc.
*/

package xpprovider

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// AlibabaCloudConfig holds provider configuration for Alibaba Cloud.
type AlibabaCloudConfig struct {
	AccessKey           string
	SecretKey           string
	SecurityToken       string
	Region              string
	ConfigurationSource string
}

// GetClient configures and returns the provider client.
func (c *AlibabaCloudConfig) GetClient(ctx context.Context, p *schema.Provider) (interface{}, diag.Diagnostics) {
	// Build the configuration map
	config := map[string]interface{}{
		"region": c.Region,
	}

	if c.AccessKey != "" {
		config["access_key"] = c.AccessKey
	}
	if c.SecretKey != "" {
		config["secret_key"] = c.SecretKey
	}
	if c.SecurityToken != "" {
		config["security_token"] = c.SecurityToken
	}
	if c.ConfigurationSource != "" {
		config["configuration_source"] = c.ConfigurationSource
	}

	// Convert config map to cty.Value
	configValues := make(map[string]cty.Value)
	for k, v := range config {
		if strVal, ok := v.(string); ok {
			configValues[k] = cty.StringVal(strVal)
		}
	}

	// Create ResourceConfig from the config map
	resourceConfig := terraform.NewResourceConfigRaw(config)

	// Configure the provider
	diags := p.Configure(ctx, resourceConfig)
	if diags.HasError() {
		return nil, diags
	}

	// Return the configured provider meta
	return p.Meta(), diags
}
