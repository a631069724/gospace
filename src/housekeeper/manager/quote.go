package manager

import (
	"housekeeper/adaptation"
	"housekeeper/api"
	"sync"
)

type Quoter struct {
	quote    map[string]adaptation.Quote
	contact  map[string]Contract
	client   api.Client
	quoteSpi *QuoteSpi
	sync.Mutex
}

type Contract struct {
	contractEName string
	shortName     string
	ctpContract   string
	money         float64
}

func NewContrct(ename string, shortname string, ctpname string, money float64) *Contract {
	return &Contract{contractEName: ename, shortName: shortname, ctpContract: ctpname, money: money}
}

func (q *Quoter) GetQuote(name string) *adaptation.Quote {
	q.Lock()
	defer q.Unlock()
	if v, ok := q.quote[name]; ok && v.HighestPrice > 0 && v.LowestPrice > 0 {
		return &v
	}
	return nil
}

func (q *Quoter) AddContract(name string, ctc *Contract) {
	q.Lock()
	defer q.Unlock()
	q.contact[name] = *ctc
}

func (q *Quoter) GetContct(name string) Contract {
	q.Lock()
	defer q.Unlock()
	return q.contact[name]
}

func (q *Quoter) SaveQuote(quote *adaptation.Quote) {
	q.Lock()
	defer q.Unlock()
	if quote.LastPrice <= 0 || quote.HighestPrice <= 0 || quote.LowestPrice <= 0 {
		return
	}
	q.quote[quote.InstrumentID] = *quote
	//	log.Println("当前行情列表:", q.quote)
}

func (q *Quoter) SubscribeMarketData(ppInstrumentID []string) {

	q.client.SubscribeMarketData(ppInstrumentID)
}

func NewQuoter(accName string, pwd string, fid string, apitype string, certinfo string, ip string, port string) *Quoter {
	var client api.Client
	quote := &Quoter{quote: make(map[string]adaptation.Quote), contact: make(map[string]Contract)}
	qs := &QuoteSpi{acc: quote}
	if apitype == "CTP" {
		client = api.NewCTPQuoteClient(certinfo, accName, pwd, "tcp://"+ip+":"+port, "tcp://"+ip+":"+port, 0, 0, qs)
	}
	quote.client = client
	quote.quoteSpi = qs
	client.StartMd()
	return quote
}
