// SPDX-FileCopyrightText: 2025 Upbound Inc. <https://upbound.io>
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/pkg/reference"
	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
)

const (
	// SelfPackagePath is the golang path for this package.
	SelfPackagePath = "github.com/crossplane-contrib/provider-upjet-alibabacloud/config/common"

	// PathIdExtractor is the golang path to
	// IdExtractor function in this package.
	PathIdExtractor                              = SelfPackagePath + ".IdExtractor()"
	PathAccountNameExtractor                     = SelfPackagePath + ".AccountNameExtractor()"
	PathRoleArnExtractor                         = SelfPackagePath + ".RoleArnExtractor()"
	PathDBEndpointIdExtractor                    = SelfPackagePath + ".DBEndpointIdExtractor()"
	PathAlarmContactGroupNameExtractor           = SelfPackagePath + ".AlarmContactGroupNameExtractor()"
	PathOssBucketCnameTokenExtractor             = SelfPackagePath + ".OssBucketCnameTokenExtractor()"
	PathAlidnsRecordDomainExtractor              = SelfPackagePath + ".AlidnsRecordDomainExtractor()"
	PathOssBucketLocationExtractor               = SelfPackagePath + ".OssBucketLocationExtractor()"
	PathPrivateLinkVpcEndpointServiceIdExtractor = SelfPackagePath + ".PrivateLinkVpcEndpointServiceIdExtractor()"
	PathFcv3FunctionArnExtractor                 = SelfPackagePath + ".Fcv3FunctionArnExtractor()"
)

// IdExtractor extracts id of the
// resources from "status.atProvider.id".
func IdExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.id")
		if err != nil {
			return ""
		}
		return r
	}
}

// AccountNameExtractor extracts id of the
// resources from "status.atProvider.account_name".
func AccountNameExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.accountName")
		if err != nil {
			return ""
		}
		return r
	}
}

// RoleArnExtractor extracts id of the
// resources from "status.atProvider.account_name".
func RoleArnExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.arn")
		if err != nil {
			return ""
		}
		return r
	}
}

// DBEndpointIdExtractor extracts id of the
// resources from "status.atProvider.account_name".
func DBEndpointIdExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.dbEndpointId")
		if err != nil {
			return ""
		}
		return r
	}
}

// AlarmContactGroupNameExtractor extracts id of the
// resources from "status.atProvider.account_name".
func AlarmContactGroupNameExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.alarmContactGroupName")
		if err != nil {
			return ""
		}
		return r
	}
}

// OssBucketCnameTokenExtractor extracts id of the
// resources from "status.atProvider.token".
func OssBucketCnameTokenExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.token")
		if err != nil {
			return ""
		}
		return r
	}
}

// AlidnsRecordDomainExtractor extracts id of the
// resources from "status.atProvider.domain".
func AlidnsRecordDomainExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.domain")
		if err != nil {
			return ""
		}
		return r
	}
}

// OssBucketLocationExtractor extracts id of the
// resources from "status.atProvider.location".
func OssBucketLocationExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.location")
		if err != nil {
			return ""
		}
		return r
	}
}

// PrivateLinkVpcEndpointServiceIdExtractor extracts id of the
// resources from "status.atProvider.service_id".
func PrivateLinkVpcEndpointServiceIdExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.service_id")
		if err != nil {
			return ""
		}
		return r
	}
}

// Fcv3FunctionArnExtractor extracts id of the
// resources from "status.atProvider.function_arn".
func Fcv3FunctionArnExtractor() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		paved, err := fieldpath.PaveObject(mg)
		if err != nil {
			return ""
		}
		r, err := paved.GetString("status.atProvider.function_arn")
		if err != nil {
			return ""
		}
		return r
	}
}
