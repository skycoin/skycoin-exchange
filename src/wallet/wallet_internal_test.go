package wallet

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustLoad(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), ".wallet1000")
	InitDir(tmpDir)
	defer func() {
		err := os.RemoveAll(tmpDir)
		assert.Nil(t, err)
	}()

	// create wallets.
	testData := []struct {
		ID   string
		Type string
		Seed string
	}{
		{"bitcoin_seed1", "bitcoin", "seed1"},
		{"bitcoin_seed2", "bitcoin", "seed2"},
		{"bitcoin_seed3", "bitcoin", "seed3"},
		{"bitcoin_seed4", "bitcoin", "seed4"},
		{"skycoin_seed1", "skycoin", "seed1"},
		{"skycoin_seed2", "skycoin", "seed2"},
		{"skycoin_seed3", "skycoin", "seed3"},
	}

	for _, d := range testData {
		if _, err := New(d.Type, d.Seed); err != nil {
			fmt.Println(d.Type, " ", d.Seed)
			t.Error(err)
			return
		}
	}

	gWallets.reset()
	gWallets.mustLoad()
	for _, d := range testData {
		if _, ok := gWallets.Value[d.ID]; !ok {
			t.Errorf("%s not loaded", d.ID)
			return
		}
	}
}
