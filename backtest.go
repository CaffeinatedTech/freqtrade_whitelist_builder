package main

import "bytes"
import "fmt"
import "os"
import "os/exec"
import "github.com/Jeffail/gabs/v2"

func backtest(args Args) {
	// Run the backtest
	// fmt.Println("Running backtest...")
	cmd := exec.Command(args.path+"/.venv/bin/python", args.path+"/freqtrade", "backtesting", "--config", args.path+"/user_data/config-whitelist-builder.json",
		"--strategy", args.strategy, "--userdir", args.path+"/user_data", "--timerange", args.timerange, "--timeframe", args.timeframe)
	execOut := new(bytes.Buffer)
	execErr := new(bytes.Buffer)
	cmd.Stdout = execOut
	cmd.Stderr = execErr
	exitCode := cmd.Run()
	if exitCode != nil {
		fmt.Println("Error running backtest")
		fmt.Println(execErr.String())
		os.Exit(2)
	}
}

func parseBacktestResults(args Args) []Pair {
	// Parse the backtest results and return all pairs that are profitable
	lastResultFilePath := args.path + "/user_data/backtest_results/.last_result.json"
	lastResultFile, err := os.Open(lastResultFilePath)
	check(err, "Error opening last result file")
	// The result file that we want to open is the only text inside lastResultFile
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(lastResultFile)
	check(err, "Error reading last result file")

	// Get the result path from the buf.  It is in JSON format in node "latest_backtest"
	jsonParsed, err := gabs.ParseJSON(buf.Bytes())
	check(err, "Error parsing json")
	resultPath := jsonParsed.Path("latest_backtest").Data().(string)
	// Close the lastResultFile and release teh jsonParsed object
	jsonParsed = nil
	lastResultFile.Close()

	resultFilePath := args.path + "/user_data/backtest_results/" + resultPath
	resultFile, err := os.Open(resultFilePath)
	check(err, "Error opening result file")
	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(resultFile)
	check(err, "Error reading result file")

	// Load the JSON from resultFile
	jsonParsedResults, err := gabs.ParseJSON(buf.Bytes())
	check(err, "Error parsing json")

	// The results are stored in ["strategy"][args.strategy]["results_per_pair"]
	results := jsonParsedResults.Path("strategy").Path(args.strategy).Path("results_per_pair")
	profitablePairs := []Pair{}
	for _, child := range results.Children() {
		if child.Path("key").Data().(string) == "TOTAL" {
			break
		}
		if child.Path("profit_sum_pct").Data().(float64) <= 0.0 {
			break
		}
		if child.Path("profit_sum_pct").Data().(float64) > 0 {
			profitablePairs = append(profitablePairs, Pair{pair: child.Path("key").Data().(string), profit: child.Path("profit_sum_pct").Data().(float64)})
		}
	}

	return profitablePairs
}
