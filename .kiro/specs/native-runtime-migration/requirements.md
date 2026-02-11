# Requirements Document

## Introduction

This document specifies the requirements for migrating the Alibaba Cloud provider from upjet v1 architecture to upjet v2 architecture. The current provider uses upjet v1.9.0, which relies on the Terraform CLI and WorkspaceStore abstraction for resource operations. The upjet v2 architecture, used by AWS, Azure, and GCP providers, embeds the Terraform provider directly into the Crossplane provider process, eliminating the need for external Terraform CLI processes and significantly improving performance and resource efficiency.

**Architecture Change:** The upjet v2 architecture compiles the Terraform provider directly into the Crossplane provider binary, removing the need for WorkspaceStore, separate provider binaries, or Terraform CLI processes. This is fundamentally different from the upjet v1 approach.

## Glossary

- **Provider**: The Crossplane provider for Alibaba Cloud that manages cloud resources
- **upjet v1**: The current architecture using WorkspaceStore and Terraform CLI processes
- **upjet v2**: The target architecture that embeds Terraform providers directly in the process
- **WorkspaceStore**: The upjet v1 component responsible for managing Terraform workspace state (removed in v2)
- **Embedded_Provider**: The upjet v2 approach where Terraform provider code is compiled into the Crossplane provider binary
- **Framework_Provider**: Terraform Plugin Framework provider interface used in upjet v2
- **SDK_Provider**: Terraform Plugin SDK provider interface used in upjet v2
- **TerraformProvider**: The compiled Terraform provider instance embedded in the process
- **SetupFn**: The function that configures Terraform provider credentials and configuration
- **Resource_Reconciliation**: The process of ensuring cloud resources match their desired state
- **xpprovider**: Package that provides the embedded Terraform provider implementation

## Requirements

### Requirement 1: upjet v2 Dependency Upgrade

**User Story:** As a developer, I want the provider to use upjet v2, so that it uses the modern embedded provider architecture.

#### Acceptance Criteria

1. WHEN upgrading dependencies, THE Provider SHALL update upjet from v1.9.0 to v2.x
2. WHEN upgrading dependencies, THE Provider SHALL update crossplane-runtime to v2.x
3. THE Provider SHALL resolve all breaking API changes from the upjet v1 to v2 migration
4. THE Provider SHALL update all imports to use the new upjet v2 package paths

### Requirement 2: Terraform Provider Embedding

**User Story:** As a developer, I want the Terraform provider embedded in the binary, so that no external Terraform CLI or provider binaries are needed.

#### Acceptance Criteria

1. WHEN building the provider, THE Build_System SHALL compile the terraform-provider-alicloud code into the provider binary
2. THE Provider SHALL use xpprovider package to expose the embedded Terraform provider
3. THE Provider SHALL initialize both Framework_Provider and SDK_Provider instances at startup
4. THE Provider SHALL NOT require external Terraform CLI installation or provider binaries

### Requirement 3: Provider Configuration Architecture

**User Story:** As a developer, I want the provider to use the upjet v2 configuration pattern, so that it follows the same architecture as AWS/Azure/GCP providers.

#### Acceptance Criteria

1. WHEN initializing the provider, THE Provider SHALL call GetProvider() with Framework_Provider and SDK_Provider instances
2. THE Provider configuration SHALL use config.NewProvider() with the embedded TerraformProvider
3. THE Provider SHALL configure resources using config.WithTerraformProvider() and config.WithTerraformPluginFrameworkProvider()
4. THE Provider SHALL remove all WorkspaceStore initialization code

### Requirement 4: Client Setup Refactoring

**User Story:** As a developer, I want the client setup to use the upjet v2 pattern, so that credentials are configured correctly for the embedded provider.

#### Acceptance Criteria

1. WHEN setting up clients, THE SetupFn SHALL receive the embedded TerraformProvider instance from SetupConfig
2. THE SetupFn SHALL configure provider credentials (access_key, secret_key, security_token, region) directly on the TerraformProvider
3. THE SetupFn SHALL return terraform.Setup with Meta and FrameworkProvider fields populated
4. THE SetupFn SHALL maintain the existing credential extraction and region resolution logic
5. THE SetupFn SHALL preserve the user agent configuration with version information

### Requirement 5: Main Entry Point Refactoring

**User Story:** As a developer, I want the main.go to follow the upjet v2 pattern, so that the provider initializes correctly.

#### Acceptance Criteria

