package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (e *ec2Implementation) getAllRegions() ([]string, error) {
	// Create a slice for all regions
	var regions []string

	// List all regions for ec2
	resp, err := e.svc.DescribeRegions(nil)
	if err != nil {
		return nil, err
	}

	// Append all regions in response to slice
	for _, r := range resp.Regions {
		regions = append(regions, *r.RegionName)
	}
	return regions, nil
}

func (e *ec2Implementation) getVpcs(region string) ([]vpcInfo, error) {
	// Create slice for vpc info
	var vpcs []vpcInfo

	// Describe all vpcs in region
	resp, err := e.svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	// if not default vpc add to vpc info struct
	for _, vpc := range resp.Vpcs {
		if *vpc.IsDefault != true {
			v := vpcInfo{
				vpcCidr: *vpc.CidrBlock,
				vpcID:   *vpc.VpcId,
				region:  region,
			}
			name := getNameTag(vpc.Tags)
			v.vpcName = name
			vpcs = append(vpcs, v)
		}
	}

	return vpcs, nil
}

func getNameTag(tag []*ec2.Tag) string {
	// Get the name tag from vpc tags
	for _, t := range tag {
		if *t.Key == "Name" {
			return *t.Value
		}
		continue
	}
	return ""
}

func (e *ec2Implementation) getSubnets(vpcID string) ([]subnetInfo, error) {
	// Create subnetInfo slice
	var subnets []subnetInfo

	// Describe all subnets for the vpcID provided
	resp, err := e.svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Add subnets to subnetInfo slice
	for _, s := range resp.Subnets {
		subs := &subnetInfo{
			subnetCidr: *s.CidrBlock,
			vpcID:      *s.VpcId,
		}
		names := getNameTag(s.Tags)
		subs.subnetName = names
		subnets = append(subnets, *subs)
	}

	return subnets, nil
}
