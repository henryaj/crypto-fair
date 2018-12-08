package main

import (
	"os"

	"github.com/henryaj/crypto-fair/clients"
)

type ExchangeClient interface {
	Connect() error
	GetOrderbook(currency string) []Order
	GetTrades(currency string) []Trade
}

type Order string
type Trade string

// SPEC: must be scalable for other currencies and exchanges
// but only flesh out getting fair value for XRP at Binance

// ---
// what currency has been selected?
// get data from exchange for our currency of choice
// calculate fair value
// output fair value in some machine-readable format
func main() {
	currency := os.Args[0]

	binanceClient := clients.NewBinanceClient()

	err := binanceClient.Connect()
	if err != nil {
		panic(err.Error())
	}

	_ = binanceClient.GetOrderbook(currency)
	_ = binanceClient.GetTrades(currency)

	// do some fancy calculations here
}
