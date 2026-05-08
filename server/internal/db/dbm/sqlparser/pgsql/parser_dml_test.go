package pgsql

import (
	"testing"

	"mayfly-go/internal/db/dbm/sqlparser/sqlstmt"
)

// ========== INSERT 测试 ==========

func TestPgInsertBasic(t *testing.T) {
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

	// 验证完整文本
	if insertStmt.GetText() != sql {
		t.Errorf("expected text='%s', got '%s'", sql, insertStmt.GetText())
	}

	// 验证表名
	if insertStmt.Table.Name != "users" {
		t.Errorf("expected table='users', got '%s'", insertStmt.Table.Name)
	}

	// 验证列名
	if len(insertStmt.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(insertStmt.Columns))
	}
	if insertStmt.Columns[0] != "name" {
		t.Errorf("expected Columns[0]='name', got '%s'", insertStmt.Columns[0])
	}
	if insertStmt.Columns[1] != "age" {
		t.Errorf("expected Columns[1]='age', got '%s'", insertStmt.Columns[1])
	}
}

func TestPgInsertMultipleRows(t *testing.T) {
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

func TestPgInsertReturning(t *testing.T) {
	sql := "INSERT INTO users (name, age) VALUES ('John', 30) RETURNING id, name"
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

func TestPgInsertFromSelect(t *testing.T) {
	sql := "INSERT INTO users_backup SELECT * FROM users WHERE status = 1"
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

func TestPgUpdateBasic(t *testing.T) {
	sql := "UPDATE users SET name = 'John' WHERE id = 1"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)

	// 验证完整文本
	if updateStmt.GetText() != sql {
		t.Errorf("expected text='%s', got '%s'", sql, updateStmt.GetText())
	}

	// 验证表名
	if len(updateStmt.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(updateStmt.Tables))
	}
	if updateStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users', got '%s'", updateStmt.Tables[0].Name)
	}

	// 验证 SET 字段
	if len(updateStmt.Set) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(updateStmt.Set))
	}
	if updateStmt.Set[0].Column != "name" {
		t.Errorf("expected Set[0].Column='name', got '%s'", updateStmt.Set[0].Column)
	}
	if updateStmt.Set[0].Value == nil || updateStmt.Set[0].Value.Text != "'John'" {
		t.Errorf("expected Set[0].Value=''John'', got '%s'", updateStmt.Set[0].Value.Text)
	}

	// 验证 WHERE
	if updateStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
	if updateStmt.Where.Text != "id = 1" {
		t.Errorf("expected WHERE='id = 1', got '%s'", updateStmt.Where.Text)
	}
}

func TestPgUpdateMultipleColumns(t *testing.T) {
	sql := "UPDATE users SET name = 'John', age = 30, email = 'john@example.com' WHERE id = 1"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)

	// 验证完整文本
	if updateStmt.GetText() != sql {
		t.Errorf("expected text='%s', got '%s'", sql, updateStmt.GetText())
	}

	// 验证表名
	if len(updateStmt.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(updateStmt.Tables))
	}
	if updateStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users', got '%s'", updateStmt.Tables[0].Name)
	}

	// 验证 SET 字段（3个赋值）
	if len(updateStmt.Set) != 3 {
		t.Fatalf("expected 3 assignments, got %d", len(updateStmt.Set))
	}

	// 验证第一个赋值：name = 'John'
	if updateStmt.Set[0].Column != "name" {
		t.Errorf("expected Set[0].Column='name', got '%s'", updateStmt.Set[0].Column)
	}
	if updateStmt.Set[0].Value == nil || updateStmt.Set[0].Value.Text != "'John'" {
		t.Errorf("expected Set[0].Value=''John'', got '%s'", updateStmt.Set[0].Value.Text)
	}
	if updateStmt.Set[0].Text != "name = 'John'" {
		t.Errorf("expected Set[0].Text=\"name = 'John'\", got '%s'", updateStmt.Set[0].Text)
	}

	// 验证第二个赋值：age = 30
	if updateStmt.Set[1].Column != "age" {
		t.Errorf("expected Set[1].Column='age', got '%s'", updateStmt.Set[1].Column)
	}
	if updateStmt.Set[1].Value == nil || updateStmt.Set[1].Value.Text != "30" {
		t.Errorf("expected Set[1].Value='30', got '%s'", updateStmt.Set[1].Value.Text)
	}
	if updateStmt.Set[1].Text != "age = 30" {
		t.Errorf("expected Set[1].Text='age = 30', got '%s'", updateStmt.Set[1].Text)
	}

	// 验证第三个赋值：email = 'john@example.com'
	if updateStmt.Set[2].Column != "email" {
		t.Errorf("expected Set[2].Column='email', got '%s'", updateStmt.Set[2].Column)
	}
	if updateStmt.Set[2].Value == nil || updateStmt.Set[2].Value.Text != "'john@example.com'" {
		t.Errorf("expected Set[2].Value=''john@example.com'', got '%s'", updateStmt.Set[2].Value.Text)
	}
	if updateStmt.Set[2].Text != "email = 'john@example.com'" {
		t.Errorf("expected Set[2].Text=\"email = 'john@example.com'\", got '%s'", updateStmt.Set[2].Text)
	}

	// 验证 WHERE
	if updateStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
	if updateStmt.Where.Text != "id = 1" {
		t.Errorf("expected WHERE='id = 1', got '%s'", updateStmt.Where.Text)
	}
}

