package check

import (
	"flag"
	"riskmanagement/conf"
)

func Run() {

	confPath := flag.String("conf", "", "Server configuration file path")
	flag.Parse()
	if *confPath == "" {
		panic("conf path is null use -conf")
	}
	conf.LoadConfig(*confPath)
}
