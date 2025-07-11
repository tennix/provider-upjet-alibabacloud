package natgateway

import (
	"github.com/crossplane-contrib/provider-upjet-alibabacloud/config/common"
	"github.com/crossplane/upjet/pkg/config"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("alicloud_nat_gateway", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "vpc"
		r.ShortGroup = string(common.NATGateway)
		r.Kind = "NATGateway"
		r.References["vpc_id"] = config.Reference{
			TerraformName: "alicloud_vpc",
		}
		// Delete deprecated fields
		delete(r.TerraformResource.Schema, "name")
		delete(r.TerraformResource.Schema, "instance_charge_type")
		delete(r.TerraformResource.Schema, "spec")
		delete(r.TerraformResource.Schema, "bandwidth_packages")
		delete(r.TerraformResource.Schema, "availability_zone")
	})
	p.AddResourceConfigurator("alicloud_forward_entry", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "vpc"
		r.ShortGroup = string(common.NATGateway)
		r.Kind = "ForwardEntry"
		r.References["vpc_id"] = config.Reference{
			TerraformName: "alicloud_vpc",
		}
		// Delete deprecated fields
		delete(r.TerraformResource.Schema, "name")
		delete(r.TerraformResource.Schema, "availability_zone")
	})
	p.AddResourceConfigurator("alicloud_snat_entry", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "vpc"
		r.ShortGroup = string(common.NATGateway)
		r.Kind = "SNATEntry"
		r.References["vswitch_id"] = config.Reference{
			TerraformName: "alicloud_vswitch",
		}
		// Delete deprecated fields
		delete(r.TerraformResource.Schema, "name")
		delete(r.TerraformResource.Schema, "availability_zone")
	})
	p.AddResourceConfigurator("alicloud_vpc_nat_ip", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "vpc"
		r.ShortGroup = string(common.NATGateway)
		r.Kind = "VPCNATIP"
		r.References["vpc_id"] = config.Reference{
			TerraformName: "alicloud_vpc",
		}
		// Delete deprecated fields
		delete(r.TerraformResource.Schema, "name")
		delete(r.TerraformResource.Schema, "instance_charge_type")
		delete(r.TerraformResource.Schema, "spec")
		delete(r.TerraformResource.Schema, "bandwidth_packages")
		delete(r.TerraformResource.Schema, "availability_zone")
	})
	p.AddResourceConfigurator("alicloud_vpc_nat_ip_cidr", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "vpc"
		r.ShortGroup = string(common.NATGateway)
		r.Kind = "VPCNATIPCIDR"
		r.References["vpc_id"] = config.Reference{
			TerraformName: "alicloud_vpc",
		}
		// Delete deprecated fields
		delete(r.TerraformResource.Schema, "name")
		delete(r.TerraformResource.Schema, "instance_charge_type")
		delete(r.TerraformResource.Schema, "spec")
		delete(r.TerraformResource.Schema, "bandwidth_packages")
		delete(r.TerraformResource.Schema, "availability_zone")
	})
}
