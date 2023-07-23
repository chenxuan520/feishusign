package model

import (
	"context"
	"fmt"
	"time"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
)

func RobotSendTextMsg(recviveID, content string) error {
	content = larkim.NewTextMsgBuilder().Text(content).Build()
	uuid := time.Now().String()
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("user_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(recviveID).
			MsgType("text").
			Content(content).
			Uuid(tools.MD5(uuid)).
			Build()).
		Build()
	resp, err := tools.GlobalLark.Im.Message.Create(context.Background(), req, larkcore.WithTenantKey(tools.GetAccessToken()))
	if err != nil {
		return err
	}
	if !resp.Success() {
		return fmt.Errorf("Error:%d %s %s", resp.Code, resp.Msg, resp.RequestId())
	}
	return nil
}
