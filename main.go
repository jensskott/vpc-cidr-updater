package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)


type ec2Implementation struct {
	session *session.Session
	svc     ec2iface.EC2API
}

type vpcInfo struct {
	vpcName string
	vpcCidr string
	region  string
	subnets []subnetInfo
}

type subnetInfo struct {
	subnetName string
	subnetCidr string
}

func main() {


	client := newClient("eu-west-1")
	regions, err := client.getAllRegions()
	if err != nil {
		log.WithField("message", "Could not describe aws regions").Fatal(err)
	}


	// Describe all vpc info
	vpcs, err := getAllVpcs(regions)
	if err != nil {
		log.WithField("message", "Could not describe aws vpcs").Fatal(err)
	}
	for _, r := range vpcs {
		out := fmt.Sprintf("VPC INFO Region: %s Name: %s, Cidr: %s \n Subnets: %v", r.region, r.vpcName, r.vpcCidr, r.subnets)
		fmt.Println(out)
	}
}

func newClient(region string) ec2Implementation {
	var ec2Client ec2Implementation
	ec2Client.session = session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client.svc = ec2.New(ec2Client.session)
	return ec2Client
}


