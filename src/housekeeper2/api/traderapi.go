package api

import (
	"errors"
	"fmt"
)

const (
	ACCOUNT_TYPE_CTP string = "CTP"
)

type TraderApi interface {
	RegistSpi(interface{})
	Connect(string, string)
	ReqUserLogin(Request)
	ReqQryInvestorPosition(Request)
	ReqQryTradingAccount(Request)
	ReqOrderInsert(Request)
	ReqUserLogout(Request)
}

func NewTraderApi(_type string) (TraderApi, error) {
	if _type == ACCOUNT_TYPE_CTP {
		return NewCtpTraderApi(), nil
	}
	return nil, errors.New(fmt.Sprintf("The Account Type :%s Not Exist", _type))
}
