package hkapi

import (
	"fmt"
	"housekeeper/manager"
	"log"
	"net/http"
	"strconv"
)

var hk *manager.HouseKeeper

func HttpApiRegist(h *manager.HouseKeeper) {
	hk = h
}

func HandleList() {
	http.HandleFunc("/addrule", AddRule)
	http.HandleFunc("/addaccount", AddRule)
}

func Run() {
	HandleList()
	http.ListenAndServe(hk.GetHttpAddr(), nil)
}

func AddRule(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	fflvl, err := strconv.ParseFloat(r.FormValue("flvl"), 64)
	if err != nil {
		fmt.Fprintln(w, "flvl err")
	}
	fblvl, err := strconv.ParseFloat(r.FormValue("blvl"), 64)
	if err != nil {
		fmt.Fprintln(w, "blvl err")
	}
	fmulilvl, err := strconv.ParseFloat(r.FormValue("mulilvl"), 64)
	if err != nil {
		fmt.Fprintln(w, "mulilvl err")
	}

	urule := manager.NewRule(r.FormValue("id"), r.FormValue("account"), r.FormValue("type"), fflvl, fblvl, fmulilvl, r.FormValue("starttime"), r.FormValue("endtime"))
	hk.AddRule("gengzhi", urule)
	fmt.Fprintf(w, `{
"code":"%s",
	"msg":"%s"
	}`, "00", "success")

}

func AddAcount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	uname := r.FormValue("uname")
	if !hk.KeeperExist(uname) {
		hk.AddKeeper(uname, manager.NewKeeper(uname))
	}

	uacc := manager.NewAccount(r.FormValue("account"), r.FormValue("pwd"), "9999", "CTP", r.FormValue("certinfo"), r.FormValue("ip"), r.FormValue("port"))

	log.Println("keeper:", uname, "Add Account", uacc)
	if uacc != nil {
		uacc.MyInit()
		uacc.SetQuoter(hk.GetQuoter())
		hk.AddAccount(uname, uacc)
		fmt.Fprintf(w, `{
				"code":"%s",
				"msg":"%s"
				}`, "00", "success")

	} else {
		fmt.Fprintf(w, `{
				"code":"%s",
				"msg":"%s"
			}`, "01", "error info")

	}
}
