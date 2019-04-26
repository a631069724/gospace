package api

import (
	"errors"
	"fmt"
)

type MdApi interface {
	RegistSpi(interface{})
	Connect(string, string)
	ReqUserLogin(Request)
	SubscribeMarketData([]string)
}

func NewMdApi(apitype string) (MdApi, error) {
	if apitype == "CTP" {
		return NewCtpMdApi(), nil
	}
	return nil, errors.New(fmt.Sprintf("The Account Type :%s Not Exist", apitype))
}
