// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connect

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	awstypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// @SDKDataSource("aws_connect_routing_profile")
func DataSourceRoutingProfile() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceRoutingProfileRead,
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_outbound_queue_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"media_concurrencies": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"concurrency": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name", "routing_profile_id"},
			},
			"queue_configs": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"queue_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"queue_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"queue_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"routing_profile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"routing_profile_id", "name"},
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceRoutingProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).ConnectClient(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	instanceID := d.Get("instance_id").(string)

	input := &connect.DescribeRoutingProfileInput{
		InstanceId: aws.String(instanceID),
	}

	if v, ok := d.GetOk("routing_profile_id"); ok {
		input.RoutingProfileId = aws.String(v.(string))
	} else if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		routingProfileSummary, err := dataSourceGetRoutingProfileSummaryByName(ctx, conn, instanceID, name)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "finding Connect Routing Profile Summary by name (%s): %s", name, err)
		}

		if routingProfileSummary == nil {
			return sdkdiag.AppendErrorf(diags, "finding Connect Routing Profile Summary by name (%s): not found", name)
		}

		input.RoutingProfileId = routingProfileSummary.Id
	}

	resp, err := conn.DescribeRoutingProfile(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "getting Connect Routing Profile: %s", err)
	}

	if resp == nil || resp.RoutingProfile == nil {
		return sdkdiag.AppendErrorf(diags, "getting Connect Routing Profile: empty response")
	}

	routingProfile := resp.RoutingProfile

	if err := d.Set("media_concurrencies", flattenRoutingProfileMediaConcurrencies(routingProfile.MediaConcurrencies)); err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	d.Set("arn", routingProfile.RoutingProfileArn)
	d.Set("default_outbound_queue_id", routingProfile.DefaultOutboundQueueId)
	d.Set("description", routingProfile.Description)
	d.Set("instance_id", instanceID)
	d.Set("name", routingProfile.Name)
	d.Set("routing_profile_id", routingProfile.RoutingProfileId)

	// getting the routing profile queues uses a separate API: ListRoutingProfileQueues
	queueConfigs, err := getRoutingProfileQueueConfigs(ctx, conn, instanceID, *routingProfile.RoutingProfileId)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "finding Connect Routing Profile Queue Configs Summary by Routing Profile ID (%s): %s", *routingProfile.RoutingProfileId, err)
	}

	d.Set("queue_configs", queueConfigs)

	if err := d.Set("tags", KeyValueTags(ctx, routingProfile.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", instanceID, aws.ToString(routingProfile.RoutingProfileId)))

	return diags
}

func dataSourceGetRoutingProfileSummaryByName(ctx context.Context, conn *connect.Client, instanceID, name string) (*awstypes.RoutingProfileSummary, error) {
	var result *awstypes.RoutingProfileSummary

	input := &connect.ListRoutingProfilesInput{
		InstanceId: aws.String(instanceID),
		MaxResults: aws.Int32(ListRoutingProfilesMaxResults),
	}

	pages := connect.NewListRoutingProfilesPaginator(conn, input)

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		for _, qs := range page.RoutingProfileSummaryList {
			if aws.ToString(qs.Name) == name {
				result = &qs
			}
		}
	}

	return result, nil
}
