package tokenizer

import (
	"strings"
	"unicode"
)

// DialectConfig 定义不同 SQL 方言的配置
type DialectConfig struct {
	// 反引号作为标识符引号（MySQL）
	BacktickAsIdentifier bool
	// 双引号作为标识符引号（PostgreSQL）
	DoubleQuoteAsIdentifier bool
	// 支持 # 行注释（MySQL）
	HashLineComment bool
	// 支持 $tag$ 风格字符串/标识符（PostgreSQL）
	DollarQuote bool
	// 额外关键字集合（合并到标准关键字中）
	ExtraKeywords map[string]bool
}

// Tokenizer 将 SQL 字符串拆分为 Token 序列
type Tokenizer struct {
	sql     string
	pos     int
	length  int
	config  DialectConfig
	Tokens  []Token
	current int
}

// New 创建一个新的 Tokenizer
func New(sql string, config DialectConfig) *Tokenizer {
	t := &Tokenizer{
		sql:    sql,
		pos:    0,
		length: len(sql),
		config: config,
		Tokens: make([]Token, 0),
	}
	t.tokenize()
	// 追加 EOF token
	t.Tokens = append(t.Tokens, Token{Type: TokenEOF, Value: "", Pos: t.length, End: t.length})
	return t
}

// tokenize 执行词法分析
func (t *Tokenizer) tokenize() {
	for t.pos < t.length {
		ch := t.sql[t.pos]

		// 跳过空白字符
		if isWhitespace(ch) {
			t.pos++
			continue
		}

		// 行注释 --
		if ch == '-' && t.pos+1 < t.length && t.sql[t.pos+1] == '-' {
			t.skipLineComment()
			continue
		}

		// MySQL # 行注释
		if t.config.HashLineComment && ch == '#' {
			t.skipLineComment()
			continue
		}

		// 块注释 /* */
		if ch == '/' && t.pos+1 < t.length && t.sql[t.pos+1] == '*' {
			t.skipBlockComment()
			continue
		}

		// 字符串字面量 '...'
		if ch == '\'' {
			t.readString('\'')
			continue
		}

		// 双引号字符串 "..."（如果不作为标识符引号）
		if ch == '"' && !t.config.DoubleQuoteAsIdentifier {
			t.readString('"')
			continue
		}

		// 双引号标识符 "..."（PostgreSQL）
		if ch == '"' && t.config.DoubleQuoteAsIdentifier {
			t.readQuotedIdentifier('"')
			continue
		}

		// 反引号标识符 `...`（MySQL）
		if ch == '`' && t.config.BacktickAsIdentifier {
			t.readQuotedIdentifier('`')
			continue
		}

		// PostgreSQL $tag$ ... $tag$
		if t.config.DollarQuote && ch == '$' {
			if t.readDollarQuote() {
				continue
			}
		}

		// 数字
		if isDigit(ch) {
			t.readNumber()
			continue
		}

		// 标识符或关键字（字母、_、@ 开头）
		if isIdentifierStart(ch) {
			t.readIdentifierOrKeyword()
			continue
		}

		// 运算符和标点符号
		if isOperatorStart(ch) {
			t.readOperator()
			continue
		}

		if isPunctuation(ch) {
			t.Tokens = append(t.Tokens, Token{
				Type:  TokenPunctuation,
				Value: string(ch),
				Pos:   t.pos,
				End:   t.pos + 1,
			})
			t.pos++
			continue
		}

		// 未知字符，跳过
		t.pos++
	}
}

// skipLineComment 跳过一个行注释（到行尾或 EOF）
func (t *Tokenizer) skipLineComment() {
	for t.pos < t.length && t.sql[t.pos] != '\n' {
		t.pos++
	}
}

// skipBlockComment 跳过一个块注释 /* */
func (t *Tokenizer) skipBlockComment() {
	t.pos += 2 // 跳过 /*
	for t.pos < t.length {
		if t.sql[t.pos] == '*' && t.pos+1 < t.length && t.sql[t.pos+1] == '/' {
			t.pos += 2
			return
		}
		t.pos++
	}
}

// readString 读取一个单引号或双引号字符串字面量
func (t *Tokenizer) readString(quote byte) {
	start := t.pos
	t.pos++ // 跳过起始引号
	for t.pos < t.length {
		ch := t.sql[t.pos]
		if ch == quote {
			// 检查是否是转义（两个连续引号）
			if t.pos+1 < t.length && t.sql[t.pos+1] == quote {
				t.pos += 2
				continue
			}
			t.pos++ // 跳过结束引号
			break
		}
		// MySQL 风格转义 \'
		if ch == '\\' && t.pos+1 < t.length {
			t.pos += 2
			continue
		}
		t.pos++
	}
	t.Tokens = append(t.Tokens, Token{
		Type:  TokenString,
		Value: t.sql[start:t.pos],
		Pos:   start,
		End:   t.pos,
	})
}

// readQuotedIdentifier 读取一个带引号的标识符（反引号或双引号）
func (t *Tokenizer) readQuotedIdentifier(quote byte) {
	start := t.pos
	t.pos++ // 跳过起始引号
	for t.pos < t.length {
		ch := t.sql[t.pos]
		if ch == quote {
			// 检查转义引号（如 "a""b" 或 ``a``b``）
			if t.pos+1 < t.length && t.sql[t.pos+1] == quote {
				t.pos += 2
				continue
			}
			t.pos++ // 跳过结束引号
			break
		}
		t.pos++
	}
	t.Tokens = append(t.Tokens, Token{
		Type:  TokenIdentifier,
		Value: t.sql[start:t.pos],
		Pos:   start,
		End:   t.pos,
	})
}

