// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs3/app/provider/v1beta1/provider_api.proto

package providerv1beta1

import (
	context "context"
	fmt "fmt"
	v1beta12 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	v1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type OpenInAppRequest struct {
	// OPTIONAL.
	// Opaque information.
	Opaque *v1beta1.Opaque `protobuf:"bytes,1,opt,name=opaque,proto3" json:"opaque,omitempty"`
	// REQUIRED.
	// The resourceInfo to be opened. The gateway grpc message has a ref instead.
	ResourceInfo *v1beta11.ResourceInfo `protobuf:"bytes,2,opt,name=resource_info,json=resourceInfo,proto3" json:"resource_info,omitempty"`
	// REQUIRED.
	// View mode.
	ViewMode ViewMode `protobuf:"varint,3,opt,name=view_mode,json=viewMode,proto3,enum=cs3.app.provider.v1beta1.ViewMode" json:"view_mode,omitempty"`
	// REQUIRED.
	// The access token this application provider will use when contacting
	// the storage provider to read and write.
	// Service implementors MUST make sure that the access token only grants
	// access to the requested resource.
	// Service implementors should use a ResourceId rather than a filepath to grant access, as
	// ResourceIds MUST NOT change when a resource is renamed.
	// The access token MUST be short-lived.
	// TODO(labkode): investigate token derivation techniques.
	AccessToken          string   `protobuf:"bytes,4,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OpenInAppRequest) Reset()         { *m = OpenInAppRequest{} }
func (m *OpenInAppRequest) String() string { return proto.CompactTextString(m) }
func (*OpenInAppRequest) ProtoMessage()    {}
func (*OpenInAppRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_c007b70b037097fe, []int{0}
}

func (m *OpenInAppRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OpenInAppRequest.Unmarshal(m, b)
}
func (m *OpenInAppRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OpenInAppRequest.Marshal(b, m, deterministic)
}
func (m *OpenInAppRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OpenInAppRequest.Merge(m, src)
}
func (m *OpenInAppRequest) XXX_Size() int {
	return xxx_messageInfo_OpenInAppRequest.Size(m)
}
func (m *OpenInAppRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OpenInAppRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OpenInAppRequest proto.InternalMessageInfo

func (m *OpenInAppRequest) GetOpaque() *v1beta1.Opaque {
	if m != nil {
		return m.Opaque
	}
	return nil
}

func (m *OpenInAppRequest) GetResourceInfo() *v1beta11.ResourceInfo {
	if m != nil {
		return m.ResourceInfo
	}
	return nil
}

func (m *OpenInAppRequest) GetViewMode() ViewMode {
	if m != nil {
		return m.ViewMode
	}
	return ViewMode_VIEW_MODE_INVALID
}

func (m *OpenInAppRequest) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

type OpenInAppResponse struct {
	// REQUIRED.
	// The response status.
	Status *v1beta12.Status `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	// OPTIONAL.
	// Opaque information.
	Opaque *v1beta1.Opaque `protobuf:"bytes,2,opt,name=opaque,proto3" json:"opaque,omitempty"`
	// REQUIRED.
	// The url that user agents will render to clients.
	// Usually the rendering happens by using HTML iframes or in separate browser tabs.
	AppUrl               *OpenInAppURL `protobuf:"bytes,3,opt,name=app_url,json=appUrl,proto3" json:"app_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *OpenInAppResponse) Reset()         { *m = OpenInAppResponse{} }
func (m *OpenInAppResponse) String() string { return proto.CompactTextString(m) }
func (*OpenInAppResponse) ProtoMessage()    {}
func (*OpenInAppResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c007b70b037097fe, []int{1}
}

func (m *OpenInAppResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OpenInAppResponse.Unmarshal(m, b)
}
func (m *OpenInAppResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OpenInAppResponse.Marshal(b, m, deterministic)
}
func (m *OpenInAppResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OpenInAppResponse.Merge(m, src)
}
func (m *OpenInAppResponse) XXX_Size() int {
	return xxx_messageInfo_OpenInAppResponse.Size(m)
}
func (m *OpenInAppResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_OpenInAppResponse.DiscardUnknown(m)
}

var xxx_messageInfo_OpenInAppResponse proto.InternalMessageInfo

func (m *OpenInAppResponse) GetStatus() *v1beta12.Status {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *OpenInAppResponse) GetOpaque() *v1beta1.Opaque {
	if m != nil {
		return m.Opaque
	}
	return nil
}

func (m *OpenInAppResponse) GetAppUrl() *OpenInAppURL {
	if m != nil {
		return m.AppUrl
	}
	return nil
}

func init() {
	proto.RegisterType((*OpenInAppRequest)(nil), "cs3.app.provider.v1beta1.OpenInAppRequest")
	proto.RegisterType((*OpenInAppResponse)(nil), "cs3.app.provider.v1beta1.OpenInAppResponse")
}

func init() {
	proto.RegisterFile("cs3/app/provider/v1beta1/provider_api.proto", fileDescriptor_c007b70b037097fe)
}

var fileDescriptor_c007b70b037097fe = []byte{
	// 438 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0x95, 0x0e, 0x15, 0xea, 0x0e, 0x18, 0xbe, 0x10, 0xaa, 0x21, 0x95, 0x1e, 0x50, 0xb5,
	0x21, 0x57, 0x6d, 0x3f, 0xc0, 0x94, 0xee, 0x54, 0x09, 0xd4, 0xc8, 0xb0, 0x1d, 0x50, 0xa5, 0xc8,
	0x73, 0xde, 0xa1, 0x88, 0x2d, 0x7e, 0x67, 0x3b, 0xa9, 0x38, 0xf1, 0x5d, 0x38, 0xf2, 0x09, 0xf8,
	0x0c, 0x7c, 0x24, 0x4e, 0x28, 0x8e, 0x13, 0x55, 0x42, 0x11, 0xbd, 0xc5, 0x7e, 0x7e, 0xcf, 0xfb,
	0xe7, 0x89, 0xc9, 0xb9, 0x34, 0xcb, 0x99, 0x40, 0x9c, 0xa1, 0x56, 0x65, 0x96, 0x82, 0x9e, 0x95,
	0xf3, 0x1b, 0xb0, 0x62, 0xde, 0x5e, 0x24, 0x02, 0x33, 0x86, 0x5a, 0x59, 0x45, 0x43, 0x69, 0x96,
	0x4c, 0x20, 0xb2, 0x46, 0x63, 0x1e, 0x1e, 0x4d, 0x3b, 0xcb, 0x68, 0x30, 0xaa, 0xd0, 0x12, 0x4c,
	0x5d, 0x63, 0x74, 0x5a, 0x91, 0x1a, 0x65, 0x0b, 0x18, 0x2b, 0x6c, 0xd1, 0xa8, 0xef, 0x2a, 0xd5,
	0x58, 0xa5, 0xc5, 0x17, 0xf8, 0x7f, 0xad, 0xd7, 0x15, 0x6d, 0xbf, 0x21, 0x98, 0x16, 0x71, 0xa7,
	0x5a, 0x9e, 0xfc, 0x09, 0xc8, 0xc9, 0x06, 0x21, 0x5f, 0xe7, 0x11, 0x22, 0x87, 0x87, 0x02, 0x8c,
	0xa5, 0x73, 0xd2, 0x57, 0x28, 0x1e, 0x0a, 0x08, 0x83, 0x71, 0x30, 0x1d, 0x2e, 0x5e, 0xb1, 0x6a,
	0xa9, 0xda, 0xe6, 0x8b, 0xb0, 0x8d, 0x03, 0xb8, 0x07, 0xe9, 0x86, 0x3c, 0x6d, 0x3a, 0x27, 0x59,
	0x7e, 0xab, 0xc2, 0x9e, 0x73, 0x9e, 0x39, 0xa7, 0x1f, 0xf6, 0x9f, 0x48, 0x18, 0xf7, 0x96, 0x75,
	0x7e, 0xab, 0xf8, 0xb1, 0xde, 0x3b, 0xd1, 0x0b, 0x32, 0x28, 0x33, 0xd8, 0x25, 0xf7, 0x2a, 0x85,
	0xf0, 0x68, 0x1c, 0x4c, 0x9f, 0x2d, 0x26, 0xac, 0x2b, 0x5b, 0x76, 0x9d, 0xc1, 0xee, 0x83, 0x4a,
	0x81, 0x3f, 0x29, 0xfd, 0x17, 0x7d, 0x43, 0x8e, 0x85, 0x94, 0x60, 0x4c, 0x62, 0xd5, 0x57, 0xc8,
	0xc3, 0x47, 0xe3, 0x60, 0x3a, 0xe0, 0xc3, 0xfa, 0xee, 0x53, 0x75, 0x35, 0xf9, 0x15, 0x90, 0x17,
	0x7b, 0xcb, 0x1b, 0x54, 0xb9, 0x01, 0x3a, 0x23, 0xfd, 0x3a, 0x6f, 0xbf, 0xfd, 0x4b, 0xd7, 0x56,
	0xa3, 0x6c, 0xbb, 0x7d, 0x74, 0x32, 0xf7, 0xd8, 0x5e, 0x5c, 0xbd, 0x43, 0xe3, 0xba, 0x20, 0x8f,
	0x05, 0x62, 0x52, 0xe8, 0x3b, 0xb7, 0xdb, 0x70, 0xf1, 0xb6, 0x7b, 0xb7, 0x76, 0xc2, 0x2b, 0xfe,
	0x9e, 0xf7, 0x05, 0xe2, 0x95, 0xbe, 0x5b, 0x18, 0x32, 0x8c, 0x3d, 0x18, 0xc5, 0x6b, 0x9a, 0x92,
	0x41, 0x8b, 0xd1, 0xb3, 0x03, 0x6a, 0xf9, 0x5f, 0x3d, 0x3a, 0x3f, 0x88, 0xad, 0x93, 0x59, 0x7d,
	0x27, 0xa7, 0x52, 0xdd, 0x77, 0x3a, 0x56, 0x27, 0xed, 0x48, 0x98, 0xc5, 0xd5, 0xf3, 0x8a, 0x83,
	0xcf, 0xcf, 0x1b, 0xca, 0x43, 0x3f, 0x7a, 0x47, 0x97, 0x51, 0xfc, 0xb3, 0x17, 0x5e, 0x9a, 0x25,
	0x8b, 0x10, 0x59, 0xe3, 0x61, 0xd7, 0xf3, 0x55, 0x05, 0xfc, 0x76, 0xd2, 0x36, 0x42, 0xdc, 0x36,
	0xd2, 0xd6, 0x4b, 0x37, 0x7d, 0xf7, 0x68, 0x97, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xb9, 0xb9,
	0x16, 0xb5, 0x92, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProviderAPIClient is the client API for ProviderAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProviderAPIClient interface {
	// Returns the App URL and all necessary info to open a resource in an online editor.
	// MUST return CODE_NOT_FOUND if the resource does not exist.
	OpenInApp(ctx context.Context, in *OpenInAppRequest, opts ...grpc.CallOption) (*OpenInAppResponse, error)
}

type providerAPIClient struct {
	cc *grpc.ClientConn
}

func NewProviderAPIClient(cc *grpc.ClientConn) ProviderAPIClient {
	return &providerAPIClient{cc}
}

func (c *providerAPIClient) OpenInApp(ctx context.Context, in *OpenInAppRequest, opts ...grpc.CallOption) (*OpenInAppResponse, error) {
	out := new(OpenInAppResponse)
	err := c.cc.Invoke(ctx, "/cs3.app.provider.v1beta1.ProviderAPI/OpenInApp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProviderAPIServer is the server API for ProviderAPI service.
type ProviderAPIServer interface {
	// Returns the App URL and all necessary info to open a resource in an online editor.
	// MUST return CODE_NOT_FOUND if the resource does not exist.
	OpenInApp(context.Context, *OpenInAppRequest) (*OpenInAppResponse, error)
}

// UnimplementedProviderAPIServer can be embedded to have forward compatible implementations.
type UnimplementedProviderAPIServer struct {
}

func (*UnimplementedProviderAPIServer) OpenInApp(ctx context.Context, req *OpenInAppRequest) (*OpenInAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenInApp not implemented")
}

func RegisterProviderAPIServer(s *grpc.Server, srv ProviderAPIServer) {
	s.RegisterService(&_ProviderAPI_serviceDesc, srv)
}

func _ProviderAPI_OpenInApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpenInAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProviderAPIServer).OpenInApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cs3.app.provider.v1beta1.ProviderAPI/OpenInApp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProviderAPIServer).OpenInApp(ctx, req.(*OpenInAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProviderAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cs3.app.provider.v1beta1.ProviderAPI",
	HandlerType: (*ProviderAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OpenInApp",
			Handler:    _ProviderAPI_OpenInApp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cs3/app/provider/v1beta1/provider_api.proto",
}
