package view

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
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
	feishuReq := larkauthen.NewCreateAccessTokenReqBuilder().Body(
		larkauthen.NewCreateAccessTokenReqBodyBuilder().
			GrantType("authorization_code").
			Code(req.Code).
			Build()).Build()
	res, err := tools.GlobalLark.Authen.AccessToken.Create(context.Background(), feishuReq, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("UserSignIn:%s", err.Error()))
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, map[string]interface{}{
		"data": res,
	})
}
