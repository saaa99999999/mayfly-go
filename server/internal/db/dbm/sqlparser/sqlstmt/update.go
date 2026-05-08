package sqlstmt

// UpdateStmt UPDATE 语句
type UpdateStmt struct {
	Base
	Tables []TableRef   // 更新的表（支持多表）
	Set    []Assignment // SET 子句
	Where  *Expr
}

func (*UpdateStmt) StmtKind() Kind { return KindUpdate }
