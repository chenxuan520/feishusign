package view

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/service"
)

type EventRoute struct {
	AdminService *service.AdminService
}

func (e *EventRoute) InitEvent() *dispatcher.EventDispatcher {
	//register event handle
	events := dispatcher.
		NewEventDispatcher(config.GlobalConfig.Feishu.Verification, config.GlobalConfig.Feishu.EncryptKey).
		OnP2MessageReceiveV1(e.MsgReceive).
		OnP1LeaveApprovalV4(e.LeaveEventApproval)
	return events
}

func (e *EventRoute) LeaveEventApproval(ctx context.Context, event *larkapproval.P1LeaveApprovalV4) error {
	if event.Event == nil {
		logger.GetLogger().Error("deal event err: nil event.Event")
		return nil
	}
	if event.Event.UserID == "" {
		logger.GetLogger().Error("deal event err: no userId")
		return nil
	}
	userId := event.Event.UserID
	if event.Event.LeaveName != "@i18n@6959807929197281283" {
		// 这个字符串代表例会假，在审批应用的管理后台，可能更改
		e.AdminService.AdminSend(userId, "请假失败，假期名称错误，请联系管理员")
		logger.GetLogger().Debug("debug: not @i18n@6959807929197281283 but " + event.Event.LeaveName)
		return nil
	}

	if err := e.AdminService.AdminDealLeave(userId, event.Event.LeaveStartTime); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("deal leave approval error: %v", err))
		e.AdminService.AdminSend(userId, "请假失败，请联系管理员查看后台日志排错")
		return nil
	}
	e.AdminService.AdminSend(userId, fmt.Sprintf("请假成功，请假时间为：%s", event.Event.LeaveStartTime))
	return nil
}

func (e *EventRoute) MsgReceive(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	if event.Event == nil || event.Event.Message == nil || event.Event.Sender == nil ||
		event.Event.Sender.SenderId.UserId == nil || event.Event.Message.MessageType == nil {
		logger.GetLogger().Error("Error:get wrong message")
		return nil
	}
	userID := *event.Event.Sender.SenderId.UserId

	text := larkim.MessagePostText{}
	err := json.Unmarshal([]byte(*event.Event.Message.Content), &text)
	if err != nil {
		e.AdminService.AdminSend(userID, err.Error())
		return nil
	}
	switch *event.Event.Message.MessageType {
	case "text":
		e.AdminService.AdminDealMsg(userID, text.Text)
	default:
		e.AdminService.AdminSend(userID, "你好，请输入文本信息")
	}
	return nil
}

func NewEventRoute() *EventRoute {
	return &EventRoute{
		AdminService: service.NewAdminService(),
	}
}
