package dm

import (
	"mayfly-go/internal/db/dbm/sqlparser/sqlstmt"
	"mayfly-go/pkg/logx"
)

type DmParser struct {
}

func (*DmParser) Parse(stmt string) (sqlstmt.Stmt, error) {
	defer func() {
		if e := recover(); e != nil {
			logx.ErrorTrace("dm sql parser err: ", e)
		}
	}()
	return NewParser(stmt).Parse()
}
