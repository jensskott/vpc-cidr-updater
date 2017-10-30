package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func writeToFile(data map[string]*allInfo, account string) error {
	// Create file to write to
	fileName := fmt.Sprintf("%s-networks", account)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a writer for csv
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write each key in map to a new line in csv
	for _, l := range data {
		var s []string
		s = append(s, l.vpcs.vpcName, l.vpcs.vpcCidr, l.vpcs.region)
		for _, subs := range l.subnets {
			s = append(s, subs.subnetName, subs.subnetCidr)
		}
		err = writer.Write(s)
		if err != nil {
			return err
		}
	}

	return nil
}
