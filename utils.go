package main

import "fmt"
import "os"

func printUsage() {
  // Print the usage message
  fmt.Println("Usage: whitelist_builder -s <strategy> [-p <path>] [-t <timerange>] [-c <config>] [-n <numPairs>] [-e <exchange>] [-tm <tradingMode>] [-nd]")
  fmt.Println("  -s, -strategy: Specify the name of the strategy to backtest")
  fmt.Println("  -p, -path: Specify the path to your working freqtrade directory (Default: current directory)")
  fmt.Println("  -t, -timerange: Specify the timerange to backtest (Default: last 30 days)")
  fmt.Println("  -tf, -timeframe: Specify the timeframe for your strategy (Default: 5m)")
  fmt.Println("  -c, -config: Specify the name of the config file to duplicate for testing, and pull the blacklist from (Default: config.json)")
  fmt.Println("  -n, -num-pairs: Specify the number of pairs to return (Default: 50)")
  fmt.Println("  -e, -exchange: Specify the exchange to pull the pairs from (Default: binance)")
  fmt.Println("  -tm, -trading-mode: Specify the trading mode (Default: spot)")
  fmt.Println("  -nd, -no-download: Do not download historical data (Default: false)")
  fmt.Println("  -ip, -inform-pairs: Specify the informative pairs for your strategy (\"BTC/USDT ETH/USDT\")")
  fmt.Println("  -it, -inform-timeframes: Specify the informative pairs timeframes for your strategy (\"1h 4h\")")
  fmt.Println("\n Example: whitelist_builder -s RSI -p /home/user/freqtrade -t 20210101-20210131 -c config.json -n 100 -e binance -tm spot")
  fmt.Print("\n")
  os.Exit(2)
}

func check(err error, message string) {
	if err != nil {
		fmt.Println(message)
		os.Exit(2)
	}
}

func cleanupOldBacktestResults(args Args) {
	// Delete everything in the backtest_results directory.

	// Open the backtest_results directory
	backtestResultsDir := args.path + "/user_data/backtest_results"
	dir, err := os.Open(backtestResultsDir)
	check(err, "Error opening backtest results directory")

	// Get the list of files in the directory
	files, err := dir.Readdir(0)
	check(err, "Error reading backtest results directory")

	// Delete each file in the directory
	for _, file := range files {
		err = os.Remove(backtestResultsDir + "/" + file.Name())
		check(err, "Error deleting file")
	}
}
