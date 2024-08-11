// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package batch

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	batch_sdkv2 "github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{
		{
			Factory: newJobDefinitionDataSource,
			Name:    "Job Definition",
		},
	}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory: newJobQueueResource,
			Name:    "Job Queue",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceComputeEnvironment,
			TypeName: "aws_batch_compute_environment",
			Name:     "Compute Environment",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  DataSourceJobQueue,
			TypeName: "aws_batch_job_queue",
			Name:     "Job Queue",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceSchedulingPolicy,
			TypeName: "aws_batch_scheduling_policy",
			Name:     "Scheduling Policy",
			Tags:     &types.ServicePackageResourceTags{},
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceComputeEnvironment,
			TypeName: "aws_batch_compute_environment",
			Name:     "Compute Environment",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceJobDefinition,
			TypeName: "aws_batch_job_definition",
			Name:     "Job Definition",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceSchedulingPolicy,
			TypeName: "aws_batch_scheduling_policy",
			Name:     "Scheduling Policy",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.Batch
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*batch_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return batch_sdkv2.NewFromConfig(cfg,
		batch_sdkv2.WithEndpointResolverV2(newEndpointResolverSDKv2()),
		withBaseEndpoint(config[names.AttrEndpoint].(string)),
	), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
