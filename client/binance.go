package client

import (
	"context"
	"fmt"
	"os"

	binance "github.com/binance-exchange/go-binance"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type BinanceClient struct {
	client binance.Binance
}

type OrderBook struct {
	Bids []Order
	Asks []Order
}

type Order struct {
	Price    float64
	Quantity float64
}

func NewBinanceClient(apiKey, apiSecret string) *BinanceClient {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowWarn())
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(apiSecret),
	}

	ctx, _ := context.WithCancel(context.Background())

	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		apiKey,
		hmacSigner,
		logger,
		ctx,
	)
	b := binance.NewBinance(binanceService)

	return &BinanceClient{client: b}
}

func (bc *BinanceClient) Connect() error {
	return bc.client.Ping()
}

func (bc *BinanceClient) GetOrderbook(fromCurrency, toCurrency string, limit int) (OrderBook, error) {
	orderBook := OrderBook{}
	currency := fromCurrency + toCurrency

	req := binance.OrderBookRequest{Symbol: currency, Limit: limit}
	res, err := bc.client.OrderBook(req)

	if err != nil {
		return OrderBook{}, err
	}

	if len(res.Bids) == 0 || len(res.Asks) == 0 {
		// market doesn't exist on Binance, or there was a problem getting it
		// try a market with the symbols reversed, then reverse the orderbook
		return bc.getReversedOrderbook(fromCurrency, toCurrency, limit)
	}

	for _, bid := range res.Bids {
		orderBook.Bids = append(orderBook.Bids, Order(*bid))
	}

	for _, ask := range res.Asks {
		orderBook.Asks = append(orderBook.Asks, Order(*ask))
	}

	return orderBook, nil
}

func (bc *BinanceClient) getReversedOrderbook(fromCurrency, toCurrency string, limit int) (OrderBook, error) {
	orderBook := OrderBook{}

	// NOTE reversed currency symbols
	currency := toCurrency + fromCurrency

	req := binance.OrderBookRequest{Symbol: currency, Limit: limit}
	res, err := bc.client.OrderBook(req)

	if err != nil {
		return OrderBook{}, err
	}

	if len(res.Bids) == 0 || len(res.Asks) == 0 {
		return OrderBook{}, fmt.Errorf("Unable to get orderbook for %s", currency)
	}

	for _, bid := range res.Bids {
		orderBook.Asks = append(orderBook.Asks, Order(*bid))
	}

	for _, ask := range res.Asks {
		orderBook.Bids = append(orderBook.Bids, Order(*ask))
	}

	return orderBook, nil
}
