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

func getAllVpcs(regions []string) ([]vpcInfo, error) {
	var vpcs []vpcInfo

	for _, r := range regions {
		client := newClient(r)

		resp, err := client.svc.DescribeVpcs(nil)
		if err != nil {
			return nil, err
		}

		for _, vpc := range resp.Vpcs {
			if *vpc.IsDefault != true {
				v := vpcInfo{
					vpcCidr: *vpc.CidrBlock,
					region:  r,
				}
				name := getNameTag(vpc.Tags)
				v.vpcName = name
				subnets, err := client.getSubnets(*vpc.VpcId)
				if err != nil {
					return nil, err
				}
				v.subnets = subnets
				vpcs = append(vpcs, v)
			}
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
func (e *ec2Implementation) getSubnets (vpcId string) ([]subnetInfo, error) {
	var subnets []subnetInfo

	resp, err := e.svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	for _, s := range resp.Subnets {
		subs := &subnetInfo{
			subnetCidr: *s.CidrBlock,
		}
		names := getNameTag(s.Tags)
		subs.subnetName = names
		subnets = append(subnets, *subs)
	}

	return subnets, nil
}
