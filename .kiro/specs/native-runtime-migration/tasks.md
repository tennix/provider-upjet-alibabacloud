# Implementation Tasks

## Phase 1: Dependency Upgrade

- [x] 1. Update go.mod dependencies
  - [x] 1.1 Upgrade upjet from v1.9.0 to v2.x
  - [x] 1.2 Upgrade crossplane-runtime to v2.x
  - [x] 1.3 Add terraform-provider-alicloud as direct dependency
  - [x] 1.4 Run `go mod tidy` to resolve dependencies

- [x] 2. Update import paths
  - [x] 2.1 Update all upjet imports to v2 paths
  - [x] 2.2 Update crossplane-runtime imports to v2 paths
  - [x] 2.3 Fix compilation errors from import changes

- [x] 3. Resolve breaking API changes
  - [x] 3.1 Fix upjet v2 API breaking changes
  - [x] 3.2 Fix crossplane-runtime v2 API breaking changes
  - [x] 3.3 Verify all packages compile successfully

## Phase 2: Provider Initialization Refactoring

- [x] 4. Create or integrate xpprovider package
  - [x] 4.1 Investigate if terraform-provider-alicloud has xpprovider package
  - [x] 4.2 If not available, create xpprovider wrapper package
  - [x] 4.3 Implement GetProvider() function to return Framework and SDK providers
  - [x] 4.4 Implement AlibabaCloudConfig struct for provider configuration
  - [x] 4.5 Implement GetClient() method for provider client initialization

- [x] 5. Update cmd/provider/main.go
  - [x] 5.1 Add xpprovider.GetProvider() call at startup
  - [x] 5.2 Update config.GetProvider() call with Framework and SDK providers
  - [x] 5.3 Remove WorkspaceStore initialization
  - [x] 5.4 Remove deprecated flags (terraform-version, terraform-provider-source, terraform-provider-version)
  - [x] 5.5 Add SetupConfig initialization with TerraformProvider
  - [x] 5.6 Update tjcontroller.Options to use Provider instead of WorkspaceStore
  - [x] 5.7 Add PollJitter and OperationTrackerStore configuration

- [x] 6. Update config/provider.go
  - [x] 6.1 Change GetProvider() signature to accept ctx, fwProvider, sdkProvider parameters
  - [x] 6.2 Add WithTerraformProvider(sdkProvider) option
  - [x] 6.3 Add WithTerraformPluginFrameworkProvider(fwProvider) option
  - [x] 6.4 Update return type to include error
  - [x] 6.5 Verify all resource configurations still work

## Phase 3: Client Setup Refactoring

- [x] 7. Update internal/clients/alibabacloud.go
  - [x] 7.1 Create SetupConfig struct with TerraformProvider and Logger fields
  - [x] 7.2 Rename TerraformSetupBuilder to SelectTerraformSetup
  - [x] 7.3 Update SelectTerraformSetup signature to accept SetupConfig
  - [x] 7.4 Remove version, providerSource, providerVersion parameters
  - [x] 7.5 Remove ps.Version and ps.Requirement initialization
  - [x] 7.6 Create configureNoForkAlibabaCloudClient function
  - [x] 7.7 Implement provider configuration with credentials in configureNoForkAlibabaCloudClient
  - [x] 7.8 Set ps.Meta with configured provider client
  - [x] 7.9 Preserve extractAndUnmarshalCredentials function (no changes)
  - [x] 7.10 Preserve getRegion function (no changes)
  - [x] 7.11 Preserve getUserAgent function (no changes)

## Phase 4: Build System Updates

- [x] 8. Update Makefile
  - [x] 8.1 Remove TERRAFORM_NATIVE_PROVIDER_BINARY variable
  - [x] 8.2 Remove TERRAFORM_PROVIDER_DOWNLOAD_NAME variable
  - [x] 8.3 Remove TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX variable
  - [x] 8.4 Keep TERRAFORM_PROVIDER_VERSION for schema generation
  - [x] 8.5 Keep TERRAFORM_PROVIDER_SOURCE for schema generation
  - [x] 8.6 Update schema generation process if needed for upjet v2

- [-] 9. Update Dockerfile
  - [x] 9.1 Remove provider binary download steps
  - [x] 9.2 Verify go.mod includes terraform-provider-alicloud
  - [x] 9.3 Ensure build process compiles embedded provider

- [x] 10. Update deployment manifests
  - [x] 10.1 Remove TERRAFORM_NATIVE_PROVIDER_PATH environment variable from package/crossplane.yaml.tmpl
  - [x] 10.2 Remove TERRAFORM_VERSION environment variable
  - [x] 10.3 Remove TERRAFORM_PROVIDER_VERSION environment variable
  - [x] 10.4 Remove TERRAFORM_PROVIDER_SOURCE environment variable
  - [x] 10.5 Remove provider binary volume mounts if any

