package dump

import (
	"fmt"

	"github.com/goinggo/mapstructure"
	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/module"
	basemodule "github.com/liangdas/mqant/module/base"
)

var Storer *OraclStorer

var Module = func() module.Module {
	this := new(Dump)
	return this
}

type Dump struct {
	basemodule.BaseModule
	options Options
}

func (this *Dump) GetType() string {
	return "Dump"
}

func (this *Dump) Version() string {
	return "1.0.0"
}

func (this *Dump) OnInit(app module.App, settings *conf.ModuleSettings) {
	this.BaseModule.OnInit(this, app, settings)

	err := mapstructure.Decode(settings.Settings["Options"], &this.options)
	if err != nil {
		fmt.Println("Options error:", settings.Settings["Options"])
	}
	Storer, err = NewOracleStorer(this.options.DbUser, this.options.DbPwd, this.options.ConnStr)
	if err != nil {
		fmt.Println("Oracle Connect error:", err.Error())
		return
	}
	OracleTableInit(Storer)
	this.GetServer().RegisterGO("httpDump", httpDump)
}

func (this *Dump) Run(closeSig chan bool) {

}

func (this *Dump) OnDestroy() {
	this.GetServer().OnDestroy()
	Storer.db.Close()
}