// readDollarQuote 读取 PostgreSQL $tag$ ... $tag$ 风格的引号内容
func (t *Tokenizer) readDollarQuote() bool {
	start := t.pos
	// 读取 $tag$
	tagEnd := t.readDollarTag()
	if tagEnd < 0 {
		return false
	}
	tag := t.sql[start : tagEnd+1] // 包含 $tag$
	// 查找结束标记（从 tag 之后开始）
	searchPos := tagEnd + 1
	for searchPos < t.length {
		if strings.HasPrefix(t.sql[searchPos:], tag) {
			searchPos += len(tag)
			t.Tokens = append(t.Tokens, Token{
				Type:  TokenString,
				Value: t.sql[start:searchPos],
				Pos:   start,
				End:   searchPos,
			})
			t.pos = searchPos
			return true
		}
		searchPos++
	}
	// 未找到结束标记，回退
	return false
}

// readDollarTag 读取 PostgreSQL $tag$ 中的 tag，返回结束 $ 的位置
func (t *Tokenizer) readDollarTag() int {
	if t.sql[t.pos] != '$' {
		return -1
	}
	pos := t.pos + 1
	for pos < t.length {
		ch := t.sql[pos]
		if ch == '$' {
			return pos
		}
		if !unicode.IsLetter(rune(ch)) && !unicode.IsDigit(rune(ch)) && ch != '_' {
			return -1
		}
		pos++
	}
	return -1
}

// readNumber 读取一个数字（整数或浮点数）
func (t *Tokenizer) readNumber() {
	start := t.pos
	for t.pos < t.length && (isDigit(t.sql[t.pos]) || t.sql[t.pos] == '.') {
		t.pos++
	}
	// 支持科学计数法 e.g. 1e10, 1.5E-3
	if t.pos < t.length && (t.sql[t.pos] == 'e' || t.sql[t.pos] == 'E') {
		t.pos++
		if t.pos < t.length && (t.sql[t.pos] == '+' || t.sql[t.pos] == '-') {
			t.pos++
		}
		for t.pos < t.length && isDigit(t.sql[t.pos]) {
			t.pos++
		}
	}
	t.Tokens = append(t.Tokens, Token{
		Type:  TokenNumber,
		Value: t.sql[start:t.pos],
		Pos:   start,
		End:   t.pos,
	})
}

// readIdentifierOrKeyword 读取标识符或关键字
func (t *Tokenizer) readIdentifierOrKeyword() {
	start := t.pos
	for t.pos < t.length && isIdentifierPart(t.sql[t.pos]) {
		t.pos++
	}
	value := t.sql[start:t.pos]
	upper := strings.ToUpper(value)

	tokType := TokenIdentifier
	if Keywords[upper] {
		tokType = TokenKeyword
	} else if t.config.ExtraKeywords[upper] {
		tokType = TokenKeyword
	}

	t.Tokens = append(t.Tokens, Token{
		Type:  tokType,
		Value: value,
		Pos:   start,
		End:   t.pos,
	})
}

// readOperator 读取运算符
func (t *Tokenizer) readOperator() {
	start := t.pos
	// 尝试读取多字符运算符
	if t.pos+1 < t.length {
		two := t.sql[t.pos : t.pos+2]
		if two == "<=" || two == ">=" || two == "<>" || two == "!=" ||
			two == "||" || two == "::" || two == "->" || two == "->>" ||
			two == "=>" || two == ".." {
			// PostgreSQL :: 类型转换, -> JSON, ->> JSON text, => key-value, .. 范围
			t.pos += 2
			// 检查 ->>（三字符）
			if two == "->" && t.pos < t.length && t.sql[t.pos] == '>' {
				t.pos++
			}
			t.Tokens = append(t.Tokens, Token{
				Type:  TokenOperator,
				Value: t.sql[start:t.pos],
				Pos:   start,
				End:   t.pos,
			})
			return
		}
	}
	t.pos++
	t.Tokens = append(t.Tokens, Token{
		Type:  TokenOperator,
		Value: t.sql[start:t.pos],
		Pos:   start,
		End:   t.pos,
	})
}

// Peek 预览当前 token（不移动位置）
func (t *Tokenizer) Peek() Token {
	if t.current >= len(t.Tokens) {
		return t.Tokens[len(t.Tokens)-1] // EOF
	}
	return t.Tokens[t.current]
}

// Next 返回当前 token 并移动到下一个
func (t *Tokenizer) Next() Token {
	if t.current >= len(t.Tokens) {
		return t.Tokens[len(t.Tokens)-1] // EOF
	}
	tok := t.Tokens[t.current]
	t.current++
	return tok
}

// Consume 消耗当前位置的 token（等同于 Next，为了可读性）
func (t *Tokenizer) Consume() Token {
	return t.Next()
}

// Pos 返回当前 token 索引
func (t *Tokenizer) Pos() int {
	return t.current
}

// SetPos 设置当前 token 索引
func (t *Tokenizer) SetPos(p int) {
	t.current = p
}

// Length 返回 token 总数
func (t *Tokenizer) Length() int {
	return len(t.Tokens)
}

// TokenAt 获取指定索引的 token
func (t *Tokenizer) TokenAt(idx int) Token {
	if idx < 0 || idx >= len(t.Tokens) {
		return t.Tokens[len(t.Tokens)-1]
	}
	return t.Tokens[idx]
}

// helper functions
func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isIdentifierStart(ch byte) bool {
	return isLetter(ch) || ch == '_' || ch == '@'
}

func isIdentifierPart(ch byte) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_' || ch == '@' || ch == '$'
}

func isOperatorStart(ch byte) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '=' ||
		ch == '<' || ch == '>' || ch == '!' || ch == '|' || ch == ':' || ch == '~'
}

func isPunctuation(ch byte) bool {
	return ch == '(' || ch == ')' || ch == ',' || ch == ';' || ch == '.'
}
