package main

import "flag"
import "fmt"
import "os"
import "sort"
import "strings"
import "time"
import "github.com/fatih/color"
import "github.com/Jeffail/gabs/v2"

// This project builds pair whitelists for your freqtrade strategies.
// You will provide the following:
// - The path to your working freqtrade install
// - The timerange to backtest
// - The name of the config file to pull the exchange and blacklist from
// - The name of the strategy to backtest

type Args struct {
	config            string
	downloadTimerange string
	exchange          string
	noDownload        bool
	numPairs          int
	path              string
	strategy          string
  timeframe         string
	timerange         string
	tradingMode       string
  blacklist         []string
  informPairs       string
  informTimeframes  string
}

type Pair struct {
	pair   string
	profit float64
}

func getArgs() Args {
	args := Args{}

	// The default timerange will be the last 30 days from today.  Dates must be in YYYYMMDD format.
	dateToday := time.Now().Format("20060102")
	sixtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("20060102")
	defaultTimerange := sixtyDaysAgo + "-" + dateToday

	flag.StringVar(&args.path, "p", "./", "Specify the path to your working freqtrade directory")
	flag.StringVar(&args.path, "path", "./", "Specify the path to your working freqtrade directory")
	flag.StringVar(&args.timerange, "t", defaultTimerange, "Specify the timerange to backtest")
	flag.StringVar(&args.timerange, "timerange", defaultTimerange, "Specify the timerange to backtest")
  flag.StringVar(&args.timeframe, "tf", "5m", "Specify the timeframe that your strategy uses")
  flag.StringVar(&args.timeframe, "timeframe", "5m", "Specify the timeframe that your strategy uses")
	flag.StringVar(&args.config, "c", "config.json", `Specify the name of the config file to pull the 
    exchange and blacklist from, Defaults to config.json`)
	flag.StringVar(&args.config, "config", "config.json", `Specify the name of the config file to pull 
    the exchange and blacklist from, Defaults to config.json`)
	flag.StringVar(&args.strategy, "s", "", "Specify the name of the strategy to backtest")
	flag.StringVar(&args.strategy, "strategy", "", "Specify the name of the strategy to backtest")
	flag.IntVar(&args.numPairs, "n", 50, "Specify the number of pairs to return")
	flag.IntVar(&args.numPairs, "num-pairs", 50, "Specify the number of pairs to return")
	flag.BoolVar(&args.noDownload, "nd", false, "Do not download historical data")
	flag.BoolVar(&args.noDownload, "no-download", false, "Do not download historical data")
	flag.StringVar(&args.exchange, "e", "binance", "Specify the exchange to pull the pairs from")
	flag.StringVar(&args.exchange, "exchange", "binance", "Specify the exchange to pull the pairs from")
	flag.StringVar(&args.tradingMode, "tm", "spot", "Specify the trading mode")
	flag.StringVar(&args.tradingMode, "trading-mode", "spot", "Specify the trading mode")
  flag.StringVar(&args.informPairs, "ip", "", "Specify the informative pairs for your strategy (\"BTC/USDT ETH/USDT\")")
  flag.StringVar(&args.informPairs, "inform-pairs", "", "Specify the informative pairs for your strategy (\"BTC/USDT ETH/USDT\")")
  flag.StringVar(&args.informTimeframes, "it", "", "Specify the informative pairs timeframes for your strategy (\"1h 4h\")")
  flag.StringVar(&args.informTimeframes, "inform-timeframes", "", "Specify the informative pairs timeframes for your strategy (\"1h 4h\")")
	flag.Parse()

	validateArgs(args)

	// Set timerange start to 30 days beforehand.
	startDateString := strings.Split(args.timerange, "-")[0]
	parsedStartDate, _ := time.Parse("20060102", startDateString)
	startDate := parsedStartDate.AddDate(0, 0, -30).Format("20060102")
	args.downloadTimerange = startDate + "-" + strings.Split(args.timerange, "-")[1]

	// Convert exchange to lowercase
	args.exchange = strings.ToLower(args.exchange)
	// Convert tradingMode to lowercase
	args.tradingMode = strings.ToLower(args.tradingMode)

	return args
}

