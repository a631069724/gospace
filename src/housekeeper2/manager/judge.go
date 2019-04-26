package manager

import (
	"errors"
	"housekeeper2/rule"
	"log"
)

func JudgeBalance(p rule.Params) (interface{}, error) {

	balance := p.GetParam("Balance").(float64)
	limitbalance := p.GetParam("ForceClose").(float64)
	if limitbalance > balance {
		return true, nil
	}
	return false, errors.New("Non-conformity")
}

func JudgeMargin(p rule.Params) (interface{}, error) {
	balance := p.GetParam("Balance").(float64)
	margin := p.GetParam("Margin").(float64)
	priofund := p.GetParam("PrioFund").(float64)
	bondmult := p.GetParam("BondMult").(float64)

	result := margin - ((balance - priofund) * bondmult)
	log.Println(margin, balance, priofund, bondmult, result)
	if result > 0 {
		return result, nil
	}
	return 0, errors.New("Non-conformity")
}

func JudgePercent(p rule.Params) (interface{}, error) {
	var positions []Position
	psts := p.GetParam("Positions").(map[string][2]Position)
	quotercln := p.GetParam("QuoterCln").(QuoterClient)
	percent := p.GetParam("UDPercent").(float64)
	for k, v := range psts {
		quoter := quotercln.GetQuoter(k)
		if quoter == nil {
			continue
		}
		if quoter.GetHighestPrice()-(quoter.GetHighestPrice()-quoter.GetOpenPrice())*percent < quoter.GetLastPrice() {
			//涨停
			positions = append(positions, v[0])
		} else if quoter.GetLowestPrice()+(quoter.GetOpenPrice()-quoter.GetLowestPrice())*percent > quoter.GetLastPrice() {
			//跌停
			positions = append(positions, v[1])
		}
	}
	if len(positions) > 0 {
		return positions, nil
	}
	return nil, errors.New("Non-conformity")
}
