package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func writeToFile(data map[string]*allInfo, account string) error {
	// Create file to write to
	fileName := fmt.Sprintf("%s-networks.csv", account)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	// Create a writer for csv
	writer := csv.NewWriter(file)
	defer writer.Flush()

	if len(data) > 0 {
		// Write each key in map to a new line in csv
		for _, l := range data {
			var s []string
			s = append(s, l.vpcs.vpcName, l.vpcs.vpcCidr, l.vpcs.region)
			for _, subs := range l.subnets {
				s = append(s, subs.subnetName, subs.subnetCidr)
			}
			err := writer.Write(s)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("no data to create csv file")
	}

	return nil
}
