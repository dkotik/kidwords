package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/dkotik/kidwords/store"
	_ "github.com/lib/pq"
)

func TestPostgresStore(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set environment variable DATABASE_URL to run this test")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(fmt.Errorf("cannot create database: %w", err))
	}
	t.Cleanup(func() { db.Close() })

	tableName := store.DefaultTableName + "_test"
	t.Cleanup(func() {
		if _, err = db.Exec(`DROP TABLE ` + tableName); err != nil {
			t.Log(`DROP TABLE ` + tableName)
			t.Fatal("could not clean up test table:", err)
		}
	})

	s, err := store.NewPostgresStore(
		db,
		store.WithTableName(tableName),
		store.WithGuaranteedTable(),
	)
	if err != nil {
		t.Fatal(err)
	}

	RunAllOperations(s)(t)
}
