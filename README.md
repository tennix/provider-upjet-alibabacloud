# provider-upjet-alibabacloud

## Build

Run the following to build the provider locally.

```
make submodules
make generate
```

## Authentication

Provider credentials are configured with a `ProviderConfig`. Static credentials
and STS session credentials are read from a Kubernetes `Secret`; AssumeRole and
AssumeRoleWithOIDC options are configured as structured fields on the
`ProviderConfig` spec.

Do not put `assume_role` or `assume_role_with_oidc` blocks in the credentials
secret. The credentials secret should contain only credential material such as
`access_key`, `secret_key`, `security_token`, and `region`.

### Static Credentials

The default static AK/SK example uses `credentials.source: Secret`:

```bash
kubectl apply -f examples/providerconfig/v1beta1/secret.yaml.tmpl
kubectl apply -f examples/providerconfig/v1beta1/providerconfig.yaml
```

The credentials secret uses JSON in the `credentials` key:

```json
{
  "access_key": "...",
  "secret_key": "...",
  "region": "cn-hangzhou"
}
```

`region_id` is still accepted as a compatibility fallback, but new examples
should use `region`. Regional managed resources can also set
`spec.forProvider.region`; that resource-level value takes precedence over the
secret region.

### STS Session Credentials

STS session credentials are supported by adding `security_token` to the secret:

```bash
kubectl apply -f examples/providerconfig/v1beta1/secret-sts-token.yaml.tmpl
kubectl apply -f examples/providerconfig/v1beta1/providerconfig-sts-token.yaml
```

### AssumeRole

For AssumeRole, use a base credential secret for the caller identity and put
the role parameters in `spec.assumeRole`:

```bash
kubectl apply -f examples/providerconfig/v1beta1/secret.yaml.tmpl
kubectl apply -f examples/providerconfig/v1beta1/providerconfig-assume-role.yaml
```

```yaml
apiVersion: alibabacloud.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: assume-role
spec:
  credentials:
    source: Secret
    secretRef:
      name: example-creds
      namespace: crossplane-system
      key: credentials
  assumeRole:
    roleARN: acs:ram::<account-id>:role/<role-name>
    sessionName: crossplane-assume-role
    sessionExpiration: 3600
```

### AssumeRoleWithOIDC

For AssumeRoleWithOIDC, configure the role, OIDC provider, and token file in
`spec.assumeRoleWithOIDC`. No static AK/SK Secret is required when
`credentials.source` is `None`:

```bash
kubectl apply -f examples/providerconfig/v1beta1/providerconfig-assume-role-with-oidc.yaml
```

```yaml
apiVersion: alibabacloud.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: assume-role-with-oidc
spec:
  credentials:
    source: None
  assumeRoleWithOIDC:
    roleARN: acs:ram::<account-id>:role/<role-name>
    oidcProviderARN: acs:ram::<account-id>:oidc-provider/<provider-name>
    oidcTokenFile: /var/run/secrets/ack.alibabacloud.com/rrsa-tokens/token
    roleSessionName: crossplane-oidc
    sessionExpiration: 3600
```

When running in ACK with RRSA enabled, mount the projected RRSA token into the
provider container and set `oidcTokenFile` to that path. The RAM role trust
policy must trust the matching Alibaba Cloud OIDC provider and ServiceAccount
subject. This authentication setup is separate from the IMS `OIDCProvider`
managed resource, which manages Alibaba Cloud OIDC provider resources.
Because `credentials.source: None` does not read a credentials Secret, regional
managed resources should set `spec.forProvider.region`.

If both `assumeRole` and `assumeRoleWithOIDC` are set, the provider forwards
both structured blocks to the Terraform provider.

After applying a `ProviderConfig`, managed resources can reference it with
`spec.providerConfigRef.name`. If omitted, resources use the `default`
`ProviderConfig`.

## Test

Add an environment variable to set the credentials for the target Alibaba
account for the tests as follows and then run `make e2e`.

```
export UPTEST_CLOUD_CREDENTIALS='{
    "access_key": "...",
    "secret_key": "...",
    "region": "us-west-1"
}'
```
