package main

import (
	"housekeeper/hkapi"
	"housekeeper/manager"
	"os"
)

func main() {
	m := manager.NewManager(os.Args[1])
	m.Init()
	hkapi.HttpApiRegist(m)
	go func() { hkapi.Run() }()
	m.Run()

}