## Phase 5: Testing and Validation

- [ ] 11. Build and test locally
  - [ ] 11.1 Run `make submodules`
  - [ ] 11.2 Run `make generate`
  - [ ] 11.3 Run `make build`
  - [ ] 11.4 Verify provider binary builds successfully
  - [ ] 11.5 Run `make local-deploy` to deploy locally
  - [ ] 11.6 Verify provider pod starts successfully
  - [ ] 11.7 Check provider logs for initialization errors

- [ ] 12. Test credential configuration
  - [ ] 12.1 Create test ProviderConfig with access_key and secret_key
  - [ ] 12.2 Verify provider can read credentials from secret
  - [ ] 12.3 Test with security_token credential
  - [ ] 12.4 Test region resolution from credentials
  - [ ] 12.5 Test region resolution from spec.forProvider.region

- [ ] 13. Run e2e tests for VPC resources
  - [ ] 13.1 Set UPTEST_CLOUD_CREDENTIALS environment variable
  - [ ] 13.2 Run `make uptest UPTEST_EXAMPLE_LIST=examples/vpc/v1alpha1/vpc.yaml`
  - [ ] 13.3 Verify VPC creation succeeds
  - [ ] 13.4 Verify VPC update works
  - [ ] 13.5 Verify VPC deletion completes
  - [ ] 13.6 Check status.atProvider is populated correctly

- [ ] 14. Run e2e tests for ECS resources
  - [ ] 14.1 Run `make uptest UPTEST_EXAMPLE_LIST=$(UPTEST_EXAMPLE_LIST_ECS)`
  - [ ] 14.2 Verify all ECS resource operations succeed
  - [ ] 14.3 Check for any regression issues

- [ ] 15. Run e2e tests for all resource groups
  - [ ] 15.1 Run tests for ACK resources
  - [ ] 15.2 Run tests for ALB resources
  - [ ] 15.3 Run tests for ALIDNS resources
  - [ ] 15.4 Run tests for OSS resources
  - [ ] 15.5 Run tests for RAM resources
  - [ ] 15.6 Run tests for remaining resource groups
  - [ ] 15.7 Document any failing tests and investigate

- [ ] 16. Performance testing
  - [ ] 16.1 Measure memory usage with 50 VPC resources (upjet v2)
  - [ ] 16.2 Compare with baseline memory usage (upjet v1 if available)
  - [ ] 16.3 Measure reconciliation latency for 50 VPC creates
  - [ ] 16.4 Compare with baseline latency
  - [ ] 16.5 Document performance improvements

- [ ] 17. Backward compatibility testing
  - [ ] 17.1 Create resources with upjet v1 (if possible)
  - [ ] 17.2 Upgrade to upjet v2
  - [ ] 17.3 Verify existing resources continue to reconcile
  - [ ] 17.4 Verify status.atProvider state is preserved
  - [ ] 17.5 Test resource updates on migrated resources
  - [ ] 17.6 Test resource deletion on migrated resources

## Phase 6: Documentation

- [ ] 18. Update README.md
  - [ ] 18.1 Remove references to Terraform CLI requirements
  - [ ] 18.2 Update build instructions for upjet v2
  - [ ] 18.3 Update test instructions
  - [ ] 18.4 Add note about embedded provider architecture

- [ ] 19. Create migration guide
  - [ ] 19.1 Document breaking changes from upjet v1 to v2
  - [ ] 19.2 Document removed environment variables
  - [ ] 19.3 Document new requirements (if any)
  - [ ] 19.4 Provide troubleshooting steps for common issues
  - [ ] 19.5 Explain state storage (no migration needed)

- [ ] 20. Update examples
  - [ ] 20.1 Review all example manifests
  - [ ] 20.2 Update any examples that reference old configuration
  - [ ] 20.3 Verify all examples work with upjet v2

## Phase 7: Final Validation and Release

- [ ] 21. Code review and cleanup
  - [ ] 21.1 Run `make lint` and fix any issues
  - [ ] 21.2 Review all code changes
  - [ ] 21.3 Remove any commented-out code
  - [ ] 21.4 Ensure code follows project conventions

- [ ] 22. Final testing
  - [ ] 22.1 Run full e2e test suite
  - [ ] 22.2 Verify all tests pass
  - [ ] 22.3 Test in clean Kubernetes cluster
  - [ ] 22.4 Verify provider installation from package

- [ ] 23. Release preparation
  - [ ] 23.1 Update CHANGELOG.md with migration notes
  - [ ] 23.2 Tag release with appropriate version
  - [ ] 23.3 Build and publish provider package
  - [ ] 23.4 Announce migration to users
