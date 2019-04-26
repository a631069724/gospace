package Tools

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func InTheTimeMins(mins string, stimes []string) bool {
	for _, stime := range stimes {
		te, err := time.Parse("150405", stime)
		if err != nil {
			log.Println("time: ", stime, "err :", err.Error())
			return false
		}
		m, _ := time.ParseDuration(fmt.Sprintf("-%sm", mins))

		ts := te.Add(m)
		tn, _ := time.Parse("150405", time.Now().Format("150405"))

		if tn.Before(te) && tn.After(ts) {
			log.Println(ts, tn, te, "TRUE")
			return true
		}
		log.Println(ts, tn, te, "FALSE")
	}

	return false
}