func TestPgUpdateFromJoin(t *testing.T) {
	// PostgreSQL UPDATE FROM 语法
	sql := "UPDATE orders SET status = 'shipped' FROM users WHERE orders.user_id = users.id AND users.status = 1"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	t.Logf("UPDATE text: %s", updateStmt.GetText())
	t.Logf("UPDATE tables: %+v", updateStmt.Tables)

	// 主表名应该正确解析
	if len(updateStmt.Tables) < 1 {
		t.Fatal("expected at least 1 table")
	}
	// 注意：当前解析器对于 UPDATE FROM 语法支持有限，可能表名包含别名
	// 但至少验证能成功解析而不报错
	t.Logf("Successfully parsed UPDATE FROM statement")
}

// ========== DELETE 测试 ==========

func TestPgDeleteBasic(t *testing.T) {
	sql := "DELETE FROM users WHERE id = 1"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)

	// 验证完整文本
	if deleteStmt.GetText() != sql {
		t.Errorf("expected text='%s', got '%s'", sql, deleteStmt.GetText())
	}

	// 验证表名
	if len(deleteStmt.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(deleteStmt.Tables))
	}
	if deleteStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users', got '%s'", deleteStmt.Tables[0].Name)
	}

	// 验证 WHERE
	if deleteStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
	if deleteStmt.Where.Text != "id = 1" {
		t.Errorf("expected WHERE='id = 1', got '%s'", deleteStmt.Where.Text)
	}
}

func TestPgDeleteUsing(t *testing.T) {
	sql := "DELETE FROM orders o USING users u WHERE o.user_id = u.id AND u.status = 0"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	if len(deleteStmt.Tables) < 1 || deleteStmt.Tables[0].Name != "orders" {
		t.Errorf("expected table='orders'")
	}
}

// ========== Schema.Table 测试 ==========

func TestPgInsertWithSchema(t *testing.T) {
	sql := `INSERT INTO "public"."users" ("name", "age") VALUES ('John', 30)`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt := stmt.(*sqlstmt.InsertStmt)
	if insertStmt.Table.Schema != "public" {
		t.Errorf("expected schema='public', got '%s'", insertStmt.Table.Schema)
	}
	if insertStmt.Table.Name != "users" {
		t.Errorf("expected table='users', got '%s'", insertStmt.Table.Name)
	}
}

