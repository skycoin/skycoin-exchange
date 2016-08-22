package skycoin_interface

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUnspentOutpts(t *testing.T) {
	addrs := []string{
		"Kb9SqqTVA3XyQjZYb4wYrBVUeZWRKEQyzZ",
		"24VUoHirWUpwJTjLxfMRkKjZZBsqESsagU9",
	}

	outpts, err := GetUnspentOutputs(addrs)
	assert.Nil(t, err)
	d, err := json.MarshalIndent(outpts, "", " ")
	assert.Nil(t, err)
	fmt.Printf(string(d))
}

func TestBroadcastTx(t *testing.T) {
	// utxos, err := GetUnspentOutputs([]string{"UsS43vk2yRqjXvgbwq12Dkjr8cHVTBxYoj"})
	// assert.Nil(t, err)
	// NewTransaction(utxos, keyMap, outs)
}
