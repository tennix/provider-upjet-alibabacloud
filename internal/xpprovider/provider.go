/*
Copyright 2021 Upbound Inc.
*/

package xpprovider

import (
	"context"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetProvider returns the framework and SDK providers for Alibaba Cloud.
// The Alibaba Cloud provider currently only uses the SDK provider (not Framework).
func GetProvider(ctx context.Context) (provider.Provider, *schema.Provider, error) {
	// Initialize the Alibaba Cloud SDK provider
	sdkProvider := alicloud.Provider()

	// Alibaba Cloud provider doesn't use Terraform Plugin Framework yet
	var fwProvider provider.Provider = nil

	return fwProvider, sdkProvider, nil
}
