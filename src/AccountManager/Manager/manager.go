package Manager

import (
	"AccountManager/ctptran"
	"fmt"
)

type MyManager struct {
	client   Client
	customer *Customer
}

func NewManager() *MyManager {
	client, err := ctptran.NewCtpClient("/home/gengzhi/gospace/src/AccountManager/Manager/confg.ini")
	if err != nil {
		fmt.Println("New Client Error", err.Error())
		return nil
	}
	return &MyManager{
		client:   client,
		customer: NewCustomer(),
	}
}

func (m *MyManager) Run() {
	c := make(chan bool)
	if !m.client.Start() {
		fmt.Println("start err")
		return
	}

	//	m.client.SubscribeMarketData(m.customer.GetMarkets())
	m.client.ReqQryInvestorPosition()
	fmt.Println("Server Running...")
	<-c
}
