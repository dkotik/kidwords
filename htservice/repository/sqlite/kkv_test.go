package sqlite

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"
)

func TestKeyKeyValueRepository(t *testing.T) {
	kkv, err := NewKeyKeyValueRepository(
	// WithConnection("file:test.sqlite3"),
	)
	if err != nil {
		t.Fatal(err)
	}

	first := []byte("first")
	second := []byte("second")
	value := []byte("value")
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	if err = kkv.Set(ctx, first, second, value); err != nil {
		t.Fatal(err)
	}

	if err = kkv.Update(ctx, first, second, func(current []byte) ([]byte, error) {
		if !bytes.Equal(current, value) {
			t.Logf("current: %s", string(current))
			t.Logf("expected: %s", string(value))
			return nil, errors.New("values do not match")
		}
		return []byte("newValue"), err
	}); err != nil {
		t.Fatal(err)
	}

	// time.Sleep(time.Second * 3)

	current, err := kkv.Get(ctx, first, second)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(current, []byte("newValue")) {
		t.Logf("current: %s", string(current))
		t.Fatal(errors.New("values do not match"))
	}

	if err = kkv.Delete(ctx, first, second); err != nil {
		t.Fatal(err)
	}
	current, err = kkv.Get(ctx, first, second)
	if err != nil {
		t.Fatal(err)
	}
	if len(current) > 0 {
		t.Fatalf("deleted entry is still present: %s", string(current))
	}
}
