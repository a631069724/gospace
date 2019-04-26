package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"housekeeper2/rule"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Unknwon/goconfig"
	_ "github.com/wendal/go-oci8"
)

var Wg sync.WaitGroup
var dblocker sync.Mutex

type Housekeeper struct {
	config           string
	cfg              *goconfig.ConfigFile
	oracleUser       string
	oraclePwd        string
	oracleSource     string
	serveraddr       string
	db               *sql.DB
	keepers          map[string]*keeper
	ctpContractTable ContractTable
	quoteClinet      QuoterClient
	logfile          string
	sync.Mutex
}

func NewHouseKeeper(config string) (*Housekeeper, error) {
	log.Println("New HouseKeeper")
	hk := &Housekeeper{
		keepers:          make(map[string]*keeper),
		ctpContractTable: NewContractTable(),
	}
	var err error

	hk.config = config
	hk.cfg, err = goconfig.LoadConfigFile(config)
	if err != nil {
		return nil, err
	}

	hk.oracleUser, err = hk.cfg.GetValue("ORACLE", "OracleUser")
	if err != nil {
		return nil, err
	}
	hk.oraclePwd, err = hk.cfg.GetValue("ORACLE", "OraclePwd")
	if err != nil {
		return nil, err
	}
	hk.oracleSource, err = hk.cfg.GetValue("ORACLE", "OracleSource")
	if err != nil {
		return nil, err
	}
	hk.db, err = sql.Open("oci8", hk.oracleUser+"/"+hk.oraclePwd+"@"+hk.oracleSource)
	if err != nil {
		return nil, err
	}
	hk.serveraddr, _ = hk.cfg.GetValue("PUBLIC", "HTTPAddr")
	hk.logfile, err = hk.cfg.GetValue("PUBLIC", "LogFile")
	if err != nil {
		log.Panic("Log File error: ", err)
		return nil, err
	}
	file, err := os.OpenFile(hk.logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		log.Panic("Open Log File error: ", err)
	}
	log.SetOutput(file)

	return hk, nil
}

func (hk *Housekeeper) GetDber() *sql.DB {
	return hk.db
}

func (hk *Housekeeper) GetHttpAddr() string {
	return hk.serveraddr
}

func (hk *Housekeeper) Start() {
	hk.InitQuoter()
	hk.InitContractTable()
	hk.InitKeepers()
	hk.InitRules()
	hk.QuoteClinetRun()
	Wg.Add(1)
	go func() {
		defer Wg.Done()
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
		<-c
		hk.Exit()
	}()

	hk.AccountRun()
	Wg.Wait()

}

func (hk *Housekeeper) Exit() {
	hk.Lock()
	defer hk.Unlock()
	keepers := hk.GetKeepers()
	for _, kp := range keepers {
		kp.AccountExit()
	}
}

func (hk *Housekeeper) QuoteClinetRun() {
	hk.quoteClinet.Run()
}

func (hk *Housekeeper) AccountRun() {
	hk.Lock()
	defer hk.Unlock()
	keepers := hk.GetKeepers()
	for _, kp := range keepers {
		kp.AccountRun()
	}
}

func (hk *Housekeeper) InitQuoter() error {

	var account string
	var pwd string
	var certinfo string
	var ip string
	var port string
	rows, err := hk.db.Query(`select 
	t.account,t.password,t.certinfo,t.ip,t.port 
	from ZHGJ_QUOTE_ACCOUNT t
	where t.is_used='1' and t.api_type='CTP'`)
	if err != nil {
		log.Println("select err:", err)
		return err
	}

	rows.Next()
	rows.Scan(&account, &pwd, &certinfo, &ip, &port)
	hk.quoteClinet = NewQuoterClient(account, pwd, "CTP", certinfo, ip, port)

	return nil

}

func (hk *Housekeeper) GetKeepers() map[string]*keeper {
	return hk.keepers
}

func (hk *Housekeeper) InitContractTable() error {
	rows, err := hk.db.Query(`select CTP,EXCHANGEID,CLOSE_TIME 
	from ZHGJ_CONTRACT`)
	if err != nil {
		log.Println("select err:", err)
		return err
	}
	for rows.Next() {
		var contract string
		var exchangeid string
		var closetime string
		rows.Scan(&contract, &exchangeid, &closetime)
		ctimes := strings.Split(closetime, ",")
		hk.ctpContractTable.SetContract(NewContract(contract, exchangeid, ctimes))
	}
	return nil
}

