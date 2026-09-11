// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"reflect"
	"strings"
	"testing"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/upjet/pkg/terraform"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/crossplane-contrib/provider-upjet-alibabacloud/apis"
	"github.com/crossplane-contrib/provider-upjet-alibabacloud/apis/v1beta1"
	vpcv1alpha1 "github.com/crossplane-contrib/provider-upjet-alibabacloud/apis/vpc/v1alpha1"
)

func TestBuildProviderConfiguration(t *testing.T) {
	cfg := buildProviderConfiguration("cn-hangzhou", map[string]any{
		"access_key":     "ak",
		"secret_key":     "sk",
		"security_token": "token",
		"ignored":        "value",
	}, nil)

	want := map[string]any{
		"region":         "cn-hangzhou",
		"access_key":     "ak",
		"secret_key":     "sk",
		"security_token": "token",
	}
	if !reflect.DeepEqual(map[string]any(cfg), want) {
		t.Fatalf("unexpected provider configuration:\n got: %#v\nwant: %#v", cfg, want)
	}
}

func TestBuildProviderConfigurationUsesTypedAssumeRole(t *testing.T) {
	roleARN := "acs:ram::1234567890123456:role/crossplane"
	sessionName := "crossplane"
	policy := "{}"
	sessionExpiration := 3600
	externalID := "external-id"

	cfg := buildProviderConfiguration("cn-hangzhou", map[string]any{
		"assume_role": map[string]any{
			"role_arn": "acs:ram::1234567890123456:role/from-secret",
		},
	}, &v1beta1.ProviderConfig{
		Spec: v1beta1.ProviderConfigSpec{
			AssumeRole: &v1beta1.AssumeRoleOptions{
				RoleARN:           roleARN,
				SessionName:       &sessionName,
				Policy:            &policy,
				SessionExpiration: &sessionExpiration,
				ExternalID:        &externalID,
			},
		},
	})

	want := []any{
		map[string]any{
			"role_arn":           roleARN,
			"session_name":       sessionName,
			"policy":             policy,
			"session_expiration": sessionExpiration,
			"external_id":        externalID,
		},
	}
	if !reflect.DeepEqual(cfg["assume_role"], want) {
		t.Fatalf("unexpected assume_role value:\n got: %#v\nwant: %#v", cfg["assume_role"], want)
	}
}

func TestBuildProviderConfigurationUsesAssumeRoleWithOIDC(t *testing.T) {
	roleARN := "acs:ram::1234567890123456:role/crossplane-oidc"
	oidcProviderARN := "acs:ram::1234567890123456:oidc-provider/ack"
	oidcTokenFile := "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token"
	roleSessionName := "crossplane"
	sessionExpiration := 3600

	cfg := buildProviderConfiguration("cn-hangzhou", map[string]any{}, &v1beta1.ProviderConfig{
		Spec: v1beta1.ProviderConfigSpec{
			AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
				RoleARN:           roleARN,
				OIDCProviderARN:   oidcProviderARN,
				OIDCTokenFile:     oidcTokenFile,
				RoleSessionName:   &roleSessionName,
				SessionExpiration: &sessionExpiration,
			},
		},
	})

	wantOIDC := []any{
		map[string]any{
			"role_arn":           roleARN,
			"oidc_provider_arn":  oidcProviderARN,
			"oidc_token_file":    oidcTokenFile,
			"role_session_name":  roleSessionName,
			"session_expiration": sessionExpiration,
		},
	}
	if !reflect.DeepEqual(cfg["assume_role_with_oidc"], wantOIDC) {
		t.Fatalf("unexpected assume_role_with_oidc value:\n got: %#v\nwant: %#v", cfg["assume_role_with_oidc"], wantOIDC)
	}
}

