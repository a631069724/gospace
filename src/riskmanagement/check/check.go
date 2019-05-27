package check

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"riskmanagement/check/fundscheck"
	"riskmanagement/check/fundsproc"
	"riskmanagement/conf"
	"riskmanagement/log"
	"syscall"
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
	fundsproc.NewFundsProcer(fchk.GetObservable(), conf.Conf.Redis)
	tick := time.Duration(conf.Conf.Tick)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	fmt.Println("=====risk management bigen running=====")
	for {
		select {
		case <-time.After(tick * time.Second):
			fchk.Check()
		case <-sigs:
			log.Info("catch exit signal")
			log.Close()
			os.Exit(0)
		}
	}

}
