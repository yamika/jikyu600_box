package togglbot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	togglToken  = "YOUR_TOGGL_API_TOKEN"
	header      = "Content-Type"
	contentType = "application/json"
)

type Workspace struct {
	ID                            int       `json:"id"`
	Name                          string    `json:"name"`
	Profile                       int       `json:"profile,omitempty"`
	IsPremium                     bool      `json:"premium,omitempty"`
	IsAdmin                       bool      `json:"admin"`
	DefaultHourlyRate             float64   `json:"default_hourly_rate,omitempty"`
	DefaultCurrency               string    `json:"default_currency,omitempty,omitempty"`
	IsOnlyAdminsMayCreateProjects bool      `json:"only_admins_may_create_projects,omitempty"`
	IsOnlyAdminsSeeBillableRates  bool      `json:"only_admins_see_billable_rates,omitempty"`
	IsOnlyAdminsSeeTeamDashboard  bool      `json:"only_admins_see_team_dashboard,omitempty"`
	IsProjectsBillableByDefault   bool      `json:"projects_billable_by_default,omitempty"`
	Rounding                      int       `json:"rounding,omitempty"`
	RoundingMinutes               int       `json:"rounding_minutes,omitempty"`
	APItoken                      string    `json:"api_token"`
	Timestamp                     time.Time `json:"at"`
	IsIcalEnabled                 bool      `json:"ical_enabled,omitempty"`
}

type Project struct {
	ID              int       `json:"id"`
	WID             int       `json:"wid"`
	Name            string    `json:"name"`
	IsBillable      bool      `json:"billable"`
	IsPrivate       bool      `json:"is_private"`
	IsActive        bool      `json:"active"`
	IsTemplate      bool      `json:"template,omitempty"`
	Timestamp       time.Time `json:"at,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	Color           string    `json:"color,omitempty"`
	IsAutoEstimates bool      `json:"auto_estimates,omitempty"`
	ActualHours     int       `json:"actual_hours,omitempty"`
	HexColor        string    `json:"hex_color,omitempty"`
}

type Tag struct {
	ID        int       `json:"id"`
	WID       int       `json:"wid"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"at,omitempty"`
}

type TimeEntry struct {
	ID         int       `json:"id"`
	UID        int       `json:"uid"`
	WID        int       `json:"wid"`
	PID        int       `json:"pid"`
	IsBillable bool      `json:"billable"`
	StartTime  time.Time `json:"start"`
	StopTime   time.Time `json:"stop"`
	Duration   int64     `json:"duration"`
	IsDuronly  bool      `json:"duronly"`
	Timestamp  time.Time `json:"at"`
	Tags       []string  `json:"tags"`
	Desc       string    `json:"description"`
}

type WorkingTime struct {
	AllWorkingTimes float64 `json:"all_working_times"`
	TimeText        string  `json:"time_text"`
}

type ResponseTimeEntry struct {
	Response TimeEntry `json:"data"`
}

type TogglManager struct {
	Workspaces          []Workspace `json:"workspaces"`
	Projects            []Project   `json:"projects"`
	Tags                []Tag       `json:"tags"`
	CurrentTimeEntry    TimeEntry   `json:"time_entry"`
	WorkingTimes        WorkingTime `json:"working_times"`
	HasTimeEntry        bool        `json:"has_time_entry"`
	HasCurrentTimeEntry bool        `json:"has_current_time_entry"`
}

func (t *TogglManager) Init() bool {
	t.HasTimeEntry = false
	t.HasCurrentTimeEntry = false
	var err = t.InitWorkspace()
	if err != nil {
		return true
	}
	log.Printf("Get workspace\n")
	err = t.InitProjects(t.Workspaces[0].ID)
	if err != nil {
		return true
	}
	log.Printf("Get projects\n")

	err = t.InitTags(t.Workspaces[0].ID)
	if err != nil {
		return true
	}
	log.Printf("Get tags\n")

	col := int(time.Now().Day())
	if col == 31{
		col = 30
	}
	data := t.GetSheetCell(col)
	if data[3].Value == "" {
		t.WorkingTimes.TimeText = "Empty"
	}else{
		t.WorkingTimes.TimeText = data[3].Value
	}
	log.Printf("Get worksheet data\n")

	return false
}

func (t *TogglManager) InitWorkspace() error {
	url := "https://www.toggl.com/api/v8/workspaces"
	req, err := http.NewRequest("GET", url, nil)

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

	initWorkspaces := &[]Workspace{}
	json.Unmarshal(bodyText, initWorkspaces)
	t.Workspaces = *initWorkspaces

	return err
}

func (t *TogglManager) InitProjects(wid int) error {
	var bufferURL bytes.Buffer
	bufferURL.WriteString("https://www.toggl.com/api/v8/workspaces/")
	bufferURL.WriteString(strconv.Itoa(wid))
	bufferURL.WriteString("/projects")

	url := bufferURL.String()
	req, err := http.NewRequest("GET", url, nil)

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

	initProjects := &[]Project{}
	json.Unmarshal(bodyText, initProjects)
	t.Projects = *initProjects

	return err
}

func (t *TogglManager) InitTags(wid int) error {
	var bufferURL bytes.Buffer
	bufferURL.WriteString("https://www.toggl.com/api/v8/workspaces/")
	bufferURL.WriteString(strconv.Itoa(wid))
	bufferURL.WriteString("/tags")

	url := bufferURL.String()
	req, err := http.NewRequest("GET", url, nil)

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

	initTags := &[]Tag{}
	json.Unmarshal(bodyText, initTags)
	t.Tags = *initTags

	return err
}

func (t *TogglManager) GetTags(wid int) string {
	var bufferURL bytes.Buffer
	bufferURL.WriteString("https://www.toggl.com/api/v8/workspaces/")
	bufferURL.WriteString(strconv.Itoa(wid))
	bufferURL.WriteString("/tags")

	url := bufferURL.String()
	req, err := http.NewRequest("GET", url, nil)

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

	var result []Tag
	json.Unmarshal(bodyText, &result)

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(bodyText), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return "タグ一覧\n" + buf.String()
}

func (t *TogglManager) GetProjects(wid int) string {
	var bufferURL bytes.Buffer
	bufferURL.WriteString("https://www.toggl.com/api/v8/workspaces/")
	bufferURL.WriteString(strconv.Itoa(wid))
	bufferURL.WriteString("/projects")

	url := bufferURL.String()
	req, err := http.NewRequest("GET", url, nil)

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

	var result []Project
	json.Unmarshal(bodyText, &result)

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(bodyText), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return "プロジェクト一覧\n" + buf.String()
}
