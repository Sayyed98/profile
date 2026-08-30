package db_test

import (
	"testing"

	"github.com/mohdhujaifa/profile/internal/db"
	"github.com/stretchr/testify/require"
)

func TestSchemaContainsCoreTables(t *testing.T) {
	raw, err := db.SchemaSQL()
	require.NoError(t, err)
	for _, table := range []string{"profiles", "skills", "experiences", "projects", "education", "contact_messages"} {
		require.Contains(t, raw, table)
	}
	stmts := db.SplitStatements(raw)
	require.GreaterOrEqual(t, len(stmts), 8)
}
