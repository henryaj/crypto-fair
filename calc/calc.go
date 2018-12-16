package calc

import (
	"github.com/henryaj/crypto-fair/client"
)

type ExchangeClient interface {
	Connect() error
	GetOrderbook(fromCurrency, toCurrency string, limit int) (client.OrderBook, error)
}

func GetFair(from, to string, client ExchangeClient) (float64, error) {
	book, err := client.GetOrderbook(from, to, 50)
	if err != nil {
		return 0, err
	}

	highestAsk, lowestBid, err := findHighestAndLowest(book)
	if err != nil {
		return 0, err
	}

	return (highestAsk + lowestBid) / 2, nil
}

func GetFairVia(from, via, to string, client ExchangeClient) (float64, error) {
	firstFair, err := GetFair(from, via, client)
	if err != nil {
		return 0, err
	}

	secondFair, err := GetFair(via, to, client)
	if err != nil {
		return 0, err
	}

	return firstFair * secondFair, nil
}

func findHighestAndLowest(book client.OrderBook) (float64, float64, error) {
	var highest float64

	for _, bid := range book.Bids {
		if bid.Price > highest {
			highest = bid.Price
		}
	}

	lowest := book.Asks[0].Price
	for _, ask := range book.Asks {
		if ask.Price < lowest {
			lowest = ask.Price
		}
	}

	return highest, lowest, nil // TODO err
}
