package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Unknwon/goconfig"
	_ "github.com/wendal/go-oci8"
)

type HouseKeeper struct {
	oracleUser   string
	oraclePwd    string
	oracleSource string
	db           *sql.DB
	httpaddr     string
	keepers      map[string]*Keeper
	quoter       *Quoter
	sync.Mutex
}

func (hk *HouseKeeper) KeeperExist(keeperName string) bool {
	hk.Lock()
	hk.Unlock()
	_, ok := hk.keepers[keeperName]
	return ok
}

/*
func (hk *HouseKeeper) Star(keeperName string) error {
	hk.Lock()
	defer hk.Unlock()
	if !hk.keepers[keeperName].Start() {
		return error.New(fmt.Sprintf("keeper %s start error", keeperName))
	}
}
*/
func (hk *HouseKeeper) AddKeeper(keeperName string, keeper *Keeper) error {
	hk.Lock()
	defer hk.Unlock()
	if _, ok := hk.keepers[keeperName]; ok {
		return errors.New(fmt.Sprintf("keeper %s exist", keeperName))
	}
	log.Println("Add Keeper: ", keeper)
	hk.keepers[keeperName] = keeper
	return nil
}

func (hk *HouseKeeper) AddRule(keeperName string, urule *Rule) error {
	hk.Lock()
	defer hk.Unlock()
	if v, ok := hk.keepers[keeperName]; ok {

		return v.AddRule(urule)
	}
	return errors.New("keeper not exist")
}

func (hk *HouseKeeper) AddAccount(keeperName string, account *Account) error {
	hk.Lock()
	defer hk.Unlock()
	log.Println("Add Account: ", account)
	if v, ok := hk.keepers[keeperName]; ok {
		v.AddAccount(account)
	}
	return nil
}

func NewManager(path string) *HouseKeeper {
	var hk HouseKeeper
	log.Println("New Manger")
	cfg, err := goconfig.LoadConfigFile(path)
	if err != nil {
		log.Println("load config: ", path, "err: ", err)
		return nil
	}
	p, err := cfg.GetValue("ORACLE", "OracleUser")
	if err != nil {
		log.Println("oracle user error", err)
		return nil
	}
	hk.oracleUser = p
	p, err = cfg.GetValue("ORACLE", "OraclePwd")
	if err != nil {
		log.Println("oracle pwd error", err)
		return nil
	}
	hk.oraclePwd = p
	p, err = cfg.GetValue("ORACLE", "OracleSource")
	if err != nil {
		log.Println("oracle source error", err)
		return nil
	}
	hk.oracleSource = p

	hk.db, err = sql.Open("oci8", hk.oracleUser+"/"+hk.oraclePwd+"@"+hk.oracleSource)
	if err != nil {
		log.Panicln("db open error:", err)
		return nil
	}
	hk.keepers = make(map[string]*Keeper)
	hk.httpaddr, _ = cfg.GetValue("PUBLIC", "HTTPAddr")
	return &hk
}

func (hk *HouseKeeper) GetHttpAddr() string {
	return hk.httpaddr
}

func (hk *HouseKeeper) QuoteInit() error {
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
	hk.quoter = NewQuoter(account, pwd, "9999", "CTP", certinfo, ip, port)

	rows, err = hk.db.Query(`select CONTRACT_ENAME,SHORT_NAME,MONEY,CTP from ZHGJ_CONTRACT`)
	if err != nil {
		log.Println("select err:", err)
		return err
	}
	for rows.Next() {
		var cename string
		var shortname string
		var money float64
		var ctpname string
		rows.Scan(&cename, &shortname, &money, &ctpname)
		hk.quoter.AddContract(ctpname, NewContrct(cename, shortname, ctpname, money))
	}

	return nil

}

func (hk *HouseKeeper) AccountInit() error {
	rows, err := hk.db.Query(`select t.u_name,t.account,t.password,r.certinfo,r.ip,r.port from ZHGJ_ACCOUNT t
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

		rows.Scan(&uname, &account, &pwd, &certinfo, &ip, &port)
		if !hk.KeeperExist(uname) {
			hk.AddKeeper(uname, NewKeeper(uname))
		}

		uacc := NewAccount(account, pwd, "9999", "CTP", certinfo, ip, port)

		log.Println("keeper:", uname, "Add Account", uacc)
		if uacc != nil {
			uacc.SetQuoter(hk.quoter)
			uacc.MyInit()
			hk.AddAccount(uname, uacc)
		}
	}
	return nil
}

func (hk *HouseKeeper) RuleInit() error {
	rows, err := hk.db.Query(`select rule_id,u_name,
	account,rule_type,fund_level,bond_level,bond_multiple,start_time,
	end_time from ZHGJ_RULE `)
	if err != nil {
		log.Println("select err:", err)
		return err
	}

	for rows.Next() {
		var ruleid string
		var uname string
		var account string
		var rtype string
		var flvl float64
		var blvl float64
		var bmlt float64
		var stime string
		var etime string

		rows.Scan(&ruleid, &uname, &account, &rtype, &flvl, &blvl, &bmlt, &stime, &etime)
		urule := NewRule(ruleid, account, rtype, flvl, blvl, bmlt, stime, etime)
		log.Println("keeper:", uname, "Add rule", urule)
		if urule != nil {
			if err := hk.AddRule(uname, urule); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func (hk *HouseKeeper) Init() {
	hk.QuoteInit()
	hk.AccountInit()
	hk.RuleInit()
}

func (hk *HouseKeeper) Run() {
	log.Println("Running...")
	for {
		for k, v := range hk.keepers {
			log.Println("kepper : ", k, "begin check")
			go func() { v.CheckRule() }()
		}
		time.Sleep(60 * time.Second)
	}
}

func (hk *HouseKeeper) GetQuoter() *Quoter {
	return hk.quoter
}
