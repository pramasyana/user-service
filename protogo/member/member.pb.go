// Code generated by protoc-gen-go. DO NOT EDIT.
// source: member.proto

/*
Package member is a generated protocol buffer package.

It is generated from these files:
	member.proto

It has these top-level messages:
	ResponseMessage
	MemberQuery
	MemberRegister
	MemberUpdate
	MemberPasswordRequest
	Address
	SocialMedia
	Member
*/
package member

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ResponseMessage struct {
	Message   string `protobuf:"bytes,1,opt,name=Message" json:"Message,omitempty"`
	Email     string `protobuf:"bytes,2,opt,name=Email" json:"Email,omitempty"`
	FirstName string `protobuf:"bytes,3,opt,name=FirstName" json:"FirstName,omitempty"`
	LastName  string `protobuf:"bytes,4,opt,name=LastName" json:"LastName,omitempty"`
}

func (m *ResponseMessage) Reset()                    { *m = ResponseMessage{} }
func (m *ResponseMessage) String() string            { return proto.CompactTextString(m) }
func (*ResponseMessage) ProtoMessage()               {}
func (*ResponseMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ResponseMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ResponseMessage) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *ResponseMessage) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *ResponseMessage) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

type MemberQuery struct {
	ID  string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	JWT string `protobuf:"bytes,2,opt,name=JWT" json:"JWT,omitempty"`
}

func (m *MemberQuery) Reset()                    { *m = MemberQuery{} }
func (m *MemberQuery) String() string            { return proto.CompactTextString(m) }
func (*MemberQuery) ProtoMessage()               {}
func (*MemberQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MemberQuery) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *MemberQuery) GetJWT() string {
	if m != nil {
		return m.JWT
	}
	return ""
}

type MemberRegister struct {
	FirstName  string `protobuf:"bytes,1,opt,name=FirstName" json:"FirstName,omitempty"`
	LastName   string `protobuf:"bytes,2,opt,name=LastName" json:"LastName,omitempty"`
	Email      string `protobuf:"bytes,3,opt,name=Email" json:"Email,omitempty"`
	Gender     string `protobuf:"bytes,4,opt,name=Gender" json:"Gender,omitempty"`
	Mobile     string `protobuf:"bytes,5,opt,name=Mobile" json:"Mobile,omitempty"`
	Birthdate  string `protobuf:"bytes,6,opt,name=Birthdate" json:"Birthdate,omitempty"`
	Password   string `protobuf:"bytes,7,opt,name=Password" json:"Password,omitempty"`
	RePassword string `protobuf:"bytes,8,opt,name=RePassword" json:"RePassword,omitempty"`
}

func (m *MemberRegister) Reset()                    { *m = MemberRegister{} }
func (m *MemberRegister) String() string            { return proto.CompactTextString(m) }
func (*MemberRegister) ProtoMessage()               {}
func (*MemberRegister) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MemberRegister) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *MemberRegister) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *MemberRegister) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *MemberRegister) GetGender() string {
	if m != nil {
		return m.Gender
	}
	return ""
}

func (m *MemberRegister) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *MemberRegister) GetBirthdate() string {
	if m != nil {
		return m.Birthdate
	}
	return ""
}

func (m *MemberRegister) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *MemberRegister) GetRePassword() string {
	if m != nil {
		return m.RePassword
	}
	return ""
}

