package wallet

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/coin"
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
		Type coin.Type
		Seed string
	}{
		{"bitcoin_seed1", coin.Bitcoin, "seed1"},
		{"bitcoin_seed2", coin.Bitcoin, "seed2"},
		{"bitcoin_seed3", coin.Bitcoin, "seed3"},
		{"bitcoin_seed4", coin.Bitcoin, "seed4"},
		{"skycoin_seed1", coin.Skycoin, "seed1"},
		{"skycoin_seed2", coin.Skycoin, "seed2"},
		{"skycoin_seed3", coin.Skycoin, "seed3"},
	}

	for _, d := range testData {
		if _, err := New(d.Type, d.Seed); err != nil {
			fmt.Println(d.Type.String(), " ", d.Seed)
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
