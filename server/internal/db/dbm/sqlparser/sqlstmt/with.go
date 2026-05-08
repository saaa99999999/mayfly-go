package sqlstmt

// WithStmt WITH 语句
type WithStmt struct {
	Base
}

func (*WithStmt) StmtKind() Kind { return KindWith }
