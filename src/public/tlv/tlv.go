package tlv

import (
	"errors"
	"fmt"
	"sync"
)

type TlvPacker struct {
	data sync.Map
}

func (this *TlvPacker) Set(name string, value []byte) {
	this.data.Store(name, value)
}

func (this *TlvPacker) Get(name string) ([]byte, error) {
	v, ok := this.data.Load(name)
	if !ok {
		return []byte{}, errors.New(fmt.Sprintf("tag: %s not exist", name))
	}
	return v.([]byte), nil
}

func (this *TlvPacker) Pack() []byte {
	tagCount, dataLen := 0, 0
	retData := make([]byte, 1024, 1024)
	retDataTmp := make([]byte, 1024, 1024)
	this.data.Range(func(key, value interface{}) bool {
		tlvdata := make([]byte, 100, 100)
		rkey := key.(string)
		rvalue := value.([]byte)
		rvlen := len(rvalue)
		tagCount++
		dataLen += len(rkey) + 5 + rvlen
		tlvdata = append(tlvdata, []byte(rkey)...)
		tlvdata[len(tlvdata)] = 0x00
		tlvdata = append(tlvdata, []byte(fmt.Sprintf("%03d", rvlen))...)
		tlvdata[len(tlvdata)] = 0x00
		tlvdata = append(tlvdata, rvalue...)
		retDataTmp = append(retDataTmp, tlvdata...)
		return true
	})
	retData = append(retData, []byte(fmt.Sprintf("%03d", tagCount))...)
	retData = append(retData, []byte(fmt.Sprintf("%05d", dataLen))...)
	retData = append(retData, retDataTmp...)
	return retData
}
