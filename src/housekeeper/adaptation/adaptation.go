package adaptation

import (
	"sync"
)

var mApi map[string]Api = make(map[string]Api)
var mProc map[string]Proc = make(map[string]Proc)
var m sync.Mutex

func RegistApi(a Api) {
	m.Lock()
	defer m.Unlock()
	mApi[a.GetName()] = a
	if v, ok := mProc[a.GetName()]; ok {
		a.SetProc(v)
	}
}

func RegisProc(p Proc) {
	m.Lock()
	defer m.Unlock()
	mProc[p.ApiName()] = p
	if v, ok := mApi[p.ApiName()]; ok {
		v.SetProc(p)
	}
}
