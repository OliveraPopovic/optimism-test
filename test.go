package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
	"golang.org/x/crypto/sha3"
	"log"
	"os"
	"strings"
)

func main() {
	GenesisAddresses, FileAddresses := loadGenesisAddresses()
	//FileAddresses := loadFileAddresses()

	numFound := 0
	numMissed := 0
	visitedGenesisAddresses := make(map[string]bool)

	// kopiranje storage adresa na hit/miss mapu
	for s, _ := range GenesisAddresses {
		visitedGenesisAddresses[s] = false
	}

	//for _, address := range FileAddresses {
	//	a := common.HexToAddress(address)
	//	b := CalcOVMETHStorageKey(a)
	//	c, _ := strings.CutPrefix(b.String(), "0x")
	//	if strings.ToLower(c) == "fa3e200bf8cec729781af9908c8b10bc70c98a770c6d8704188eb42d2f416cfd" {
	//		panic("")
	//	}
	//	if address != strings.ToLower(address) {
	//		fmt.Println()
	//	}
	//}

	for _, addressHex := range FileAddresses {
		address := common.HexToAddress(addressHex)
		addressStorage := CalcOVMETHStorageKey(address)
		addressStorageWithoutPrefix, _ := strings.CutPrefix(addressStorage.String(), "0x")

		// TODO remove second condition, only a test
		if GenesisAddresses[addressStorageWithoutPrefix] != nil || GenesisAddresses[strings.ToLower(addressStorageWithoutPrefix)] != nil {
			visitedGenesisAddresses[addressStorageWithoutPrefix] = true
		}
	}

	for _, b := range visitedGenesisAddresses {
		if b == false {
			numMissed++
		} else {
			numFound++
		}
	}

	color.Green("%d", numFound)
	color.Red("%d", numMissed)
}

func loadGenesisAddresses() (map[string]interface{}, []string) {
	// Read the JSON file
	data, err := os.ReadFile("/Users/oliverapopovic/Downloads/genesis-berlin.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create a map to hold the parsed JSON data
	var jsonData map[string]interface{}

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	var addresses map[string]interface{}
	var otherAddresses []string
	allocField, ok := jsonData["alloc"]

	if ok {
		OvmEthContract, ok := allocField.(map[string]interface{})["DeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000"]
		if ok {
			addresses, _ = OvmEthContract.(map[string]interface{})["storage"].(map[string]interface{})
			for key, _ := range allocField.(map[string]interface{}) {
				// TODO remove, only a test
				// doesn't seem to change anything
				key = strings.ToLower(key)

				otherAddresses = append(otherAddresses, "0x"+key)
			}
		}
	}

	normalizedAddresses := make(map[string]interface{})

	for key, value := range addresses {
		if key != strings.ToLower(key) {
			fmt.Println(key)
		}

		normalizedAddresses[strings.ToLower(key)] = value
	}

	return normalizedAddresses, otherAddresses
}

func loadFileAddresses() []string {
	data, err := os.ReadFile("/Users/oliverapopovic/Repositories/optimism/packages/migration-data/data/ovm-addresses.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create a map to hold the parsed JSON data
	var jsonData []string

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	return jsonData

	//var fakeMap = make(map[string]interface{})
	//
	//for _, datum := range jsonData {
	//	fakeMap[datum] = nil
	//}
	//
	//return fakeMap
}

// BytesBacked is a re-export of the same interface in Geth,
// which is unfortunately private.
type BytesBacked interface {
	Bytes() []byte
}

// CalcOVMETHStorageKey calculates the storage key of an OVM ETH balance.
func CalcOVMETHStorageKey(addr common.Address) common.Hash {
	return CalcStorageKey(addr, common.Big0)
}

// CalcStorageKey is a helper method to calculate storage keys.
func CalcStorageKey(a, b BytesBacked) common.Hash {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(common.LeftPadBytes(a.Bytes(), 32))
	hasher.Write(common.LeftPadBytes(b.Bytes(), 32))
	digest := hasher.Sum(nil)
	return common.BytesToHash(digest)
}