func TestBuildProviderConfigurationOmitsUnsetOptionalAuthFields(t *testing.T) {
	cfg := buildProviderConfiguration("cn-hangzhou", map[string]any{}, &v1beta1.ProviderConfig{
		Spec: v1beta1.ProviderConfigSpec{
			AssumeRole: &v1beta1.AssumeRoleOptions{
				RoleARN: "acs:ram::1234567890123456:role/crossplane",
			},
			AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
				RoleARN:         "acs:ram::1234567890123456:role/crossplane-oidc",
				OIDCProviderARN: "acs:ram::1234567890123456:oidc-provider/ack",
				OIDCTokenFile:   "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
			},
		},
	})

	wantAssumeRole := []any{
		map[string]any{
			"role_arn": "acs:ram::1234567890123456:role/crossplane",
		},
	}
	wantOIDC := []any{
		map[string]any{
			"role_arn":          "acs:ram::1234567890123456:role/crossplane-oidc",
			"oidc_provider_arn": "acs:ram::1234567890123456:oidc-provider/ack",
			"oidc_token_file":   "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
		},
	}
	if !reflect.DeepEqual(cfg["assume_role"], wantAssumeRole) {
		t.Fatalf("unexpected assume_role:\n got: %#v\nwant: %#v", cfg["assume_role"], wantAssumeRole)
	}
	if !reflect.DeepEqual(cfg["assume_role_with_oidc"], wantOIDC) {
		t.Fatalf("unexpected assume_role_with_oidc:\n got: %#v\nwant: %#v", cfg["assume_role_with_oidc"], wantOIDC)
	}
}

