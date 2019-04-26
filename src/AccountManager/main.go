package main

import (
	"AccountManager/Manager"
	"fmt"
)

func main() {
	m := Manager.NewManager()
	m.Run()
	fmt.Println("!!!!!Server Start OK!!!!!")
}
