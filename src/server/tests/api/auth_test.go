package api

// func TestAuthBadPubkey(t *testing.T) {
// 	svr := tests.FakeServer{
// 		A: &tests.FakeAccount{
// 			ID:      cli_pubkey,
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 	}
//
// 	cr := {
// 		AccountID: errPubkey,
// 		Nonce:     cipher.RandByte(xchacha20.NonceSize),
// 	}
//
// 	req := cr.MustJson()
//
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", bytes.NewBuffer(req)))
// 	assert.Equal(t, w.Code, 400)
// }

//
// func TestAuthBadKey(t *testing.T) {
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      cli_pubkey,
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 	}
//
// 	cr := ContentRequest{
// 		AccountID: cli_pubkey,
// 		Nonce:     cipher.RandByte(xchacha20.NonceSize),
// 	}
//
// 	req := cr.MustJson()
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", bytes.NewBuffer(req)))
// 	assert.Equal(t, w.Code, 401)
// }
//
// func TestAuthUnauthorized(t *testing.T) {
// 	svr := FakeServer{
// 		A: &FakeAccount{
// 			ID:      cli_pubkey,
// 			Addr:    "16VV1EbKHK7e3vJu4rhq2dJwegDcbaCcma",
// 			Balance: uint64(0),
// 		},
// 		Seckey: cipher.MustSecKeyFromHex(server_seckey),
// 	}
//
// 	cr := ContentRequest{
// 		AccountID: cli_pubkey,
// 		Nonce:     cipher.RandByte(xchacha20.NonceSize),
// 	}
// 	req := cr.MustJson()
//
// 	w := MockServer(&svr, HttpRequestCase("POST", "/api/v1/accounts", bytes.NewBuffer(req)))
// 	assert.Equal(t, w.Code, 401)
// }