func validateArgs(args Args) {
	// Check for required arguments
	required := []string{"strategy"}
	for _, req := range required {
		if flag.Lookup(req).Value.String() == "" {
      printUsage()
			os.Exit(0)
		}
	}
	// Validate path
	if _, err := os.Stat(args.path); os.IsNotExist(err) {
		fmt.Println("Path does not exist:", args.path)
		os.Exit(2)
	}
	// Make sure freqtrade is at that path
	if _, err := os.Stat(args.path + "/freqtrade/freqtradebot.py"); os.IsNotExist(err) {
		fmt.Println("freqtrade not found at path:", args.path)
    printUsage()
		os.Exit(2)
	}
	// Check if the config file exists
	if _, err := os.Stat(args.path + "/user_data/" + args.config); os.IsNotExist(err) {
		fmt.Println("Config file not found:", args.path+"/user_data/config/"+args.config)
		os.Exit(2)
	}
	// Validate timerange
	startDateString := strings.Split(args.timerange, "-")[0]
	endDateString := strings.Split(args.timerange, "-")[1]
	_, startDateErr := time.Parse("20060102", startDateString)
	_, endDateErr := time.Parse("20060102", endDateString)
	if len(strings.Split(args.timerange, "-")) != 2 || startDateErr != nil || endDateErr != nil {
		fmt.Println("Invalid timerange:", args.timerange, ". Must be in the format YYYYMMDD-YYYYMMDD")
		os.Exit(2)
	}
	// Validate tradingMode
	if args.tradingMode != "spot" && args.tradingMode != "futures" {
		fmt.Println("Invalid trading mode:", args.tradingMode, ". Must be spot or futures")
		os.Exit(2)
	}
}

func main() {
	// Parse command line arguments.
	args := getArgs()

	green := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

	fmt.Printf("Calculating the top %s profitable pairs for strategy %s in timerange %s\n\n", green(args.numPairs), green(args.strategy), green(args.timerange))

	// Cleanup old backtest results
	cleanupOldBacktestResults(args)

	// We need to fetch the pairs list from the exchange, but the ccxt library is not available in Go.
	marketPairs := getMarketPairs(args)

	// Check and setup backtest config file
	createBacktestConfig(args)
  // Get the blacklist from the config file
  args.blacklist = getBlacklist(args)
	updateBacktestConfig(args, marketPairs)
	// Download the historical data
	if !args.noDownload {
		downloadData(args)
	}

	results := []Pair{}

	for len(marketPairs) > 0 {
		fmt.Print("\r\033[K\r")
		fmt.Printf("Running backtests... %d pairs remaining.", len(marketPairs))

		thisBatch := marketPairs[:min(args.numPairs, len(marketPairs))]
		marketPairs = marketPairs[min(args.numPairs, len(marketPairs)):]
		// Update the backtest config file.
		updateBacktestConfig(args, thisBatch)
		// Run the backtest
		backtest(args)
		// Parse the backtest results
		thisBatchResults := parseBacktestResults(args)
		results = append(results, thisBatchResults...)
	}

	// Sort the results by profit
	sort.Slice(results, func(i, j int) bool { return results[i].profit > results[j].profit })

	// Print the results
	fmt.Print("\r\033[K\r")
	fmt.Printf("Found %d profitable pairs.  Here are the top %d:\n\n", len(results), min(args.numPairs, len(results)))

	// Print the first 50 pairs of the results slice as a JSON array whitelist using gabs
	topPairs := results[:min(args.numPairs, len(results))]
	jsonParsed, err := gabs.ParseJSON([]byte("[]"))
	if err != nil {
		fmt.Println("Error parsing json")
		os.Exit(2)
	}
	for _, pair := range topPairs {
		jsonParsed.ArrayAppend(pair.pair)
	}

	fmt.Println(jsonParsed.String())

	fmt.Println("\nDone")
}
