package dump

import (
	"fmt"
	"testing"
)

func TestStruct2InserSql(t *testing.T) {
	fmt.Println(Struct2InserSql("table_name", new(ConsumeReq)))
}
