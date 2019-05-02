package check

import (
	"flag"
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
	log.InitLog(conf.Conf.Log.Path)
	fchk := fundscheck.NewFundsChecker(conf.Conf)
	fundsproc.NewFundsProcer(fchk.GetObservable())
	ticker := time.NewTicker(time.Duration(conf.Conf.Tick) * time.Second)
	for _ = range ticker.C {
		fchk.Check()
	}
}
