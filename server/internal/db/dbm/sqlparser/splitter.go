package sqlparser

import (
	"io"

	"mayfly-go/internal/pkg/utils"
)

// SQLSplitter SQL 切割器接口
// 不同方言有不同的切割规则，例如：
//   - WITH CTE 语句不能被分号切割
//   - PostgreSQL 的 DO $$ ... $$ 块
//   - MySQL 的 DELIMITER 命令
//   - 存储过程/函数中的分号
type SQLSplitter interface {
	// SplitSQL 切割 SQL 语句
	//  - r: 读取器
	//  - callback: 回调函数，每条完整的 SQL 语句调用一次
	SplitSQL(r io.Reader, callback utils.StmtCallback) error
}

// DefaultSplitter 默认 SQL 切割器
// 适用于大多数标准 SQL 场景
type DefaultSplitter struct {
	delimiter rune
}

// NewDefaultSplitter 创建默认切割器
func NewDefaultSplitter(delimiter ...rune) *DefaultSplitter {
	delim := rune(';')
	if len(delimiter) > 0 {
		delim = delimiter[0]
	}
	return &DefaultSplitter{delimiter: delim}
}

// SplitSQL 实现 SQLSplitter 接口
func (s *DefaultSplitter) SplitSQL(r io.Reader, callback utils.StmtCallback) error {
	return utils.SplitStmts(r, s.delimiter, callback)
}
