// Code generated by protoc-gen-go.
// source: pp.deposit.proto
// DO NOT EDIT!

package pp

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type GetDepositAddrReq struct {
	AccountId        *string `protobuf:"bytes,1,opt,name=account_id" json:"account_id,omitempty"`
	CoinType         *string `protobuf:"bytes,2,opt,name=coin_type" json:"coin_type,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GetDepositAddrReq) Reset()                    { *m = GetDepositAddrReq{} }
func (m *GetDepositAddrReq) String() string            { return proto.CompactTextString(m) }
func (*GetDepositAddrReq) ProtoMessage()               {}
func (*GetDepositAddrReq) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *GetDepositAddrReq) GetAccountId() string {
	if m != nil && m.AccountId != nil {
		return *m.AccountId
	}
	return ""
}

func (m *GetDepositAddrReq) GetCoinType() string {
	if m != nil && m.CoinType != nil {
		return *m.CoinType
	}
	return ""
}

type GetDepositAddrRes struct {
	AccountId        *string `protobuf:"bytes,10,opt,name=account_id" json:"account_id,omitempty"`
	CoinType         *string `protobuf:"bytes,11,opt,name=coin_type" json:"coin_type,omitempty"`
	Address          *string `protobuf:"bytes,12,opt,name=address" json:"address,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GetDepositAddrRes) Reset()                    { *m = GetDepositAddrRes{} }
func (m *GetDepositAddrRes) String() string            { return proto.CompactTextString(m) }
func (*GetDepositAddrRes) ProtoMessage()               {}
func (*GetDepositAddrRes) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *GetDepositAddrRes) GetAccountId() string {
	if m != nil && m.AccountId != nil {
		return *m.AccountId
	}
	return ""
}

func (m *GetDepositAddrRes) GetCoinType() string {
	if m != nil && m.CoinType != nil {
		return *m.CoinType
	}
	return ""
}

func (m *GetDepositAddrRes) GetAddress() string {
	if m != nil && m.Address != nil {
		return *m.Address
	}
	return ""
}

func init() {
	proto.RegisterType((*GetDepositAddrReq)(nil), "pp.GetDepositAddrReq")
	proto.RegisterType((*GetDepositAddrRes)(nil), "pp.GetDepositAddrRes")
}

func init() { proto.RegisterFile("pp.deposit.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 126 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x28, 0x28, 0xd0, 0x4b,
	0x49, 0x2d, 0xc8, 0x2f, 0xce, 0x2c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x28,
	0x50, 0xb2, 0xe2, 0x12, 0x74, 0x4f, 0x2d, 0x71, 0x81, 0x88, 0x3b, 0xa6, 0xa4, 0x14, 0x05, 0xa5,
	0x16, 0x0a, 0x09, 0x71, 0x71, 0x25, 0x26, 0x27, 0xe7, 0x97, 0xe6, 0x95, 0xc4, 0x67, 0xa6, 0x48,
	0x30, 0x2a, 0x30, 0x6a, 0x70, 0x0a, 0x09, 0x72, 0x71, 0x26, 0xe7, 0x67, 0xe6, 0xc5, 0x97, 0x54,
	0x16, 0xa4, 0x4a, 0x30, 0x81, 0x84, 0x94, 0xbc, 0x31, 0xf5, 0x16, 0xa3, 0xe9, 0xe5, 0xc2, 0xd4,
	0xcb, 0x0d, 0x16, 0xe2, 0xe7, 0x62, 0x4f, 0x04, 0xea, 0x48, 0x2d, 0x2e, 0x96, 0xe0, 0x01, 0x09,
	0x00, 0x02, 0x00, 0x00, 0xff, 0xff, 0x6c, 0x36, 0x5d, 0xf1, 0x9f, 0x00, 0x00, 0x00,
}