func (hk *Housekeeper) GetContractTable() ContractTable {
	return hk.ctpContractTable
}

func (hk *Housekeeper) InitRules() error {
	rows, err := hk.db.Query(`select rule_id,u_name,
	account,rule_type from ZHGJ_ACCOUNT_RULE `)
	if err != nil {
		log.Println("select err:", err)
		return err
	}

	for rows.Next() {
		var ruleid string
		var uname string
		var account string
		var rtype string

		rows.Scan(&ruleid, &uname, &account, &rtype)
		kp := hk.GetKeeper(uname)
		if kp == nil {
			log.Println("keeper: ", uname, " not exist")
			continue
		}
		acc := kp.GetAccount(account)
		if acc == nil {
			log.Println("account: ", account, " not exist")
			continue
		}
		var ruler rule.Ruler
		if rtype == "01" {
			ruler = rule.NewMyRule(rtype, JudgeBalance)
		} else if rtype == "02" {
			ruler = rule.NewMyRule(rtype, JudgeMargin)
		} else if rtype == "03" {
			ruler = rule.NewMyRule(rtype, JudgePercent)
		} else {
			continue
		}
		acc.AddRuler(rtype, ruler)

		log.Println("keeper: ", uname, "Account: ", account, "Add rule: ", rtype)

	}
	return nil
}

func (hk *Housekeeper) InitKeepers() error {
	rows, err := hk.db.Query(`select t.u_name,t.account,t.password,r.certinfo,r.ip,r.port,t.EXPIRATION,
	nvl(t.PRIORITY_FUND,0),nvl(t.AFTER_FUND,0),nvl(t.BOND_MULTIPLE,0),nvl(t.UDPERCENT,0),nvl(t.PRE_MINUTES,0),nvl(t.FORCE_CLOSE,0),t.CLOSE_YT
	from ZHGJ_ACCOUNT t
	left join ZHGJ_FUTURE_COMPANY r
	on t.future_id=r.future_id
	where t.is_used='1'`)
	if err != nil {
		log.Println("select err:", err)
		return err
	}

	for rows.Next() {
		var uname string
		var account string
		var pwd string
		var certinfo string
		var ip string
		var port string
		var expriration string
		var bmult float64
		var priofund float64
		var aftfund float64
		var udpect float64
		var premin string
		var forceclose float64
		var closeflag byte

		rows.Scan(&uname, &account, &pwd, &certinfo, &ip, &port, &expriration, &priofund, &aftfund, &bmult, &udpect, &premin, &forceclose, &closeflag)
		if !hk.KeeperExist(uname) {
			if err := hk.AddKeeper(NewKeeper(uname)); err != nil {
				log.Println("cannot add the keeper: ", err.Error())
			}
		}
		uacc, err := NewAccount(true, certinfo, account, pwd, ip, port, expriration, "CTP", priofund, aftfund, bmult, udpect, premin, forceclose, hk.ctpContractTable, closeflag)
		if err != nil {
			log.Println("New account error: ", err.Error())
			return err
		}
		uacc.SetDber(hk.GetDber())
		kp := hk.GetKeeper(uname)
		if kp == nil {
			continue
		}
		uacc.SetQuoteClient(hk.GetQuoteClient())
		log.Println("keeper:", kp.Name(), "Add Account", uacc)
		if err := kp.AddAccount(uacc); err != nil {
			log.Println("cannot add the account: ", err.Error())
		}
	}
	return nil
}

func (hk *Housekeeper) GetQuoteClient() QuoterClient {
	return hk.quoteClinet
}

func (hk *Housekeeper) KeeperExist(name string) bool {
	hk.Lock()
	defer hk.Unlock()
	_, ok := hk.keepers[name]
	return ok
}

func (hk *Housekeeper) AddKeeper(kp *keeper) error {
	hk.Lock()
	defer hk.Unlock()
	if kp == nil {
		return errors.New("keeper is nil")
	}
	name := kp.Name()
	if name == "" {
		return errors.New("keeper name is nil")
	}
	if _, ok := hk.keepers[name]; ok {
		return errors.New(fmt.Sprintf("keeper %s has exist", name))
	}
	hk.keepers[name] = kp
	return nil
}

func (hk *Housekeeper) GetKeeper(name string) *keeper {
	hk.Lock()
	defer hk.Unlock()
	keepers := hk.GetKeepers()
	return keepers[name]
}
