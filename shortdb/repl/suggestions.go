package repl

// SuggestionTexts returns autocomplete entries for the REPL input (Bubble Tea textinput).
// Textinput matches suggestions as a case-insensitive prefix of the whole line.
func SuggestionTexts() []string {
	return []string{
		".help",
		".tables",
		".open",
		".save",
		".close",
		"SELECT name, columns FROM 'shortdb_tables'",
		"CREATE TABLE users ( id integer, name text );",
		"CREATE TABLE users ( id int, title string, active boolean );",
		"CREATE TABLE items ( sku text, qty integer, in_stock bool );",
		"create table",
		"DROP TABLE users;",
		"drop table",
		"SELECT",
		"UPDATE",
		"INSERT INTO users ( id, name, active ) VALUES ( '1', 'Alice', 'true' );",
		"INSERT INTO users ( 'id', 'name', 'active' ) VALUES ( '2', 'Bob', 'false' );",
		"INSERT INTO items ( sku, qty ) VALUES ( 'ABC', '10' );",
		"INSERT INTO",
		"insert into",
		"DELETE",
	}
}
