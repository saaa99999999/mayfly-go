package sqlstmt

// DeleteStmt DELETE 语句
type DeleteStmt struct {
	Base
	Tables []TableRef // 删除的表
	Where  *Expr
}

func (*DeleteStmt) StmtKind() Kind { return KindDelete }
