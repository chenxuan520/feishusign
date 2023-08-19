package model

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
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
	if !res.Success() || res.Data.UserId == nil || res.Data.Name == nil {
		err := res.Error()
		logger.GetLogger().Error(err)
		return "", "", fmt.Errorf(err)
	}
	return *res.Data.UserId, *res.Data.Name, nil
}

func GetUserPartByID(userID string) ([]string, error) {
	messageReq := larkcontact.NewGetUserReqBuilder().UserId(userID).UserIdType("user_id").Build()
	resMsg, err := tools.GlobalLark.Contact.User.Get(context.Background(), messageReq, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		return nil, err
	}
	if !resMsg.Success() {
		err := resMsg.Error()
		logger.GetLogger().Error(err)
		return nil, fmt.Errorf(err)
	}
	var parts []string
	for _, v := range resMsg.Data.User.DepartmentIds {
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
	var members [][2]string
	hasMore := true
	pageToken := ""

	for hasMore {
		req := larkim.NewGetChatMembersReqBuilder().ChatId(charID).MemberIdType("user_id").
			PageSize(100).PageToken(pageToken).Build()
		res, err := tools.GlobalLark.Im.ChatMembers.Get(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
		if err != nil {
			return nil, err
		}
		if !res.Success() {
			return nil, fmt.Errorf("%v", res.Error())
		}
		for _, v := range res.Data.Items {
			members = append(members, [2]string{*v.MemberId, *v.Name})
		}
		hasMore = *res.Data.HasMore
		pageToken = *res.Data.PageToken
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
	for _, v := range res.Data.Items {
		if *v.Name == "Dian团队在站队员交流群" {
			//TODO 写进配置
			return *v.ChatId, nil
		}
	}
	return "", fmt.Errorf("dian Group not found")
}

func GetUsernameById(userId string) (string, error) {
	req := larkcontact.NewGetUserReqBuilder().UserIdType("user_id").UserId(userId).Build()
	res, err := tools.GlobalLark.Contact.User.Get(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		return "", err
	}
	if !res.Success() {
		return "", fmt.Errorf("%v", res.Error())
	}
	return *res.Data.User.Name, nil
}
