package api

import (
	"log"

	"github.com/qerio/goctp"
)

type CtpMdSpi struct {
	Spi MdSpi
}

func NewCtpMdSpi(spi MdSpi) interface{} {
	return goctp.NewDirectorCThostFtdcMdSpi(&CtpMdSpi{Spi: spi})
}

func (ctpspi *CtpMdSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {
	// 如果ErrorID != 0, 说明收到了错误的响应
	bResult := (pRspInfo.GetErrorID() != 0)
	if bResult {
		log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), pRspInfo.GetErrorMsg())
	}
	return bResult
}

func (ctpspi *CtpMdSpi) OnFrontConnected() {
	log.Println("CtpMdSpi.OnFrontConnected")
	ctpspi.Spi.OnConnect()
}

func (ctpspi *CtpMdSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

	log.Println("CtpMdSpi.OnRspUserLogin.")
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

func (ctpspi *CtpMdSpi) OnRspSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
	log.Printf("GoCThostFtdcMdSpi.OnRspSubMarketData: %#v %#v %#v\n", pSpecificInstrument.GetInstrumentID(), nRequestID, bIsLast)
}

func (ctpspi *CtpMdSpi) OnRtnDepthMarketData(pDepthMarketData goctp.CThostFtdcDepthMarketDataField) {

	quotedata := NewQuoter()

	quotedata.SetTradingDay(pDepthMarketData.GetTradingDay())
	quotedata.SetInstrumentID(pDepthMarketData.GetInstrumentID())
	quotedata.SetExchangeId(pDepthMarketData.GetExchangeID())
	quotedata.SetExchangeInstId(pDepthMarketData.GetExchangeInstID())
	quotedata.SetLastPrice(pDepthMarketData.GetLastPrice())
	quotedata.SetPreSettlementPrice(pDepthMarketData.GetPreSettlementPrice())
	quotedata.SetPreClosePrice(pDepthMarketData.GetPreClosePrice())
	quotedata.SetPreOpenInterest(pDepthMarketData.GetPreOpenInterest())
	quotedata.SetOpenPrice(pDepthMarketData.GetOpenPrice())
	quotedata.SetHighestPrice(pDepthMarketData.GetHighestPrice())
	quotedata.SetLowestPrice(pDepthMarketData.GetLowestPrice())
	quotedata.SetVolume(pDepthMarketData.GetVolume())

	ctpspi.Spi.OnRtnDepthMarketData(quotedata)

}
