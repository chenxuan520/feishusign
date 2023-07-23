package model

import (
	"context"
	"fmt"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
)

type User struct {
	Name     string
	Status   Status
	Part     []string
	SignTime string
}

func GetUserMsgByCode(code string) (string, string, error) {
	openIDReq := larkauthen.NewCreateAccessTokenReqBuilder().Body(
		larkauthen.NewCreateAccessTokenReqBodyBuilder().
			GrantType("authorization_code").
			Code(code).
			Build()).Build()
	res, err := tools.GlobalLark.Authen.AccessToken.Create(context.Background(), openIDReq, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		return "", "", err
	}
	if res.Data.UserId == nil || res.Data.Name == nil {
		return "", "", fmt.Errorf("get id wrong " + res.Msg)
	}
	return *res.Data.UserId, *res.Data.Name, nil
}

func GetUserPartByID(userID string) ([]string, error) {
	messageReq := larkcontact.NewGetUserReqBuilder().UserId(userID).UserIdType("user_id").Build()
	resmsg, err := tools.GlobalLark.Contact.User.Get(context.Background(), messageReq, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		return nil, err
	}
	parts := []string{}
	for _, v := range resmsg.Data.User.DepartmentIds {
		partReq := larkcontact.NewGetDepartmentReqBuilder().
			DepartmentId(v).
			Build()
		resPart, err := tools.GlobalLark.Contact.Department.Get(context.Background(), partReq)
		if err != nil {
			return nil, err
		}
		parts = append(parts, *resPart.Data.Department.Name)
	}
	return parts, nil
}
