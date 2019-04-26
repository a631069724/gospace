package fundscheck

import (
	"riskmanagement/position"
	"riskmanagement/subject"
)

type FundsChecker struct {
	/*观察者对象 用于通知事件变化 fundschecker对象用此通知需要处理的持仓*/
	subject.Observable
	position.Position
}

func NewFundsChecker(oraclCon string) *FundsChecker {
	return &FundsChecker{
		Observable: subject.NewObservable(),
	}
}

func (this *FundsChecker) Check() {
	this.SetChanged()
	this.NotifyObserver()
}
