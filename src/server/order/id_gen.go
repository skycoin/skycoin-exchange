package order

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/util/file"
)

// id generator
type IDGenerator struct {
	IDC  chan uint64
	Path string
}

func newIDGenerator(cp string) *IDGenerator {
	name := strings.Replace(cp, "/", "_", 1)
	return &IDGenerator{
		IDC:  make(chan uint64),
		Path: filepath.Join(orderDir, name+"."+idExt),
	}
}

func (ig IDGenerator) Run(closing chan bool) {
	id := struct {
		ID uint64 `json:"id"`
	}{}
	if _, err := os.Stat(ig.Path); !os.IsNotExist(err) {
		if err := file.LoadJSON(ig.Path, &id); err != nil {
			panic(err)
		}
	}
	for {
		select {
		case <-closing:
			return
		default:
			id.ID += 1
			ig.IDC <- id.ID

			if err := file.SaveJSON(ig.Path, id, 0600); err != nil {
				panic(err)
			}
		}
	}
}

func (ig IDGenerator) GetID() uint64 {
	return <-ig.IDC
}