1. WHEN the provider starts, THE Provider SHALL initialize Framework_Provider and SDK_Provider using xpprovider.GetProvider()
2. THE Provider SHALL pass both provider instances to config.GetProvider()
3. THE Provider SHALL configure tjcontroller.Options with the Provider instance (not WorkspaceStore)
4. THE Provider SHALL remove all deprecated Terraform-related flags (terraform-version, terraform-provider-version, etc.)

### Requirement 6: Build System Updates

**User Story:** As a developer, I want the build system updated for upjet v2, so that the provider compiles correctly.

#### Acceptance Criteria

1. WHEN building the provider, THE Build_System SHALL include terraform-provider-alicloud as a Go dependency in go.mod
2. THE Build_System SHALL NOT download or package separate Terraform provider binaries in Dockerfile
3. THE Build_System SHALL remove TERRAFORM_NATIVE_PROVIDER_BINARY variable from Makefile
4. THE Build_System SHALL remove TERRAFORM_PROVIDER_DOWNLOAD_NAME and TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX variables
5. THE Build_System SHALL maintain TERRAFORM_PROVIDER_VERSION for schema generation compatibility
6. THE Build_System SHALL update the schema generation process to work with upjet v2

### Requirement 7: Deployment Configuration Simplification

**User Story:** As a deployment engineer, I want simplified deployment manifests, so that deployments are easier to manage.

#### Acceptance Criteria

1. WHEN deploying the provider, THE Deployment_Configuration SHALL NOT include TERRAFORM_NATIVE_PROVIDER_PATH environment variable
2. THE Deployment_Configuration SHALL NOT include TERRAFORM_VERSION environment variable
3. THE Deployment_Configuration SHALL NOT include TERRAFORM_PROVIDER_VERSION environment variable
4. THE Deployment_Configuration SHALL NOT mount separate provider binaries

### Requirement 8: Backward Compatibility

**User Story:** As a platform operator, I want the migration to preserve existing functionality, so that current deployments continue to work.

#### Acceptance Criteria

1. WHEN using upjet v2, THE Provider SHALL maintain the same external API for resource management
2. THE Provider SHALL support the same credential configuration methods as before
3. THE Provider SHALL maintain compatibility with existing ProviderConfig resources
4. THE Provider SHALL maintain compatibility with existing managed resources

### Requirement 9: Error Handling and Diagnostics

**User Story:** As a platform operator, I want clear error messages, so that I can quickly diagnose and fix issues.

#### Acceptance Criteria

1. IF provider initialization fails, THEN THE Provider SHALL log descriptive errors with context
2. IF credential configuration fails, THEN THE Provider SHALL log the failure reason
3. WHEN the embedded provider encounters errors, THE Provider SHALL log detailed error information
4. THE Provider SHALL include provider status in health check responses

### Requirement 10: Resource Efficiency

**User Story:** As a platform operator, I want upjet v2 to reduce resource consumption, so that I can manage more resources with the same infrastructure.

#### Acceptance Criteria

1. WHEN reconciling resources, THE Provider SHALL use the embedded provider instance instead of spawning processes
2. THE Provider SHALL reuse the same provider instance across multiple resource operations
3. WHEN compared to upjet v1, THE Provider SHALL demonstrate reduced memory footprint
4. WHEN compared to upjet v1, THE Provider SHALL demonstrate reduced CPU overhead for resource operations

### Requirement 11: Testing and Validation

**User Story:** As a developer, I want comprehensive tests for upjet v2, so that I can verify the migration works correctly.

#### Acceptance Criteria

1. THE Test_Suite SHALL include e2e tests using uptest framework that verify upjet v2 initialization
2. THE Test_Suite SHALL include e2e tests that verify resource CRUD operations work with upjet v2
3. THE Test_Suite SHALL include e2e tests that verify credential configuration works correctly
4. THE Test_Suite SHALL verify existing example manifests continue to work with upjet v2
5. THE Test_Suite SHALL use the existing uptest framework and test patterns (setup.sh, UPTEST_CLOUD_CREDENTIALS)

### Requirement 12: Documentation

**User Story:** As a developer or operator, I want documentation explaining the upjet v2 migration, so that I understand the changes and how to troubleshoot issues.

#### Acceptance Criteria

1. THE Documentation SHALL explain the differences between upjet v1 and upjet v2 architectures
2. THE Documentation SHALL provide migration guide for upgrading from v1 to v2
3. THE Documentation SHALL include troubleshooting guidance for common upjet v2 issues
4. THE Documentation SHALL document the removed environment variables and configuration options
