// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: commercionetwork/government/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	//grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MsgSetGovAddress struct {
}

func (m *MsgSetGovAddress) Reset()         { *m = MsgSetGovAddress{} }
func (m *MsgSetGovAddress) String() string { return proto.CompactTextString(m) }
func (*MsgSetGovAddress) ProtoMessage()    {}
func (*MsgSetGovAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_bef64d04997e18a4, []int{0}
}
func (m *MsgSetGovAddress) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetGovAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetGovAddress.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetGovAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetGovAddress.Merge(m, src)
}
func (m *MsgSetGovAddress) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetGovAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetGovAddress.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetGovAddress proto.InternalMessageInfo

type MsgSetGovAddressResponse struct {
}

func (m *MsgSetGovAddressResponse) Reset()         { *m = MsgSetGovAddressResponse{} }
func (m *MsgSetGovAddressResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSetGovAddressResponse) ProtoMessage()    {}
func (*MsgSetGovAddressResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_bef64d04997e18a4, []int{1}
}
func (m *MsgSetGovAddressResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetGovAddressResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetGovAddressResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetGovAddressResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetGovAddressResponse.Merge(m, src)
}
func (m *MsgSetGovAddressResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetGovAddressResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetGovAddressResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetGovAddressResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgSetGovAddress)(nil), "commercionetwork.commercionetwork.government.MsgSetGovAddress")
	proto.RegisterType((*MsgSetGovAddressResponse)(nil), "commercionetwork.commercionetwork.government.MsgSetGovAddressResponse")
}

func init() {
	proto.RegisterFile("commercionetwork/government/tx.proto", fileDescriptor_bef64d04997e18a4)
}

var fileDescriptor_bef64d04997e18a4 = []byte{
	// 192 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x49, 0xce, 0xcf, 0xcd,
	0x4d, 0x2d, 0x4a, 0xce, 0xcc, 0xcf, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x4f, 0xcf, 0x2f,
	0x4b, 0x2d, 0xca, 0xcb, 0x4d, 0xcd, 0x2b, 0xd1, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0xd2, 0x41, 0x57, 0xa5, 0x87, 0x21, 0x80, 0xd0, 0xa6, 0x24, 0xc4, 0x25, 0xe0, 0x5b, 0x9c,
	0x1e, 0x9c, 0x5a, 0xe2, 0x9e, 0x5f, 0xe6, 0x98, 0x92, 0x52, 0x94, 0x5a, 0x5c, 0xac, 0x24, 0xc5,
	0x25, 0x81, 0x2e, 0x16, 0x94, 0x5a, 0x5c, 0x90, 0x9f, 0x57, 0x9c, 0x6a, 0x34, 0x9f, 0x91, 0x8b,
	0xd9, 0xb7, 0x38, 0x5d, 0x68, 0x3a, 0x23, 0x17, 0x2f, 0x8a, 0x0a, 0x21, 0x3b, 0x3d, 0x52, 0x2c,
	0xd6, 0x43, 0xb7, 0x41, 0xca, 0x8d, 0x32, 0xfd, 0x30, 0x17, 0x3a, 0x45, 0x9e, 0x78, 0x24, 0xc7,
	0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c,
	0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x7d, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x12, 0xc8, 0x58,
	0x7d, 0x8c, 0xa0, 0xc4, 0x10, 0xa8, 0x40, 0x09, 0xdd, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0x70,
	0x08, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xce, 0x2b, 0x3d, 0x7b, 0x89, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	SetGovAddress(ctx context.Context, in *MsgSetGovAddress, opts ...grpc.CallOption) (*MsgSetGovAddressResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SetGovAddress(ctx context.Context, in *MsgSetGovAddress, opts ...grpc.CallOption) (*MsgSetGovAddressResponse, error) {
	out := new(MsgSetGovAddressResponse)
	err := c.cc.Invoke(ctx, "/commercionetwork.commercionetwork.government.Msg/SetGovAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SetGovAddress(context.Context, *MsgSetGovAddress) (*MsgSetGovAddressResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SetGovAddress(ctx context.Context, req *MsgSetGovAddress) (*MsgSetGovAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetGovAddress not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SetGovAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetGovAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetGovAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/commercionetwork.commercionetwork.government.Msg/SetGovAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetGovAddress(ctx, req.(*MsgSetGovAddress))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "commercionetwork.commercionetwork.government.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetGovAddress",
			Handler:    _Msg_SetGovAddress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "commercionetwork/government/tx.proto",
}

func (m *MsgSetGovAddress) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetGovAddress) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetGovAddress) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgSetGovAddressResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetGovAddressResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetGovAddressResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgSetGovAddress) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgSetGovAddressResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgSetGovAddress) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSetGovAddress: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetGovAddress: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgSetGovAddressResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSetGovAddressResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetGovAddressResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
