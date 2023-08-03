package model

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
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

func GetUsersByChat(charID string) ([][2]string, error) {
	req := larkim.NewGetChatMembersReqBuilder().ChatId(charID).MemberIdType("user_id").Build()
	res, err := tools.GlobalLark.Im.ChatMembers.Get(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		return nil, err
	}
	if !res.Success() {
		return nil, fmt.Errorf("%v", res.Error())
	}
	var members [][2]string
	for _, v := range res.Data.Items {
		members = append(members, [2]string{*v.MemberId, *v.Name})
	}
	return members, nil
}

func GetChatID() (string, error) {
	req := larkim.NewListChatReqBuilder().UserIdType("user_id").SortType("ByCreateTimeAsc").Build()
	res, err := tools.GlobalLark.Im.Chat.List(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		return "", err
	}
	if !res.Success() {
		return "", fmt.Errorf("%v", res.Error())
	}
	// TODO 机器人会在多个群中，需要获取特定的群
	return *res.Data.Items[0].ChatId, nil
}
