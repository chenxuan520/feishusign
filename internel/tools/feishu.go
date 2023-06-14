package tools

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

var GlobalLark *lark.Client = nil

//feishu accessToken,use by GetAccessToken
type feishuToken struct {
	mut   sync.RWMutex
	token string
}

const updateTokenTime = time.Hour + 40*time.Minute

var accessToken feishuToken

func InitLarkClient(appID, appSecret string) error {
	GlobalLark = lark.NewClient(appID, appSecret)
	if GlobalLark == nil {
		return fmt.Errorf("InitLarkClient wrong")
	}
	go loopUpdateToken(appID, appSecret)
	return nil
}

func loopUpdateToken(appID, appSecret string) {
	for {
		select {
		case <-time.After(updateTokenTime):
			err := updateAccessToken(appID, appSecret)
			if err != nil {
				log.Println(fmt.Sprintf("update token wrong %s", err.Error()))
			}
		}
	}
}

func updateAccessToken(appID, appSecret string) error {
	accessToken.mut.Lock()
	defer accessToken.mut.Unlock()

	req := &larkcore.SelfBuiltTenantAccessTokenReq{
		AppID:     appID,
		AppSecret: appSecret,
	}
	res, err := GlobalLark.GetTenantAccessTokenBySelfBuiltApp(context.Background(), req)
	if err != nil {
		accessToken.token = ""
		return err
	}

	accessToken.token = res.TenantAccessToken
	return nil
}

func GetAccessToken() string {
	accessToken.mut.RLock()
	defer accessToken.mut.RUnlock()

	result := accessToken.token
	return result
}