func TestValidateProviderConfigRejectsDynamicAuthInSecret(t *testing.T) {
	err := validateProviderConfig(&v1beta1.ProviderConfig{}, map[string]any{
		"assume_role": map[string]any{
			"role_arn": "acs:ram::1234567890123456:role/from-secret",
		},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateProviderConfigRequiresNonBlankAuthFields(t *testing.T) {
	tests := map[string]struct {
		pc      *v1beta1.ProviderConfig
		wantErr string
	}{
		"AssumeRoleRoleARNWhitespace": {
			pc: &v1beta1.ProviderConfig{
				Spec: v1beta1.ProviderConfigSpec{
					AssumeRole: &v1beta1.AssumeRoleOptions{
						RoleARN: " \t\n",
					},
				},
			},
			wantErr: "spec.assumeRole.roleARN is required",
		},
		"OIDCRoleARNWhitespace": {
			pc: &v1beta1.ProviderConfig{
				Spec: v1beta1.ProviderConfigSpec{
					AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
						RoleARN:         " ",
						OIDCProviderARN: "acs:ram::1234567890123456:oidc-provider/ack",
						OIDCTokenFile:   "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
					},
				},
			},
			wantErr: "spec.assumeRoleWithOIDC.roleARN is required",
		},
		"OIDCProviderARNWhitespace": {
			pc: &v1beta1.ProviderConfig{
				Spec: v1beta1.ProviderConfigSpec{
					AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
						RoleARN:         "acs:ram::1234567890123456:role/crossplane-oidc",
						OIDCProviderARN: "\t",
						OIDCTokenFile:   "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
					},
				},
			},
			wantErr: "spec.assumeRoleWithOIDC.oidcProviderARN is required",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateProviderConfig(tt.pc, nil)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: got %v, want substring %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidateProviderConfigRequiresOIDCTokenFile(t *testing.T) {
	roleARN := "acs:ram::1234567890123456:role/crossplane-oidc"
	oidcProviderARN := "acs:ram::1234567890123456:oidc-provider/ack"

	err := validateProviderConfig(&v1beta1.ProviderConfig{
		Spec: v1beta1.ProviderConfigSpec{
			AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
				RoleARN:         roleARN,
				OIDCProviderARN: oidcProviderARN,
			},
		},
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateProviderConfigSessionExpirationBounds(t *testing.T) {
	tests := map[string]struct {
		pc      *v1beta1.ProviderConfig
		wantOK  bool
		wantErr string
	}{
		"AssumeRoleNil": {
			pc:     providerConfigWithAssumeRoleExpiration(nil),
			wantOK: true,
		},
		"AssumeRoleMinimumInclusive": {
			pc:     providerConfigWithAssumeRoleExpiration(intPtr(900)),
			wantOK: true,
		},
		"AssumeRoleMaximumInclusive": {
			pc:     providerConfigWithAssumeRoleExpiration(intPtr(3600)),
			wantOK: true,
		},
		"AssumeRoleBelowMinimum": {
			pc:      providerConfigWithAssumeRoleExpiration(intPtr(899)),
			wantErr: "spec.assumeRole.sessionExpiration must be between 900 and 3600",
		},
		"AssumeRoleAboveMaximum": {
			pc:      providerConfigWithAssumeRoleExpiration(intPtr(3601)),
			wantErr: "spec.assumeRole.sessionExpiration must be between 900 and 3600",
		},
		"OIDCNil": {
			pc:     providerConfigWithOIDCExpiration(nil),
			wantOK: true,
		},
		"OIDCMinimumInclusive": {
			pc:     providerConfigWithOIDCExpiration(intPtr(900)),
			wantOK: true,
		},
		"OIDCMaximumInclusive": {
			pc:     providerConfigWithOIDCExpiration(intPtr(43200)),
			wantOK: true,
		},
		"OIDCBelowMinimum": {
			pc:      providerConfigWithOIDCExpiration(intPtr(899)),
			wantErr: "spec.assumeRoleWithOIDC.sessionExpiration must be between 900 and 43200",
		},
		"OIDCAboveMaximum": {
			pc:      providerConfigWithOIDCExpiration(intPtr(43201)),
			wantErr: "spec.assumeRoleWithOIDC.sessionExpiration must be between 900 and 43200",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateProviderConfig(tt.pc, nil)
			if tt.wantOK {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: got %v, want substring %q", err, tt.wantErr)
			}
		})
	}
}

func TestTerraformSetupBuilderUsesOIDCWithoutCredentialSecret(t *testing.T) {
	ctx := context.Background()
	region := "cn-shanghai"
	roleSessionName := "crossplane-oidc"
	sessionExpiration := 3600
	pc := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "oidc"},
		Spec: v1beta1.ProviderConfigSpec{
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceNone,
			},
			AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
				RoleARN:           "acs:ram::1234567890123456:role/crossplane-oidc",
				OIDCProviderARN:   "acs:ram::1234567890123456:oidc-provider/ack",
				OIDCTokenFile:     "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
				RoleSessionName:   &roleSessionName,
				SessionExpiration: &sessionExpiration,
			},
		},
	}
	mg := managedVPC("example", "oidc", &region)

	c := fakeClient(t, pc, mg)
	setup, err := TerraformSetupBuilder("v1.0.0", "registry.example/alicloud", "1.260.0", nil)(ctx, c, mg)
	if err != nil {
		t.Fatalf("unexpected setup error: %v", err)
	}

	wantOIDC := []any{
		map[string]any{
			"role_arn":           "acs:ram::1234567890123456:role/crossplane-oidc",
			"oidc_provider_arn":  "acs:ram::1234567890123456:oidc-provider/ack",
			"oidc_token_file":    "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
			"role_session_name":  roleSessionName,
			"session_expiration": sessionExpiration,
		},
	}
	if got := setup.Configuration["region"]; got != region {
		t.Fatalf("unexpected region: got %v, want %v", got, region)
	}
	if !reflect.DeepEqual(setup.Configuration["assume_role_with_oidc"], wantOIDC) {
		t.Fatalf("unexpected assume_role_with_oidc:\n got: %#v\nwant: %#v", setup.Configuration["assume_role_with_oidc"], wantOIDC)
	}
	if _, ok := setup.Configuration["access_key"]; ok {
		t.Fatal("did not expect access_key with credentials.source None")
	}
	if _, ok := setup.Configuration["secret_key"]; ok {
		t.Fatal("did not expect secret_key with credentials.source None")
	}
	if setup.Scheduler != nil {
		t.Fatal("unexpected scheduler")
	}
	assertProviderConfigUsage(t, c, mg, "oidc")
}

func TestTerraformSetupBuilderUsesInjectedIdentityWithOIDCWithoutCredentialSecret(t *testing.T) {
	ctx := context.Background()
	region := "cn-shanghai"
	roleSessionName := "crossplane-oidc"
	sessionExpiration := 3600
	pc := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "oidc"},
		Spec: v1beta1.ProviderConfigSpec{
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceInjectedIdentity,
			},
			AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
				RoleARN:           "acs:ram::1234567890123456:role/crossplane-oidc",
				OIDCProviderARN:   "acs:ram::1234567890123456:oidc-provider/ack",
				OIDCTokenFile:     "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
				RoleSessionName:   &roleSessionName,
				SessionExpiration: &sessionExpiration,
			},
		},
	}
	mg := managedVPC("example", "oidc", &region)
	c := fakeClient(t, pc, mg)

	setup, err := TerraformSetupBuilder("v1.0.0", "registry.example/alicloud", "1.260.0", nil)(ctx, c, mg)
	if err != nil {
		t.Fatalf("unexpected setup error: %v", err)
	}
	if _, ok := setup.Configuration["access_key"]; ok {
		t.Fatal("did not expect access_key with credentials.source InjectedIdentity")
	}
	if _, ok := setup.Configuration["secret_key"]; ok {
		t.Fatal("did not expect secret_key with credentials.source InjectedIdentity")
	}
	if _, ok := setup.Configuration["assume_role_with_oidc"]; !ok {
		t.Fatalf("expected assume_role_with_oidc in provider configuration: %#v", setup.Configuration)
	}
	assertProviderConfigUsage(t, c, mg, "oidc")
}

