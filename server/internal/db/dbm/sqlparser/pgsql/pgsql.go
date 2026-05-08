package pgsql

import (
	"mayfly-go/internal/db/dbm/sqlparser/sqlstmt"
	"mayfly-go/pkg/gox"
	"mayfly-go/pkg/logx"
)

type PgsqlParser struct {
}

func (*PgsqlParser) Parse(stmt string) (sqlstmt.Stmt, error) {
	defer gox.Recover(func(e error) {
		logx.ErrorTrace("postgres sql parser err: ", e)
	})
	return NewParser(stmt).Parse()
}
