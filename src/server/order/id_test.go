package order

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/stretchr/testify/assert"
)

func TestIdGeneratorEmpty(t *testing.T) {
	idg := newIDGenerator("test/sky")
	closing := make(chan bool)
	go idg.Run(closing)
	id := idg.GetID()
	if id != 1 {
		t.Fatal("id error")
	}
	close(closing)
	// remove file
	path := filepath.Join(orderDir, "test_sky."+idExt)
	time.Sleep(500 * time.Millisecond)
	err := os.Remove(path)
	assert.Nil(t, err)
}

func TestIdGeneratorNoneEmpty(t *testing.T) {
	// prepare data first.
	id := struct {
		ID uint64 `json:"id"`
	}{4}
	// save the id into local disk
	path := filepath.Join(orderDir, "test1_sky."+idExt)
	err := file.SaveJSON(path, id, 0600)
	assert.Nil(t, err)

	// remove file
	defer os.RemoveAll(filepath.Join(orderDir, "test1_sky."+idExt))

	idg := newIDGenerator("test1/sky")
	closing := make(chan bool)
	go idg.Run(closing)
	for i := 1; i < 10; i++ {
		nid := idg.GetID()
		expected := uint64(4 + i)
		if nid != expected {
			t.Fatal("id error")
		}
	}
	close(closing)
	// check the file value.
	time.Sleep(500 * time.Millisecond)
	err = file.LoadJSON(path, &id)
	assert.Nil(t, err)
	if id.ID != 13 {
		t.Fatal("sync file failed")
	}
}
