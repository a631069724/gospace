package main

import (
	"fmt"
	"housekeeper2/manager"
	"housekeeper2/serverapi"
	"os"
)

func main() {
	hk, err := manager.NewHouseKeeper(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go serverapi.Start(hk)
	hk.Start()
}