type MemberUpdate struct {
	ID            string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	FirstName     string `protobuf:"bytes,2,opt,name=FirstName" json:"FirstName,omitempty"`
	LastName      string `protobuf:"bytes,3,opt,name=LastName" json:"LastName,omitempty"`
	Gender        string `protobuf:"bytes,4,opt,name=Gender" json:"Gender,omitempty"`
	Birthdate     string `protobuf:"bytes,5,opt,name=Birthdate" json:"Birthdate,omitempty"`
	Mobile        string `protobuf:"bytes,6,opt,name=Mobile" json:"Mobile,omitempty"`
	Phone         string `protobuf:"bytes,7,opt,name=Phone" json:"Phone,omitempty"`
	Ext           string `protobuf:"bytes,8,opt,name=Ext" json:"Ext,omitempty"`
	Province      string `protobuf:"bytes,9,opt,name=Province" json:"Province,omitempty"`
	ProvinceID    string `protobuf:"bytes,10,opt,name=ProvinceID" json:"ProvinceID,omitempty"`
	City          string `protobuf:"bytes,11,opt,name=City" json:"City,omitempty"`
	CityID        string `protobuf:"bytes,12,opt,name=CityID" json:"CityID,omitempty"`
	District      string `protobuf:"bytes,13,opt,name=District" json:"District,omitempty"`
	DistrictID    string `protobuf:"bytes,14,opt,name=DistrictID" json:"DistrictID,omitempty"`
	SubDistrict   string `protobuf:"bytes,15,opt,name=SubDistrict" json:"SubDistrict,omitempty"`
	SubDistrictID string `protobuf:"bytes,16,opt,name=SubDistrictID" json:"SubDistrictID,omitempty"`
	ZipCode       string `protobuf:"bytes,17,opt,name=ZipCode" json:"ZipCode,omitempty"`
	Street1       string `protobuf:"bytes,18,opt,name=Street1" json:"Street1,omitempty"`
	Street2       string `protobuf:"bytes,19,opt,name=Street2" json:"Street2,omitempty"`
	Status        string `protobuf:"bytes,20,opt,name=Status" json:"Status,omitempty"`
}

func (m *MemberUpdate) Reset()                    { *m = MemberUpdate{} }
func (m *MemberUpdate) String() string            { return proto.CompactTextString(m) }
func (*MemberUpdate) ProtoMessage()               {}
func (*MemberUpdate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MemberUpdate) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *MemberUpdate) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *MemberUpdate) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *MemberUpdate) GetGender() string {
	if m != nil {
		return m.Gender
	}
	return ""
}

func (m *MemberUpdate) GetBirthdate() string {
	if m != nil {
		return m.Birthdate
	}
	return ""
}

func (m *MemberUpdate) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *MemberUpdate) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *MemberUpdate) GetExt() string {
	if m != nil {
		return m.Ext
	}
	return ""
}

func (m *MemberUpdate) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *MemberUpdate) GetProvinceID() string {
	if m != nil {
		return m.ProvinceID
	}
	return ""
}

func (m *MemberUpdate) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *MemberUpdate) GetCityID() string {
	if m != nil {
		return m.CityID
	}
	return ""
}

func (m *MemberUpdate) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *MemberUpdate) GetDistrictID() string {
	if m != nil {
		return m.DistrictID
	}
	return ""
}

func (m *MemberUpdate) GetSubDistrict() string {
	if m != nil {
		return m.SubDistrict
	}
	return ""
}

func (m *MemberUpdate) GetSubDistrictID() string {
	if m != nil {
		return m.SubDistrictID
	}
	return ""
}

func (m *MemberUpdate) GetZipCode() string {
	if m != nil {
		return m.ZipCode
	}
	return ""
}

func (m *MemberUpdate) GetStreet1() string {
	if m != nil {
		return m.Street1
	}
	return ""
}

func (m *MemberUpdate) GetStreet2() string {
	if m != nil {
		return m.Street2
	}
	return ""
}

func (m *MemberUpdate) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type MemberPasswordRequest struct {
	JWT 		string `protobuf:"bytes,2,opt,name=JWT" json:"JWT,omitempty"`
	ID          string `protobuf:"bytes,2,opt,name=ID" json:"ID,omitempty"`
	OldPassword string `protobuf:"bytes,3,opt,name=OldPassword" json:"OldPassword,omitempty"`
	NewPassword string `protobuf:"bytes,4,opt,name=NewPassword" json:"NewPassword,omitempty"`
	RePassword  string `protobuf:"bytes,5,opt,name=RePassword" json:"RePassword,omitempty"`
}

func (m *MemberPasswordRequest) Reset()                    { *m = MemberPasswordRequest{} }
func (m *MemberPasswordRequest) String() string            { return proto.CompactTextString(m) }
func (*MemberPasswordRequest) ProtoMessage()               {}
func (*MemberPasswordRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *MemberPasswordRequest) GetJWT() string {
	if m != nil {
		return m.JWT
	}
	return ""
}

