package pvtz

import (
	"github.com/crossplane-contrib/provider-upjet-alibabacloud/config/common"
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("alicloud_pvtz_zone", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "pvtz"
		r.ShortGroup = string(common.PVTZ)
		r.Kind = "Zone"
		// name is deprecated
		delete(r.TerraformResource.Schema, "name")
	})
	p.AddResourceConfigurator("alicloud_pvtz_zone_attachment", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "pvtz"
		r.ShortGroup = string(common.PVTZ)
		r.Kind = "ZoneAttachment"
		r.References["zone_id"] = config.Reference{
			TerraformName: "alicloud_pvtz_zone",
		}
		r.References["vpc_id"] = config.Reference{
			TerraformName: "alicloud_vpc",
		}
		// vpc_ids is conflict with vpcs, and vpcs is recommended
		delete(r.TerraformResource.Schema, "vpc_ids")
	})
	p.AddResourceConfigurator("alicloud_pvtz_zone_record", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "pvtz"
		r.ShortGroup = string(common.PVTZ)
		r.Kind = "ZoneRecord"
		r.References["zone_id"] = config.Reference{
			TerraformName: "alicloud_pvtz_zone",
		}
		// resource_record is deprecated
		delete(r.TerraformResource.Schema, "resource_record")
	})
	p.AddResourceConfigurator("alicloud_pvtz_endpoint", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "pvtz"
		r.ShortGroup = string(common.PVTZ)
		r.Kind = "Endpoint"
		r.References["vpc_id"] = config.Reference{
			TerraformName: "alicloud_vpc",
		}
		r.References["security_group_id"] = config.Reference{
			TerraformName: "alicloud_security_group",
		}
	})
	p.AddResourceConfigurator("alicloud_pvtz_rule", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "pvtz"
		r.ShortGroup = string(common.PVTZ)
		r.Kind = "Rule"
		r.References["endpoint_id"] = config.Reference{
			TerraformName: "alicloud_pvtz_endpoint",
		}
	})
	p.AddResourceConfigurator("alicloud_pvtz_rule_attachment", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "pvtz"
		r.ShortGroup = string(common.PVTZ)
		r.Kind = "RuleAttachment"
		r.References["rule_id"] = config.Reference{
			TerraformName: "alicloud_pvtz_rule",
		}
		r.References["vpc_id"] = config.Reference{
			TerraformName: "alicloud_vpc",
		}
	})
	p.AddResourceConfigurator("alicloud_pvtz_user_vpc_authorization", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "pvtz"
		r.ShortGroup = string(common.PVTZ)
		r.Kind = "UserVPCAuthorization"
	})
}