func TestTerraformSetupBuilderUsesSecretCredentialsAndRegionFallback(t *testing.T) {
	ctx := context.Background()
	pc := providerConfigWithSecret("default")
	secret := credentialSecret("example-creds", `{
		"access_key": "ak",
		"secret_key": "sk",
		"security_token": "token",
		"region_id": "cn-hangzhou"
	}`)
	mg := managedVPC("example", "default", nil)
	scheduler := terraform.NewNoOpProviderScheduler()
	c := fakeClient(t, pc, secret, mg)

	setup, err := TerraformSetupBuilder("v1.0.0", "registry.example/alicloud", "1.260.0", scheduler)(ctx, c, mg)
	if err != nil {
		t.Fatalf("unexpected setup error: %v", err)
	}

	want := map[string]any{
		"region":               "cn-hangzhou",
		"access_key":           "ak",
		"secret_key":           "sk",
		"security_token":       "token",
		"configuration_source": getUserAgent(),
	}
	if !reflect.DeepEqual(map[string]any(setup.Configuration), want) {
		t.Fatalf("unexpected setup configuration:\n got: %#v\nwant: %#v", setup.Configuration, want)
	}
	if setup.Scheduler != scheduler {
		t.Fatal("scheduler was not preserved")
	}
	assertProviderConfigUsage(t, c, mg, "default")
}

func TestTerraformSetupBuilderUsesManagedResourceRegionOverSecretRegion(t *testing.T) {
	ctx := context.Background()
	region := "cn-shanghai"
	pc := providerConfigWithSecret("default")
	secret := credentialSecret("example-creds", `{
		"access_key": "ak",
		"secret_key": "sk",
		"region": "cn-hangzhou"
	}`)
	mg := managedVPC("example", "default", &region)

	setup, err := TerraformSetupBuilder("v1.0.0", "registry.example/alicloud", "1.260.0", nil)(ctx, fakeClient(t, pc, secret, mg), mg)
	if err != nil {
		t.Fatalf("unexpected setup error: %v", err)
	}
	if got := setup.Configuration["region"]; got != region {
		t.Fatalf("unexpected region: got %v, want %v", got, region)
	}
}

func TestTerraformSetupBuilderOIDCProviderHandlesAreStableAndDistinct(t *testing.T) {
	policyA := `{"Statement":[{"Action":"vpc:*","Effect":"Allow","Resource":"*"}]}`
	policyB := `{"Statement":[{"Action":"ecs:*","Effect":"Allow","Resource":"*"}]}`

	base := setupForOIDCProviderConfig(t,
		"base",
		"acs:ram::1234567890123456:role/crossplane-oidc",
		"/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
		&policyA)
	same := setupForOIDCProviderConfig(t,
		"same",
		"acs:ram::1234567890123456:role/crossplane-oidc",
		"/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
		&policyA)
	differentRole := setupForOIDCProviderConfig(t,
		"different-role",
		"acs:ram::1234567890123456:role/crossplane-oidc-2",
		"/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
		&policyA)
	differentTokenFile := setupForOIDCProviderConfig(t,
		"different-token-file",
		"acs:ram::1234567890123456:role/crossplane-oidc",
		"/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/other-token",
		&policyA)
	differentPolicy := setupForOIDCProviderConfig(t,
		"different-policy",
		"acs:ram::1234567890123456:role/crossplane-oidc",
		"/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
		&policyB)

	baseHandle, err := base.Configuration.ToProviderHandle()
	if err != nil {
		t.Fatalf("cannot build base provider handle: %v", err)
	}
	baseHandleAgain, err := base.Configuration.ToProviderHandle()
	if err != nil {
		t.Fatalf("cannot rebuild base provider handle: %v", err)
	}
	if baseHandle != baseHandleAgain {
		t.Fatalf("provider handle is not deterministic: got %q then %q", baseHandle, baseHandleAgain)
	}
	sameHandle, err := same.Configuration.ToProviderHandle()
	if err != nil {
		t.Fatalf("cannot build same provider handle: %v", err)
	}
	if baseHandle != sameHandle {
		t.Fatalf("identical OIDC configurations produced different handles: got %q want %q", sameHandle, baseHandle)
	}
	for name, setup := range map[string]terraform.Setup{
		"different role":       differentRole,
		"different token file": differentTokenFile,
		"different policy":     differentPolicy,
	} {
		handle, err := setup.Configuration.ToProviderHandle()
		if err != nil {
			t.Fatalf("cannot build provider handle for %s: %v", name, err)
		}
		if handle == baseHandle {
			t.Fatalf("%s produced the same provider handle %q", name, handle)
		}
	}
}

