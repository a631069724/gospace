package dump

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/liangdas/mqant/log"
)

func httpDump(moduletype string, reqdump []byte, respdump []byte) (string, string) {
	req, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(reqdump)))
	reqBody, _ := ioutil.ReadAll(req.Body)
	resp, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(respdump)))
	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Info("httpDump:reqBody[%s]", string(reqBody))
	log.Info("httpDump:respBody[%s]", string(respBody))

	var reqData, respData interface{}
	var dumptableReq, dumptableResp string
	trantype := req.URL.Path
	if trantype == "/pay/consume" {
		reqData = new(ConsumeReq)
		respData = new(ConsumeResp)
		dumptableReq = "pay_consume_req"
		dumptableResp = "pay_consume_resp"
	} else if trantype == "/pay/notify/consume" {
		reqData = new(ConsumeNotifyReq)
		respData = new(ConsumeNotifyResp)
		dumptableReq = "pay_notify_consume_req"
		dumptableResp = "pay_notify_consume_resp"
	} else if trantype == "/pay/reverse" {
		reqData = new(ReverseReq)
		respData = new(ReverseResp)
		dumptableReq = "pay_reverse_req"
		dumptableResp = "pay_reverse_resp"
	} else if trantype == "/pay/cancel" {
		reqData = new(CancelReq)
		respData = new(CancelResp)
		dumptableReq = "pay_cancel_req"
		dumptableResp = "pay_cancel_resp"
	} else if trantype == "/pay/notify/cancel" {
		reqData = new(CancelNotifyReq)
		respData = new(CancelNotifyResp)
		dumptableReq = "pay_notify_cancel_req"
		dumptableResp = "pay_notify_cancel_resp"
	} else if trantype == "/pay/return" {
		reqData = new(ReturnReq)
		respData = new(ReturnResp)
		dumptableReq = "pay_return_req"
		dumptableResp = "pay_return_resp"
	} else if trantype == "/pay/notify/return" {
		reqData = new(ReturnNotifyReq)
		respData = new(ReturnNotifyResp)
		dumptableReq = "pay_notify_return_req"
		dumptableResp = "pay_notify_return_resp"
	} else if trantype == "/pay/query" {
		reqData = new(QueryReq)
		respData = new(QueryResp)
		dumptableReq = "pay_notify_query_req"
		dumptableResp = "pay_notify_query_resp"
	} else {
		log.Info("undefined [%s] dump", trantype)
		return "", ""
	}
	err := json.Unmarshal(reqBody, reqData)

	if err != nil {
		log.Error("Unmarshal ReqData[%s] error:%s", reqBody, err.Error())
		return string(reqBody), err.Error()
	}
	if err := Storer.DataStore(dumptableReq, reqData); err != nil {
		log.Error("oracle store data err:%s", err.Error())
	}
	err = json.Unmarshal(respBody, respData)
	if err != nil {
		log.Error("Unmarshal ReqData[%s] error:%s", respBody, err.Error())
		return string(respBody), err.Error()
	}
	if err := Storer.DataStore(dumptableResp, respData); err != nil {
		log.Error("oracle store data err:%s", err.Error())
	}
	return "", ""
}
