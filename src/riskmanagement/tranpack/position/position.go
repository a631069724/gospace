package position

import "fmt"

type Position interface {
	GetUserId() string
	GetContract() string
	GetPositionId() string
	GetCloseBuyOrSell() string
	GetPositionNum() string
}

type DefaultPosition struct {
	UserId       string
	ContrctEname string
	ContractDate string
	Num          int
	BuyOrSell    int
	PositionId   int64
}

func (this *DefaultPosition) GetUserId() string {
	return this.UserId
}

func (this *DefaultPosition) GetContract() string {
	return this.ContrctEname + " " + this.ContractDate
}

func (this *DefaultPosition) GetPositionId() string {
	return fmt.Sprintf("%d", this.PositionId)

}

func (this *DefaultPosition) GetCloseBuyOrSell() string {
	if this.BuyOrSell == 1 {
		return "-1"
	}
	return "1"
}

func (this *DefaultPosition) GetPositionNum() string {
	return fmt.Sprintf("%d", this.Num)

}
