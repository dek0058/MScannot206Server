package table_test

import (
	"MScannot206/shared/table"
	"path/filepath"
	"testing"
)

func TestRepository(t *testing.T) {
	r := &table.Repository{}

	relativePath := "../../data"
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	if err := r.Load(absolutePath); err != nil {
		t.Fatalf("failed to load repository: %v", err)
	}

	_, ok := r.Item.Get("1")
	if !ok {
		t.Log("item with key '1' not found in item table")
	}

	item, ok := r.Item.Get("hair-43")
	if !ok {
		t.Log("item with key 'hair-43' not found in item table")
	} else {
		t.Logf("item with key 'hair-43' found: %+v", item)
	}
}
