package api

import (
	"housekeeper/adaptation"
)

type Client interface {
	StartTrader() bool
	StartMd()
	ReqQryInvestorPosition(int)
	ReqQryTradingAccount(int)
	GetTraderRequestID() int
	ReqOrderInsert(int)
	ClearPosition()
	SubscribeMarketData([]string)
	ClosePosition(adaptation.PositionField, int)
}
