package main

import "bytes"
import "fmt"
import "io"
import "os"
import "github.com/Jeffail/gabs/v2"

func createBacktestConfig(args Args) {
	// Create the backtest config file
	// We need to copy the existing config file
	// fmt.Println("Updating backtest config file...")
	// Copy the existing config file
	original := args.path + "/user_data/" + args.config
	testConfig := args.path + "/user_data/config-whitelist-builder.json"

	// If the testConfig file doesn't exist, copy it from the original
	if _, err := os.Stat(testConfig); os.IsNotExist(err) {
		fmt.Println("Creating backtest config file...")

		originalFile, err := os.Open(original)
		check(err, "Error opening original config file")
		defer originalFile.Close()

		testConfigFile, err := os.Create(testConfig)
		check(err, "Error creating test config file")
		defer testConfigFile.Close()

		_, err = io.Copy(testConfigFile, originalFile)
		check(err, "Error copying original config file")
	}

}

func getBlacklist(args Args) []string {
  // Get the blacklist from the config file
  // Open the testConfig file
  config := args.path + "/user_data/" + args.config
  configFile, err := os.Open(config)
  check(err, "Error opening config file")
  defer configFile.Close()
  // Read the testConfig file
  buf := new(bytes.Buffer)
  _, err = buf.ReadFrom(configFile)
  check(err, "Error reading test config file")
  // Convert the buffer to a JSON object
  jsonParsed, err := gabs.ParseJSON(buf.Bytes())
  check(err, "Error parsing json")
  // Get the blacklist from the JSON object
  blacklist := []string{}
  blacklistArray := jsonParsed.Path("exchange").Path("pair_blacklist").Children()
  for _, child := range blacklistArray {
    blacklist = append(blacklist, child.Data().(string))
  }
  return blacklist
}

func updateBacktestConfig(args Args, pairsList []string) {
	// Update the backtest config file
	// fmt.Println("Updating backtest config file...")
	// Open the testConfig file
	testConfig := args.path + "/user_data/config-whitelist-builder.json"
	testConfigFile, err := os.Open(testConfig)
	check(err, "Error opening test config file")
	defer testConfigFile.Close()
	// Read the testConfig file
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(testConfigFile)
	check(err, "Error reading test config file")
	// Convert the buffer to a JSON object
	jsonParsed, err := gabs.ParseJSON(buf.Bytes())
	check(err, "Error parsing json")
	// Update the JSON object
  jsonParsed.Set(args.exchange, "exchange", "name")
	jsonParsed.Set(pairsList, "exchange", "pair_whitelist")
  jsonParsed.Set(args.numPairs, "max_open_trades")
  jsonParsed.Set("unlimited", "stake_amount")
  jsonParsed.Set(args.numPairs * 100, "dry_run_wallet")
  jsonParsed.Set(args.tradingMode, "trading_mode")
  marginMode := ""
  if args.tradingMode == "futures" {
    marginMode = "isolated"
  }
  jsonParsed.Set(marginMode, "margin_mode")
  jsonParsed.Set(args.blacklist, "exchange", "pair_blacklist")

	// Update the buffer with the new JSON object
	buf = new(bytes.Buffer)
	_, err = buf.Write(jsonParsed.Bytes())
	check(err, "Error writing to buffer")

	// Write the JSON object back to the file
	testConfigFile, err = os.Create(testConfig)
	check(err, "Error creating test config file")

	_, err = testConfigFile.Write(buf.Bytes())
	check(err, "Error writing to test config file")

}
