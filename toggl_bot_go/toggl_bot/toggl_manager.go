package toggl_bot

import (
  "time"
  "log"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "strconv"
)

const (
	togglToken = "YOUR_TOGGL_API_TOKEN"
  header = "Content-Type"
  content_type = "application/json"
)

type Workspace struct{
  Id int `json:"id"`
  Name string `json:"name"`
  Profile int `json:"profile,omitempty"`
  IsPremium bool `json:"premium,omitempty"`
  IsAdmin bool `json:"admin"`
  DefaultHourlyRate float64 `json:"default_hourly_rate,omitempty"`
  DefaultCurrency  string `json:"default_currency,omitempty,omitempty"`
  IsOnlyAdminsMayCreateProjects bool `json:"only_admins_may_create_projects,omitempty"`
  IsOnlyAdminsSeeBillableRates bool `json:"only_admins_see_billable_rates,omitempty"`
  IsOnlyAdminsSeeTeamDashboard bool `json:"only_admins_see_team_dashboard,omitempty"`
  IsProjectsBillableByDefault bool `json:"projects_billable_by_default,omitempty"`
  Rounding int `json:"rounding,omitempty"`
  RoundingMinutes int `json:"rounding_minutes,omitempty"`
  APItoken string `json:"api_token"`
  Timestamp time.Time `json:"at"`
  IsIcalEnabled bool `json:"ical_enabled,omitempty"`
}

type Project struct {
    Id int `json:"id"`
    Wid int `json:"wid"`
    Name string `json:"name"`
    IsBillable bool `json:"billable"`
    IsPrivate bool  `json:"is_private"`
    IsActive bool `json:"active"`
    IsTemplate bool `json:"template,omitempty"`
    Timestamp time.Time `json:"at,omitempty"`
    CreatedAt time.Time `json:"created_at,omitempty"`
    Color string `json:"color,omitempty"`
    IsAutoEstimates bool `json:"auto_estimates,omitempty"`
    ActualHours int `json:"actual_hours,omitempty"`
    HexColor string `json:"hex_color,omitempty"`
}

type Tag struct {
  Id int `json:"id"`
  Wid int `json:"wid"`
  Name string `json:"name"`
  Timestamp time.Time `json:"at,omitempty"`
}

type TimeEntry struct {
  Id int `json:"id"`
  Uid int `json:"uid"`
  Wid int `json:"wid"`
  Pid int `json:"pid"`
  IsBillable bool `json:"billable"`
  StartTime time.Time `json:"start"`
  StopTime time.Time `json:"stop"`
  Duration int64 `json:"duration"`
  IsDuronly bool `json:"duronly"`
  Timestamp time.Time `json:"at"`
  Tags []string `json:"tags"`
  Desc string `json:"description"`
}

type WorkingTime struct {
  AllWorkingTimes float64 `json:"all_working_times"`
  TimeText string `json:"time_text"`
}

type ResponseTimeEntry struct{
  Response TimeEntry `json:"data"`
}

type Toggl_manager struct {
    Workspaces []Workspace `json:"workspaces"`
    Projects []Project `json:"projects"`
    Tags []Tag `json:"tags"`
    CurrentTimeEntry TimeEntry `json:"time_entry"`
    WorkingTimes WorkingTime `json:"working_times"`
    HasTimeEntry bool `json:"has_time_entry"`
    HasCurrentTimeEntry bool `json:"has_current_time_entry"`
}

func (t *Toggl_manager) Init() bool {
  t.HasTimeEntry = false
  t.HasCurrentTimeEntry = false
  var err = t.Init_workspace()
  if err != nil{
    return true
  }
  log.Printf("Get workspace\n")
  err = t.Init_projects(t.Workspaces[0].Id)
  if err != nil{
    return true
  }
  log.Printf("Get projects\n")

  err = t.Init_tags(t.Workspaces[0].Id)
  if err != nil{
    return true
  }
  log.Printf("Get tags\n")

  col := int(time.Now().Day())
  data := t.Get_sheet_cell(col)
  if data[3].Value == ""{
    t.WorkingTimes.TimeText  = "Empty"
  }
  t.WorkingTimes.TimeText  = data[3].Value
  log.Printf("Get worksheet data\n")

  return false
}

