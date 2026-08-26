package test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/dkotik/kidwords/store"
	_ "modernc.org/sqlite"
)

func TestSQliteStore(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:?cache=shared&mode=rwc")
	if err != nil {
		t.Fatal(fmt.Errorf("cannot create database: %w", err))
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { db.Close() })

	s, err := store.NewSQLiteStore(db, store.WithGuaranteedTable())
	if err != nil {
		t.Fatal(err)
	}
	RunAllOperations(s)(t)
}
