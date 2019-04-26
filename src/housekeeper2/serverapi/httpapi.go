package serverapi

import (
	"fmt"
	"housekeeper2/manager"
	"housekeeper2/rule"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var hk *manager.Housekeeper

func HandleList() {
	http.HandleFunc("/addrule", AddRule)
	http.HandleFunc("/addaccount", AddAccount)
	http.HandleFunc("/delrule", DelRule)
	http.HandleFunc("/delaccount", DelAccount)
	http.HandleFunc("/startaccount", StartAccount)
	http.HandleFunc("/stopaccount", StopAccount)
	http.HandleFunc("/getaccinfo", GetAccountInfo)
	http.HandleFunc("/renewals", Renewals)
}

func Start(h *manager.Housekeeper) {
	hk = h
	HandleList()
	http.ListenAndServe(hk.GetHttpAddr(), nil)

}

func DelRulee(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	keeperName := r.FormValue("umane")
	accountName := r.FormValue("account")
	if keeperName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不能为空")
		return
	}
	if accountName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不能为空")
		return
	}

	keeper := hk.GetKeeper(keeperName)
	if keeper == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不存在")
		return
	}
	account := keeper.GetAccount(accountName)
	if account == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不存在")
		return
	}
	ttypes := strings.Split(r.FormValue("types"), ",")
	for _, ttype := range ttypes {
		account.DelRuler(ttype)
	}
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":"%s"
			}`, "00", "success")
}

func AddRule(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	ttypes := strings.Split(r.FormValue("types"), ",")

	keeperName := r.FormValue("umane")
	accountName := r.FormValue("account")
	if keeperName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不能为空")
		return
	}
	if accountName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不能为空")
		return
	}

	keeper := hk.GetKeeper(keeperName)
	if keeper == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不存在")
		return
	}
	account := keeper.GetAccount(accountName)
	if account == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不存在")
		return
	}

	for _, ttype := range ttypes {
		var ruler rule.Ruler
		if ttype == "01" {
			ruler = rule.NewMyRule(ttype, manager.JudgeBalance)
			fflvl, err := strconv.ParseFloat(r.FormValue("flvl"), 64)
			if err != nil {
				fmt.Fprintf(w, `{
					"code":"%s",
						"msg":"%s"
						}`, "01", "flvl err")
				return
			}
			account.SetForceClose(fflvl)
		} else if ttype == "02" {
			ruler = rule.NewMyRule(ttype, manager.JudgeMargin)
			fpriofund, err := strconv.ParseFloat(r.FormValue("priofund"), 64)
			if err != nil {
				fmt.Fprintf(w, `{
	   			"code":"%s",
	   				"msg":"%s"
	   				}`, "01", "priofund err")
				return
			}
			account.SetPrioritypFund(fpriofund)
			fmulilvl, err := strconv.ParseFloat(r.FormValue("mulilvl"), 64)
			if err != nil {
				fmt.Fprintf(w, `{
	   			"code":"%s",
	   				"msg":"%s"
	   				}`, "01", "mulilvl err")
				return
			}
			account.SetBondMult(fmulilvl)
			fpremin := r.FormValue("premin")
			if fpremin == "" {
				fmt.Fprintf(w, `{
	   			"code":"%s",
	   				"msg":"%s"
	   				}`, "01", "premin err")
				return
			}
			account.SetPreMin(fpremin)
		} else if ttype == "03" {
			ruler = rule.NewMyRule(ttype, manager.JudgePercent)
			fudpercent, err := strconv.ParseFloat(r.FormValue("udpercent"), 64)
			if err != nil {
				fmt.Fprintf(w, `{
	   			"code":"%s",
	   				"msg":"%s"
	   				}`, "01", "udpercent err")
				return
			}
			account.SetUDPercent(fudpercent)
		} else {
			continue
		}
		account.AddRuler(ttype, ruler)
	}
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":"%s"
			}`, "00", "success")

}

