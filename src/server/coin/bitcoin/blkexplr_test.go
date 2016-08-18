package bitcoin_interface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlkExplrUtxos(t *testing.T) {
	_, err := GetUnspentOutputs([]string{"19EC57DDAtTCVcKENVcd5tbRXk7yKSKvGK"})
	assert.Nil(t, err)
}
