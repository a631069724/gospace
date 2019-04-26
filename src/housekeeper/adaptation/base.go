/*
 * @Description: 适配接口
 * @Author: gdd
 * @Date: 2018-11-07 09:52:18
 * @LastEditTime: 2018-11-28 14:03:57
 * @LastEditors: Please set LastEditors
 */
package adaptation

/**
 * @msg: 封装注册第三方服务接口
 * @param {type}
 * @return:
 */
type Api interface {
	GetName() string
	SetProc(Proc)
}

/**
 * @msg: 封装内部业务逻辑接口
 * @param {type}
 * @return:
 */
type Proc interface {
	ApiName() string
	OnRspQryInvestorPosition(*PositionField, bool)
	OnRspQryTradingAccount(*TradingAccount)
	SetSessionId(int)
	SetFrontId(int)
	OnRtnOrder()
	OnRtnTrade(*PositionField)
	OnRspOrderInsert(*PositionField, bool)
}

/**
 * @msg: 行情回调拆处理
 * @param {type}
 * @return:
 */
type QuoteSpi interface {
	OnRspUserLogin()
	OnRtnDepthMarketData(*Quote)
}
