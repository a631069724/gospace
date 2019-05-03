package check

import (
	"flag"
	"fmt"
	"riskmanagement/check/fundscheck"
	"riskmanagement/check/fundsproc"
	"riskmanagement/conf"
	"riskmanagement/log"
	"time"
)

func Run() {

	confPath := flag.String("conf", "", "Server configuration file path")
	flag.Parse()
	if *confPath == "" {
		panic("conf path is null use -conf")
	}
	conf.LoadConfig(*confPath)
	fmt.Println("Log File:", conf.Conf.Log.Path)
	log.InitLog(conf.Conf.Log.Path)
	fchk := fundscheck.NewFundsChecker(conf.Conf)
	fproc := fundsproc.NewFundsProcer(fchk.GetObservable(), conf.Conf.Redis)
	fmt.Println("funds procer:", fproc)
	ticker := time.NewTicker(time.Duration(conf.Conf.Tick) * time.Second)
	for _ = range ticker.C {
		fchk.Check()
	}
}