func (m *MemberPasswordRequest) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *MemberPasswordRequest) GetOldPassword() string {
	if m != nil {
		return m.OldPassword
	}
	return ""
}

func (m *MemberPasswordRequest) GetNewPassword() string {
	if m != nil {
		return m.NewPassword
	}
	return ""
}

func (m *MemberPasswordRequest) GetRePassword() string {
	if m != nil {
		return m.RePassword
	}
	return ""
}

type Address struct {
	Province      string `protobuf:"bytes,1,opt,name=Province" json:"Province,omitempty"`
	ProvinceID    string `protobuf:"bytes,2,opt,name=ProvinceID" json:"ProvinceID,omitempty"`
	City          string `protobuf:"bytes,3,opt,name=City" json:"City,omitempty"`
	CityID        string `protobuf:"bytes,4,opt,name=CityID" json:"CityID,omitempty"`
	District      string `protobuf:"bytes,5,opt,name=District" json:"District,omitempty"`
	DistrictID    string `protobuf:"bytes,6,opt,name=DistrictID" json:"DistrictID,omitempty"`
	SubDistrict   string `protobuf:"bytes,7,opt,name=SubDistrict" json:"SubDistrict,omitempty"`
	SubDistrictID string `protobuf:"bytes,8,opt,name=SubDistrictID" json:"SubDistrictID,omitempty"`
	ZipCode       string `protobuf:"bytes,9,opt,name=ZipCode" json:"ZipCode,omitempty"`
	Street1       string `protobuf:"bytes,10,opt,name=Street1" json:"Street1,omitempty"`
	Street2       string `protobuf:"bytes,11,opt,name=Street2" json:"Street2,omitempty"`
}

func (m *Address) Reset()                    { *m = Address{} }
func (m *Address) String() string            { return proto.CompactTextString(m) }
func (*Address) ProtoMessage()               {}
func (*Address) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Address) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *Address) GetProvinceID() string {
	if m != nil {
		return m.ProvinceID
	}
	return ""
}

func (m *Address) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *Address) GetCityID() string {
	if m != nil {
		return m.CityID
	}
	return ""
}

func (m *Address) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *Address) GetDistrictID() string {
	if m != nil {
		return m.DistrictID
	}
	return ""
}

func (m *Address) GetSubDistrict() string {
	if m != nil {
		return m.SubDistrict
	}
	return ""
}

func (m *Address) GetSubDistrictID() string {
	if m != nil {
		return m.SubDistrictID
	}
	return ""
}

func (m *Address) GetZipCode() string {
	if m != nil {
		return m.ZipCode
	}
	return ""
}

func (m *Address) GetStreet1() string {
	if m != nil {
		return m.Street1
	}
	return ""
}

func (m *Address) GetStreet2() string {
	if m != nil {
		return m.Street2
	}
	return ""
}

type SocialMedia struct {
	FacebookID string `protobuf:"bytes,1,opt,name=FacebookID" json:"FacebookID,omitempty"`
	GoogleID   string `protobuf:"bytes,2,opt,name=GoogleID" json:"GoogleID,omitempty"`
	AzureID    string `protobuf:"bytes,3,opt,name=AzureID" json:"AzureID,omitempty"`
}

func (m *SocialMedia) Reset()                    { *m = SocialMedia{} }
func (m *SocialMedia) String() string            { return proto.CompactTextString(m) }
func (*SocialMedia) ProtoMessage()               {}
func (*SocialMedia) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *SocialMedia) GetFacebookID() string {
	if m != nil {
		return m.FacebookID
	}
	return ""
}

func (m *SocialMedia) GetGoogleID() string {
	if m != nil {
		return m.GoogleID
	}
	return ""
}

func (m *SocialMedia) GetAzureID() string {
	if m != nil {
		return m.AzureID
	}
	return ""
}

