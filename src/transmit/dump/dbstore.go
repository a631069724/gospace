package dump

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/liangdas/mqant/log"
	_ "github.com/wendal/go-oci8"
)

var (
	insertSqlFormat = "insert into %s (%s) values (%s)"
)

type OraclStorer struct {
	db *sql.DB
}

func NewOracleStorer(usr, pwd, connstr string) (*OraclStorer, error) {
	db, err := sql.Open("oci8", usr+"/"+pwd+"@"+connstr)

	return &OraclStorer{
		db: db,
	}, err
}

func columsAndvaluesStr(data interface{}, dataTag string) (string, string) {
	var columns string
	var values string
	s := reflect.TypeOf(data).Elem()
	v := reflect.ValueOf(data).Elem()
	columns = s.Field(0).Tag.Get(dataTag)
	values = `"` + v.Field(0).Interface().(string) + `"`

	for i := 1; i < s.NumField(); i++ {
		columns += "," + s.Field(i).Tag.Get(ORACLE_TAG)
		values += `,"` + v.Field(i).Interface().(string) + `"`
	}
	return columns, values
}

func Struct2InserSql(table string, data interface{}, dataTag string) string {
	var runSql string
	column, values := columsAndvaluesStr(data, dataTag)
	runSql = fmt.Sprintf(insertSqlFormat, table, column, values)
	return runSql
}

func (this *OraclStorer) DataStore(table string, data interface{}) error {
	runSql := Struct2InserSql(table, data, ORACLE_TAG)

	return this.RunSql(runSql)
}

func (this *OraclStorer) RunSql(runSql string) error {
	log.Info("RunSql [%s]", runSql)
	_, err := this.db.Exec(runSql)
	if err != nil {
		return err
	}
	this.db.Exec("commit")
	return nil
}
