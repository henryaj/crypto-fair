package main

import (
	"fmt"
	"os"

	"github.com/henryaj/crypto-fair/clients"
	flags "github.com/jessevdk/go-flags"
)

type ExchangeClient interface {
	Connect() error
	GetOrderbook(currency string) client.OrderBook
	// GetTrades(currency string) []Trade
}

// get data from exchange for our currency of choice
// calculate fair value
// output fair value in some machine-readable format
func main() {
	var opts struct {
		// Slice of bool will append 'true' each time the option
		// is encountered (can be set multiple times, like -vvv)
		Currency string `short:"c" long:"currency" description:"Ticker name of asset to calculate fair value, e.g. 'XMRBTC'" required:"true"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	binanceClient := client.NewBinanceClient(
		os.Getenv("BINANCE_API_KEY"),
		os.Getenv("BINANCE_SECRET_KEY"),
	)

	if err := binanceClient.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to exchange")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	book, err := binanceClient.GetOrderBook(opts.Currency, 100)
	if err != nil {
		panic(err)
	}

	fmt.Println(book)
	// _ = binanceClient.GetTrades(currency)

	// do some fancy calculations here
}
