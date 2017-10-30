package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	account = kingpin.Arg("account", "Account number.").Required().String()
)

// Ec2 Implementation
type ec2Implementation struct {
	session *session.Session
	svc     ec2iface.EC2API
}

// Info about vpc
type vpcInfo struct {
	vpcName string
	vpcCidr string
	vpcID   string
	region  string
}

// Info about subnet
type subnetInfo struct {
	subnetName string
	subnetCidr string
	vpcID      string
}

// Combine vpc and subnet info
type allInfo struct {
	vpcs    vpcInfo
	subnets []subnetInfo
}

func main() {
	// Parse kingpin variables
	kingpin.Version("0.0.1")
	kingpin.Parse()

	// Create map for vpc data
	var data map[string]*allInfo
	data = make(map[string]*allInfo)

	// Get all regions associated with EC2 ("eu-west-1") is endpoint for looking up regions only
	client := newClient("eu-west-1")

	regions, err := client.getAllRegions()
	if err != nil {
		log.WithField("message", "Could not describe aws regions").Fatal(err)
	}

	// Run program for all regions
	for _, r := range regions {
		client := newClient(r)
		// Describe all vpc info
		vpcs, err := client.getVpcs(r)
		if err != nil {
			log.WithField("message", "Could not describe aws vpcs").Fatal(err)
		}
		// Just append data where vpcs are listed
		if len(vpcs) > 0 {
			for _, s := range vpcs {
				subnets, err := client.getSubnets(s.vpcID)
				if err != nil {
					log.WithField("message", "Could not describe aws subnets").Fatal(err)
				}
				// Add vpc and subnet info to map
				data[s.vpcID] = &allInfo{
					subnets: subnets,
					vpcs: vpcInfo{
						vpcName: s.vpcName,
						vpcCidr: s.vpcCidr,
						region:  r,
					},
				}
			}
		}
	}
	// Write vpc and subnet info to csv
	err = writeToFile(data, *account)
	if err != nil {
		log.WithField("message", "Could not write to cvs file").Fatal(err)
	}
}

// Ec2 implementation
func newClient(region string) ec2Implementation {
	var ec2Client ec2Implementation
	ec2Client.session = session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	ec2Client.svc = ec2.New(ec2Client.session)
	return ec2Client
}
