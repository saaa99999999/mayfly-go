package oracle

import (
	"mayfly-go/internal/db/dbm/sqlparser/sqlstmt"
	"mayfly-go/pkg/logx"
)

type OracleParser struct {
}

func (*OracleParser) Parse(stmt string) (sqlstmt.Stmt, error) {
	defer func() {
		if e := recover(); e != nil {
			logx.ErrorTrace("oracle sql parser err: ", e)
		}
	}()
	return NewParser(stmt).Parse()
}
