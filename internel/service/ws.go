package service

import (
	"encoding/json"
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

var DefaultWsService *WsService = nil

type WsService struct {
	hashmap *sync.Map
}

type WsConn struct {
	MeetingID string
	UserID    string
	mux       sync.RWMutex
	url       string
	conn      *websocket.Conn
	Exit      chan struct{}
}

type MeetingMsg struct {
	MeetingID string `json:"meeting_id"`
	Code      string `json:"code"`
}

func (w *WsConn) Serve() {
	go w.downStream()
	go w.upStream()
}

func (w *WsConn) downStream() {
	//update first
	url := w.updateUrl()
	err := w.conn.WriteMessage(websocket.TextMessage, []byte(url))
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintln("Error:write msg wrong ", err.Error()))
		return
	}
	//loop write
	for {
		select {
		case <-w.Exit:
			return
		case <-time.After(ChangeTime):
			url := w.updateUrl()
			err := w.conn.WriteMessage(websocket.TextMessage, []byte(url))
			if err != nil {
				logger.GetLogger().Error(fmt.Sprintln("Error:write msg wrong ", err.Error()))
				return
			}
		}
	}
}

func (w *WsConn) upStream() {
	for {
		//check if close
		_, _, err := w.conn.ReadMessage()
		if err != nil {
			close(w.Exit)
			w.conn.Close()
			if DefaultWsService != nil {
				DefaultWsService.DelWsConn(w.MeetingID)
			}
			return
		}
	}
}

func (w *WsConn) updateUrl() string {
	//lock
	w.mux.Lock()
	defer w.mux.Unlock()

	//calc adn store md5
	val := tools.MD5(time.Now().String() + HashSalt)
	w.url = val

	msg := MeetingMsg{
		MeetingID: w.MeetingID,
		Code:      val,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Error: json encode %s", err.Error()))
		return ""
	}
	base64 := tools.Base64Encode(data)
	str := url.QueryEscape(config.GlobalConfig.Server.RedirectURL)
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=%s&app_id=%s&state=%s", str, config.GlobalConfig.Feishu.AppID, base64)
	return url
}

func (w *WsService) AddWsConn(conn *websocket.Conn, userID, meetingID string) error {
	wsconn := &WsConn{
		MeetingID: meetingID,
		UserID:    userID,
		mux:       sync.RWMutex{},
		url:       "",
		conn:      conn,
		Exit:      make(chan struct{}),
	}
	_, ok := w.hashmap.LoadOrStore(meetingID, wsconn)
	if ok {
		logger.GetLogger().Error(fmt.Sprintln("Error: exist meeting ", meetingID))
		close(wsconn.Exit)
		conn.Close()
		return fmt.Errorf("Error: add conn wrong ")
	}
	go wsconn.Serve()
	return nil
}

func (w *WsService) GetMeetingUrl(meeting string) (string, error) {
	//check conn
	val, ok := w.hashmap.Load(meeting)
	if !ok {
		logger.GetLogger().Error(fmt.Sprintf("Error: empty %s", meeting))
		return "", fmt.Errorf("empty meeting")
	}
	//assert interface
	conn, ok := val.(*WsConn)
	if !ok || conn == nil {
		logger.GetLogger().Error(fmt.Sprintf("Error: conn %s", meeting))
		return "", fmt.Errorf("conn interface wrong")
	}
	//lock
	conn.mux.RLock()
	defer conn.mux.RUnlock()
	return conn.url, nil
}

func (w *WsService) DelWsConn(meetingID string) {
	w.hashmap.Delete(meetingID)
}

func NewWsService() *WsService {
	DefaultWsService = &WsService{
		hashmap: &sync.Map{},
	}
	return DefaultWsService
}
