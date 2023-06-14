package view

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/request"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/response"
)

func UserSignIn(c *gin.Context) {
	req := request.ReqSignin{
		Code:  "",
		State: "",
	}
	req.Code = c.Query("code")
	req.State = c.Query("state")
	if req.Code == "" || req.State == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("please log first"))
		return
	}
	openIDReq := larkauthen.NewCreateAccessTokenReqBuilder().Body(
		larkauthen.NewCreateAccessTokenReqBodyBuilder().
			GrantType("authorization_code").
			Code(req.Code).
			Build()).Build()
	res, err := tools.GlobalLark.Authen.AccessToken.Create(context.Background(), openIDReq, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("UserSignIn:%s", err.Error()))
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	if res.Data.UserId == nil {
		response.ErrorDetail(c, http.StatusBadRequest, map[string]interface{}{"data": res}, fmt.Errorf("get empty userid"))
		return
	}
	userID := *res.Data.UserId
	fmt.Println(userID)
	messageReq := larkcontact.NewGetUserReqBuilder().UserId(userID).UserIdType("user_id").Build()
	resmsg, err := tools.GlobalLark.Contact.User.Get(context.Background(), messageReq, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("UserSignInresMsg:%s", err.Error()))
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	parts := []string{}
	for _, v := range resmsg.Data.User.DepartmentIds {
		partReq := larkcontact.NewGetDepartmentReqBuilder().
			DepartmentId(v).
			Build()
		// 发起请求
		resPart, err := tools.GlobalLark.Contact.Department.Get(context.Background(), partReq)
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("GetPart:%s", err.Error()))
			response.Error(c, http.StatusBadRequest, err)
			return
		}
		parts = append(parts, *resPart.Data.Department.Name)
	}
	response.Success(c, map[string]interface{}{
		"data": resmsg,
		"part": parts,
	})
}
