package gomsp

import (
	"fmt"
	"testing"
)

func TestMsp(t *testing.T) {
	MspAttach(99)
	putbuf := "test go msp"
	MspPut([]byte(putbuf), 99)
	getbuf, len, id, ret := MspGet(60)
	fmt.Println(getbuf, len, id, ret)
}
