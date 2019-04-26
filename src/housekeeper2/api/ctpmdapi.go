package api

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. 
#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import (
	"log"
	"unsafe"

	"github.com/qerio/goctp"
)

type CtpMdApi struct {
	ctpmdapi  goctp.CThostFtdcMdApi
}

func NewCtpMdApi() MdApi {
	return &CtpMdApi{
		ctpmdapi: goctp.CThostFtdcMdApiCreateFtdcMdApi(),
	}

}


func (ctp *CtpMdApi) RegistSpi(spi interface{}) {
	ctp.ctpmdapi.RegisterSpi(spi.(goctp.CThostFtdcMdSpi))
}


func (ctp *CtpMdApi) Connect(ip string, port string) {
	
	ctp.ctpmdapi.RegisterFront("tcp://" + ip + ":" + port)
	ctp.ctpmdapi.Init()
}

func (ctp *CtpMdApi) ReqUserLogin(reqdata Request) {
	log.Println("CtpMdApi.ReqUserLogin.")
	
	req := goctp.NewCThostFtdcReqUserLoginField()
	req.SetBrokerID(reqdata.BrokeId())
	req.SetUserID(reqdata.UserName())
	req.SetPassword(reqdata.Pwd())

	iResult := ctp.ctpmdapi.ReqUserLogin(req, reqdata.ReqId())

	if iResult != 0 {
		log.Println("发送用户登录请求: 失败.")
	} else {
		log.Println("发送用户登录请求: 成功.")
	}
}

func (ctp *CtpMdApi) SubscribeMarketData(name []string) {

	log.Println("CtpMdApi.SubscribeMarketData")

	arg := make([](*C.char), 0)
	l := len(name)
	for i, _ := range name {
		cchar := C.CString(name[i])
		defer C.free(unsafe.Pointer(cchar)) //释放内存
		strptr := (*C.char)(unsafe.Pointer(cchar))
		arg = append(arg, strptr) //将char指针加入到arg切片
	}

	iResult := ctp.ctpmdapi.SubscribeMarketData((*string)(unsafe.Pointer(&arg[0])), int(l))

	if iResult != 0 {
		log.Println("发送行情订阅请求: 失败.")
	} else {
		log.Println("发送行情订阅请求: 成功.")
	}
}