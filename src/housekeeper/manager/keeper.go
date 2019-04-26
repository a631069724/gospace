package manager

import (
	"errors"
	"log"
	"sync"
)

type Keeper struct {
	uname    string
	accounts map[string]*Account

	sync.Mutex
}

func NewKeeper(uname string) *Keeper {
	return &Keeper{uname: uname, accounts: make(map[string]*Account)}
}

func (k *Keeper) AddAccount(account *Account) {
	k.Lock()
	defer k.Unlock()
	if _, ok := k.accounts[account.account]; !ok {
		k.accounts[account.account] = account
	}
}

func (k *Keeper) SetQuoter(acc string, q *Quoter) {
	k.Lock()
	defer k.Unlock()
	if v, ok := k.accounts[acc]; ok {
		v.SetQuoter(q)
	}
}

/*
func (k *Keeper)DelAccount(){
	k.accounts
}*/
/*
func (k *Keeper) AddRule(urule *Rule) {
	k.Lock()
	defer k.Unlock()
	if _, ok := k.rules[urule.id]; !ok {
		k.rules[urule.id] = urule
	}
}*/

func (k *Keeper) AddRule(urule *Rule) error {
	k.Lock()
	defer k.Unlock()

	if v, ok := k.accounts[urule.account]; ok {
		v.AddRule(urule)
		return nil
	}
	return errors.New("account not exist")
}

func (k *Keeper) CheckRule() {

	k.Lock()
	defer k.Unlock()
	for _, v := range k.accounts {
		log.Println(v.account, " checking...")
		v.ClearPosition()
		v.client.ReqQryInvestorPosition(v.client.GetTraderRequestID())
		<-v.wait
		v.client.ReqQryTradingAccount(v.client.GetTraderRequestID())
		<-v.wait
		v.CheckRule()
	}

}
