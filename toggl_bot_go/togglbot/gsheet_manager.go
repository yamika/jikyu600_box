package togglbot

import (
	"bytes"
	"encoding/json"
	"github.com/Iwark/spreadsheet"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

// GSheet has your google spreadsheet columns
type GSheet struct {
	Date         string `json:"日付,omitempty"`
	AllWorkTimes string `json:"総時間,omitempty"`
	WorkDesc     string `json:"内容,omitempty"`
	WorkTimeText string `json:"その他,omitempty"`
}

func (t *TogglManager) GetSheetCell(col int) []spreadsheet.Cell {
	data, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		panic(err.Error())
	}

	conf, err := google.JWTConfigFromJSON(data, "https://spreadsheets.google.com/feeds")
	if err != nil {
		panic(err.Error())
	}

	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)

	spreadsheet, err := service.FetchSpreadsheet("YOUR_SPREADSHEET_ID")
	if err != nil {
		panic(err.Error())
	}

	sheet, err := spreadsheet.SheetByIndex(0)
	if err != nil {
		panic(err.Error())
	}
	return sheet.Rows[2+col]
}

func (t *TogglManager) GetWorksheet(col int) string {
	data := t.GetSheetCell(col)
	var gsheet GSheet
	gsheet.Date = data[0].Value
	gsheet.AllWorkTimes = data[1].Value
	gsheet.WorkDesc = data[2].Value
	gsheet.WorkTimeText = data[3].Value

	b, _ := json.Marshal(&gsheet)
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(b), "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func (t *TogglManager) SendWorkingTime(desc string) string {
	col := int(time.Now().Day())
	if col == 31{
		col = 30
	}
	data, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		panic(err.Error())
	}

	conf, err := google.JWTConfigFromJSON(data, "https://spreadsheets.google.com/feeds")
	if err != nil {
		panic(err.Error())
	}

	client := conf.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)

	spreadsheet, err := service.FetchSpreadsheet("YOUR_SPREADSHEET_ID")
	if err != nil {
		panic(err.Error())
	}

	sheet, err := spreadsheet.SheetByIndex(0)
	if err != nil {
		panic(err.Error())
	}
	if desc == "" {
		sheet.Update(2+col, 1, strconv.FormatFloat(t.WorkingTimes.AllWorkingTimes, 'f', 2, 64))
		sheet.Update(2+col, 3, t.WorkingTimes.TimeText)
		err := sheet.Synchronize()
		if err != nil {
			panic(err.Error())
		}
		return "更新したぜ\n" + t.GetWorksheet(col)
	}
	sheet.Update(2+col, 2, desc)
	err = sheet.Synchronize()
	if err != nil {
		panic(err.Error())
	}
	return "説明を更新したぜ\n" + t.GetWorksheet(col)
}
