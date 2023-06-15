package view

import (
	"context"
	"fmt"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
)

type EventRoute struct {
}

func (e *EventRoute) InitEvent() *dispatcher.EventDispatcher {
	//register event handle
	return dispatcher.NewEventDispatcher(config.GlobalConfig.Feishu.Verification, config.GlobalConfig.Feishu.EncryptKey).OnP2MessageReceiveV1(e.MsgReceive).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	})
}

func (*EventRoute) MsgReceive(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	fmt.Println(larkcore.Prettify(event))
	fmt.Println(event.RequestId())
	return nil
}

func NewEventRoute() *EventRoute {
	return &EventRoute{}
}
