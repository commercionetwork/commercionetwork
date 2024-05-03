// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: commercionetwork/did/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
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

type QueryResolveIdentityRequest struct {
	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (m *QueryResolveIdentityRequest) Reset()         { *m = QueryResolveIdentityRequest{} }
func (m *QueryResolveIdentityRequest) String() string { return proto.CompactTextString(m) }
func (*QueryResolveIdentityRequest) ProtoMessage()    {}
func (*QueryResolveIdentityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_82b34d99399b1471, []int{0}
}
func (m *QueryResolveIdentityRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryResolveIdentityRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryResolveIdentityRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryResolveIdentityRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResolveIdentityRequest.Merge(m, src)
}
func (m *QueryResolveIdentityRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryResolveIdentityRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResolveIdentityRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResolveIdentityRequest proto.InternalMessageInfo

func (m *QueryResolveIdentityRequest) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

type QueryResolveIdentityResponse struct {
	Identity *Identity `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
}

func (m *QueryResolveIdentityResponse) Reset()         { *m = QueryResolveIdentityResponse{} }
func (m *QueryResolveIdentityResponse) String() string { return proto.CompactTextString(m) }
func (*QueryResolveIdentityResponse) ProtoMessage()    {}
func (*QueryResolveIdentityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_82b34d99399b1471, []int{1}
}
func (m *QueryResolveIdentityResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryResolveIdentityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryResolveIdentityResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryResolveIdentityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResolveIdentityResponse.Merge(m, src)
}
func (m *QueryResolveIdentityResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryResolveIdentityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResolveIdentityResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResolveIdentityResponse proto.InternalMessageInfo

func (m *QueryResolveIdentityResponse) GetIdentity() *Identity {
	if m != nil {
		return m.Identity
	}
	return nil
}

type QueryResolveIdentityHistoryRequest struct {
	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (m *QueryResolveIdentityHistoryRequest) Reset()         { *m = QueryResolveIdentityHistoryRequest{} }
func (m *QueryResolveIdentityHistoryRequest) String() string { return proto.CompactTextString(m) }
func (*QueryResolveIdentityHistoryRequest) ProtoMessage()    {}
func (*QueryResolveIdentityHistoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_82b34d99399b1471, []int{2}
}
func (m *QueryResolveIdentityHistoryRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryResolveIdentityHistoryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryResolveIdentityHistoryRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryResolveIdentityHistoryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResolveIdentityHistoryRequest.Merge(m, src)
}
func (m *QueryResolveIdentityHistoryRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryResolveIdentityHistoryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResolveIdentityHistoryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResolveIdentityHistoryRequest proto.InternalMessageInfo

func (m *QueryResolveIdentityHistoryRequest) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

type QueryResolveIdentityHistoryResponse struct {
	Identities []*Identity `protobuf:"bytes,1,rep,name=identities,proto3" json:"identities,omitempty"`
}

func (m *QueryResolveIdentityHistoryResponse) Reset()         { *m = QueryResolveIdentityHistoryResponse{} }
func (m *QueryResolveIdentityHistoryResponse) String() string { return proto.CompactTextString(m) }
func (*QueryResolveIdentityHistoryResponse) ProtoMessage()    {}
func (*QueryResolveIdentityHistoryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_82b34d99399b1471, []int{3}
}
func (m *QueryResolveIdentityHistoryResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryResolveIdentityHistoryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryResolveIdentityHistoryResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryResolveIdentityHistoryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResolveIdentityHistoryResponse.Merge(m, src)
}
func (m *QueryResolveIdentityHistoryResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryResolveIdentityHistoryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResolveIdentityHistoryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResolveIdentityHistoryResponse proto.InternalMessageInfo

func (m *QueryResolveIdentityHistoryResponse) GetIdentities() []*Identity {
	if m != nil {
		return m.Identities
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryResolveIdentityRequest)(nil), "commercionetwork.commercionetwork.did.QueryResolveIdentityRequest")
	proto.RegisterType((*QueryResolveIdentityResponse)(nil), "commercionetwork.commercionetwork.did.QueryResolveIdentityResponse")
	proto.RegisterType((*QueryResolveIdentityHistoryRequest)(nil), "commercionetwork.commercionetwork.did.QueryResolveIdentityHistoryRequest")
	proto.RegisterType((*QueryResolveIdentityHistoryResponse)(nil), "commercionetwork.commercionetwork.did.QueryResolveIdentityHistoryResponse")
}

func init() { proto.RegisterFile("commercionetwork/did/query.proto", fileDescriptor_82b34d99399b1471) }

var fileDescriptor_82b34d99399b1471 = []byte{
	// 390 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0xc1, 0x4e, 0xea, 0x40,
	0x14, 0x86, 0x19, 0xc8, 0xbd, 0xe1, 0xce, 0x4d, 0xee, 0x4d, 0x66, 0x45, 0x2a, 0x69, 0x48, 0x09,
	0xd1, 0x98, 0xd0, 0x09, 0xa8, 0x89, 0x6b, 0x64, 0x61, 0x75, 0x61, 0xec, 0xd2, 0x5d, 0x4b, 0x27,
	0x65, 0x02, 0xf4, 0x94, 0xce, 0x80, 0x12, 0xe3, 0xc6, 0x27, 0x30, 0xf1, 0x3d, 0x7c, 0x08, 0x57,
	0x2e, 0x49, 0xdc, 0xb8, 0xd3, 0x80, 0x0f, 0x62, 0x68, 0x0b, 0x2a, 0x54, 0x42, 0xc4, 0xed, 0xe9,
	0x7f, 0xbe, 0xff, 0xff, 0x7b, 0x5a, 0x5c, 0x68, 0x40, 0xa7, 0xc3, 0x82, 0x06, 0x07, 0x8f, 0xc9,
	0x73, 0x08, 0x5a, 0xd4, 0xe1, 0x0e, 0xed, 0xf6, 0x58, 0x30, 0xd0, 0xfd, 0x00, 0x24, 0x90, 0xd2,
	0xbc, 0x42, 0x5f, 0x18, 0x38, 0xdc, 0x51, 0xf2, 0x2e, 0x80, 0xdb, 0x66, 0xd4, 0xf2, 0x39, 0xb5,
	0x3c, 0x0f, 0xa4, 0x25, 0x39, 0x78, 0x22, 0x82, 0x28, 0xdb, 0x0d, 0x10, 0x1d, 0x10, 0xd4, 0xb6,
	0x04, 0x8b, 0xe8, 0xb4, 0x5f, 0xb1, 0x99, 0xb4, 0x2a, 0xd4, 0xb7, 0x5c, 0xee, 0x85, 0xe2, 0x58,
	0x5b, 0x4c, 0x8c, 0xc4, 0x1d, 0xe6, 0x49, 0x2e, 0xe3, 0x54, 0x5a, 0x19, 0x6f, 0x9c, 0x4e, 0x30,
	0x26, 0x13, 0xd0, 0xee, 0x33, 0x23, 0x7e, 0x6a, 0xb2, 0x6e, 0x8f, 0x09, 0x49, 0xfe, 0xe1, 0xb4,
	0x51, 0xcf, 0xa1, 0x02, 0xda, 0xfa, 0x63, 0xa6, 0x8d, 0xba, 0xd6, 0xc2, 0xf9, 0x64, 0xb9, 0xf0,
	0xc1, 0x13, 0x8c, 0x1c, 0xe3, 0xec, 0xd4, 0x20, 0xdc, 0xfa, 0x5b, 0xa5, 0xfa, 0x4a, 0xbd, 0xf5,
	0x19, 0x6a, 0x06, 0xd0, 0x76, 0xb1, 0x96, 0x64, 0x76, 0xc8, 0x85, 0x84, 0xe0, 0xcb, 0x88, 0x7d,
	0x5c, 0x5c, 0xba, 0x15, 0x27, 0x3d, 0xc1, 0x38, 0x36, 0xe2, 0x4c, 0xe4, 0x50, 0x21, 0xf3, 0x9d,
	0xac, 0x1f, 0x10, 0xd5, 0xbb, 0x0c, 0xfe, 0x15, 0x1a, 0x93, 0x7b, 0x84, 0xb3, 0x53, 0x09, 0xa9,
	0xad, 0xc8, 0x5c, 0x72, 0x05, 0xe5, 0x60, 0x2d, 0x46, 0x54, 0x58, 0x2b, 0x5f, 0x3f, 0xbe, 0xde,
	0xa6, 0x37, 0x49, 0x89, 0x26, 0x7e, 0x17, 0x97, 0x46, 0xfd, 0x8a, 0xbe, 0xd7, 0x21, 0xcf, 0x08,
	0xff, 0x9f, 0x7b, 0x77, 0xc4, 0x58, 0x23, 0xc7, 0xe7, 0xab, 0x29, 0x47, 0x3f, 0x81, 0x8a, 0x9b,
	0xed, 0x85, 0xcd, 0x28, 0x29, 0xaf, 0xd4, 0x8c, 0x36, 0xa3, 0xf5, 0x9a, 0xf9, 0x30, 0x52, 0xd1,
	0x70, 0xa4, 0xa2, 0x97, 0x91, 0x8a, 0x6e, 0xc6, 0x6a, 0x6a, 0x38, 0x56, 0x53, 0x4f, 0x63, 0x35,
	0x75, 0xb6, 0xef, 0x72, 0xd9, 0xec, 0xd9, 0x93, 0x44, 0x8b, 0xc8, 0x85, 0xc1, 0x45, 0xe8, 0x22,
	0x07, 0x3e, 0x13, 0xf6, 0xef, 0xf0, 0xaf, 0xda, 0x79, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xd2, 0x37,
	0xd7, 0x11, 0x0f, 0x04, 0x00, 0x00,
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
	// Queries a DidDocument by id.
	Identity(ctx context.Context, in *QueryResolveIdentityRequest, opts ...grpc.CallOption) (*QueryResolveIdentityResponse, error)
	IdentityHistory(ctx context.Context, in *QueryResolveIdentityHistoryRequest, opts ...grpc.CallOption) (*QueryResolveIdentityHistoryResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Identity(ctx context.Context, in *QueryResolveIdentityRequest, opts ...grpc.CallOption) (*QueryResolveIdentityResponse, error) {
	out := new(QueryResolveIdentityResponse)
	err := c.cc.Invoke(ctx, "/commercionetwork.commercionetwork.did.Query/Identity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) IdentityHistory(ctx context.Context, in *QueryResolveIdentityHistoryRequest, opts ...grpc.CallOption) (*QueryResolveIdentityHistoryResponse, error) {
	out := new(QueryResolveIdentityHistoryResponse)
	err := c.cc.Invoke(ctx, "/commercionetwork.commercionetwork.did.Query/IdentityHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries a DidDocument by id.
	Identity(context.Context, *QueryResolveIdentityRequest) (*QueryResolveIdentityResponse, error)
	IdentityHistory(context.Context, *QueryResolveIdentityHistoryRequest) (*QueryResolveIdentityHistoryResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Identity(ctx context.Context, req *QueryResolveIdentityRequest) (*QueryResolveIdentityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Identity not implemented")
}
func (*UnimplementedQueryServer) IdentityHistory(ctx context.Context, req *QueryResolveIdentityHistoryRequest) (*QueryResolveIdentityHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IdentityHistory not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Identity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryResolveIdentityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Identity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/commercionetwork.commercionetwork.did.Query/Identity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Identity(ctx, req.(*QueryResolveIdentityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_IdentityHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryResolveIdentityHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).IdentityHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/commercionetwork.commercionetwork.did.Query/IdentityHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).IdentityHistory(ctx, req.(*QueryResolveIdentityHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "commercionetwork.commercionetwork.did.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Identity",
			Handler:    _Query_Identity_Handler,
		},
		{
			MethodName: "IdentityHistory",
			Handler:    _Query_IdentityHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "commercionetwork/did/query.proto",
}

func (m *QueryResolveIdentityRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryResolveIdentityRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryResolveIdentityRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryResolveIdentityResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryResolveIdentityResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryResolveIdentityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Identity != nil {
		{
			size, err := m.Identity.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryResolveIdentityHistoryRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryResolveIdentityHistoryRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryResolveIdentityHistoryRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryResolveIdentityHistoryResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryResolveIdentityHistoryResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryResolveIdentityHistoryResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Identities) > 0 {
		for iNdEx := len(m.Identities) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Identities[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
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
func (m *QueryResolveIdentityRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryResolveIdentityResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Identity != nil {
		l = m.Identity.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryResolveIdentityHistoryRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryResolveIdentityHistoryResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Identities) > 0 {
		for _, e := range m.Identities {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryResolveIdentityRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryResolveIdentityRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryResolveIdentityRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
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
			m.ID = string(dAtA[iNdEx:postIndex])
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
func (m *QueryResolveIdentityResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryResolveIdentityResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryResolveIdentityResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Identity", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Identity == nil {
				m.Identity = &Identity{}
			}
			if err := m.Identity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *QueryResolveIdentityHistoryRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryResolveIdentityHistoryRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryResolveIdentityHistoryRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
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
			m.ID = string(dAtA[iNdEx:postIndex])
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
func (m *QueryResolveIdentityHistoryResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryResolveIdentityHistoryResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryResolveIdentityHistoryResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Identities", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Identities = append(m.Identities, &Identity{})
			if err := m.Identities[len(m.Identities)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
