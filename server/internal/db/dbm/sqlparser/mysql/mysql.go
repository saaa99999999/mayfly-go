package mysql

import (
	"mayfly-go/internal/db/dbm/sqlparser/sqlstmt"
	"mayfly-go/pkg/gox"
	"mayfly-go/pkg/logx"
)

type MysqlParser struct {
}

func (*MysqlParser) Parse(stmt string) (sqlstmt.Stmt, error) {
	defer gox.Recover(func(e error) {
		logx.ErrorTrace("mysql sql parser err: ", e)
	})
	return NewParser(stmt).Parse()
}
