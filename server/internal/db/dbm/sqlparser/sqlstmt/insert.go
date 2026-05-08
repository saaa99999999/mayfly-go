package sqlstmt

// InsertStmt INSERT 语句
type InsertStmt struct {
	Base
	Table   TableRef
	Columns []string
	Values  string // VALUES 部分原始文本
}

func (*InsertStmt) StmtKind() Kind { return KindInsert }
