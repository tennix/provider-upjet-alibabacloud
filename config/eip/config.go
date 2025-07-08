package eip

import (
	"github.com/crossplane-contrib/provider-upjet-alibabacloud/config/common"
	"github.com/crossplane/upjet/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("alicloud_eip", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "vpc"
		r.ShortGroup = string(common.EIP)
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{
				"name",
			},
		}
	})
	p.AddResourceConfigurator("alicloud_eip_address", func(r *config.Resource) {
		r.Kind = "EIPAddress"
		r.ShortGroup = string(common.EIP)
		delete(r.TerraformResource.Schema, "name")
		delete(r.TerraformResource.Schema, "instance_charge_type")
	})
	p.AddResourceConfigurator("alicloud_eip_association", func(r *config.Resource) {
		r.Kind = "EIPAssociation"
		r.ShortGroup = string(common.EIP)
		delete(r.TerraformResource.Schema, "name")
	})
}
