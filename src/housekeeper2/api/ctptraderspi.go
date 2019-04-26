package api

import (
	"log"
	"reflect"

	"github.com/qerio/goctp"
)

type CtpTraderSpi struct {
	Spi TraderSpi
}

func NewCtpTraderSpi(spi TraderSpi) interface{} {
	return goctp.NewDirectorCThostFtdcTraderSpi(&CtpTraderSpi{Spi: spi})
}

func (ctpspi *CtpTraderSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {
	// 如果ErrorID != 0, 说明收到了错误的响应
	bResult := (pRspInfo.GetErrorID() != 0)
	if bResult {
		log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
	}
	return bResult
}

func (ctpspi *CtpTraderSpi) OnFrontConnected() {
	log.Println("ctptraderspi.OnFrontConnected")
	ctpspi.Spi.OnConnect()
}

func (ctpspi *CtpTraderSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("CtpTraderSpi.OnRspUserLogin.")
	if bIsLast && !ctpspi.IsErrorRspInfo(pRspInfo) {

		//log.Printf("获取当前交易日  : %#v\n", p.Client.TraderApi.GetTradingDay())
		log.Printf("获取用户登录信息: %#v %#v %#v\n", pRspUserLogin.GetFrontID(), pRspUserLogin.GetSessionID(), pRspUserLogin.GetMaxOrderRef())

		// // 保存会话参数
		// FRONT_ID = pRspUserLogin->FrontID;
		// SESSION_ID = pRspUserLogin->SessionID;
		// int iNextOrderRef = atoi(pRspUserLogin->MaxOrderRef);
		// iNextOrderRef++;
		// sprintf(ORDER_REF, "%d", iNextOrderRef);
		// sprintf(EXECORDER_REF, "%d", 1);
		// sprintf(FORQUOTE_REF, "%d", 1);
		// sprintf(QUOTE_REF, "%d", 1);
		// ///获取当前交易日
		// cerr << "获取当前交易日 = " << pMdApi->GetTradingDay() << endl;
		// ///投资者结算结果确认
		//p.ReqSettlementInfoConfirm()
		respdata := NewRespData()
		ctpspi.Spi.OnRspUserLogin(respdata, pRspInfo, nRequestID, bIsLast)

	} else {
		log.Printf("errid=%s msg=%s\n", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
	}
}

func (ctpspi *CtpTraderSpi) OnRspQryInvestorPosition(pInvestorPosition goctp.CThostFtdcInvestorPositionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("CtpTraderSpi.OnRspQryInvestorPosition.", bIsLast)
	v := reflect.ValueOf(pInvestorPosition)
	if !v.IsValid() {
		return
	}

	respdata := NewRespData()

	respdata.SetCloseProfit(pInvestorPosition.GetCloseProfit())

	if pInvestorPosition.GetPosiDirection() == 0x32 {
		respdata.SetPosiDirection(0x31)
	} else {
		respdata.SetPosiDirection(0x30)
	}
	respdata.SetInstrumentID(pInvestorPosition.GetInstrumentID())
	respdata.SetPosition(pInvestorPosition.GetPosition())
	respdata.SetYdPosition(pInvestorPosition.GetYdPosition())
	respdata.SetOpenAmount(pInvestorPosition.GetOpenAmount())
	ctpspi.Spi.OnRspQryInvestorPosition(respdata, pRspInfo, nRequestID, bIsLast)

}

func (ctpspi *CtpTraderSpi) OnRspQryTradingAccount(pTradingAccount goctp.CThostFtdcTradingAccountField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("CtpTraderSpi.OnRspQryTradingAccount.")
	respdata := NewRespData()
	respdata.SetBalance(pTradingAccount.GetAvailable() + pTradingAccount.GetCurrMargin())
	respdata.SetAccountMargin(pTradingAccount.GetCurrMargin())
	ctpspi.Spi.OnRspQryTradingAccount(respdata, pRspInfo, nRequestID, bIsLast)

}

func (ctpspi *CtpTraderSpi) OnRtnOrder(pOrder goctp.CThostFtdcOrderField) {
	//log.Println("CtpTraderSpi.OnRtnOrder.")
}

func (ctpspi *CtpTraderSpi) OnRtnTrade(pTrade goctp.CThostFtdcTradeField) {
	log.Println("CtpTraderSpi.OnRtnTrade.")
	respdata := NewRespData()

	respdata.SetInstrumentID(pTrade.GetInstrumentID())
	respdata.SetOpenOrClose(pTrade.GetOffsetFlag())
	if pTrade.GetOffsetFlag() == '0' { //开仓
		if pTrade.GetDirection() == '0' { //多
			respdata.SetPosiDirection('1')
		} else { //空
			respdata.SetPosiDirection('0')
		}
	} else { //平仓
		if pTrade.GetDirection() == '1' { //多
			respdata.SetPosiDirection('1')
		} else { //空
			respdata.SetPosiDirection('0')
		}
	}
	respdata.SetPosition(pTrade.GetVolume())
	ctpspi.Spi.OnRtnTrade(respdata)
}

func (ctpspi *CtpTraderSpi) OnRspOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("CtpTraderSpi.OnRspOrderInsert.")
	log.Println("错误信息 : ", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())

}

func (ctpspi *CtpTraderSpi) OnErrRtnOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, pRspInfo goctp.CThostFtdcRspInfoField) {
	log.Println("CtpTraderSpi.OnErrRtnOrderInsert.")
	log.Println("错误信息 : ", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())

}

func (ctpspi *CtpTraderSpi) OnRspUserLogout(pUserLogout goctp.CThostFtdcUserLogoutField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Println("CtpTraderSpi.OnRspUserLogout.", pUserLogout.GetUserID())

	ctpspi.Spi.OnRspUserLogout()
}
