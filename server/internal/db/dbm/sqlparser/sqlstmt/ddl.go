package sqlstmt

// DdlStmt DDL 语句
type DdlStmt struct {
	Base
	DdlKind string // CREATE_TABLE, DROP_TABLE, ALTER_TABLE 等
}

func (*DdlStmt) StmtKind() Kind { return KindDdl }
