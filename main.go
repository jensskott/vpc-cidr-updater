package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	log "github.com/sirupsen/logrus"
)

type ec2Implementation struct {
	session *session.Session
	svc     ec2iface.EC2API
}

type vpcInfo struct {
	vpcName string
	vpcCidr string
	vpcID   string
	region  string
}

type subnetInfo struct {
	subnetName string
	subnetCidr string
	vpcID      string
}

type allInfo struct {
	vpcs    vpcInfo
	subnets []subnetInfo
}

func main() {
	var data []allInfo

	client := newClient("eu-west-1")
	regions, err := client.getAllRegions()
	if err != nil {
		log.WithField("message", "Could not describe aws regions").Fatal(err)
	}

	for _, r := range regions {
		client := newClient(r)
		// Describe all vpc info
		vpcs, err := client.getVpcs(r)
		if err != nil {
			log.WithField("message", "Could not describe aws vpcs").Fatal(err)
		}

		for _, s := range vpcs {
			d := &allInfo{
				vpcs: s,
			}
			subnets, err := client.getSubnets(s.vpcID)
			if err != nil {
				log.WithField("message", "Could not describe aws subnets").Fatal(err)
			}
			d.subnets = subnets

			data = append(data, *d)
		}
	}
	fmt.Println(data)
}

func newClient(region string) ec2Implementation {
	var ec2Client ec2Implementation
	ec2Client.session = session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client.svc = ec2.New(ec2Client.session)
	return ec2Client
}
