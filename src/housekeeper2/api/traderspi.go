package api

type TraderSpi interface {
	OnConnect()
	OnRspUserLogin(Response, RespInfo, int, bool)
	OnRspQryInvestorPosition(Response, RespInfo, int, bool)
	OnRspQryTradingAccount(Response, RespInfo, int, bool)
	OnRtnTrade(Response)
	OnRspUserLogout()
}

func NewTraderSpi(spi TraderSpi, _type string) interface{} {
	if _type == ACCOUNT_TYPE_CTP {
		return NewCtpTraderSpi(spi)
	}
	return nil
}
