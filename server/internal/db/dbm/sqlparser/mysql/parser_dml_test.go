package mysql

import (
	"testing"

	"mayfly-go/internal/db/dbm/sqlparser/sqlstmt"
)

// ========== INSERT 测试 ==========

func TestInsertBasic(t *testing.T) {
	sql := "INSERT INTO users (name, age) VALUES ('John', 30)"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt, ok := stmt.(*sqlstmt.InsertStmt)
	if !ok {
		t.Fatal("expected InsertStmt")
	}
	if insertStmt.GetText() != sql {
		t.Errorf("expected text='%s', got '%s'", sql, insertStmt.GetText())
	}
}

func TestInsertMultipleRows(t *testing.T) {
	sql := "INSERT INTO users (name, age) VALUES ('John', 30), ('Jane', 25), ('Bob', 35)"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt := stmt.(*sqlstmt.InsertStmt)
	if insertStmt.GetText() != sql {
		t.Errorf("expected text='%s'", sql)
	}
}

// ========== UPDATE 测试 ==========

func TestUpdateBasic(t *testing.T) {
	sql := "UPDATE users SET name = 'John' WHERE id = 1"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)

	if updateStmt.GetText() != sql {
		t.Errorf("expected text='%s', got '%s'", sql, updateStmt.GetText())
	}
	if len(updateStmt.Tables) != 1 || updateStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users'")
	}
	if updateStmt.Where == nil || updateStmt.Where.Text != "id = 1" {
		t.Errorf("expected WHERE='id = 1'")
	}
}

func TestUpdateMultipleColumns(t *testing.T) {
	sql := "UPDATE users SET name = 'John', age = 30, email = 'john@example.com' WHERE id = 1"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	if updateStmt.Where == nil || updateStmt.Where.Text != "id = 1" {
		t.Errorf("expected WHERE='id = 1'")
	}
}

func TestUpdateWithOrderByLimit(t *testing.T) {
	sql := "UPDATE users SET status = 0 WHERE status = 1 ORDER BY id LIMIT 10"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)

	if updateStmt.Where == nil || updateStmt.Where.Text != "status = 1" {
		t.Errorf("expected WHERE='status = 1'")
	}
}

// ========== DELETE 测试 ==========

func TestDeleteBasic(t *testing.T) {
	sql := "DELETE FROM users WHERE id = 1"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)

	if deleteStmt.GetText() != sql {
		t.Errorf("expected text='%s', got '%s'", sql, deleteStmt.GetText())
	}
	if len(deleteStmt.Tables) != 1 || deleteStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users'")
	}
	if deleteStmt.Where == nil || deleteStmt.Where.Text != "id = 1" {
		t.Errorf("expected WHERE='id = 1'")
	}
}

func TestDeleteMultipleTables(t *testing.T) {
	sql := "DELETE t1, t2 FROM users t1 JOIN orders t2 ON t1.id = t2.user_id WHERE t1.status = 0"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)

	if len(deleteStmt.Tables) < 1 {
		t.Errorf("expected at least 1 table")
	}
}

func TestDeleteWithOrderByLimit(t *testing.T) {
	sql := "DELETE FROM users WHERE status = 0 ORDER BY id LIMIT 100"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)

	if deleteStmt.Where == nil || deleteStmt.Where.Text != "status = 0" {
		t.Errorf("expected WHERE='status = 0'")
	}
}

// ========== DDL 测试 ==========

func TestDDLCreate(t *testing.T) {
	sql := "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50), age INT)"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	ddlStmt := stmt.(*sqlstmt.DdlStmt)

	if ddlStmt.DdlKind != "CREATE" {
		t.Errorf("expected DdlKind='CREATE', got '%s'", ddlStmt.DdlKind)
	}
	if ddlStmt.GetText() != sql {
		t.Errorf("expected text='%s'", sql)
	}
}

func TestDDLDrop(t *testing.T) {
	sql := "DROP TABLE users"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	ddlStmt := stmt.(*sqlstmt.DdlStmt)

	if ddlStmt.DdlKind != "DROP" {
		t.Errorf("expected DdlKind='DROP', got '%s'", ddlStmt.DdlKind)
	}
}

func TestDDLAlter(t *testing.T) {
	sql := "ALTER TABLE users ADD COLUMN email VARCHAR(255)"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	ddlStmt := stmt.(*sqlstmt.DdlStmt)

	if ddlStmt.DdlKind != "ALTER" {
		t.Errorf("expected DdlKind='ALTER', got '%s'", ddlStmt.DdlKind)
	}
}

func TestDDLTruncate(t *testing.T) {
	sql := "TRUNCATE TABLE users"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// TRUNCATE 可能返回 DdlStmt 或其他类型
	if stmt == nil {
		t.Fatalf("expected stmt not nil")
	}
	t.Logf("TRUNCATE stmt type: %T", stmt)
	t.Logf("TRUNCATE text: %s", stmt.GetText())
}
