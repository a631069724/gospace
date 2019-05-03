package fundscheck

import (
	"database/sql"
	"fmt"
	"riskmanagement/conf"
	"riskmanagement/log"
	"riskmanagement/subject"
	"riskmanagement/tranpack/position"

	_ "github.com/wendal/go-oci8"
)

type FundsChecker struct {
	oracle         *conf.Oracle
	runPositionSql *conf.RunSql
	db             *sql.DB
	/*观察者对象 用于通知事件变化 fundschecker对象用此通知需要处理的持仓*/
	subject.Observable
	position.Position
}

func NewFundsChecker(c conf.Config) *FundsChecker {

	db, err := sql.Open("oci8", c.Oracle.UserName+"/"+c.Oracle.Pwd+"@"+c.Oracle.Source)
	if err != nil {
		panic(fmt.Sprintf("oracle open error %v", err))
	}
	return &FundsChecker{
		oracle:         c.Oracle,
		runPositionSql: c.RunSql,
		db:             db,
		Observable:     subject.NewObservable(),
	}
}

func (this *FundsChecker) Check() error {
	for _, sqlstr := range this.runPositionSql.Sql {
		fmt.Println("Run Sql:", sqlstr)
		log.Info("Run Sql:%s", sqlstr)
		rows, err := this.db.Query(sqlstr)
		if err != nil {
			log.Info("select err:", err)
			return err
		}

		for rows.Next() {
			pst := position.DefaultPosition{}
			rows.Scan(&pst.UserId, &pst.ContrctEname, &pst.ContractDate, &pst.PositionId,
				&pst.Num, &pst.BuyOrSell)
			this.Position = &pst
			this.updateFlag(pst.PositionId)
			this.deleteFollowRelation(pst.UserId)
			log.Info("Position:%d Notify", pst.PositionId)
			this.SetChanged()
			this.NotifyObserver(pst)

		}
	}
	return nil
}

func (this *FundsChecker) GetObservable() subject.Observable {
	return this.Observable
}

func (this *FundsChecker) updateFlag(PositionId int64) {
	this.db.Exec(`update yishenglong_trad_position set SFORCE_CLOSE_FLAG='21' where position_id=%d`, PositionId)
	this.db.Exec("commit")
}

func (this *FundsChecker) deleteFollowRelation(fwUid string) {
	this.db.Exec(`delete yishenglong_follow_relation where fw_u_id=%s`, fwUid)
	this.db.Exec("commit")
}
