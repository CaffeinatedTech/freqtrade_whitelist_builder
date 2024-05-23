package main

import "bytes"
import "fmt"
import "os"
import "os/exec"
import "strings"
import "github.com/Jeffail/gabs/v2"

// Fetch the market pairs from the exchange, and return them as a JSON string
func getMarketPairs(args Args) []string {
	// Get the market pairs from the exchange
	fmt.Print("Fetching market pairs...")
	cmd := exec.Command(args.path+"/.venv/bin/python", args.path+"/freqtrade", "list-pairs", "--exchange", args.exchange,
		"--trading-mode", args.tradingMode, "--quote", "USDT", "--print-json", "--userdir", args.path+"/user_data")
	var execOut bytes.Buffer
	var execErr bytes.Buffer
	cmd.Stdout = &execOut
	cmd.Stderr = &execErr
	exitCode := cmd.Run()
	if exitCode != nil {
		fmt.Println("Error fetching market pairs")
		fmt.Print(execErr.String())
		os.Exit(2)
	}

	// convert the string to a JSON object
	jsonParsed, err := gabs.ParseJSON([]byte("{\"array\":" + execOut.String() + "}"))
	check(err, "Error parsing json")

	// JSON array to list
	var pairsList []string
	for _, child := range jsonParsed.S("array").Children() {
		pairsList = append(pairsList, child.Data().(string))
		// pairsList.PushBack(child.Data().(string))
	}

	// Print the list length
	fmt.Print("\r\033[K\r")
	fmt.Printf("Exchange %s has %d %s pairs\n", args.exchange, len(pairsList), strings.ToUpper(args.tradingMode))

	return pairsList
}

func downloadData(args Args) {
	// Download the historical data

	// We need to run a python script inside of its venv freqtrade download-data
	// We need to pass the path, timerange, and config file

	fmt.Print("Downloading historical data...")

	cmd := exec.Command(args.path+"/.venv/bin/python", args.path+"/freqtrade", "download-data", "--userdir", args.path+"/user_data",
		"--timerange", args.downloadTimerange, "-c", args.path + "/user_data/config-whitelist-builder.json", "--timeframes", args.timeframe)
  var execOut bytes.Buffer
  var execErr bytes.Buffer
	cmd.Stdout = &execOut
	cmd.Stderr = &execErr

	exitCode := cmd.Run()
	if exitCode != nil {
		fmt.Println("Error downloading historical data")
    fmt.Print(execErr.String())
		os.Exit(2)
	}

  if args.informPairs != "" && args.informTimeframes != "" {
    timeframesArray := strings.Split(args.informTimeframes, " ")
    fmt.Print("\r\033[K\r")
    fmt.Print("Downloading informative pairs historical data...")
    cmdArgs := []string{args.path+"/freqtrade", "download-data", "--userdir", args.path+"/user_data", 
      "--timerange", args.downloadTimerange, "-c", args.path + "/user_data/config-whitelist-builder.json", "--pairs", args.informPairs, "--timeframes", args.informTimeframes}
    cmdArgs = append(cmdArgs, timeframesArray...)
    cmd := exec.Command(args.path+"/.venv/bin/python", cmdArgs...)
    var execOut bytes.Buffer
    var execErr bytes.Buffer
    cmd.Stdout = &execOut
    cmd.Stderr = &execErr
    exitCode := cmd.Run()
    if exitCode != nil {
      fmt.Println("Error downloading informative pairs historical data")
      fmt.Print(execErr.String())
      os.Exit(2)
    }
  }
}