func (t *Toggl_manager) Init_workspace() error {
  url := "https://www.toggl.com/api/v8/workspaces"
  req,err := http.NewRequest("GET",url,nil)

  if err != nil {
    panic(err.Error)
  }

  req.Header.Set(header,content_type)
  req.SetBasicAuth(togglToken,"api_token")
  client := &http.Client{}
  resp,err := client.Do(req)
  if err != nil{
    panic(err.Error)
  }
  defer resp.Body.Close()

  bodyText, err := ioutil.ReadAll(resp.Body)
  if err != nil{
    panic(err.Error)
  }

  initWorkspaces := &[]Workspace{}
  json.Unmarshal(bodyText,initWorkspaces)
  t.Workspaces = *initWorkspaces

  return err
}

func (t *Toggl_manager) Init_projects(wid int) error {
  var buffer_url bytes.Buffer
  buffer_url.WriteString("https://www.toggl.com/api/v8/workspaces/")
  buffer_url.WriteString(strconv.Itoa(wid))
  buffer_url.WriteString("/projects")

  url := buffer_url.String()
  req,err := http.NewRequest("GET",url,nil)

  if err != nil {
    panic(err.Error)
  }

  req.Header.Set(header,content_type)
  req.SetBasicAuth(togglToken,"api_token")
  client := &http.Client{}
  resp,err := client.Do(req)
  if err != nil{
    panic(err.Error)
  }
  defer resp.Body.Close()

  bodyText, err := ioutil.ReadAll(resp.Body)
  if err != nil{
    panic(err.Error)
  }

  initProjects := &[]Project{}
  json.Unmarshal(bodyText,initProjects)
  t.Projects = *initProjects

  return err
}

func (t *Toggl_manager) Init_tags(wid int) error {
  var buffer_url bytes.Buffer
  buffer_url.WriteString("https://www.toggl.com/api/v8/workspaces/")
  buffer_url.WriteString(strconv.Itoa(wid))
  buffer_url.WriteString("/tags")

  url := buffer_url.String()
  req,err := http.NewRequest("GET",url,nil)

  if err != nil {
    panic(err.Error)
  }

  req.Header.Set(header,content_type)
  req.SetBasicAuth(togglToken,"api_token")
  client := &http.Client{}
  resp,err := client.Do(req)
  if err != nil{
    panic(err.Error)
  }
  defer resp.Body.Close()

  bodyText, err := ioutil.ReadAll(resp.Body)
  if err != nil{
    panic(err.Error)
  }

  initTags := &[]Tag{}
  json.Unmarshal(bodyText,initTags)
  t.Tags = *initTags

  return err
}

func (t *Toggl_manager) Get_Tags(wid int) string{
  var buffer_url bytes.Buffer
  buffer_url.WriteString("https://www.toggl.com/api/v8/workspaces/")
  buffer_url.WriteString(strconv.Itoa(wid))
  buffer_url.WriteString("/tags")

  url := buffer_url.String()
  req,err := http.NewRequest("GET",url,nil)

  if err != nil {
    panic(err.Error)
  }

  req.Header.Set(header,content_type)
  req.SetBasicAuth(togglToken,"api_token")
  client := &http.Client{}
  resp,err := client.Do(req)
  if err != nil{
    panic(err.Error)
  }
  defer resp.Body.Close()

  bodyText, err := ioutil.ReadAll(resp.Body)
  if err != nil{
   panic(err.Error)
  }

  var result []Tag
  json.Unmarshal(bodyText,&result)

  var buf bytes.Buffer
  err = json.Indent(&buf,[]byte(bodyText),"","  ")
  if err != nil{
    log.Fatal(err)
  }

  return "タグ一覧\n"+buf.String()
}

func (t *Toggl_manager) Get_Projects(wid int) string{
  var buffer_url bytes.Buffer
  buffer_url.WriteString("https://www.toggl.com/api/v8/workspaces/")
  buffer_url.WriteString(strconv.Itoa(wid))
  buffer_url.WriteString("/projects")

  url := buffer_url.String()
  req,err := http.NewRequest("GET",url,nil)

  if err != nil {
    panic(err.Error)
  }

  req.Header.Set(header,content_type)
  req.SetBasicAuth(togglToken,"api_token")
  client := &http.Client{}
  resp,err := client.Do(req)
  if err != nil{
    panic(err.Error)
  }
  defer resp.Body.Close()

  bodyText, err := ioutil.ReadAll(resp.Body)
  if err != nil{
    panic(err.Error)
  }

  var result []Project
  json.Unmarshal(bodyText,&result)

  var buf bytes.Buffer
  err = json.Indent(&buf,[]byte(bodyText),"","  ")
  if err != nil{
    log.Fatal(err)
  }

  return "プロジェクト一覧\n"+buf.String()
}
