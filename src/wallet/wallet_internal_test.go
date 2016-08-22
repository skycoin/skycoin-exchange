package wallet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/coin"
)

func TestMustLoad(t *testing.T) {
	RegisterCreator(coin.Bitcoin, NewBtcWltCreator())
	tmpDir := filepath.Join(os.TempDir(), ".wallet1000")
	InitDir(tmpDir)

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
		{"bitcoin_seed5", coin.Bitcoin, "seed5"},
	}

	for _, d := range testData {
		if _, err := New(d.Type, d.Seed); err != nil {
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
