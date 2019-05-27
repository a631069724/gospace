package main

import (
	"transmit/dump"
	"transmit/proxy"

	"github.com/liangdas/mqant"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/registry"
	"github.com/nats-io/go-nats"
)

func main() {
	//	rs := registry.NewRegistry(registry.Addrs(""))
	rs := registry.DefaultRegistry
	nc, err := nats.Connect(nats.DefaultURL, nats.MaxReconnects(10000))
	if err != nil {

	}
	app := mqant.CreateApp(
		module.Nats(nc),
		module.Registry(rs),
	)
	app.Run(false, proxy.Module(), dump.Module())

}
