package v1beta1

import (
	"reflect"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestProviderConfigDeepCopyAuthOptions(t *testing.T) {
	const input = `
spec:
  credentials:
    source: None
  assumeRole:
    roleARN: acs:ram::123456789:role/test
    sessionName: role-session
    policy: role-policy
    sessionExpiration: 900
    externalID: external-id
  assumeRoleWithOIDC:
    roleARN: acs:ram::123456789:role/oidc
    oidcProviderARN: acs:ram::123456789:oidc-provider/test
    oidcTokenFile: /var/run/oidc/token
    roleSessionName: oidc-session
    policy: oidc-policy
    sessionExpiration: 3600
`
	var original, expected ProviderConfig
	if err := yaml.UnmarshalStrict([]byte(input), &original); err != nil {
		t.Fatal(err)
	}
	if err := yaml.UnmarshalStrict([]byte(input), &expected); err != nil {
		t.Fatal(err)
	}
	copied := original.DeepCopy()
	if !reflect.DeepEqual(&original, copied) {
		t.Fatal("deep copy lost auth options")
	}
	const changed = "changed"
	copied.Spec.AssumeRole.RoleARN = changed
	*copied.Spec.AssumeRole.SessionName = changed
	*copied.Spec.AssumeRole.Policy = changed
	*copied.Spec.AssumeRole.SessionExpiration = 1800
	*copied.Spec.AssumeRole.ExternalID = changed
	copied.Spec.AssumeRoleWithOIDC.RoleARN = changed
	*copied.Spec.AssumeRoleWithOIDC.RoleSessionName = changed
	*copied.Spec.AssumeRoleWithOIDC.Policy = changed
	*copied.Spec.AssumeRoleWithOIDC.SessionExpiration = 1800
	if !reflect.DeepEqual(original, expected) {
		t.Fatal("mutating copied auth options modified the original ProviderConfig")
	}
}
