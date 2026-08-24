package gorocksdb_test

import (
	"testing"

	"gorocksdb/pkg/gorocksdb"
)

func TestScanAfterFlushReturnsPersistedRows(t *testing.T) {
	db, err := gorocksdb.Open(gorocksdb.Options{
		Dir:     t.TempDir(),
		Profile: "test",
		Sync:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	records := []struct {
		key   string
		value string
	}{
		{key: "accounts/1001", value: "active"},
		{key: "accounts/1002", value: "locked"},
		{key: "accounts/1003", value: "pending"},
	}
	for _, record := range records {
		if err := db.Put([]byte(record.key), []byte(record.value)); err != nil {
			t.Fatalf("put %q: %v", record.key, err)
		}
	}
	if err := db.Flush(); err != nil {
		t.Fatalf("flush: %v", err)
	}

	value, err := db.Get([]byte(records[1].key))
	if err != nil {
		t.Fatalf("get persisted row: %v", err)
	}
	if string(value) != records[1].value {
		t.Fatalf("get persisted row = %q, want %q", value, records[1].value)
	}

	rows, err := db.Scan([]byte("accounts/"), []byte("accounts0"), len(records))
	if err != nil {
		t.Fatalf("scan persisted rows: %v", err)
	}
	if len(rows) != len(records) {
		t.Fatalf("scan returned %d rows, want %d", len(rows), len(records))
	}
	for i, record := range records {
		if string(rows[i].Key) != record.key || string(rows[i].Value) != record.value {
			t.Fatalf("row %d = %q/%q, want %q/%q", i, rows[i].Key, rows[i].Value, record.key, record.value)
		}
	}
}
