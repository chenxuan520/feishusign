package service

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
)

// TODO: add to config //
const (
	ChangeTime = time.Second * 7
	HashKey    = "dian"
	HashSalt   = "sign"
)

type WsService struct {
	conn    *websocket.Conn
	hashmap sync.Map
	Exit    chan struct{}
}

func (w *WsService) DownStream() {
	for {
		select {
		case <-w.Exit:
			w.conn.Close()
			return
		case <-time.After(ChangeTime):
			url := w.updateUrl()
			err := w.conn.WriteMessage(websocket.TextMessage, []byte(url))
			if err != nil {
				logger.GetLogger().Error(fmt.Sprintln("Error:write msg wrong ", err.Error()))
				w.conn.Close()
				return
			}
		}
	}
}

func (w *WsService) UpStream() {
	for {
		_, _, err := w.conn.ReadMessage()
		if err != nil {
			//TODO
		}
	}
}

func (w *WsService) updateUrl() string {
	val := tools.SHA1(time.Now().String() + HashSalt)
	w.hashmap.Store(HashKey, val)
	str := url.QueryEscape(config.GlobalConfig.Server.RedirectURL)
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=%s&app_id=%s&state=%s", str, config.GlobalConfig.Feishu.AppID, "hello")
	return url
}

func NewWsService(conn *websocket.Conn) *WsService {
	ws := &WsService{
		conn:    conn,
		hashmap: sync.Map{},
		Exit:    make(chan struct{}),
	}

	return ws
}
