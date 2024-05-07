// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: commercionetwork/vbr/vbr.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type VbrPool struct {
	Amount github_com_cosmos_cosmos_sdk_types.DecCoins `protobuf:"bytes,1,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.DecCoins" json:"amount" yaml:"amount"`
}

func (m *VbrPool) Reset()         { *m = VbrPool{} }
func (m *VbrPool) String() string { return proto.CompactTextString(m) }
func (*VbrPool) ProtoMessage()    {}
func (*VbrPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_f2f117998702db8c, []int{0}
}
func (m *VbrPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *VbrPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_VbrPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *VbrPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VbrPool.Merge(m, src)
}
func (m *VbrPool) XXX_Size() int {
	return m.Size()
}
func (m *VbrPool) XXX_DiscardUnknown() {
	xxx_messageInfo_VbrPool.DiscardUnknown(m)
}

var xxx_messageInfo_VbrPool proto.InternalMessageInfo

func (m *VbrPool) GetAmount() github_com_cosmos_cosmos_sdk_types.DecCoins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func init() {
	proto.RegisterType((*VbrPool)(nil), "commercionetwork.vbr.VbrPool")
}

func init() { proto.RegisterFile("commercionetwork/vbr/vbr.proto", fileDescriptor_f2f117998702db8c) }

var fileDescriptor_f2f117998702db8c = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xce, 0xcf, 0xcd,
	0x4d, 0x2d, 0x4a, 0xce, 0xcc, 0xcf, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x2f, 0x4b, 0x2a,
	0x02, 0x61, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x11, 0x74, 0x79, 0xbd, 0xb2, 0xa4, 0x22,
	0x29, 0x91, 0xf4, 0xfc, 0xf4, 0x7c, 0xb0, 0x02, 0x7d, 0x10, 0x0b, 0xa2, 0x56, 0x4a, 0x2e, 0x39,
	0xbf, 0x38, 0x37, 0xbf, 0x58, 0x3f, 0x29, 0xb1, 0x38, 0x55, 0xbf, 0xcc, 0x30, 0x29, 0xb5, 0x24,
	0xd1, 0x50, 0x3f, 0x39, 0x3f, 0x33, 0x0f, 0x22, 0xaf, 0xd4, 0xca, 0xc8, 0xc5, 0x1e, 0x96, 0x54,
	0x14, 0x90, 0x9f, 0x9f, 0x23, 0x54, 0xc5, 0xc5, 0x96, 0x98, 0x9b, 0x5f, 0x9a, 0x57, 0x22, 0xc1,
	0xa8, 0xc0, 0xac, 0xc1, 0x6d, 0x24, 0xa3, 0x07, 0xd1, 0xac, 0x07, 0xd2, 0xac, 0x07, 0xd5, 0xac,
	0xe7, 0x92, 0x9a, 0xec, 0x9c, 0x9f, 0x99, 0xe7, 0xe4, 0x72, 0xe2, 0x9e, 0x3c, 0xc3, 0xa7, 0x7b,
	0xf2, 0xbc, 0x95, 0x89, 0xb9, 0x39, 0x56, 0x4a, 0x10, 0x9d, 0x4a, 0xab, 0xee, 0xcb, 0x6b, 0xa7,
	0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x43, 0x6d, 0x87, 0x50, 0xba, 0xc5,
	0x29, 0xd9, 0xfa, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x30, 0x43, 0x8a, 0x83, 0xa0, 0x36, 0x3a, 0x05,
	0x9d, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb,
	0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x05, 0x8a, 0x71, 0x68, 0x01,
	0x83, 0x21, 0x50, 0x01, 0x0e, 0x2b, 0xb0, 0x25, 0x49, 0x6c, 0x60, 0x2f, 0x1a, 0x03, 0x02, 0x00,
	0x00, 0xff, 0xff, 0xd7, 0x0b, 0xf9, 0xbe, 0x50, 0x01, 0x00, 0x00,
}

func (m *VbrPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VbrPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *VbrPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintVbr(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintVbr(dAtA []byte, offset int, v uint64) int {
	offset -= sovVbr(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *VbrPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovVbr(uint64(l))
		}
	}
	return n
}

func sovVbr(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVbr(x uint64) (n int) {
	return sovVbr(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *VbrPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVbr
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: VbrPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VbrPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVbr
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthVbr
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVbr
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.DecCoin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVbr(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVbr
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipVbr(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVbr
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowVbr
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowVbr
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthVbr
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVbr
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVbr
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVbr        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVbr          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVbr = fmt.Errorf("proto: unexpected end of group")
)
