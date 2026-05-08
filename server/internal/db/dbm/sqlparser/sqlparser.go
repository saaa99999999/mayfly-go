package sqlparser

import (
	"mayfly-go/internal/db/dbm/sqlparser/sqlstmt"
)

type SqlParser interface {

	// Parse 解析单条 SQL 语句
	//  - 返回: 解析后的 Stmt 对象
	Parse(stmt string) (sqlstmt.Stmt, error)
}
