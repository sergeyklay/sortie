package persistence

import _ "embed"

// Migration represents a single numbered schema migration. Migrations are
// applied sequentially by [Store.Migrate]. The SQL field may contain multiple
// statements separated by semicolons.
type Migration struct {
	Version     int
	Description string
	SQL         string
}

//go:embed sql/001_initial.sql
var migration001SQL string

var migrations = []Migration{
	{Version: 1, Description: "core persistence tables", SQL: migration001SQL},
}
