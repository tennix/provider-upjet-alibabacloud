package ims

import (
	"github.com/crossplane-contrib/provider-upjet-alibabacloud/config/common"
	"github.com/crossplane/upjet/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("alicloud_ims_oidc_provider", func(r *config.Resource) {
		r.ShortGroup = string(common.IMS)
		r.Kind = "OIDCProvider"
	})
}
