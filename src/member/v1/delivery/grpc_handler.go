package delivery

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Bhinneka/user-service/helper"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Bhinneka/user-service/protogo/member"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
)

const (
	messageJwtCannotBeEmpty = "jwt arg cannot be empty"
	scopeCheckJWT           = "check_jwt_from_arg"
)

// GRPCHandler data structure
type GRPCHandler struct {
	MemberUseCase usecase.MemberUseCase
	QPublisher    service.QPublisher
	Topic         string
	PublicKey     *rsa.PublicKey
}

// NewGRPCHandler function for initializing grpc handler object
func NewGRPCHandler(memberUseCase usecase.MemberUseCase,
	qPublisher service.QPublisher, topic string, publicKey *rsa.PublicKey) *GRPCHandler {
	return &GRPCHandler{
		MemberUseCase: memberUseCase,
		QPublisher:    qPublisher,
		Topic:         topic,
		PublicKey:     publicKey,
	}
}

func (h *GRPCHandler) GetMember(ctxReq context.Context, memberID string) <-chan model.GetMemberResult {
	output := make(chan model.GetMemberResult)
	go func() {
		memberResult := <-h.MemberUseCase.GetDetailMemberByID(ctxReq, memberID)
		if memberResult.Error != nil {
			if memberResult.Error == fmt.Errorf(helper.ErrorDataNotFound, "member") {
				memberResult.Error = errors.New(helper.ErrorUnauthorized)
				output <- model.GetMemberResult{Error: memberResult.Error, Scope: "get_detail_member"}
				return
			}
			output <- model.GetMemberResult{Error: memberResult.Error, Scope: "get_detail_member"}
			return
		}

		member, ok := memberResult.Result.(model.Member)
		if !ok {
			err := errors.New("result is not member")
			output <- model.GetMemberResult{Error: err, Scope: "parse_detail_member"}
			return
		}
		output <- model.GetMemberResult{Result: member}
	}()

	return output

}

// Register function for saving member data
func (h *GRPCHandler) Register(c context.Context, arg *pb.MemberRegister) (*pb.ResponseMessage, error) {
	ctx := "MemberPresenter-Register"

	member := &model.Member{}
	member.FirstName = strings.Trim(arg.FirstName, " \t")
	member.LastName = strings.Trim(arg.LastName, " \t")
	member.Email = strings.Trim(arg.Email, " \t")
	member.NewPassword = strings.Trim(arg.Password, " \t")
	member.RePassword = strings.Trim(arg.RePassword, " \t")
	member.GenderString = strings.Trim(arg.Gender, " \t")
	member.BirthDateString = strings.Trim(arg.Birthdate, " \t")
	member.Mobile = strings.Trim(arg.Mobile, " \t")
	member.Type = "register"

	saveResult := <-h.MemberUseCase.RegisterMember(c, member)
	if saveResult.Error != nil {
		return nil, status.Error(codes.Internal, saveResult.Error.Error())
	}

	result, ok := saveResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New("result is not proper response")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// send the data to NSQ before being sent to dolphin after registration
	dolphinData := serviceModel.MemberDolphin{
		ID:        result.ID,
		Email:     member.Email,
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Gender:    strings.ToUpper(member.GenderString),
		DOB:       member.BirthDateString,
		Mobile:    member.Mobile,
		Created:   time.Now().Format(time.RFC3339),
	}

	NSQPayload := serviceModel.DolphinPayloadNSQ{
		EventType: "register",
		Counter:   0,
		Payload:   dolphinData,
	}

	// prepare to send to nsq
	payloadJSON, _ := json.Marshal(NSQPayload)
	if err := h.QPublisher.Publish(c, h.Topic, shared.MemberRegistration, payloadJSON); err != nil {
		helper.SendErrorLog(c, ctx, "publish_payload", err, NSQPayload)
	}

	msg := pb.ResponseMessage{
		Message:   helper.SuccessMessage,
		Email:     member.Email,
		FirstName: member.FirstName,
		LastName:  member.LastName,
	}

	return &msg, nil
}

// Update function for updating member data
func (h *GRPCHandler) Update(c context.Context, arg *pb.MemberUpdate) (*pb.ResponseMessage, error) {
	memberID := strings.Trim(arg.ID, " \t")

	// get member data
	memberResult := <-h.GetMember(c, memberID)
	if memberResult.Error != nil {
		return nil, status.Error(codes.Internal, memberResult.Error.Error())
	}

	member := memberResult.Result

	// append data from GRPC
	ma := model.Address{}

	member.FirstName = strings.Trim(arg.FirstName, " \t")
	member.LastName = strings.Trim(arg.LastName, " \t")
	member.Mobile = strings.Trim(arg.Mobile, " \t")
	member.Phone = strings.Trim(arg.Phone, " \t")
	member.Ext = strings.Trim(arg.Ext, " \t")
	member.GenderString = strings.Trim(arg.Gender, " \t")
	member.BirthDateString = strings.Trim(arg.Birthdate, " \t")

	ma.Street1 = strings.Trim(arg.Street1, " \t")
	ma.Street2 = strings.Trim(arg.Street2, " \t")
	ma.ZipCode = strings.Trim(arg.ZipCode, " \t")
	ma.SubDistrictID = strings.Trim(arg.SubDistrictID, " \t")
	ma.SubDistrict = strings.Trim(arg.SubDistrict, " \t")
	ma.DistrictID = strings.Trim(arg.DistrictID, " \t")
	ma.District = strings.Trim(arg.District, " \t")
	ma.CityID = strings.Trim(arg.CityID, " \t")
	ma.City = strings.Trim(arg.City, " \t")
	ma.ProvinceID = strings.Trim(arg.ProvinceID, " \t")
	ma.Province = strings.Trim(arg.ProvinceID, " \t")

	member.Address = ma
	member.Type = "update"

	saveResult := <-h.MemberUseCase.UpdateDetailMemberByID(c, member)
	if saveResult.Error != nil {
		return nil, status.Error(codes.Internal, saveResult.Error.Error())
	}

	msg := pb.ResponseMessage{
		Message: helper.SuccessMessage,
	}

	return &msg, nil
}

