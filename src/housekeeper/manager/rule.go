package manager

type Rule struct {
	id           string
	account      string
	rtype        string
	fundlevel    float64
	bondlevel    float64
	bondmultiple float64
	starttime    string
	endtime      string
}

func NewRule(id string, account string, rtype string, flvl float64, blvl float64, bmult float64, stime string, etime string) *Rule {
	return &Rule{
		id:           id,
		account:      account,
		rtype:        rtype,
		fundlevel:    flvl,
		bondlevel:    blvl,
		bondmultiple: bmult,
		starttime:    stime,
		endtime:      etime}
}
