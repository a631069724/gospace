package manager

import (
	"errors"
	"housekeeper2/api"
	"log"
	"sync"
	"time"
)

type QuoterTable interface {
	AddQuoter(api.Quoter)
	GetQuoter(string) api.Quoter
	ContractExist(string) bool
}

type QuoterClient interface {
	Run() error
	OnConnect()
	OnRspUserLogin(api.Response, api.RespInfo, int, bool)
	OnRtnDepthMarketData(api.Quoter)
	SubscribeMarketData([]string)
	GetQuoter(string) api.Quoter
	QuoterExist(string) bool
}

type MyQuoterTable map[string]api.Quoter

type MyQuoterClient struct {
	reqid       int
	relogin     bool
	result      chan error
	accName     string
	pwd         string
	certinfo    string
	apitype     string
	ip          string
	port        string
	mdapi       api.MdApi
	quoterTable QuoterTable
	sync.Mutex
}

func NewQuoterClient(accName string, pwd string, apitype string, certinfo string, ip string, port string) QuoterClient {
	a, _ := api.NewMdApi("CTP")
	return &MyQuoterClient{
		accName:     accName,
		pwd:         pwd,
		apitype:     apitype,
		certinfo:    certinfo,
		ip:          ip,
		port:        port,
		mdapi:       a,
		quoterTable: NewQuoterTable(),
		result:      make(chan error),
	}
}

func (c *MyQuoterClient) Name() string {
	return c.accName
}

func (c *MyQuoterClient) Pwd() string {
	return c.pwd
}

func (c *MyQuoterClient) Certinfo() string {
	return c.certinfo
}

func (c *MyQuoterClient) Run() error {
	c.RegistSpi(c)
	if err := c.Connect(); err != nil {
		log.Println("acount: ", c.Name(), "connect error: ", err.Error())
		return err
	}
	if err := c.ReqUserLogin(); err != nil {
		log.Println("acount: ", c.Name(), "login error: ", err.Error())
		return err
	}
	c.SetReLogin(true)
	return nil
}

func (c *MyQuoterClient) QuoterExist(name string) bool {
	return c.quoterTable.ContractExist(name)
}

func (c *MyQuoterClient) RegistSpi(spi api.MdSpi) error {
	if spi == nil {
		return errors.New("The spi is NULL")
	}
	c.mdapi.RegistSpi(api.NewMdSpi(spi, c.AccType()))
	return nil
}

func (c *MyQuoterClient) Connect() error {
	c.mdapi.Connect(c.ip, c.port)
	return c.Result()
}

func (c *MyQuoterClient) OnConnect() {
	log.Println("account.OnConnect")
	if c.ReLogin() {
		c.ReqUserLogin()
		return
	}
	c.Return(nil)

}

func (c *MyQuoterClient) ReLogin() bool {
	return c.relogin
}

func (c *MyQuoterClient) AccType() string {
	return c.apitype
}

func (c *MyQuoterClient) GetRequestID() int {
	c.Lock()
	defer c.Unlock()
	c.reqid += 1
	if c.reqid >= 999999 {
		c.reqid = 0
	}
	return c.reqid
}

func (c *MyQuoterClient) SetReLogin(b bool) {
	c.relogin = b
}

func (c *MyQuoterClient) ReqUserLogin() error {
	reqdata := api.NewReqData()
	reqdata.SetName(c.Name())
	reqdata.SetPwd(c.Pwd())
	reqdata.SetId(c.GetRequestID())
	reqdata.SetBrokeId(c.Certinfo())
	c.mdapi.ReqUserLogin(reqdata)
	if c.ReLogin() {
		return nil
	}
	return c.Result()

}

func (c *MyQuoterClient) OnRspUserLogin(resp api.Response, info api.RespInfo, reqid int, islast bool) {
	log.Println("MyQuoterClient.OnRspUserLogin")
	c.Return(nil)
}

func (c *MyQuoterClient) Result() error {
	select {
	case b := <-c.result:
		return b
	case <-time.After(5 * time.Second):
		return errors.New("over time")
	}
}

func (c *MyQuoterClient) Return(err error) bool {
	select {
	case <-time.After(5 * time.Second):
		return false
	case c.result <- err:
		return true
	}

}

func (c *MyQuoterClient) SubscribeMarketData(names []string) {
	log.Println("MyQuoterClient.SubscribeMarketData")
	c.mdapi.SubscribeMarketData(names)
}

func (c *MyQuoterClient) OnRtnDepthMarketData(q api.Quoter) {
	//log.Println("MyQuoterClient.OnRtnDepthMarketData")
	c.Lock()
	defer c.Unlock()
	c.quoterTable.AddQuoter(q)

}

func (c *MyQuoterClient) GetQuoter(name string) api.Quoter {
	c.Lock()
	defer c.Unlock()
	return c.quoterTable.GetQuoter(name)
}

func NewQuoterTable() QuoterTable {
	qt := make(MyQuoterTable)
	return qt
}

func (qt MyQuoterTable) AddQuoter(q api.Quoter) {

	qt[q.GetInstrumentID()] = q
}

func (qt MyQuoterTable) GetQuoter(name string) api.Quoter {

	return qt[name]
}

func (qt MyQuoterTable) ContractExist(name string) bool {

	_, ok := qt[name]
	return ok
}