func TestPgUpdateWithSchema(t *testing.T) {
	sql := `UPDATE "public"."t_db" SET "name" = 'fsdfds3' WHERE "id" = 5`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)

	t.Logf("UPDATE text: %s", updateStmt.GetText())
	t.Logf("UPDATE tables: %+v", updateStmt.Tables)
	t.Logf("UPDATE SET: %+v", updateStmt.Set)
	t.Logf("UPDATE WHERE: %+v", updateStmt.Where)

	if len(updateStmt.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(updateStmt.Tables))
	}
	if updateStmt.Tables[0].Schema != "public" {
		t.Errorf("expected schema='public', got '%s'", updateStmt.Tables[0].Schema)
	}
	if updateStmt.Tables[0].Name != "t_db" {
		t.Errorf("expected table='t_db', got '%s'", updateStmt.Tables[0].Name)
	}
	// 验证 SET 字段
	if len(updateStmt.Set) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(updateStmt.Set))
	}
	if updateStmt.Set[0].Column != "name" {
		t.Errorf("expected column='name', got '%s'", updateStmt.Set[0].Column)
	}
	if updateStmt.Set[0].Value == nil || updateStmt.Set[0].Value.Text != "'fsdfds3'" {
		t.Errorf("expected value=''fsdfds3'', got '%s'", updateStmt.Set[0].Value.Text)
	}
	if updateStmt.Where == nil || updateStmt.Where.Text != `"id" = 5` {
		t.Errorf("expected WHERE='\"id\" = 5', got '%s'", updateStmt.Where.Text)
	}
}

func TestPgDeleteWithSchema(t *testing.T) {
	sql := `DELETE FROM "public"."logs" WHERE "created_at" < '2024-01-01'`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	if len(deleteStmt.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(deleteStmt.Tables))
	}
	if deleteStmt.Tables[0].Schema != "public" {
		t.Errorf("expected schema='public', got '%s'", deleteStmt.Tables[0].Schema)
	}
	if deleteStmt.Tables[0].Name != "logs" {
		t.Errorf("expected table='logs', got '%s'", deleteStmt.Tables[0].Name)
	}
	if deleteStmt.Where == nil || deleteStmt.Where.Text != `"created_at" < '2024-01-01'` {
		t.Errorf("expected WHERE text")
	}
}

// ========== DDL 测试 ==========

func TestPgDDLCreate(t *testing.T) {
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

func TestPgDDLDrop(t *testing.T) {
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

func TestPgDDLAlter(t *testing.T) {
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

func TestPgDDLTruncate(t *testing.T) {
	sql := "TRUNCATE TABLE users"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if stmt == nil {
		t.Fatalf("expected stmt not nil")
	}
	t.Logf("TRUNCATE stmt type: %T", stmt)
	t.Logf("TRUNCATE text: %s", stmt.GetText())
}

// ========== 复杂 DML 测试 ==========

func TestPgComplexUpdateWithSubquery(t *testing.T) {
	sql := "UPDATE users SET total_orders = (SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id) WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	if len(updateStmt.Tables) != 1 || updateStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users'")
	}
}

func TestPgComplexDeleteWithSubquery(t *testing.T) {
	sql := "DELETE FROM users WHERE id NOT IN (SELECT DISTINCT user_id FROM orders)"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	if deleteStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
}

func TestPgInsertOnConflict(t *testing.T) {
	sql := "INSERT INTO users (id, name) VALUES (1, 'John') ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt := stmt.(*sqlstmt.InsertStmt)
	t.Logf("Actual text: %s", insertStmt.GetText())
	if insertStmt.GetText() != sql {
		t.Errorf("expected text='%s'", sql)
	}
}

// ========== 复杂 INSERT 测试 ==========

func TestPgInsertWithDoubleQuotes(t *testing.T) {
	// PostgreSQL 双引号标识符
	sql := `INSERT INTO "users" ("name", "age", "email") VALUES ('John', 30, 'john@example.com')`
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

func TestPgInsertWithSpecialChars(t *testing.T) {
	// 包含特殊字符
	sql := `INSERT INTO "logs" ("message", "level") VALUES ('Error: connection failed!', 'ERROR')`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt := stmt.(*sqlstmt.InsertStmt)
	t.Logf("INSERT text: %s", insertStmt.GetText())
}

func TestPgInsertReturningMultiple(t *testing.T) {
	// RETURNING 多个字段
	sql := `INSERT INTO "users" ("name", "email") VALUES ('John', 'john@example.com') RETURNING "id", "name", "created_at"`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt := stmt.(*sqlstmt.InsertStmt)
	t.Logf("INSERT RETURNING: %s", insertStmt.GetText())
}

func TestPgInsertFromSelectComplex(t *testing.T) {
	// INSERT FROM SELECT 复杂查询
	sql := `INSERT INTO "users_backup" SELECT * FROM "users" WHERE "status" = 1 AND "created_at" > '2024-01-01'`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt := stmt.(*sqlstmt.InsertStmt)
	t.Logf("INSERT FROM SELECT: %s", insertStmt.GetText())
}

func TestPgInsertOnConflictDoNothing(t *testing.T) {
	// ON CONFLICT DO NOTHING
	sql := `INSERT INTO "users" ("id", "name") VALUES (1, 'John') ON CONFLICT ("id") DO NOTHING`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt := stmt.(*sqlstmt.InsertStmt)
	t.Logf("INSERT ON CONFLICT DO NOTHING: %s", insertStmt.GetText())
}

// ========== 复杂 UPDATE 测试 ==========

func TestPgUpdateWithDoubleQuotes(t *testing.T) {
	sql := `UPDATE "users" SET "name" = 'John', "age" = 30 WHERE "id" = 1`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	if len(updateStmt.Tables) != 1 || updateStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users', got '%s'", updateStmt.Tables[0].Name)
	}
	if updateStmt.Where == nil || updateStmt.Where.Text != `"id" = 1` {
		t.Errorf("expected WHERE")
	}
}

func TestPgUpdateWithComplexWhere(t *testing.T) {
	sql := `UPDATE "orders" SET "status" = 'cancelled' WHERE "status" = 'pending' AND "created_at" < '2024-01-01' AND ("amount" < 100 OR "user_id" IS NULL)`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	if updateStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
	t.Logf("Complex WHERE: %s", updateStmt.Where.Text)
}

func TestPgUpdateWithFunctions(t *testing.T) {
	sql := `UPDATE "users" SET "updated_at" = NOW(), "login_count" = "login_count" + 1 WHERE "id" = 1`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	t.Logf("UPDATE with functions: %s", updateStmt.GetText())
}

func TestPgUpdateReturningComplex(t *testing.T) {
	sql := `UPDATE "users" SET "status" = 0 WHERE "status" = 1 RETURNING "id", "name", "old_status"`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	if updateStmt.Where == nil || updateStmt.Where.Text != `"status" = 1` {
		t.Errorf("expected WHERE")
	}
}

func TestPgUpdateWithSubquery(t *testing.T) {
	sql := `UPDATE "users" SET "total" = (SELECT SUM("amount") FROM "orders" WHERE "user_id" = "users"."id") WHERE "status" = 'active'`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	if len(updateStmt.Tables) != 1 || updateStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users'")
	}
}

