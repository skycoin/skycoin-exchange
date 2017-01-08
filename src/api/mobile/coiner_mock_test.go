/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package mobile

import (
	"fmt"
	mock "github.com/stretchr/testify/mock"

	coin "github.com/skycoin/skycoin-exchange/src/coin"
)

// CoinerMock mock
type CoinerMock struct {
	mock.Mock
}

func NewCoinerMock() *CoinerMock {
	return &CoinerMock{}
}

// BroadcastTx mocked method
func (m *CoinerMock) BroadcastTx(p0 string) (string, error) {

	ret := m.Called(p0)

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// CreateRawTx mocked method
func (m *CoinerMock) CreateRawTx(p0 []coin.TxIn, p1 coin.GetPrivKey, p2 interface{}) (string, error) {

	ret := m.Called(p0, p1, p2)

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetBalance mocked method
func (m *CoinerMock) GetBalance(p0 []string) (uint64, error) {

	ret := m.Called(p0)

	var r0 uint64
	switch res := ret.Get(0).(type) {
	case nil:
	case uint64:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetNodeAddr mocked method
func (m *CoinerMock) GetNodeAddr() string {

	ret := m.Called()

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// GetOutputByID mocked method
func (m *CoinerMock) GetOutputByID(p0 string) (string, error) {

	ret := m.Called(p0)

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// GetTransactionByID mocked method
func (m *CoinerMock) GetTransactionByID(p0 string) (string, error) {

	ret := m.Called(p0)

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// Name mocked method
func (m *CoinerMock) Name() string {

	ret := m.Called()

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}

// PrepareTx mocked method
func (m *CoinerMock) PrepareTx(p0 interface{}) ([]coin.TxIn, interface{}, error) {

	ret := m.Called(p0)

	var r0 []coin.TxIn
	switch res := ret.Get(0).(type) {
	case nil:
	case []coin.TxIn:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 interface{}
	switch res := ret.Get(1).(type) {
	case nil:
	case interface{}:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r2 error
	switch res := ret.Get(2).(type) {
	case nil:
	case error:
		r2 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1, r2

}

// Send mocked method
func (m *CoinerMock) Send(p0 string, p1 string, p2 string, p3 ...Option) (string, error) {

	ret := m.Called(p0, p1, p2, p3)

	var r0 string
	switch res := ret.Get(0).(type) {
	case nil:
	case string:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	var r1 error
	switch res := ret.Get(1).(type) {
	case nil:
	case error:
		r1 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0, r1

}

// ValidateAddr mocked method
func (m *CoinerMock) ValidateAddr(p0 string) error {

	ret := m.Called(p0)

	var r0 error
	switch res := ret.Get(0).(type) {
	case nil:
	case error:
		r0 = res
	default:
		panic(fmt.Sprintf("unexpected type: %v", res))
	}

	return r0

}