func TestTerraformSetupBuilderRejectsInvalidCredentialJSON(t *testing.T) {
	ctx := context.Background()
	pc := providerConfigWithSecret("default")
	secret := credentialSecret("example-creds", `{`)
	mg := managedVPC("example", "default", nil)

	_, err := TerraformSetupBuilder("v1.0.0", "registry.example/alicloud", "1.260.0", nil)(ctx, fakeClient(t, pc, secret, mg), mg)
	if err == nil {
		t.Fatal("expected setup error")
	}
	if strings.Contains(err.Error(), "ak") || strings.Contains(err.Error(), "sk") {
		t.Fatalf("error leaked credential material: %v", err)
	}
}

func TestTerraformSetupBuilderReturnsExpectedSetupErrors(t *testing.T) {
	tests := map[string]struct {
		mg              *vpcv1alpha1.VPC
		objs            []client.Object
		wantErr         string
		wantUsage       bool
		providerRefName string
	}{
		"NoProviderConfigReference": {
			mg:      managedVPCWithoutProviderConfig("example"),
			wantErr: errNoProviderConfig,
		},
		"MissingProviderConfig": {
			mg:              managedVPC("example", "missing", nil),
			wantErr:         errGetProviderConfig,
			wantUsage:       true,
			providerRefName: "missing",
		},
		"MissingSecret": {
			mg:              managedVPC("example", "default", nil),
			objs:            []client.Object{providerConfigWithSecret("default")},
			wantErr:         errExtractCredentials,
			wantUsage:       true,
			providerRefName: "default",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			objs := tt.objs
			if len(objs) == 0 && tt.mg != nil {
				objs = []client.Object{tt.mg}
			}
			c := fakeClient(t, objs...)
			_, err := TerraformSetupBuilder("v1.0.0", "registry.example/alicloud", "1.260.0", nil)(context.Background(), c, tt.mg)
			if err == nil {
				t.Fatal("expected setup error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("unexpected error: got %v, want substring %q", err, tt.wantErr)
			}
			if tt.wantUsage {
				assertProviderConfigUsage(t, c, tt.mg, tt.providerRefName)
			}
		})
	}
}

func TestTerraformSetupBuilderRejectsRoleSettingsInSecret(t *testing.T) {
	ctx := context.Background()
	pc := providerConfigWithSecret("default")
	secret := credentialSecret("example-creds", `{
		"access_key": "ak",
		"secret_key": "sk",
		"assume_role": {"role_arn": "acs:ram::1234567890123456:role/from-secret"}
	}`)
	mg := managedVPC("example", "default", nil)

	_, err := TerraformSetupBuilder("v1.0.0", "registry.example/alicloud", "1.260.0", nil)(ctx, fakeClient(t, pc, secret, mg), mg)
	if err == nil {
		t.Fatal("expected setup error")
	}
	if !strings.Contains(err.Error(), "spec.assumeRole") {
		t.Fatalf("expected spec.assumeRole guidance, got: %v", err)
	}
}

func fakeClient(t *testing.T, objs ...client.Object) client.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("cannot add core scheme: %v", err)
	}
	if err := apis.AddToScheme(scheme); err != nil {
		t.Fatalf("cannot add provider scheme: %v", err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func providerConfigWithSecret(name string) *v1beta1.ProviderConfig {
	return &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: v1beta1.ProviderConfigSpec{
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv1.CommonCredentialSelectors{
					SecretRef: &xpv1.SecretKeySelector{
						SecretReference: xpv1.SecretReference{
							Name:      "example-creds",
							Namespace: "crossplane-system",
						},
						Key: "credentials",
					},
				},
			},
		},
	}
}

