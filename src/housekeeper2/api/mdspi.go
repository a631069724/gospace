package api

type MdSpi interface {
	OnConnect()
	OnRspUserLogin(Response, RespInfo, int, bool)
	OnRtnDepthMarketData(Quoter)
}

func NewMdSpi(spi MdSpi, _type string) interface{} {
	if _type == ACCOUNT_TYPE_CTP {
		return NewCtpMdSpi(spi)
	}
	return nil
}
