## Freqtrade Whitelist Builder
A handy tool for freqtrade that will build a set of pairs for you
to use in your whitelist.  Using backtesting we can find the most
profitable coins in the past couple of weeks, and compile a list
ready for you to paste into your whitelist.

### How it works
Using freqtrade in the background, we can fetch all pairs from
your exchange, download the required amount of historical data,
and run backtesting on every pair.  We then build a list of
profitable pairs sorted in descending order, and supply the top
X amount.

### Requirements
You'll need a working and configured installation of freqtrade, complete
with a strategy.

Additionally if you will be building from source, you will need
a working install of Go version 1.21+

### Installation
You can either clone the repo, and build it yourself, or just
download a binary from the releases page.

**Binary Method:**  
* Download and unzip the appropriate binary from the releases page
* Optionally move the binary to your freqtrade directory

**Build Method:**  
* Clone the repo `git clone https://github.com/CaffeinatedTech/freqtrade_whitelist_builder.git`
* Change into the directory `cd whitelist_builder`
* Install requirements `go install .`
* Build `go build`
* Optinally move the binary to your freqtrade directory

### Usage
You'll need to tell whitelist_builder where your freqtrade is
installed if you haven't placed it in the same directory.

At minimum you should supply the strategy to use.  
`./whitelist_builder -s ttV8dca`

| argument | description | example |
|----------|-------------|---------|
| -p --path | The path to your working freqtrade directory if not the current directory | -p /home/user/freqtrade |
| -t --timerange | OPTIONAL - The timerange used for backtesting (defaults to the last 30 days) | -t 20230112-20240125 |
| -tf --timeframes | The timeframe used for backtesting - one only (defaults to 5m) | -tf 1h |
| -s --strategy | Your strategy | -s ttV8dca |
| -c --config | The name of the config file to pull the exchange and blacklist from (defaults to config.json) | -c config_bybit.json |
| -e --exchange | The exchange to use, otherwise it will be fetched from the current config.json | -f bybit |
| -n --num-pairs | The number of pairs to return for your whitelist (defaults to 50) | -n 75 |
| -ip --inform-pairs | The informative pairs used by your strategy or 'ALL' - it will also downlaod the data for these | -ip "BTC/USDT ETH/USDT" |
| -it --inform-timeframes | The timeframes for your informative pairs | -it "1h 4h" |
| -tm --trading-mode | Are you using SPOT or FUTURES? (defaults to SPOT) | -tm FUTURES |
| -nd --no-download | Skip downloading historical data - if you already have it | -nd |

**Examples:**  
* `./whitelist_builder -s EI3v2 -e binance -n 20` Using strategy EI3v2 on binance, fetch the top 20 pairs. (binary is located within the freqtrade directory)
* `./whitelist_builder -p /user/home/adam/freqtrade --strategy ttV8dca -tm FUTURES -t 20240101-20240401 -ip "BTC/USDT ETH/USDT" -it "1h 4h"` Get the top 50 FUTURES pairs from the timerange 20240101-20240401 for ttV8dca using informative pairs.
