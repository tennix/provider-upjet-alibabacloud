# Build the native plugin from the same upstream version as the generated schema.
# The local patch adds tags to the EIP/vSwitch create requests; CRDs are unchanged.
FROM --platform=$BUILDPLATFORM golang:1.24.1-alpine AS terraform-provider-build
RUN apk --no-cache add ca-certificates git patch
ARG TERRAFORM_PROVIDER_VERSION
WORKDIR /src
ADD https://github.com/aliyun/terraform-provider-alicloud/archive/refs/tags/v${TERRAFORM_PROVIDER_VERSION}.tar.gz /tmp/provider.tar.gz
RUN tar -xzf /tmp/provider.tar.gz --strip-components=1 -C /src
COPY terraform-provider-patches/ /patches/
RUN patch --batch --fuzz=0 -p1 < /patches/network-create-tags.patch
COPY terraform-provider-patches/network_create_tags_test.go.txt /src/alicloud/network_create_tags_test.go
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go test -p 2 ./alicloud -run '^TestNetworkCreateRequestTags$' -count=1
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -p 2 -ldflags="-s -w" -o /out/terraform-provider-alicloud .

FROM alpine:3.20.3
RUN apk --no-cache add ca-certificates bash

ARG TARGETOS
ARG TARGETARCH

ADD "bin/${TARGETOS}_${TARGETARCH}/provider" /usr/local/bin/provider

ENV USER_ID=65532

# Setup Terraform environment

## Provider-dependent configuration
ARG TERRAFORM_VERSION
ARG TERRAFORM_PROVIDER_SOURCE
ARG TERRAFORM_PROVIDER_VERSION
ARG TERRAFORM_PROVIDER_DOWNLOAD_NAME
ARG TERRAFORM_NATIVE_PROVIDER_BINARY
ARG TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX

## End of - Provider-dependent configuration

ENV PLUGIN_DIR=/terraform/provider-mirror/registry.terraform.io/${TERRAFORM_PROVIDER_SOURCE}/${TERRAFORM_PROVIDER_VERSION}/${TARGETOS}_${TARGETARCH}
ENV TF_CLI_CONFIG_FILE=/terraform/.terraformrc
ENV TF_FORK=0

RUN mkdir -p ${PLUGIN_DIR}

ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${TARGETOS}_${TARGETARCH}.zip /tmp
ADD terraformrc.hcl ${TF_CLI_CONFIG_FILE}

RUN unzip /tmp/terraform_${TERRAFORM_VERSION}_${TARGETOS}_${TARGETARCH}.zip -d /usr/local/bin \
  && chmod +x /usr/local/bin/terraform \
  && rm /tmp/terraform_${TERRAFORM_VERSION}_${TARGETOS}_${TARGETARCH}.zip \
  && chown -R ${USER_ID}:${USER_ID} /terraform
COPY --from=terraform-provider-build --chown=65532:65532 /out/terraform-provider-alicloud ${PLUGIN_DIR}/${TERRAFORM_NATIVE_PROVIDER_BINARY}
# End of - Setup Terraform environment

# Provider controller needs these environment variable at runtime
ENV TERRAFORM_VERSION=${TERRAFORM_VERSION}
ENV TERRAFORM_PROVIDER_SOURCE=${TERRAFORM_PROVIDER_SOURCE}
ENV TERRAFORM_PROVIDER_VERSION=${TERRAFORM_PROVIDER_VERSION}
ENV TERRAFORM_NATIVE_PROVIDER_PATH=${PLUGIN_DIR}/${TERRAFORM_NATIVE_PROVIDER_BINARY}

USER ${USER_ID}
EXPOSE 8080

ENTRYPOINT ["provider"]