type Member struct {
	ID           string       `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	FirstName    string       `protobuf:"bytes,2,opt,name=FirstName" json:"FirstName,omitempty"`
	LastName     string       `protobuf:"bytes,3,opt,name=LastName" json:"LastName,omitempty"`
	Email        string       `protobuf:"bytes,4,opt,name=Email" json:"Email,omitempty"`
	Gender       string       `protobuf:"bytes,5,opt,name=Gender" json:"Gender,omitempty"`
	Mobile       string       `protobuf:"bytes,6,opt,name=Mobile" json:"Mobile,omitempty"`
	Phone        string       `protobuf:"bytes,7,opt,name=Phone" json:"Phone,omitempty"`
	Ext          string       `protobuf:"bytes,8,opt,name=Ext" json:"Ext,omitempty"`
	Birthdate    string       `protobuf:"bytes,9,opt,name=Birthdate" json:"Birthdate,omitempty"`
	Address      *Address     `protobuf:"bytes,10,opt,name=Address" json:"Address,omitempty"`
	JobTitle     string       `protobuf:"bytes,11,opt,name=JobTitle" json:"JobTitle,omitempty"`
	Department   string       `protobuf:"bytes,12,opt,name=Department" json:"Department,omitempty"`
	Status       string       `protobuf:"bytes,13,opt,name=Status" json:"Status,omitempty"`
	SocialMedia  *SocialMedia `protobuf:"bytes,14,opt,name=SocialMedia" json:"SocialMedia,omitempty"`
	IsAdmin      bool         `protobuf:"varint,15,opt,name=IsAdmin" json:"IsAdmin,omitempty"`
	IsStaff      bool         `protobuf:"varint,16,opt,name=IsStaff" json:"IsStaff,omitempty"`
	SignUpFrom   string       `protobuf:"bytes,17,opt,name=SignUpFrom" json:"SignUpFrom,omitempty"`
	HasPassword  bool         `protobuf:"varint,18,opt,name=HasPassword" json:"HasPassword,omitempty"`
	LastLogin    string       `protobuf:"bytes,19,opt,name=LastLogin" json:"LastLogin,omitempty"`
	Version      int32        `protobuf:"varint,20,opt,name=Version" json:"Version,omitempty"`
	Created      string       `protobuf:"bytes,21,opt,name=Created" json:"Created,omitempty"`
	LastModified string       `protobuf:"bytes,22,opt,name=LastModified" json:"LastModified,omitempty"`
	LastBlocked  string       `protobuf:"bytes,23,opt,name=LastBlocked" json:"LastBlocked,omitempty"`
}

func (m *Member) Reset()                    { *m = Member{} }
func (m *Member) String() string            { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()               {}
func (*Member) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Member) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Member) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *Member) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *Member) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Member) GetGender() string {
	if m != nil {
		return m.Gender
	}
	return ""
}

func (m *Member) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

func (m *Member) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Member) GetExt() string {
	if m != nil {
		return m.Ext
	}
	return ""
}

func (m *Member) GetBirthdate() string {
	if m != nil {
		return m.Birthdate
	}
	return ""
}

func (m *Member) GetAddress() *Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Member) GetJobTitle() string {
	if m != nil {
		return m.JobTitle
	}
	return ""
}

func (m *Member) GetDepartment() string {
	if m != nil {
		return m.Department
	}
	return ""
}

func (m *Member) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Member) GetSocialMedia() *SocialMedia {
	if m != nil {
		return m.SocialMedia
	}
	return nil
}

func (m *Member) GetIsAdmin() bool {
	if m != nil {
		return m.IsAdmin
	}
	return false
}

func (m *Member) GetIsStaff() bool {
	if m != nil {
		return m.IsStaff
	}
	return false
}

func (m *Member) GetSignUpFrom() string {
	if m != nil {
		return m.SignUpFrom
	}
	return ""
}

func (m *Member) GetHasPassword() bool {
	if m != nil {
		return m.HasPassword
	}
	return false
}

func (m *Member) GetLastLogin() string {
	if m != nil {
		return m.LastLogin
	}
	return ""
}

func (m *Member) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *Member) GetCreated() string {
	if m != nil {
		return m.Created
	}
	return ""
}

func (m *Member) GetLastModified() string {
	if m != nil {
		return m.LastModified
	}
	return ""
}

func (m *Member) GetLastBlocked() string {
	if m != nil {
		return m.LastBlocked
	}
	return ""
}

func init() {
	proto.RegisterType((*ResponseMessage)(nil), "member.ResponseMessage")
	proto.RegisterType((*MemberQuery)(nil), "member.MemberQuery")
	proto.RegisterType((*MemberRegister)(nil), "member.MemberRegister")
	proto.RegisterType((*MemberUpdate)(nil), "member.MemberUpdate")
	proto.RegisterType((*MemberPasswordRequest)(nil), "member.MemberPasswordRequest")
	proto.RegisterType((*Address)(nil), "member.Address")
	proto.RegisterType((*SocialMedia)(nil), "member.SocialMedia")
	proto.RegisterType((*Member)(nil), "member.Member")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for MemberService service

type MemberServiceClient interface {
	Register(ctx context.Context, in *MemberRegister, opts ...grpc.CallOption) (*ResponseMessage, error)
	Update(ctx context.Context, in *MemberUpdate, opts ...grpc.CallOption) (*ResponseMessage, error)
	FindByID(ctx context.Context, in *MemberQuery, opts ...grpc.CallOption) (*Member, error)
	GetMe(ctx context.Context, in *MemberQuery, opts ...grpc.CallOption) (*Member, error)
	UpdatePassword(ctx context.Context, in *MemberPasswordRequest, opts ...grpc.CallOption) (*ResponseMessage, error)
}

type memberServiceClient struct {
	cc *grpc.ClientConn
}

func NewMemberServiceClient(cc *grpc.ClientConn) MemberServiceClient {
	return &memberServiceClient{cc}
}

func (c *memberServiceClient) Register(ctx context.Context, in *MemberRegister, opts ...grpc.CallOption) (*ResponseMessage, error) {
	out := new(ResponseMessage)
	err := grpc.Invoke(ctx, "/member.MemberService/Register", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberServiceClient) Update(ctx context.Context, in *MemberUpdate, opts ...grpc.CallOption) (*ResponseMessage, error) {
	out := new(ResponseMessage)
	err := grpc.Invoke(ctx, "/member.MemberService/Update", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberServiceClient) FindByID(ctx context.Context, in *MemberQuery, opts ...grpc.CallOption) (*Member, error) {
	out := new(Member)
	err := grpc.Invoke(ctx, "/member.MemberService/FindByID", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberServiceClient) GetMe(ctx context.Context, in *MemberQuery, opts ...grpc.CallOption) (*Member, error) {
	out := new(Member)
	err := grpc.Invoke(ctx, "/member.MemberService/GetMe", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *memberServiceClient) UpdatePassword(ctx context.Context, in *MemberPasswordRequest, opts ...grpc.CallOption) (*ResponseMessage, error) {
	out := new(ResponseMessage)
	err := grpc.Invoke(ctx, "/member.MemberService/UpdatePassword", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for MemberService service

type MemberServiceServer interface {
	Register(context.Context, *MemberRegister) (*ResponseMessage, error)
	Update(context.Context, *MemberUpdate) (*ResponseMessage, error)
	FindByID(context.Context, *MemberQuery) (*Member, error)
	GetMe(context.Context, *MemberQuery) (*Member, error)
	UpdatePassword(context.Context, *MemberPasswordRequest) (*ResponseMessage, error)
}

func RegisterMemberServiceServer(s *grpc.Server, srv MemberServiceServer) {
	s.RegisterService(&_MemberService_serviceDesc, srv)
}

func _MemberService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberRegister)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.MemberService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServiceServer).Register(ctx, req.(*MemberRegister))
	}
	return interceptor(ctx, in, info, handler)
}

func _MemberService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.MemberService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServiceServer).Update(ctx, req.(*MemberUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _MemberService_FindByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServiceServer).FindByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.MemberService/FindByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServiceServer).FindByID(ctx, req.(*MemberQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _MemberService_GetMe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServiceServer).GetMe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.MemberService/GetMe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServiceServer).GetMe(ctx, req.(*MemberQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _MemberService_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MemberPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MemberServiceServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/member.MemberService/UpdatePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MemberServiceServer).UpdatePassword(ctx, req.(*MemberPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MemberService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "member.MemberService",
	HandlerType: (*MemberServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _MemberService_Register_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _MemberService_Update_Handler,
		},
		{
			MethodName: "FindByID",
			Handler:    _MemberService_FindByID_Handler,
		},
		{
			MethodName: "GetMe",
			Handler:    _MemberService_GetMe_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _MemberService_UpdatePassword_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "member.proto",
}

func init() { proto.RegisterFile("member.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 888 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xdd, 0x6e, 0x1b, 0x45,
	0x14, 0x96, 0x7f, 0xd6, 0x3f, 0xc7, 0x89, 0x53, 0x26, 0x69, 0x3a, 0x8a, 0x00, 0x45, 0x2b, 0x2e,
	0x40, 0x42, 0xad, 0x30, 0x42, 0x5c, 0x70, 0x95, 0xc4, 0x4d, 0xea, 0xaa, 0x2e, 0x65, 0xdd, 0x82,
	0xc4, 0xdd, 0xda, 0x7b, 0xe2, 0x8e, 0x6a, 0xef, 0x98, 0x99, 0x71, 0x4b, 0xe1, 0x11, 0x78, 0x09,
	0xc4, 0x03, 0xf0, 0x12, 0xbc, 0x0d, 0x4f, 0x81, 0xe6, 0x67, 0x77, 0x67, 0xdc, 0x6e, 0x5a, 0xa1,
	0x5e, 0x79, 0xbe, 0xef, 0xcc, 0xf8, 0xfc, 0x7c, 0x67, 0xce, 0x0e, 0xec, 0xad, 0x71, 0x3d, 0x47,
	0x71, 0x77, 0x23, 0xb8, 0xe2, 0xa4, 0x63, 0x51, 0xfc, 0x3b, 0x1c, 0x24, 0x28, 0x37, 0x3c, 0x97,
	0x38, 0x45, 0x29, 0xd3, 0x25, 0x12, 0x0a, 0x5d, 0xb7, 0xa4, 0x8d, 0xd3, 0xc6, 0xe7, 0xfd, 0xa4,
	0x80, 0xe4, 0x08, 0xa2, 0xfb, 0xeb, 0x94, 0xad, 0x68, 0xd3, 0xf0, 0x16, 0x90, 0x8f, 0xa1, 0x7f,
	0xc9, 0x84, 0x54, 0x8f, 0xd3, 0x35, 0xd2, 0x96, 0xb1, 0x54, 0x04, 0x39, 0x81, 0xde, 0xa3, 0xd4,
	0x19, 0xdb, 0xc6, 0x58, 0xe2, 0xf8, 0x1e, 0x0c, 0xa6, 0x26, 0x8c, 0x1f, 0xb6, 0x28, 0x5e, 0x93,
	0x21, 0x34, 0x27, 0x63, 0xe7, 0xb3, 0x39, 0x19, 0x93, 0x5b, 0xd0, 0x7a, 0xf8, 0xd3, 0x53, 0xe7,
	0x4c, 0x2f, 0xe3, 0x7f, 0x1b, 0x30, 0xb4, 0x27, 0x12, 0x5c, 0x32, 0xa9, 0x50, 0x84, 0xde, 0x1b,
	0x37, 0x79, 0x6f, 0x86, 0xde, 0xab, 0x6c, 0x5a, 0x7e, 0x36, 0xc7, 0xd0, 0xb9, 0xc2, 0x3c, 0x43,
	0xe1, 0xa2, 0x75, 0x48, 0xf3, 0x53, 0x3e, 0x67, 0x2b, 0xa4, 0x91, 0xe5, 0x2d, 0xd2, 0xfe, 0xcf,
	0x99, 0x50, 0xcf, 0xb3, 0x54, 0x21, 0xed, 0x58, 0xff, 0x25, 0xa1, 0xfd, 0x3f, 0x49, 0xa5, 0x7c,
	0xc5, 0x45, 0x46, 0xbb, 0xd6, 0x7f, 0x81, 0xc9, 0xa7, 0x00, 0x09, 0x96, 0xd6, 0x9e, 0xb1, 0x7a,
	0x4c, 0xfc, 0x67, 0x1b, 0xf6, 0x6c, 0xb2, 0xcf, 0x36, 0xe6, 0xcf, 0x76, 0xeb, 0x13, 0xa4, 0xde,
	0xbc, 0x29, 0xf5, 0xd6, 0x4e, 0xea, 0x75, 0x49, 0x06, 0xc9, 0x44, 0xbb, 0xc9, 0x54, 0x25, 0xe8,
	0x04, 0x25, 0x38, 0x82, 0xe8, 0xc9, 0x73, 0x9e, 0xa3, 0xcb, 0xd0, 0x02, 0xad, 0xde, 0xfd, 0x5f,
	0x95, 0xcb, 0x4b, 0x2f, 0x4d, 0x31, 0x04, 0x7f, 0xc9, 0xf2, 0x05, 0xd2, 0xbe, 0x2b, 0x86, 0xc3,
	0xba, 0x18, 0xc5, 0x7a, 0x32, 0xa6, 0x60, 0x8b, 0x51, 0x31, 0x84, 0x40, 0xfb, 0x82, 0xa9, 0xd7,
	0x74, 0x60, 0x2c, 0x66, 0xad, 0xe3, 0xd1, 0xbf, 0x93, 0x31, 0xdd, 0xb3, 0xf1, 0x58, 0xa4, 0xfd,
	0x8c, 0x99, 0x54, 0x82, 0x2d, 0x14, 0xdd, 0xb7, 0x7e, 0x0a, 0xac, 0xfd, 0x14, 0xeb, 0xc9, 0x98,
	0x0e, 0xad, 0x9f, 0x8a, 0x21, 0xa7, 0x30, 0x98, 0x6d, 0xe7, 0xe5, 0xf1, 0x03, 0xb3, 0xc1, 0xa7,
	0xc8, 0x67, 0xb0, 0xef, 0xc1, 0xc9, 0x98, 0xde, 0x32, 0x7b, 0x42, 0x52, 0x5f, 0xa2, 0x9f, 0xd9,
	0xe6, 0x82, 0x67, 0x48, 0x3f, 0xb2, 0x97, 0xc8, 0x41, 0x6d, 0x99, 0x29, 0x81, 0xa8, 0xbe, 0xa2,
	0xc4, 0x5a, 0x1c, 0xac, 0x2c, 0x23, 0x7a, 0xe8, 0x5b, 0x46, 0x3a, 0xd3, 0x99, 0x4a, 0xd5, 0x56,
	0xd2, 0x23, 0x9b, 0xa9, 0x45, 0xf1, 0x1f, 0x0d, 0xb8, 0x6d, 0x5b, 0xa4, 0xe8, 0x9a, 0x04, 0x7f,
	0xd9, 0xa2, 0x54, 0x6f, 0xf4, 0xca, 0x29, 0x0c, 0xbe, 0x5f, 0x65, 0x65, 0xb7, 0xd9, 0x6e, 0xf1,
	0x29, 0xbd, 0xe3, 0x31, 0xbe, 0x2a, 0x77, 0xd8, 0x96, 0xf1, 0xa9, 0x9d, 0x86, 0x6d, 0xbf, 0xd1,
	0xb0, 0xff, 0x34, 0xa1, 0x7b, 0x96, 0x65, 0x02, 0xa5, 0x0c, 0xb4, 0x6e, 0xdc, 0xa8, 0x75, 0xb3,
	0x56, 0xeb, 0xd6, 0x5b, 0xb5, 0x6e, 0xd7, 0x6a, 0x1d, 0xdd, 0xa8, 0x75, 0xe7, 0x5d, 0x5a, 0x77,
	0xdf, 0x43, 0xeb, 0xde, 0x3b, 0xb4, 0xee, 0xd7, 0x6a, 0x0d, 0xb5, 0x5a, 0x0f, 0x02, 0xad, 0xe3,
	0x05, 0x0c, 0x66, 0x7c, 0xc1, 0xd2, 0xd5, 0x14, 0x33, 0x96, 0xea, 0x24, 0x2e, 0xd3, 0x05, 0xce,
	0x39, 0x7f, 0x51, 0x0a, 0xea, 0x31, 0xba, 0x00, 0x57, 0x9c, 0x2f, 0x57, 0x55, 0x29, 0x4b, 0xac,
	0x9d, 0x9c, 0xfd, 0xb6, 0x15, 0xda, 0x64, 0x6b, 0x59, 0xc0, 0xf8, 0xaf, 0x08, 0x3a, 0xb6, 0x71,
	0x3e, 0xe0, 0x54, 0x29, 0x07, 0x6a, 0xfb, 0xed, 0x03, 0x35, 0xaa, 0x19, 0xa8, 0xff, 0x6f, 0x9a,
	0x04, 0xb3, 0xaa, 0xbf, 0x3b, 0xab, 0xbe, 0x28, 0x5b, 0xd1, 0x54, 0x7e, 0x30, 0x3a, 0xb8, 0xeb,
	0xbe, 0x7f, 0x8e, 0x4e, 0xfc, 0x56, 0x7d, 0xc8, 0xe7, 0x4f, 0x99, 0x5a, 0xa1, 0xd3, 0xa2, 0xc4,
	0xa6, 0x85, 0x70, 0x93, 0x0a, 0xb5, 0xc6, 0x5c, 0xb9, 0x31, 0xe3, 0x31, 0xde, 0xc5, 0xdc, 0xf7,
	0x2f, 0x26, 0xf9, 0x26, 0x10, 0xd1, 0xcc, 0x99, 0xc1, 0xe8, 0xb0, 0x08, 0xc1, 0x33, 0x25, 0x81,
	0xd8, 0x14, 0xba, 0x13, 0x79, 0x96, 0xad, 0x59, 0x6e, 0x26, 0x4f, 0x2f, 0x29, 0xa0, 0xb5, 0xcc,
	0x54, 0x7a, 0x7d, 0x6d, 0xe6, 0x8d, 0xb1, 0x18, 0xa8, 0x43, 0x9c, 0xb1, 0x65, 0xfe, 0x6c, 0x73,
	0x29, 0xf8, 0xda, 0x0d, 0x1b, 0x8f, 0xd1, 0x5d, 0xfe, 0x20, 0x95, 0xe5, 0xb5, 0x25, 0xe6, 0xb4,
	0x4f, 0xe9, 0x4a, 0x6a, 0x0d, 0x1f, 0xf1, 0x25, 0xcb, 0xdd, 0xe4, 0xa9, 0x08, 0xed, 0xf9, 0x47,
	0x14, 0x92, 0xf1, 0xdc, 0x0c, 0x9f, 0x28, 0x29, 0xa0, 0xb6, 0x5c, 0x08, 0x4c, 0x15, 0x66, 0xf4,
	0xb6, 0x6d, 0x2f, 0x07, 0x49, 0x0c, 0x7b, 0xfa, 0x0f, 0xa6, 0x3c, 0x63, 0xd7, 0x0c, 0x33, 0x7a,
	0x6c, 0xcc, 0x01, 0xa7, 0xe3, 0xd2, 0xf8, 0x7c, 0xc5, 0x17, 0x2f, 0x30, 0xa3, 0x77, 0xec, 0xed,
	0xf3, 0xa8, 0xd1, 0xdf, 0x4d, 0xd8, 0xb7, 0x4d, 0x3a, 0x43, 0xf1, 0x92, 0x2d, 0x90, 0x7c, 0x07,
	0xbd, 0xf2, 0xc3, 0x7f, 0x5c, 0x54, 0x33, 0x7c, 0x10, 0x9c, 0xdc, 0x29, 0xf8, 0xdd, 0x77, 0xcd,
	0xb7, 0xd0, 0x71, 0x1f, 0xd2, 0xa3, 0xf0, 0xa8, 0x65, 0xeb, 0x0f, 0xde, 0x83, 0xde, 0x25, 0xcb,
	0xb3, 0x73, 0x3d, 0x6f, 0x0e, 0xc3, 0xa3, 0xe6, 0xe1, 0x72, 0x32, 0x0c, 0x49, 0xf2, 0x25, 0x44,
	0x57, 0xa8, 0xa6, 0xf8, 0x7e, 0xbb, 0x1f, 0xc0, 0xd0, 0x46, 0x50, 0x0a, 0xf2, 0x49, 0xb8, 0x63,
	0x67, 0xb6, 0xd7, 0x06, 0x3a, 0xef, 0x98, 0xb7, 0xdd, 0xd7, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff,
	0x97, 0x7b, 0x76, 0x07, 0xeb, 0x09, 0x00, 0x00,
}