package api

import (
	"log"
	"time"

	"github.com/qerio/goctp"
)

type CtpTraderApi struct {
	ctptraderapi goctp.CThostFtdcTraderApi
}

func NewCtpTraderApi() TraderApi {
	return &CtpTraderApi{
		ctptraderapi: goctp.CThostFtdcTraderApiCreateFtdcTraderApi(),
	}

}

func (ctp *CtpTraderApi) Connect(ip string, port string) {
	ctp.ctptraderapi.SubscribePublicTopic(2 /*THOST_TERT_QUICK*/)  // 注册公有流
	ctp.ctptraderapi.SubscribePrivateTopic(2 /*THOST_TERT_QUICK*/) // 注册私有流
	ctp.ctptraderapi.RegisterFront("tcp://" + ip + ":" + port)
	ctp.ctptraderapi.Init()
}

func (ctp *CtpTraderApi) RegistSpi(spi interface{}) {
	ctp.ctptraderapi.RegisterSpi(spi.(goctp.CThostFtdcTraderSpi))
}

func (ctp *CtpTraderApi) ReqUserLogin(reqdata Request) {
	log.Println("CtpTraderApi.ReqUserLogin.")

	req := goctp.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(reqdata.BrokeId())
	req.SetUserID(reqdata.UserName())
	req.SetPassword(reqdata.Pwd())

	iResult := ctp.ctptraderapi.ReqUserLogin(req, reqdata.ReqId())

	if iResult != 0 {
		log.Println("发送用户登录请求: 失败.")
	} else {
		log.Println("发送用户登录请求: 成功.")
	}
}

func (ctp *CtpTraderApi) ReqUserLogout(reqdata Request) {
	log.Println("CtpTraderApi.ReqUserLogout.")

	req := goctp.NewCThostFtdcUserLogoutField()
	req.SetBrokerID(reqdata.BrokeId())
	req.SetUserID(reqdata.UserName())

	ctp.ctptraderapi.ReqUserLogout(req, reqdata.ReqId())
}

func (ctp *CtpTraderApi) IsFlowControl(iResult int) bool {
	return ((iResult == -2) || (iResult == -3))
}

func (ctp *CtpTraderApi) ReqQryInvestorPosition(reqdata Request) {
	log.Println("CtpTraderApi.ReqQryInvestorPosition.")
	req := goctp.NewCThostFtdcQryInvestorPositionField()
	req.SetBrokerID(reqdata.BrokeId())
	req.SetInvestorID(reqdata.UserName())
	req.SetInstrumentID(reqdata.InstrumentID())
	for {
		iResult := ctp.ctptraderapi.ReqQryInvestorPosition(req, reqdata.ReqId())

		if !ctp.IsFlowControl(iResult) {
			log.Printf("--->>> 请求查询投资者持仓: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询投资者持仓: 受到流控 %#v\n", iResult)
			time.Sleep(1 * time.Second)
		}
	}
}

func (ctp *CtpTraderApi) ReqQryTradingAccount(reqdata Request) {
	log.Println("CtpTraderApi.ReqQryTradingAccount.")
	req := goctp.NewCThostFtdcQryTradingAccountField()
	req.SetBrokerID(reqdata.BrokeId())
	req.SetInvestorID(reqdata.UserName())
	req.SetCurrencyID(reqdata.Cur())
	for {
		iResult := ctp.ctptraderapi.ReqQryTradingAccount(req, reqdata.ReqId())

		if !ctp.IsFlowControl(iResult) {
			log.Printf("--->>> 请求查询资金账户: 成功 %#v\n", iResult)
			break
		} else {
			log.Printf("--->>> 请求查询资金账户: 受到流控 %#v\n", iResult)
		}
		time.Sleep(1 * time.Second)
	}
}

func (ctp *CtpTraderApi) ReqOrderInsert(reqdata Request) {
	log.Println("CtpTraderApi.ReqOrderInsert.")
	req := goctp.NewCThostFtdcInputOrderField()
	req.SetBrokerID(reqdata.BrokeId())
	req.SetInvestorID(reqdata.UserName())
	req.SetInstrumentID(reqdata.InstrumentID())
	req.SetUserID(reqdata.UserName())
	req.SetOrderPriceType(reqdata.GetOrderPriceType())
	req.SetDirection(reqdata.GetDirection())
	req.SetCombOffsetFlag(reqdata.GetCombOffSetFlag())
	req.SetCombHedgeFlag(reqdata.GetCombHedgeFlag())
	req.SetVolumeTotalOriginal(reqdata.GetNum())
	req.SetVolumeCondition('1')
	req.SetMinVolume(1)
	req.SetContingentCondition('1')
	req.SetForceCloseReason('0')
	req.SetIsAutoSuspend(0)
	req.SetUserForceClose(0)
	req.SetTimeCondition('3')
	req.SetLimitPrice(reqdata.GetLimitPrice())
	log.Println(reqdata.BrokeId(),
		reqdata.UserName(),
		reqdata.InstrumentID(),
		reqdata.UserName(),
		reqdata.GetOrderPriceType(),
		reqdata.GetDirection(),
		reqdata.GetCombOffSetFlag(),
		reqdata.GetCombHedgeFlag(),
		reqdata.GetNum(),
		reqdata.GetLimitPrice(),
		reqdata.ReqId(),
	)
	ctp.ctptraderapi.ReqOrderInsert(req, reqdata.ReqId())

	/*
		req.SetLimitPrice(pst.HighestPrice)
	*/
}
