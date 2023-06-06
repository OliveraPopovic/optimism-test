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

type Address struct {
	Address string `json:"address"`
}

const OVMETHContractAddress = "DeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000"

var VisitedStorageAddresses map[string]bool
var VisitedAccountAddresses map[string]bool
var OVMETHStorageAddresses map[string]interface{}
var GenesisAddresses []string

func main() {
	OVMETHStorageAddresses, GenesisAddresses = loadGenesisAddresses()

	VisitedStorageAddresses = make(map[string]bool)
	VisitedAccountAddresses = make(map[string]bool)

	fmt.Println(CalcOVMETHStorageKey(common.HexToAddress("70b6910626fb2e3ba27528236895659a791ac9fc")))
	// kopiranje storage adresa na hit/miss mapu
	for s, _ := range OVMETHStorageAddresses {
		VisitedStorageAddresses[s] = false
	}

	checkNumberOfMatchedAddresses()

	newGenesis := addAddressesToGenesis(false)
	filePath := "./new_genesis.json"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	s, _ := json.MarshalIndent(newGenesis, "", "  ")
	file.Write(s)
}

func addAddressesToGenesis(printChanges bool) *Genesis {
	originalGenesisData := loadFullGenesisData()
	addressesWithBalances := loadStoredAddressesForTest()

	for _, address := range addressesWithBalances {
		accHex := CalcOVMETHStorageKey(common.HexToAddress(address))
		addressStorageWithoutPrefix, _ := strings.CutPrefix(accHex.String(), "0x")
		if OVMETHStorageAddresses[addressStorageWithoutPrefix] == nil {
			panic("shouldn't happen ever")
		}

		balance := originalGenesisData.Alloc[OVMETHContractAddress].Storage[addressStorageWithoutPrefix]
		newAccount := GenesisAccount{
			Balance: balance,
			Nonce:   "0",
		}

		oldAccount := originalGenesisData.Alloc[address]
		if printChanges {
			color.Magenta("Address: %s", address)
			color.White("Previous { balance: %s, nonce %s }", originalGenesisData.Alloc[address].Balance, originalGenesisData.Alloc[address].Nonce)
		}
		if oldAccount.Nonce == "" {
			originalGenesisData.Alloc[address] = newAccount
		} else {
			oldAccount.Balance = balance
			originalGenesisData.Alloc[address] = oldAccount
		}
		if printChanges {
			color.Cyan("New { balance: %s, nonce %s }", originalGenesisData.Alloc[address].Balance, originalGenesisData.Alloc[address].Nonce)
		}
	}

	return originalGenesisData
}

func checkNumberOfMatchedAddresses() {
	GenesisAddresses = loadStoredAddressesForTest()
	DuneDumpAddresses := loadDuneDump()
	PreJuneAddresses := loadPreJuneAddresses()
	GoerliDumpAddresses := loadGoerliAddresses()

	goThroughAddresses(GenesisAddresses)
	printCurrentStatus("After Genesis addresses")

	goThroughAddresses(PreJuneAddresses)
	printCurrentStatus("After adding pre-June data")

	goThroughAddresses(GoerliDumpAddresses)
	printCurrentStatus("After Goerli dump addresses")

	goThroughAddresses(DuneDumpAddresses)
	printCurrentStatus("After Dune addresses")

	saveAddressesWithStorageToFile()
}

func loadStoredAddressesForTest() []string {
	data, err := os.ReadFile("./optimism_addresses_w_initial_balance.json")
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
}

func goThroughAddresses(addresses []string) {
	//for _, address := range addresses {
	//	if strings.ToLower(address) == "0x70b6910626fb2e3ba27528236895659a791ac9fc" || strings.ToLower(address) == "70b6910626fb2e3ba27528236895659a791ac9fc" {
	//		panic("THIS SHOULD HAPPEN FFS")
	//	}
	//}

	for _, addressHex := range addresses {
		//address := common.HexToAddress(addressHex)
		//addressStorage := CalcOVMETHStorageKey(address)
		//addressStorageWithoutPrefix, _ := strings.CutPrefix(addressStorage.String(), "0x")
		addressStorageWithoutPrefix := addressHex

		// TODO remove second condition, only a test
		if OVMETHStorageAddresses[addressStorageWithoutPrefix] != nil || OVMETHStorageAddresses[strings.ToLower(addressStorageWithoutPrefix)] != nil {
			VisitedAccountAddresses[addressHex] = true
			VisitedStorageAddresses[addressStorageWithoutPrefix] = true
		}
	}
}