func providerConfigWithAssumeRoleExpiration(expiration *int) *v1beta1.ProviderConfig {
	return &v1beta1.ProviderConfig{
		Spec: v1beta1.ProviderConfigSpec{
			AssumeRole: &v1beta1.AssumeRoleOptions{
				RoleARN:           "acs:ram::1234567890123456:role/crossplane",
				SessionExpiration: expiration,
			},
		},
	}
}

func providerConfigWithOIDCExpiration(expiration *int) *v1beta1.ProviderConfig {
	return &v1beta1.ProviderConfig{
		Spec: v1beta1.ProviderConfigSpec{
			AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
				RoleARN:           "acs:ram::1234567890123456:role/crossplane-oidc",
				OIDCProviderARN:   "acs:ram::1234567890123456:oidc-provider/ack",
				OIDCTokenFile:     "/var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token",
				SessionExpiration: expiration,
			},
		},
	}
}

func intPtr(v int) *int {
	return &v
}

func credentialSecret(name, credentials string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "crossplane-system",
		},
		Data: map[string][]byte{
			"credentials": []byte(credentials),
		},
	}
}

func oidcProviderConfig(name, roleARN, oidcTokenFile string, policy *string) *v1beta1.ProviderConfig {
	roleSessionName := "crossplane-oidc"
	sessionExpiration := 3600
	return &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: v1beta1.ProviderConfigSpec{
			Credentials: v1beta1.ProviderCredentials{
				Source: xpv1.CredentialsSourceNone,
			},
			AssumeRoleWithOIDC: &v1beta1.AssumeRoleWithOIDCOptions{
				RoleARN:           roleARN,
				OIDCProviderARN:   "acs:ram::1234567890123456:oidc-provider/ack",
				OIDCTokenFile:     oidcTokenFile,
				RoleSessionName:   &roleSessionName,
				Policy:            policy,
				SessionExpiration: &sessionExpiration,
			},
		},
	}
}

func setupForOIDCProviderConfig(t *testing.T, name, roleARN, oidcTokenFile string, policy *string) terraform.Setup {
	t.Helper()
	region := "cn-shanghai"
	mg := managedVPC(name, name, &region)
	c := fakeClient(t, oidcProviderConfig(name, roleARN, oidcTokenFile, policy), mg)

	setup, err := TerraformSetupBuilder("v1.0.0", "registry.example/alicloud", "1.260.0", nil)(context.Background(), c, mg)
	if err != nil {
		t.Fatalf("unexpected setup error for %s: %v", name, err)
	}
	assertProviderConfigUsage(t, c, mg, name)
	return setup
}

func managedVPC(name, providerConfigName string, region *string) *vpcv1alpha1.VPC {
	return &vpcv1alpha1.VPC{
		TypeMeta: metav1.TypeMeta{
			APIVersion: vpcv1alpha1.VPC_GroupVersionKind.GroupVersion().String(),
			Kind:       vpcv1alpha1.VPC_Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			UID:  types.UID(name + "-uid"),
		},
		Spec: vpcv1alpha1.VPCSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{Name: providerConfigName},
			},
			ForProvider: vpcv1alpha1.VPCParameters{
				Region: region,
			},
		},
	}
}

func managedVPCWithoutProviderConfig(name string) *vpcv1alpha1.VPC {
	return &vpcv1alpha1.VPC{
		TypeMeta: metav1.TypeMeta{
			APIVersion: vpcv1alpha1.VPC_GroupVersionKind.GroupVersion().String(),
			Kind:       vpcv1alpha1.VPC_Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			UID:  types.UID(name + "-uid"),
		},
	}
}

func assertProviderConfigUsage(t *testing.T, c client.Client, mg *vpcv1alpha1.VPC, providerConfigName string) {
	t.Helper()
	usage := &v1beta1.ProviderConfigUsage{}
	if err := c.Get(context.Background(), types.NamespacedName{Name: string(mg.GetUID())}, usage); err != nil {
		t.Fatalf("expected ProviderConfigUsage to be created: %v", err)
	}
	if got := usage.ProviderConfigReference.Name; got != providerConfigName {
		t.Fatalf("unexpected ProviderConfigUsage provider ref: got %q want %q", got, providerConfigName)
	}
	if got := usage.ResourceReference.Name; got != mg.GetName() {
		t.Fatalf("unexpected ProviderConfigUsage resource ref: got %q want %q", got, mg.GetName())
	}
}
