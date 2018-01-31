package togglbot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (t *TogglManager) StartTimer(wid int) string {
	url := "https://www.toggl.com/api/v8/time_entries/start"
	var r struct {
		Parameter struct {
			ID          int    `json:"id,omitempty"`
			WID         int    `json:"wid"`
			StartTime   string `json:"start"`
			CreatedWith string `json:"created_with,omitempty"`
		} `json:"time_entry"`
	}
	r.Parameter.WID = wid
	r.Parameter.StartTime = time.Now().String()
	r.Parameter.CreatedWith = "slack"

	param, err := json.Marshal(r)

	if err != nil {
		panic(err.Error)
	}

	payload := bytes.NewBuffer([]byte(string(param)))
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		panic(err.Error)
	}

	req.Header.Set(header, contentType)
	req.SetBasicAuth(togglToken, "api_token")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error)
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error)
	}

	var result ResponseTimeEntry
	json.Unmarshal(bodyText, &result)
	t.CurrentTimeEntry = result.Response
	t.HasCurrentTimeEntry = true

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(bodyText), "", "  ")
	if err != nil {
		panic(err.Error)
	}

	return "スタートしたぜ\n" + buf.String()
}

func (t *TogglManager) StopTimer(id int) string {
	var bufferURL bytes.Buffer
	bufferURL.WriteString("https://www.toggl.com/api/v8/time_entries/")
	bufferURL.WriteString(strconv.Itoa(id))
	bufferURL.WriteString("/stop")
	url := bufferURL.String()

	req, err := http.NewRequest("PUT", url, nil)

	if err != nil {
		panic(err.Error)
	}

	req.Header.Set(header, contentType)
	req.SetBasicAuth(togglToken, "api_token")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error)
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error)
	}

	var result ResponseTimeEntry
	json.Unmarshal(bodyText, &result)
	t.CurrentTimeEntry = result.Response
	t.HasTimeEntry = true
	t.HasCurrentTimeEntry = false

	t.CalcTimeFromTimeEntries()

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(bodyText), "", "  ")
	if err != nil {
		panic(err.Error)
	}

	return "総作業時間 : " + t.WorkingTimes.TimeText + "\n"
}

func (t *TogglManager) CalcTimeFromTimeEntries() {
	// UTCをJSTに変換
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	start := t.CurrentTimeEntry.StartTime.In(jst)
	stop := t.CurrentTimeEntry.StopTime.In(jst)
	duration := stop.Sub(start)

	hours0 := int(duration.Hours())
	hours := float64(hours0 % 24)
	mins0 := int(duration.Minutes()) % 60
	mins := 0.
	if mins0 > 29 {
		mins = 0.5
	}

	t.WorkingTimes.AllWorkingTimes += (hours + mins)
	starttime := strconv.Itoa(int(start.Hour())) + ":" + strconv.Itoa(int(start.Minute()))
	stoptime := strconv.Itoa(int(stop.Hour())) + ":" + strconv.Itoa(int(stop.Minute()))

	if t.WorkingTimes.TimeText == "Empty" {
		t.WorkingTimes.TimeText = starttime + "~" + stoptime
	} else {
		t.WorkingTimes.TimeText += (", " + starttime + "~" + stoptime)
	}
}
