package Manager

type Client interface {
	Start() bool
	SubscribeMarketData([]string)
	ReqQryInvestorPosition()
}
