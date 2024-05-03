// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: commercionetwork/government/query.proto

package types

import (
	context "context"
	fmt "fmt"
	//grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type QueryGovernmentAddrRequest struct {
}

func (m *QueryGovernmentAddrRequest) Reset()         { *m = QueryGovernmentAddrRequest{} }
func (m *QueryGovernmentAddrRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGovernmentAddrRequest) ProtoMessage()    {}
func (*QueryGovernmentAddrRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_88c43dcc4e5acf65, []int{0}
}
func (m *QueryGovernmentAddrRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGovernmentAddrRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGovernmentAddrRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGovernmentAddrRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGovernmentAddrRequest.Merge(m, src)
}
func (m *QueryGovernmentAddrRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGovernmentAddrRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGovernmentAddrRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGovernmentAddrRequest proto.InternalMessageInfo

type QueryGovernmentAddrResponse struct {
	GovernmentAddress string `protobuf:"bytes,1,opt,name=governmentAddress,proto3" json:"governmentAddress,omitempty"`
}

func (m *QueryGovernmentAddrResponse) Reset()         { *m = QueryGovernmentAddrResponse{} }
func (m *QueryGovernmentAddrResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGovernmentAddrResponse) ProtoMessage()    {}
func (*QueryGovernmentAddrResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_88c43dcc4e5acf65, []int{1}
}
func (m *QueryGovernmentAddrResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGovernmentAddrResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGovernmentAddrResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGovernmentAddrResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGovernmentAddrResponse.Merge(m, src)
}
func (m *QueryGovernmentAddrResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGovernmentAddrResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGovernmentAddrResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGovernmentAddrResponse proto.InternalMessageInfo

func (m *QueryGovernmentAddrResponse) GetGovernmentAddress() string {
	if m != nil {
		return m.GovernmentAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*QueryGovernmentAddrRequest)(nil), "commercionetwork.commercionetwork.government.QueryGovernmentAddrRequest")
	proto.RegisterType((*QueryGovernmentAddrResponse)(nil), "commercionetwork.commercionetwork.government.QueryGovernmentAddrResponse")
}

func init() {
	proto.RegisterFile("commercionetwork/government/query.proto", fileDescriptor_88c43dcc4e5acf65)
}

var fileDescriptor_88c43dcc4e5acf65 = []byte{
	// 255 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4f, 0xce, 0xcf, 0xcd,
	0x4d, 0x2d, 0x4a, 0xce, 0xcc, 0xcf, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x4f, 0xcf, 0x2f,
	0x4b, 0x2d, 0xca, 0xcb, 0x4d, 0xcd, 0x2b, 0xd1, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xd2, 0x41, 0x57, 0xa8, 0x87, 0x21, 0x80, 0xd0, 0x29, 0x25, 0x93, 0x9e,
	0x9f, 0x9f, 0x9e, 0x93, 0xaa, 0x9f, 0x58, 0x90, 0xa9, 0x9f, 0x98, 0x97, 0x97, 0x5f, 0x92, 0x58,
	0x92, 0x99, 0x9f, 0x57, 0x0c, 0x31, 0x4b, 0x49, 0x86, 0x4b, 0x2a, 0x10, 0x64, 0xb4, 0x3b, 0x5c,
	0x83, 0x63, 0x4a, 0x4a, 0x51, 0x50, 0x6a, 0x61, 0x69, 0x6a, 0x71, 0x89, 0x92, 0x37, 0x97, 0x34,
	0x56, 0xd9, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x21, 0x1d, 0x2e, 0xc1, 0x74, 0x14, 0x99, 0xd4,
	0xe2, 0x62, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x4c, 0x09, 0xa3, 0xe7, 0x8c, 0x5c, 0xac,
	0x60, 0xd3, 0x84, 0xee, 0x32, 0x72, 0xf1, 0xa1, 0x1a, 0x29, 0xe4, 0xa1, 0x47, 0x8a, 0xa7, 0xf4,
	0x70, 0xbb, 0x59, 0xca, 0x93, 0x0a, 0x26, 0x41, 0xfc, 0xa7, 0x64, 0xd6, 0x74, 0xf9, 0xc9, 0x64,
	0x26, 0x03, 0x21, 0x3d, 0x7d, 0x7c, 0x51, 0x83, 0xe1, 0x53, 0xa7, 0xc8, 0x13, 0x8f, 0xe4, 0x18,
	0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5,
	0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0xb2, 0x4f, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0x02, 0xb9, 0x08,
	0xd3, 0x4c, 0x0c, 0x81, 0x0a, 0x64, 0x6b, 0x4a, 0x2a, 0x0b, 0x52, 0x8b, 0x93, 0xd8, 0xc0, 0xd1,
	0x66, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x05, 0xbc, 0xa2, 0x8d, 0x2d, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Queries the government address.
	GovernmentAddr(ctx context.Context, in *QueryGovernmentAddrRequest, opts ...grpc.CallOption) (*QueryGovernmentAddrResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) GovernmentAddr(ctx context.Context, in *QueryGovernmentAddrRequest, opts ...grpc.CallOption) (*QueryGovernmentAddrResponse, error) {
	out := new(QueryGovernmentAddrResponse)
	err := c.cc.Invoke(ctx, "/commercionetwork.commercionetwork.government.Query/GovernmentAddr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries the government address.
	GovernmentAddr(context.Context, *QueryGovernmentAddrRequest) (*QueryGovernmentAddrResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) GovernmentAddr(ctx context.Context, req *QueryGovernmentAddrRequest) (*QueryGovernmentAddrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GovernmentAddr not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_GovernmentAddr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGovernmentAddrRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GovernmentAddr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/commercionetwork.commercionetwork.government.Query/GovernmentAddr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GovernmentAddr(ctx, req.(*QueryGovernmentAddrRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "commercionetwork.commercionetwork.government.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GovernmentAddr",
			Handler:    _Query_GovernmentAddr_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "commercionetwork/government/query.proto",
}

func (m *QueryGovernmentAddrRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGovernmentAddrRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGovernmentAddrRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryGovernmentAddrResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGovernmentAddrResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGovernmentAddrResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.GovernmentAddress) > 0 {
		i -= len(m.GovernmentAddress)
		copy(dAtA[i:], m.GovernmentAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.GovernmentAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryGovernmentAddrRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryGovernmentAddrResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.GovernmentAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryGovernmentAddrRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryGovernmentAddrRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGovernmentAddrRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryGovernmentAddrResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryGovernmentAddrResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGovernmentAddrResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GovernmentAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GovernmentAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
