package skycoin_interface

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/skycoin/skycoin-exchange/src/pp"
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

func TestGetOutput(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		args    args
		want    *pp.Output
		wantErr error
	}{
		// TODO: Add test cases.
		{
			"normal",
			args{
				"a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c5",
			},
			&pp.Output{
				Time:          pp.PtrUint64(1477037552),
				SrcBlockSeq:   pp.PtrUint64(443),
				SrcTx:         pp.PtrString("b8ca61c0788bd711c89563f9bc60add172ee01b543ea5dcb1955c51bbfcbbaa2"),
				OwnerAddress:  pp.PtrString("cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW"),
				Coins:         pp.PtrUint64(1000000),
				Hours:         pp.PtrUint64(7),
				SpentBlockSeq: pp.PtrUint64(450),
				SpentTx:       pp.PtrString("b1481d614ffcc27408fe2131198d9d2821c78601a0aa23d8e9965b2a5196edc0"),
			},
			nil,
		},
		{
			"invalid uxhash len",
			args{
				"123",
			},
			nil,
			errors.New("invalid output hash, encoding/hex: odd length hex string"),
		},
		{
			"invalid uxhash ",
			args{
				"a57c038591f862b8fada57e496ef948183b153348d7932921f865a8541a477c7",
			},
			nil,
			errors.New("not found\n"),
		},
	}
	for _, tt := range tests {
		got, err := GetOutput(tt.args.hash)
		if !reflect.DeepEqual(err, tt.wantErr) {
			t.Errorf("%q. GetOutput() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetOutput() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
