package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/henryaj/crypto-fair/calc"
	"github.com/henryaj/crypto-fair/client"
)

func main() {
	binanceClient := client.NewBinanceClient(
		os.Getenv("BINANCE_API_KEY"),
		os.Getenv("BINANCE_SECRET_KEY"),
	)

	if err := binanceClient.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to exchange")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("~~~~~~ XRP:BTC fair price ~~~~~~")

	for {
		// calculate fairs
		// XRP:BTC
		var wg sync.WaitGroup
		var results []float64

		resultChan := make(chan float64, 1)

		wg.Add(3)

		go func() {
			xrpBtcFair, err := calc.GetFair("XRP", "BTC", binanceClient)
			if err != nil {
				panic(err)
			}

			resultChan <- xrpBtcFair
		}()

		// XRP:ETH, ETH:BTC
		go func() {
			xrpBtcFairViaEth, err := calc.GetFairVia("XRP", "ETH", "BTC", binanceClient)
			if err != nil {
				panic(err)
			}

			resultChan <- xrpBtcFairViaEth
		}()

		// XRP:BNB, BNB:BTC
		go func() {
			xrpBtcFairViaBnb, err := calc.GetFairVia("XRP", "BNB", "BTC", binanceClient)
			if err != nil {
				panic(err)
			}

			resultChan <- xrpBtcFairViaBnb
		}()

		go func() {
			for val := range resultChan {
				results = append(results, val)
				wg.Done()
			}
		}()

		wg.Wait()

		averageFair := (results[0] + results[1] + results[2]) / 3
		fmt.Printf("%f\n", averageFair)
	}
}
