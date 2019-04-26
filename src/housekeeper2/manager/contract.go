package manager

import "regexp"

type Contract interface {
	//SetContractName()
	GetContractName() string
	//	SetExchangeId()
	GetExchangeId() string
	//	SetCloseTime()
	GetCloseTime() []string
	SetMargin(string, float64)
	GetMargin(string) float64
	IsSHFE() bool
}

type ContractTable interface {
	SetContract(Contract)
	GetContract(string) Contract
}

type MyContractTable map[string]Contract

func NewContractTable() ContractTable {
	table := make(MyContractTable)
	return table
}

func (ct MyContractTable) SetContract(c Contract) {
	ct[c.GetContractName()] = c
}

func (ct MyContractTable) GetContract(name string) Contract {
	re, _ := regexp.Compile("[^0-9]+")
	sname := re.FindString(name)
	return ct[sname]
}

type MyContract struct {
	contractName string
	exchangeId   string
	closeTime    []string
	margin       map[string]float64
}

func NewContract(cname string, exid string, cltime []string) Contract {
	return &MyContract{
		contractName: cname,
		exchangeId:   exid,
		closeTime:    cltime,
		margin:       make(map[string]float64),
	}
}

func (c *MyContract) SetMargin(name string, f float64) {
	c.margin[name] = f
}

func (c *MyContract) GetMargin(name string) float64 {
	return c.margin[name]
}

func (c *MyContract) GetContractName() string {
	return c.contractName
}

func (c *MyContract) GetExchangeId() string {
	return c.exchangeId
}

func (c *MyContract) GetCloseTime() []string {
	return c.closeTime
}

func (c *MyContract) IsSHFE() bool {
	return c.exchangeId == "SHFE"
}
