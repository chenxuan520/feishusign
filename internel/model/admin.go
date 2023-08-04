package model

import (
	"context"
	"fmt"
	larksheets "github.com/larksuite/oapi-sdk-go/v3/service/sheets/v3"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"net/http"
	"strings"
	"time"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
)

func RobotSendTextMsg(receiveID, content string) error {
	content = larkim.NewTextMsgBuilder().Text(content).Build()
	uuid := time.Now().String()
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("user_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(receiveID).
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

func CreateSpreadSheet(date string) (string, error) {
	spreadSheet := larksheets.NewSpreadsheetBuilder().Title(date + "签到情况").FolderToken("MJaMfgcSnlNr99dFSrocKf7nnjb").Build()
	req := larksheets.NewCreateSpreadsheetReqBuilder().Spreadsheet(spreadSheet).Build()
	res, err := tools.GlobalLark.Sheets.Spreadsheet.Create(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))

	if err != nil {
		return "", fmt.Errorf("send spreadsheet create req err : %v", err)
	}
	if !res.Success() {
		return "", fmt.Errorf("create spreadsheet err : %v", res.Error())
	}

	spreadsheetToken := *res.Data.Spreadsheet.SpreadsheetToken

	sheetid, err := GetFirstSheetId(spreadsheetToken)
	if err != nil {
		return "", err
	}
	url := *res.Data.Spreadsheet.Url

	if err := UpdateMeetingUrl(date, url); err != nil {
		// TODO 如果更新数据库失败原则上需要将已经创建了的表格删除，但飞书似乎没有提供删除文档的接口
		return "", fmt.Errorf("update meetingUrl err: %v", err)
	}

	chatId, err := GetChatID()
	if err != nil {
		return "", fmt.Errorf("get chatid err: %v", err)
	}
	members, err := GetUsersByChat(chatId)
	if err != nil {
		return "", fmt.Errorf("get members err: %v", err)
	}
	var values [][]string
	values = append(values, []string{"姓名", "签到情况", "项目组"})
	for _, info := range members {
		id := info[0]
		name := info[1]
		status, err := GetSignStatusById(id, date)
		var s string
		if err != nil {
			// 没找到
			s = "缺席"
		} else {
			if status == Scan {
				s = "已签到"
			} else {
				s = "请假"
			}
		}
		parts, err := GetUserPartByID(id)
		if err != nil {
			return "", fmt.Errorf("get user part err : %v", err)
		}
		values = append(values, []string{name, s, parts[0]})
	}

	err = InsertItem(spreadsheetToken, sheetid, values)
	if err != nil {
		return "", fmt.Errorf("insert item err : %v", err)
	}
	return url, nil
}

func InsertItem(sheetToken, sheetId string, values [][]string) error {

	body := map[string]interface{}{
		"valueRange": map[string]interface{}{
			"range":  sheetId,
			"values": values,
		},
	}

	// 发起请求
	resp, err := tools.GlobalLark.Do(context.Background(),
		&larkcore.ApiReq{
			HttpMethod:                http.MethodPost,
			ApiPath:                   "https://open.feishu.cn/open-apis/sheets/v2/spreadsheets/:spreadsheetToken/values_prepend",
			Body:                      body,
			QueryParams:               nil,
			PathParams:                larkcore.PathParams{"spreadsheetToken": sheetToken},
			SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
		},
		larkcore.WithTenantAccessToken(tools.GetAccessToken()),
	)

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("request err : %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		// TODO 打印错误信息
		//log.Println(string(resp.RawBody))
		return fmt.Errorf("response statusCode : %v", resp.StatusCode)
	}
	return nil
}

func GetFirstSheetId(spreadsheetToken string) (string, error) {
	req := larksheets.NewQuerySpreadsheetSheetReqBuilder().SpreadsheetToken(spreadsheetToken).Build()
	res, err := tools.GlobalLark.Sheets.SpreadsheetSheet.Query(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		return "", fmt.Errorf("send get sheetid req err : %v", err)
	}
	if !res.Success() {
		return "", fmt.Errorf("get sheetid err: %v", res.Error())
	}
	return *res.Data.Sheets[0].SheetId, nil
}

func CheckSpreadSheetIfExist(url string) (bool, error) {
	tmp := strings.Split(url, "/")
	token := tmp[len(tmp)-1]
	req := larksheets.NewGetSpreadsheetReqBuilder().SpreadsheetToken(token).Build()
	res, err := tools.GlobalLark.Sheets.Spreadsheet.Get(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		return false, fmt.Errorf("get spreadsheet info err: %v", err)
	}
	if res.Success() {
		return true, nil
	}
	if res.Code == 1310249 {
		// 该错误码表示文件被删除
		return false, nil
	}
	logger.GetLogger().Error(res.Error())
	return false, fmt.Errorf("%v", res.Error())
}
