package clients

type BinanceClient struct{}

func NewBinanceClient() BinanceClient {
	return BinanceClient{}
}

func (b *BinanceClient) Connect() error {
	return nil
}

func (b *BinanceClient) GetOrderbook(currency string) []string {
	return nil
}

func (b *BinanceClient) GetTrades(currency string) []string {
	return nil
}
