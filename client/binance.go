package client

import (
	"context"
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
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(apiSecret),
	}

	ctx, _ := context.WithCancel(context.Background())
	// use second return value for cancelling request when shutting down the app

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

func (bc *BinanceClient) GetOrderBook(currency string, limit int) (OrderBook, error) {
	orderBook := OrderBook{}

	req := binance.OrderBookRequest{Symbol: currency, Limit: limit}
	res, err := bc.client.OrderBook(req)

	if err != nil {
		return OrderBook{}, err
	}

	for _, bid := range res.Bids {
		orderBook.Bids = append(orderBook.Bids, Order(*bid))
	}

	for _, ask := range res.Asks {
		orderBook.Asks = append(orderBook.Asks, Order(*ask))
	}

	return orderBook, nil
}