func DelRule(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	ttypes := strings.Split(r.FormValue("types"), ",")
	keeperName := r.FormValue("umane")
	accountName := r.FormValue("account")
	if keeperName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不能为空")
		return
	}
	if accountName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不能为空")
		return
	}

	keeper := hk.GetKeeper(keeperName)
	if keeper == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不存在")
		return
	}
	account := keeper.GetAccount(accountName)
	if account == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不存在")
		return
	}
	for _, ttype := range ttypes {
		if ttype == "01" || ttype == "02" || ttype == "03" {
			account.DelRuler(ttype)
		}
	}
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":"%s"
			}`, "00", "success")

}

func AddAccount(w http.ResponseWriter, r *http.Request) {
	keeperName := r.FormValue("umane")
	if !hk.KeeperExist(keeperName) {
		if err := hk.AddKeeper(manager.NewKeeper(keeperName)); err != nil {
			//log.Println("cannot add the keeper: ", err.Error())
			fmt.Fprintf(w, `{
				"code":"%s",
					"msg":"cannot add the keeper %s"
					}`, "01", err.Error())
			return
		}
	}
	fpriofund, _ := strconv.ParseFloat(r.FormValue("priofund"), 64)
	fafterfund, _ := strconv.ParseFloat(r.FormValue("afterfund"), 64)
	bmult, _ := strconv.ParseFloat(r.FormValue("mulilvl"), 64)
	udpect, _ := strconv.ParseFloat(r.FormValue("udpercent"), 64)
	forceclose, _ := strconv.ParseFloat(r.FormValue("flvl"), 64)

	uacc, err := manager.NewAccount(false,
		r.FormValue("certinfo"),
		r.FormValue("account"),
		r.FormValue("pwd"),
		r.FormValue("ip"),
		r.FormValue("port"),
		r.FormValue("expriration"),
		r.FormValue("apitype"),
		fpriofund, fafterfund, bmult, udpect, r.FormValue("premin"), forceclose, hk.GetContractTable(), '0')
	if err != nil {
		//log.Println("New account error: ", err.Error())
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"New account error %s"
				}`, "01", err.Error())
		return
	}
	uacc.SetDber(hk.GetDber())
	kp := hk.GetKeeper(keeperName)
	uacc.SetQuoteClient(hk.GetQuoteClient())
	log.Println("keeper:", kp.Name(), "Add Account", uacc)
	if err := kp.AddAccount(uacc); err != nil {
		//	log.Println("cannot add the account: ", err.Error())
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"cannot add the account: %s"
				}`, "01", err.Error())
		return
	}
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":"%s"
			}`, "00", "success")
}

func StartAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	keeperName := r.FormValue("umane")
	accountName := r.FormValue("account")
	if keeperName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不能为空")
		return
	}
	if accountName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不能为空")
		return
	}

	keeper := hk.GetKeeper(keeperName)
	if keeper == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不存在")
		return
	}
	account := keeper.GetAccount(accountName)
	if account == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不存在")
		return
	}
	manager.Wg.Add(1)
	go func() {
		defer manager.Wg.Done()
		account.SetUsed(true)
		if err := account.Run(); err != nil {
			log.Println("Account Running error: ", err.Error())
		}
	}()
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":"%s"
			}`, "00", "success")
}

func StopAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	keeperName := r.FormValue("umane")
	accountName := r.FormValue("account")
	if keeperName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不能为空")
		return
	}
	if accountName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不能为空")
		return
	}

	keeper := hk.GetKeeper(keeperName)
	if keeper == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不存在")
		return
	}
	account := keeper.GetAccount(accountName)
	if account == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不存在")
		return
	}
	account.SetUsed(false)
	account.Exit()
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":"%s"
			}`, "00", "success")
}

func DelAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	keeperName := r.FormValue("umane")
	accountName := r.FormValue("account")
	if keeperName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不能为空")
		return
	}
	if accountName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不能为空")
		return
	}

	keeper := hk.GetKeeper(keeperName)
	if keeper == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不存在")
		return
	}
	account := keeper.GetAccount(accountName)
	if account == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不存在")
		return
	}
	account.Exit()
	keeper.DelAccount(accountName)
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":"%s"
			}`, "00", "success")
}

func GetAccountInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	keeperName := r.FormValue("umane")
	accountName := r.FormValue("account")
	if keeperName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不能为空")
		return
	}
	if accountName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不能为空")
		return
	}

	keeper := hk.GetKeeper(keeperName)
	if keeper == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不存在")
		return
	}
	account := keeper.GetAccount(accountName)
	if account == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不存在")
		return
	}
	msg := account.GetAccountInfo()
	if msg == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":%s
				}`, "01", "json msg error")
		return
	}
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":%s
			}`, "00", msg)
}

func Renewals(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)

	keeperName := r.FormValue("umane")
	accountName := r.FormValue("account")
	if keeperName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不能为空")
		return
	}
	if accountName == "" {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不能为空")
		return
	}

	keeper := hk.GetKeeper(keeperName)
	if keeper == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "客户号不存在")
		return
	}
	account := keeper.GetAccount(accountName)
	if account == nil {
		fmt.Fprintf(w, `{
			"code":"%s",
				"msg":"%s"
				}`, "01", "账号不存在")
		return
	}
	account.SetExpir(r.FormValue("expriration"))
	fmt.Fprintf(w, `{
		"code":"%s",
			"msg":%s
			}`, "00", "success")
}
