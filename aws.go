package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (e *ec2Implementation) getAllRegions() ([]string, error) {
	var regions []string
	resp, err := e.svc.DescribeRegions(nil)
	if err != nil {
		return nil, err
	}
	for _, r := range resp.Regions {
		regions = append(regions, *r.RegionName)
	}
	return regions, nil
}

func (e *ec2Implementation) getVpcs(region string) ([]vpcInfo, error) {
	var vpcs []vpcInfo

	resp, err := e.svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

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
	for _, t := range tag {
		return *t.Value
	}

	return ""
}

// TODO: Get all subnets of a vpc here
func (e *ec2Implementation) getSubnets(vpcID string) ([]subnetInfo, error) {
	var subnets []subnetInfo

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
