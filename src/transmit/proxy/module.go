package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	argsutil "github.com/liangdas/mqant/rpc/util"

	"github.com/goinggo/mapstructure"

	"github.com/liangdas/mqant/conf"
	"github.com/liangdas/mqant/log"
	"github.com/liangdas/mqant/module"
	basemodule "github.com/liangdas/mqant/module/base"
)

var Module = func() module.Module {
	this := new(HttpProxy)
	return this
}

type HttpProxy struct {
	basemodule.BaseModule
	options Options
}

func (this *HttpProxy) GetType() string {
	return "HttpProxy"
}

func (this *HttpProxy) Version() string {
	return "1.0.0"
}

func (this *HttpProxy) OnInit(app module.App, settings *conf.ModuleSettings) {
	this.BaseModule.OnInit(this, app, settings)
	/*this.options.Port = settings.Settings["Port"].(string)
	this.options.Raddr = settings.Settings["Raddr"].(string)
	this.options.Monitor = settings.Settings["Monitor"].(string) == "true"*/
	mapstructure.Decode(settings.Settings["Options"], &this.options)
	fmt.Println(this.options)

}

func (this *HttpProxy) Run(closeSig chan bool) {

	go func() {
		log.Info("HttpProxy listening On: %s", this.options.Port)
		http.ListenAndServe(this.options.Port, this)
		<-closeSig
		log.Info("HttpProxy stop!!!!")
	}()
}

func (this *HttpProxy) OnDestroy() {
	this.GetServer().OnDestroy()
}

func (this *HttpProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !this.options.Dump {
		remote, err := url.Parse("http://" + this.options.Raddr)
		if err != nil {
			log.Debug("Parse url err:%s", err.Error())
		}
		p := httputil.NewSingleHostReverseProxy(remote)
		p.ServeHTTP(w, r)

	} else {

		reqDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Debug("DumpRequest error:%s", err.Error())
			remote, err := url.Parse("http://" + this.options.Raddr)
			if err != nil {
				log.Debug("Parse url err:%s", err.Error())
			}
			p := httputil.NewSingleHostReverseProxy(remote)
			p.ServeHTTP(w, r)
			return
		}

		connIn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			log.Error("hijack error:%s", err.Error())
		}
		defer connIn.Close()
		connOut, err := net.DialTimeout("tcp", this.options.Raddr, time.Second*30)
		if err != nil {
			log.Error("dail to:%s err:%s", this.options.Raddr, err.Error())
			return
		}
		if err = r.Write(connOut); err != nil {
			log.Error("send to server error:%s", err.Error())
			return
		}
		respOut, err := http.ReadResponse(bufio.NewReader(connOut), r)
		if err != nil && err != io.EOF {
			log.Error("read response error:%s", err.Error())
			return
		}
		if respOut == nil {
			log.Error("respOut is nil")
			return
		}
		respDump, err := httputil.DumpResponse(respOut, true)
		if err != nil {
			log.Error("respDump error:%s", err.Error())
		}
		_, err = connIn.Write(respDump)
		if err != nil {
			log.Error("connIn write error:%s", err.Error())
		}
		log.Info("reqDump [%s]", reqDump)
		log.Info("respDump [%s]", respDump)

		this.RpcInvokeNRArgs("Dump", "httpDump",
			[]string{argsutil.STRING, argsutil.BYTES, argsutil.BYTES},
			[][]byte{[]byte(this.GetType()), reqDump, reqDump})
	}
}