// GetMe function for finding detail member by its JWT
func (h *GRPCHandler) GetMe(c context.Context, arg *pb.MemberQuery) (*pb.Member, error) {
	ctx := "MemberGRPCPresenter-GetMe"

	accessToken := strings.Trim(arg.JWT, " \t")
	if accessToken == "" {
		return nil, status.Error(codes.InvalidArgument, messageJwtCannotBeEmpty)
	}

	claims, err := shared.JWTExtract(h.PublicKey, accessToken)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	memberID := claims.Subject

	msg, err := h.generateMsgMember(c, ctx, memberID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &msg, nil
}

// FindByID function for finding detail member by ID
func (h *GRPCHandler) FindByID(c context.Context, arg *pb.MemberQuery) (*pb.Member, error) {
	ctx := "MemberPresenter-FindByID"

	memberID := strings.Trim(arg.ID, " \t")

	msg, err := h.generateMsgMember(c, ctx, memberID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &msg, nil
}

// UpdatePassword function for updating password
func (h *GRPCHandler) UpdatePassword(c context.Context, arg *pb.MemberPasswordRequest) (*pb.ResponseMessage, error) {

	memberID := strings.Trim(arg.ID, " \t")

	oldPassword := strings.Trim(arg.OldPassword, " \t")
	newPassword := strings.Trim(arg.NewPassword, " \t")

	accessToken := strings.Trim(arg.JWT, " \t")
	if accessToken == "" {
		return nil, status.Error(codes.InvalidArgument, messageJwtCannotBeEmpty)
	}

	var token string
	if split := strings.Split(accessToken, " "); len(split) > 1 {
		token = split[1]
	}

	passResult := <-h.MemberUseCase.UpdatePassword(c, token, memberID, oldPassword, newPassword)

	if passResult.Error != nil {
		return nil, status.Error(codes.Internal, passResult.Error.Error())
	}

	res, ok := passResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New("result is not proper response")
		return nil, status.Error(codes.Internal, err.Error())
	}

	msg := pb.ResponseMessage{
		Message:   helper.SuccessMessage,
		Email:     res.Email,
		FirstName: res.FirstName,
		LastName:  res.LastName,
	}

	return &msg, nil
}

// generateMsgMember function for generate msg detail member
func (h *GRPCHandler) generateMsgMember(c context.Context, ctx, memberID string) (pb.Member, error) {
	msg := pb.Member{}
	// get member data
	memberResult := <-h.GetMember(c, memberID)
	if memberResult.Error != nil {
		return msg, memberResult.Error
	}

	member := memberResult.Result

	address := &pb.Address{
		Province:      member.Address.Province,
		ProvinceID:    member.Address.ProvinceID,
		City:          member.Address.City,
		CityID:        member.Address.CityID,
		District:      member.Address.District,
		DistrictID:    member.Address.DistrictID,
		SubDistrict:   member.Address.SubDistrict,
		SubDistrictID: member.Address.SubDistrictID,
		ZipCode:       member.Address.ZipCode,
		Street1:       member.Address.Street1,
		Street2:       member.Address.Street2,
	}

	socialMedia := &pb.SocialMedia{
		FacebookID: member.SocialMedia.FacebookID,
		GoogleID:   member.SocialMedia.GoogleID,
		AzureID:    member.SocialMedia.AzureID,
	}

	msg = pb.Member{
		ID:           member.ID,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		Email:        member.Email,
		Gender:       member.Gender.String(),
		Mobile:       member.Mobile,
		Phone:        member.Phone,
		Ext:          member.Ext,
		Birthdate:    member.BirthDateString,
		Address:      address,
		JobTitle:     member.JobTitle,
		Department:   member.Department,
		Status:       member.Status.String(),
		SocialMedia:  socialMedia,
		IsAdmin:      member.IsAdmin,
		IsStaff:      member.IsStaff,
		SignUpFrom:   member.SignUpFrom,
		HasPassword:  member.HasPassword,
		LastLogin:    member.LastLoginString,
		Version:      int32(member.Version),
		Created:      member.CreatedString,
		LastBlocked:  member.LastBlockedString,
		LastModified: member.LastModifiedString,
	}

	return msg, nil
}
