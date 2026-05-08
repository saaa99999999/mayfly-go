package sqlstmt

// OtherStmt 其他语句
type OtherStmt struct {
	Base
}

func (*OtherStmt) StmtKind() Kind { return KindOther }
