package main

import (
	"testing"

	"encoding/csv"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestWriteToFile(t *testing.T) {
	// Make a map with data to test
	var data map[string]*allInfo
	data = make(map[string]*allInfo)

	data["vpc-7d5e1719"] = &allInfo{
		vpcs: vpcInfo{
			vpcName: "test_name",
			vpcCidr: "10.1.0.0/16",
			region:  "us-west-1",
		},
		subnets: []subnetInfo{
			{
				subnetName: "test_name_subnet",
				subnetCidr: "10.1.1.0/24",
			},
		},
	}

	// Write data to file
	err := writeToFile(data, "1234")
	assert.NoError(t, err)

	// Read file for test
	file, err := os.Open("1234-networks.csv")
	assert.NoError(t, err)
	r := csv.NewReader(file)
	r.Comma = ','
	// Read all lines and assert
	lines, _ := r.ReadAll()
	for _, l := range lines {
		assert.Equal(t, "test_name", l[0])
		assert.Equal(t, "10.1.1.0/24", l[4])
	}
	// Remove the test file
	os.Remove("1234-networks.csv")
}

func TestWriteToFileError(t *testing.T) {
	// Make a map with empty data to test
	var data map[string]*allInfo
	data = make(map[string]*allInfo)

	err := writeToFile(data, "1234")
	assert.Error(t, err)

	// Remove the test file
	os.Remove("1234-networks.csv")
}
