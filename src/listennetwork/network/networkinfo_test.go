package network

import (
	"fmt"
	"testing"
)

func TestIps(t *testing.T) {
	ips, err := Ips()
	if err != nil {
		t.Error("test error:", err)
	}
	for k, v := range ips {
		fmt.Println("key:", k, "value:", v)
	}
}
