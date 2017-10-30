package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/jensskott/vpc-cidr-updater/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAllVpcs(t *testing.T) {
	resp := &ec2.DescribeVpcsOutput{
		Vpcs: []*ec2.Vpc{
			{
				CidrBlock: aws.String("10.1.0.0/16"),
				VpcId:     aws.String("vpc-7d5e1719"),
				IsDefault: aws.Bool(false),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-vpc"),
					},
				},
			},
			{
				CidrBlock: aws.String("172.13.0.1/16"),
				VpcId:     aws.String("vpc-a74ab7c2"),
				IsDefault: aws.Bool(true),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("default"),
					},
				},
			},
		},
	}
	// Setup gomock controller with data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockEC2API(ctrl)
	mockSvc.EXPECT().DescribeVpcs(gomock.Any()).Return(resp, nil)

	e := ec2Implementation{
		svc: mockSvc,
	}

	testResp, err := e.getVpcs("us-west-1")
	assert.Equal(t, "test", testResp)
	assert.NoError(t, err)
}

func TestGetSubnets(t *testing.T) {
	resp := &ec2.DescribeSubnetsOutput{
		Subnets: []*ec2.Subnet{
			{
				CidrBlock: aws.String("10.1.1.0/24"),
				VpcId:     aws.String("vpc-7d5e1719"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("test-vpc-subnet"),
					},
				},
			},
		},
	}
	// Setup gomock controller with data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSvc := mocks.NewMockEC2API(ctrl)
	mockSvc.EXPECT().DescribeSubnets(gomock.Any()).Return(resp, nil)

	e := ec2Implementation{
		svc: mockSvc,
	}

	testResp, err := e.getSubnets("vpc-7d5e1719")
	assert.Equal(t, "test", testResp)
	assert.NoError(t, err)
}
