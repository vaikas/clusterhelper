// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/experimental/authorization_config.proto

package api // import "google.golang.org/genproto/googleapis/api"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Configuration of authorization.
//
// This section determines the authorization provider, if unspecified, then no
// authorization check will be done.
//
// Example:
//
//     experimental:
//       authorization:
//         provider: firebaserules.googleapis.com
type AuthorizationConfig struct {
	// The name of the authorization provider, such as
	// firebaserules.googleapis.com.
	Provider             string   `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthorizationConfig) Reset()         { *m = AuthorizationConfig{} }
func (m *AuthorizationConfig) String() string { return proto.CompactTextString(m) }
func (*AuthorizationConfig) ProtoMessage()    {}
func (*AuthorizationConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_authorization_config_201ba24923dedcb5, []int{0}
}
func (m *AuthorizationConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthorizationConfig.Unmarshal(m, b)
}
func (m *AuthorizationConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthorizationConfig.Marshal(b, m, deterministic)
}
func (dst *AuthorizationConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthorizationConfig.Merge(dst, src)
}
func (m *AuthorizationConfig) XXX_Size() int {
	return xxx_messageInfo_AuthorizationConfig.Size(m)
}
func (m *AuthorizationConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthorizationConfig.DiscardUnknown(m)
}

var xxx_messageInfo_AuthorizationConfig proto.InternalMessageInfo

func (m *AuthorizationConfig) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func init() {
	proto.RegisterType((*AuthorizationConfig)(nil), "google.api.AuthorizationConfig")
}

func init() {
	proto.RegisterFile("google/api/experimental/authorization_config.proto", fileDescriptor_authorization_config_201ba24923dedcb5)
}

var fileDescriptor_authorization_config_201ba24923dedcb5 = []byte{
	// 180 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x4a, 0xcf, 0xcf, 0x4f,
	0xcf, 0x49, 0xd5, 0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xad, 0x28, 0x48, 0x2d, 0xca, 0xcc, 0x4d, 0xcd,
	0x2b, 0x49, 0xcc, 0xd1, 0x4f, 0x2c, 0x2d, 0xc9, 0xc8, 0x2f, 0xca, 0xac, 0x4a, 0x2c, 0xc9, 0xcc,
	0xcf, 0x8b, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2,
	0x82, 0xe8, 0xd1, 0x4b, 0x2c, 0xc8, 0x54, 0x32, 0xe4, 0x12, 0x76, 0x44, 0x56, 0xe9, 0x0c, 0x56,
	0x28, 0x24, 0xc5, 0xc5, 0x51, 0x50, 0x94, 0x5f, 0x96, 0x99, 0x92, 0x5a, 0x24, 0xc1, 0xa8, 0xc0,
	0xa8, 0xc1, 0x19, 0x04, 0xe7, 0x3b, 0x25, 0x71, 0xf1, 0x25, 0xe7, 0xe7, 0xea, 0x21, 0x0c, 0x71,
	0x92, 0xc0, 0x62, 0x44, 0x00, 0xc8, 0xaa, 0x00, 0xc6, 0x28, 0x5d, 0xa8, 0xba, 0xf4, 0xfc, 0x9c,
	0xc4, 0xbc, 0x74, 0xbd, 0xfc, 0xa2, 0x74, 0xfd, 0xf4, 0xd4, 0x3c, 0xb0, 0x43, 0xf4, 0x21, 0x52,
	0x89, 0x05, 0x99, 0xc5, 0x20, 0xf7, 0x5b, 0x27, 0x16, 0x64, 0x2e, 0x62, 0x62, 0x71, 0x77, 0x0c,
	0xf0, 0x4c, 0x62, 0x03, 0x2b, 0x30, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x52, 0x27, 0x0c, 0xba,
	0xdf, 0x00, 0x00, 0x00,
}
