package view

import (
	"context"
	"encoding/json"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
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
	events := dispatcher.NewEventDispatcher(config.GlobalConfig.Feishu.Verification, config.GlobalConfig.Feishu.EncryptKey)
	events.OnP2MessageReceiveV1(e.MsgReceive)
	return events
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
		e.AdminService.AdminSend(userID, "happy new year")
	}
	return nil
}

func (*EventRoute) sendRobotMsg() {
	// TODO : send csv file
}

func NewEventRoute() *EventRoute {
	return &EventRoute{
		AdminService: service.NewAdminService(),
	}
}
