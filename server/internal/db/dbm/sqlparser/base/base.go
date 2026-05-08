package base

// SingleSQL is a separate SQL split from multi-SQL.
type SingleSQL struct {
	Text string
	// BaseLine is the line number of the first line of the SQL in the original SQL.
	// HINT: ZERO based.
	BaseLine int
	// FirstStatementLine is the line number of the first non-comment and non-blank line of the SQL in the original SQL.
	// HINT: ZERO based.
	FirstStatementLine int
	// FirstStatementColumn is the column number of the first non-comment and non-blank line of the SQL in the original SQL.
	// HINT: ZERO based.
	FirstStatementColumn int
	// LastLine is the line number of the last line of the SQL in the original SQL.
	// HINT: ZERO based.
	LastLine int
	// LastColumn is the column number of the last line of the SQL in the original SQL.
	// HINT: ZERO based.
	LastColumn int
	// The sql is empty, such as `/* comments */;` or just `;`.
	Empty bool

	// ByteOffsetStart is the start position of the sql.
	// This field may not be present for every engine.
	// ByteOffsetStart is intended for sql execution log display. It may not represent the actual sql that is sent to the database.
	ByteOffsetStart int
	// ByteOffsetEnd is the end position of the sql.
	// This field may not be present for every engine.
	// ByteOffsetEnd is intended for sql execution log display. It may not represent the actual sql that is sent to the database.
	ByteOffsetEnd int
}

// SyntaxError is a syntax error.
type SyntaxError struct {
	Line    int
	Column  int
	Message string
}

// Error returns the error message.
func (e *SyntaxError) Error() string {
	return e.Message
}

func FilterEmptySQL(list []SingleSQL) []SingleSQL {
	var result []SingleSQL
	for _, sql := range list {
		if !sql.Empty {
			result = append(result, sql)
		}
	}
	return result
}

func FilterEmptySQLWithIndexes(list []SingleSQL) ([]SingleSQL, []int32) {
	var result []SingleSQL
	var originalIndex []int32
	for i, sql := range list {
		if !sql.Empty {
			result = append(result, sql)
			originalIndex = append(originalIndex, int32(i))
		}
	}
	return result, originalIndex
}

func GetOffsetLength(total int) int {
	length := 1
	for {
		if total < 10 {
			return length
		}
		total /= 10
		length++
	}
}