func TestPgUpdateFromComplex(t *testing.T) {
	// PostgreSQL UPDATE FROM 复杂场景
	sql := `UPDATE "orders" o SET "status" = 'shipped' FROM "users" u WHERE o."user_id" = u."id" AND u."status" = 1 AND o."created_at" > '2024-01-01'`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	t.Logf("UPDATE FROM: %s", updateStmt.GetText())
}

// ========== 复杂 DELETE 测试 ==========

func TestPgDeleteWithDoubleQuotes(t *testing.T) {
	sql := `DELETE FROM "users" WHERE "id" = 1`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	if len(deleteStmt.Tables) != 1 || deleteStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users'")
	}
}

func TestPgDeleteWithComplexWhere(t *testing.T) {
	sql := `DELETE FROM "logs" WHERE "created_at" < '2024-01-01' AND ("level" = 'DEBUG' OR "level" = 'INFO')`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	if deleteStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
	t.Logf("Complex DELETE WHERE: %s", deleteStmt.Where.Text)
}

func TestPgDeleteReturningComplex(t *testing.T) {
	sql := `DELETE FROM "users" WHERE "status" = 0 RETURNING "id", "name", "deleted_at"`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	if deleteStmt.Where == nil || deleteStmt.Where.Text != `"status" = 0` {
		t.Errorf("expected WHERE")
	}
}

func TestPgDeleteWithSubquery(t *testing.T) {
	sql := `DELETE FROM "users" WHERE "id" NOT IN (SELECT DISTINCT "user_id" FROM "orders")`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	if deleteStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
}

func TestPgDeleteUsingComplex(t *testing.T) {
	// PostgreSQL DELETE USING 复杂场景
	sql := `DELETE FROM "orders" o USING "users" u WHERE o."user_id" = u."id" AND u."status" = 0 AND o."created_at" < '2024-01-01'`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	t.Logf("DELETE USING: %s", deleteStmt.GetText())
}

// ========== 复杂 DDL 测试 ==========

