package gorocksdb_test

import (
	"bytes"
	"testing"

	"gorocksdb/pkg/gorocksdb"
)

func TestPutPreservesSharedInputBuffer(t *testing.T) {
	db, err := gorocksdb.Open(gorocksdb.Options{Dir: t.TempDir(), Profile: "test"})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close: %v", err)
		}
	})

	const keyText = "ledger/42"
	const valueText = "settled-balance"
	frameLen := len(keyText) + len(valueText)
	frame := make([]byte, frameLen, frameLen+8)
	copy(frame, keyText+valueText)
	key := frame[:len(keyText)]
	value := frame[len(keyText):]
	wantFrame := bytes.Clone(frame)
	wantValue := bytes.Clone(value)

	if err := db.Put(key, value); err != nil {
		t.Fatalf("put: %v", err)
	}
	if !bytes.Equal(frame, wantFrame) {
		t.Errorf("Put changed its input frame: got %v, want %v", frame, wantFrame)
	}

	got, err := db.Get([]byte(keyText))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !bytes.Equal(got, wantValue) {
		t.Errorf("Get returned %v, want %v", got, wantValue)
	}
}
