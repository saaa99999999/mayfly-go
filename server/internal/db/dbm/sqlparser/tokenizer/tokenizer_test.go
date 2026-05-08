package tokenizer

import (
	"testing"
)

func TestTokenizerBasic(t *testing.T) {
	sql := "SELECT id, name FROM users WHERE age > 18;"
	tok := New(sql, DialectConfig{})
	if len(tok.Tokens) < 5 {
		t.Fatalf("expected at least 5 tokens, got %d", len(tok.Tokens))
	}
	// 检查第一个 token 是 SELECT 关键字
	if !tok.Tokens[0].IsKeyword("SELECT") {
		t.Fatalf("expected first token to be SELECT, got %s", tok.Tokens[0].Value)
	}
}

func TestTokenizerMySQLBacktick(t *testing.T) {
	sql := "SELECT `id`, `name` FROM `users`"
	tok := New(sql, DialectConfig{BacktickAsIdentifier: true})
	found := false
	for _, token := range tok.Tokens {
		if token.Type == TokenIdentifier && token.Value == "`id`" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to find backtick identifier `id`")
	}
}

func TestTokenizerPostgresDoubleQuote(t *testing.T) {
	sql := `SELECT "id", "name" FROM "users"`
	tok := New(sql, DialectConfig{DoubleQuoteAsIdentifier: true})
	found := false
	for _, token := range tok.Tokens {
		if token.Type == TokenIdentifier && token.Value == `"id"` {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to find double-quote identifier \"id\"")
	}
}

func TestTokenizerComments(t *testing.T) {
	sql := `SELECT /* comment */ id FROM users -- line comment
WHERE age = 1`
	tok := New(sql, DialectConfig{})
	for _, token := range tok.Tokens {
		if token.Value == "comment" || token.Value == "line" {
			t.Fatalf("comment content should be skipped, got token: %s", token.Value)
		}
	}
}

func TestTokenizerMySQLHashComment(t *testing.T) {
	sql := "SELECT id FROM users # this is a comment\nWHERE age = 1"
	tok := New(sql, DialectConfig{HashLineComment: true})
	for _, token := range tok.Tokens {
		if token.Value == "this" || token.Value == "comment" {
			t.Fatalf("hash comment content should be skipped, got token: %s", token.Value)
		}
	}
}

func TestTokenizerString(t *testing.T) {
	sql := "SELECT * FROM users WHERE name = 'hello world'"
	tok := New(sql, DialectConfig{})
	found := false
	for _, token := range tok.Tokens {
		if token.Type == TokenString && token.Value == "'hello world'" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to find string literal 'hello world'")
	}
}

func TestTokenizerStringEscape(t *testing.T) {
	sql := "SELECT 'It''s a test'"
	tok := New(sql, DialectConfig{})
	found := false
	for _, token := range tok.Tokens {
		if token.Type == TokenString && token.Value == "'It''s a test'" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected to find escaped string literal")
	}
}

func TestTokenizerNumbers(t *testing.T) {
	sql := "SELECT * FROM users LIMIT 10 OFFSET 20"
	tok := New(sql, DialectConfig{})
	found10 := false
	found20 := false
	for _, token := range tok.Tokens {
		if token.Type == TokenNumber && token.Value == "10" {
			found10 = true
		}
		if token.Type == TokenNumber && token.Value == "20" {
			found20 = true
		}
	}
	if !found10 || !found20 {
		t.Fatalf("expected to find numbers 10 and 20")
	}
}
