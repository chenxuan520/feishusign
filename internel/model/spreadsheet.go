package model

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larksheets "github.com/larksuite/oapi-sdk-go/v3/service/sheets/v3"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
	"net/http"
	"strings"
)

func CreateSpreadSheet(date string) (string, error) {
	spreadSheet := larksheets.NewSpreadsheetBuilder().Title(date + "签到情况").FolderToken(config.GlobalConfig.Sign.FolderToken).Build()
	req := larksheets.NewCreateSpreadsheetReqBuilder().Spreadsheet(spreadSheet).Build()
	res, err := tools.GlobalLark.Sheets.Spreadsheet.Create(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))

	if err != nil {
		return "", fmt.Errorf("send spreadsheet create req err : %v", err)
	}
	if !res.Success() {
		return "", fmt.Errorf("create spreadsheet err : %v", res.Error())
	}

	spreadsheetToken := *res.Data.Spreadsheet.SpreadsheetToken

	url := *res.Data.Spreadsheet.Url

	if err := UpdateContent(date, spreadsheetToken); err != nil {
		return "", err
	}

	if err := UpdateMeetingUrl(date, url); err != nil {
		// TODO 如果更新数据库失败原则上需要将已经创建了的表格删除，但飞书似乎没有提供删除文档的接口
		return "", fmt.Errorf("update meetingUrl err: %v", err)
	}

	return url, nil
}

func UpdateContent(date string, spreadsheetToken string) error {
	chatId, err := GetChatID()
	if err != nil {
		return fmt.Errorf("get chatid err: %v", err)
	}
	members, err := GetUsersByChat(chatId)
	if err != nil {
		return fmt.Errorf("get chat members err: %v", err)
	}

	signLogs, err := BatchSignLogByMeeting(date)
	if err != nil {
		return fmt.Errorf("batch sign logs err: %v", err)
	}

	var values [][]string
	values = append(values, []string{"姓名", "签到情况", "项目组"})
deal:
	for _, info := range members {
		id := info[0]
		name := info[1]
		parts, err := GetUserPartByID(id)
		if err != nil {
			return fmt.Errorf("get user parts err : %v", err)
		}
		for _, part := range parts {
			if part == "导师组" {
				continue deal
			}
		}

		s := "缺席"
		for _, log := range signLogs {
			if id == log.UserID {
				if log.Status == Scan {
					s = "已签到"
				} else {
					s = "请假"
				}
				break
			}
		}
		values = append(values, append([]string{name, s}, parts...))
	}

	sheetId, err := GetFirstSheetId(spreadsheetToken)
	if err != nil {
		return err
	}

	err = InsertItem(spreadsheetToken, sheetId, values)
	if err != nil {
		return fmt.Errorf("insert item err : %v", err)
	}

	return nil
}

func InsertItem(sheetToken, sheetId string, values [][]string) error {
	// 组装请求体
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
		//TODO 打印错误信息，需要将响应体反序列化再获取errMsg，有点麻烦
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

func CheckSpreadSheetIfExist(token string) (bool, error) {
	req := larksheets.NewGetSpreadsheetReqBuilder().SpreadsheetToken(token).Build()
	res, err := tools.GlobalLark.Sheets.Spreadsheet.Get(context.Background(), req, larkcore.WithTenantAccessToken(tools.GetAccessToken()))
	if err != nil {
		logger.GetLogger().Error(err.Error())
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

func GetSpreadsheetTokenByUrl(url string) string {
	tmp := strings.Split(url, "/")
	token := tmp[len(tmp)-1]
	return token
}
