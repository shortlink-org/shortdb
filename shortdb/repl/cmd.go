package repl

import (
	"fmt"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
)

const HISTORY_LIMIT = 100

func (r *Repl) init() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	path := os.TempDir() + "/repl.history"

	// create file if not exist
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm) // #nosec
	if err != nil {
		return fmt.Errorf("failed to open history: %w", err)
	}

	defer func() {
		_ = file.Close() //nolint:errcheck // #nosec
	}()

	// read file
	payload, err := os.ReadFile(path) // #nosec
	if err != nil {
		return fmt.Errorf("failed to read history: %w", err)
	}

	if len(payload) != 0 {
		err = proto.Unmarshal(payload, r.session)
		if err != nil {
			return fmt.Errorf("failed to unmarshal history: %w", err)
		}
	}

	return nil
}

func (r *Repl) helpString() string {
	return fmt.Sprintf(`
ShortDB repl
Enter ".help" for usage hints.
Connected to a transient in-memory database.
Use ".open DATABASENAME" to reopen on a persistent database.

Discovery (SQL)
  Overview of user tables (virtual read-only catalog):
    SELECT name, columns FROM 'shortdb_tables';
    SELECT name FROM 'shortdb_tables' LIMIT 10;
  Same catalog without quotes: FROM shortdbcatalog (parser token must be [a-zA-Z0-9]+).
  Columns are listed as "col type, col type" (integer, string, boolean).
  Catalog WHERE: only name = 'table_name' (quoted literal) is supported.

  Inspect rows (use quoted table names if the name contains underscores):
    SELECT field1, field2 FROM mytable LIMIT 20;
    SELECT id, title FROM 'books_backup' LIMIT 20;

DDL — CREATE TABLE
  CREATE TABLE name ( column type [, column type ...] );
  Column types: integer or int, text or string, boolean or bool (use these spellings).
  Table and column names: letters, digits, underscore; first character not a digit.
  Names shortdb_tables and shortdbcatalog are reserved (catalog only).
  End the statement with a semicolon. Tab completes example templates.

DDL — insert rows
  INSERT INTO table ( col1, col2 ) VALUES ( 'val1', 'val2' );
  Use single-quoted values; column names are identifiers (or quoted like 'id'). End with a semicolon.

Workspace
  f1 — CLI tab (SQL shell + transcript)
  f2 — Observable tab (SELECT + table: catalog, queries, row enter opens table)
  On Observable: click the tab labels ( CLI │ Observable ) to switch (mouse is off on CLI so text selection works; use f1/f2 there)

Dot commands
  .help    — this text
  .tables  — open Observable tab (same as f2)
  .open    — switch database
  .save    — flush / close engine
  .close   — save history and exit

current database: %s
`, r.session.GetCurrentDatabase())
}

func (r *Repl) save() error {
	err := r.engine.Close()
	if err != nil {
		return fmt.Errorf("failed to close engine: %w", err)
	}

	return nil
}

func (r *Repl) close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	path := os.TempDir() + "/repl.history"

	// create a file if not exist
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, os.ModePerm) // #nosec
	if err != nil {
		return fmt.Errorf("failed to open history: %w", err)
	}

	defer func() {
		_ = file.Close() //nolint:errcheck // #nosec
	}()

	// Save last 100 record
	if len(r.session.GetHistory()) > HISTORY_LIMIT {
		r.session.History = r.session.GetHistory()[HISTORY_LIMIT:]
	}

	payload, err := proto.Marshal(r.session)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	_, err = file.Write(payload)
	if err != nil {
		return fmt.Errorf("failed to write history: %w", err)
	}

	return nil
}

func (r *Repl) open(t string) error {
	s := strings.Split(t, " ")
	if len(s) != 2 { //nolint:mnd,goerr113 // ignore
		return ErrStatus
	}

	r.session.CurrentDatabase = s[1]

	return nil
}