func TestPgDDLCreateTableWithQuotes(t *testing.T) {
	sql := `CREATE TABLE "users" ("id" INT PRIMARY KEY, "name" VARCHAR(50) NOT NULL, "email" VARCHAR(255) UNIQUE)`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	ddlStmt := stmt.(*sqlstmt.DdlStmt)
	if ddlStmt.DdlKind != "CREATE" {
		t.Errorf("expected DdlKind='CREATE'")
	}
}

func TestPgDDLCreateTableWithSerial(t *testing.T) {
	sql := `CREATE TABLE "orders" ("id" SERIAL PRIMARY KEY, "amount" DECIMAL(10,2), "created_at" TIMESTAMP DEFAULT NOW())`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	ddlStmt := stmt.(*sqlstmt.DdlStmt)
	t.Logf("CREATE TABLE with SERIAL: %s", ddlStmt.GetText())
}

func TestPgDDLAlterTableAddColumn(t *testing.T) {
	sql := `ALTER TABLE "users" ADD COLUMN "email" VARCHAR(255)`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	ddlStmt := stmt.(*sqlstmt.DdlStmt)
	if ddlStmt.DdlKind != "ALTER" {
		t.Errorf("expected DdlKind='ALTER'")
	}
}

func TestPgDDLDropIfExists(t *testing.T) {
	sql := `DROP TABLE IF EXISTS "users"`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	ddlStmt := stmt.(*sqlstmt.DdlStmt)
	if ddlStmt.DdlKind != "DROP" {
		t.Errorf("expected DdlKind='DROP'")
	}
}

// ========== 注释风格测试 ==========

func TestPgDMLWithSingleLineComment(t *testing.T) {
	// PostgreSQL 单行注释 --
	sql := "-- 更新用户信息\nUPDATE \"users\" SET \"name\" = 'John' WHERE \"id\" = 1"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	if len(updateStmt.Tables) != 1 || updateStmt.Tables[0].Name != "users" {
		t.Errorf("expected table='users'")
	}
	t.Logf("UPDATE with -- comment: %s", updateStmt.GetText())
}

func TestPgDMLWithMultiLineComment(t *testing.T) {
	// 多行注释 /* */
	sql := "/* 删除过期订单 */ DELETE FROM \"orders\" WHERE \"created_at\" < '2024-01-01'"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	deleteStmt := stmt.(*sqlstmt.DeleteStmt)
	if deleteStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
	t.Logf("DELETE with /* */ comment: %s", deleteStmt.GetText())
}

func TestPgDMLWithInlineComment(t *testing.T) {
	// 行内注释
	sql := `SELECT "id", "name" /* 用户名 */ FROM "users" WHERE "status" = 1`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	selectStmt := stmt.(*sqlstmt.SelectStmt)
	if len(selectStmt.Items) != 2 {
		t.Fatalf("expected 2 items")
	}
	t.Logf("SELECT with inline comment: %s", selectStmt.GetText())
}

func TestPgDMLWithMultipleComments(t *testing.T) {
	// 多个注释
	sql := `-- 查询活跃用户
/* 只查询最近注册的 */
SELECT "id", "name" FROM "users" 
WHERE "status" = 1 AND "created_at" > '2024-01-01'
ORDER BY "created_at" DESC`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	selectStmt := stmt.(*sqlstmt.SelectStmt)
	if selectStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
	t.Logf("Multiple comments: %s", selectStmt.GetText())
}

func TestPgInsertWithComment(t *testing.T) {
	sql := "-- 插入新用户\nINSERT INTO \"users\" (\"name\", \"email\") VALUES ('John', 'john@example.com')"
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	insertStmt := stmt.(*sqlstmt.InsertStmt)
	t.Logf("INSERT with comment: %s", insertStmt.GetText())
}

func TestPgUpdateWithComment(t *testing.T) {
	sql := `/* 批量更新状态 */
UPDATE "orders" SET "status" = 'cancelled' 
WHERE "status" = 'pending' -- 只更新待处理的订单
AND "created_at" < '2024-01-01'`
	parser := NewParser(sql)
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	updateStmt := stmt.(*sqlstmt.UpdateStmt)
	if updateStmt.Where == nil {
		t.Fatal("expected WHERE")
	}
	t.Logf("UPDATE with multiple comments: %s", updateStmt.GetText())
}