func printCurrentStatus(text string) {
	numFound := 0
	numMissed := 0
	for _, b := range VisitedStorageAddresses {
		if b == false {
			numMissed++
		} else {
			numFound++
		}
	}

	color.Cyan(text)
	color.Green("%d", numFound)
	color.Red("%d", numMissed)
}

func saveAddressesWithStorageToFile() {

	keys := make([]string, 0, len(VisitedAccountAddresses))
	for key := range VisitedAccountAddresses {
		if strings.HasPrefix(key, "0x") {
			key, _ = strings.CutPrefix(key, "0x")
		}
		keys = append(keys, key)
	}

	// Open the file for writing
	filePath := "./optimism_addresses_w_initial_balance.json"
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Encode the keys as JSON and write to the file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(keys); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("Keys have been written to", filePath)
}

func loadFullGenesisData() *Genesis {
	data, err := os.ReadFile("/Users/oliverapopovic/Downloads/genesis-berlin.json")
	if err != nil {
		log.Fatal(err)
	}
	var genesis Genesis

	err = json.Unmarshal(data, &genesis)
	if err != nil {
		log.Fatal(err)
	}

	return &genesis
}

func loadGenesisAddresses() (map[string]interface{}, []string) {
	// Read the JSON file
	data, err := os.ReadFile("/Users/oliverapopovic/Downloads/genesis-berlin.json")
	//data, err := os.ReadFile("./new_genesis.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create a map to hold the parsed JSON data
	var jsonData Genesis

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	var OVMETHStorage = make(map[string]interface{})
	var accountAddresses []string

	for key, value := range jsonData.Alloc[OVMETHContractAddress].Storage {
		OVMETHStorage[key] = value
	}

	for key, _ := range jsonData.Alloc {
		accountAddresses = append(accountAddresses, key)
	}
	//allocField, ok := jsonData["alloc"]

	//if ok {
	//	OvmEthContract, ok := allocField.(map[string]interface{})["DeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000"]
	//	if ok {
	//		addresses, _ = OvmEthContract.(map[string]interface{})["storage"].(map[string]interface{})
	//		for key, _ := range allocField.(map[string]interface{}) {
	//			// TODO remove, only a test
	//			// doesn't seem to change anything
	//			key = strings.ToLower(key)
	//
	//			otherAddresses = append(otherAddresses, "0x"+key)
	//		}
	//	}
	//}

	//normalizedAddresses := make(map[string]interface{})
	//
	//for key, value := range addresses {
	//	if key != strings.ToLower(key) {
	//		fmt.Println(key)
	//	}
	//
	//	normalizedAddresses[strings.ToLower(key)] = value
	//}

	return OVMETHStorage, accountAddresses
}

func loadGoerliAddresses() []string {
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

func loadDuneDump() []string {
	data, err := os.ReadFile("/Users/oliverapopovic/Downloads/01H1RM0J3XN1V4YEE6W0Y1470X.csv")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")

	data, err = os.ReadFile("/Users/oliverapopovic/Downloads/01H1RK10FFXHDCV9FP51RC9CDN.csv")
	if err != nil {
		log.Fatal(err)
	}

	lines2 := strings.Split(string(data), "\n")

	toReturn := append(lines, lines2...)

	return toReturn
}

func loadPreJuneAddresses() []string {
	data, err := os.ReadFile("/Users/oliverapopovic/Downloads/addr.json")
	if err != nil {
		log.Fatal(err)
	}

	// Create a map to hold the parsed JSON data
	var addresses []Address

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(data, &addresses)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the address values
	var addressValues []string
	for _, addr := range addresses {
		addressValues = append(addressValues, addr.Address)
	}

	return addressValues
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
