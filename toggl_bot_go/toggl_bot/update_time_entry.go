package toggl_bot

import (
  "log"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "bytes"
  "strconv"
)

func (t *Toggl_manager) Find_projectid_from_Projects(pname string) int {
  for _,p := range t.Projects{
    if pname == p.Name{
      return p.Id
    }
  }
  return -1
}

func (t *Toggl_manager) Has_tags(tags []string) bool {
  IsValid := false
  for _,p := range tags{
    for _,tag := range t.Tags {
      if p == tag.Name {
        IsValid = true
        break
      }
    }
    if !IsValid {
      break
    }
  }
  return IsValid
}

func (t *Toggl_manager) Update_current_time_entry_project(id int,pname string) string{
  var buffer_url bytes.Buffer
  buffer_url.WriteString("https://www.toggl.com/api/v8/time_entries/")
  buffer_url.WriteString(strconv.Itoa(id))
  url := buffer_url.String()

  var r struct{
    Parameter struct {
      ProjectId int `json:"pid"`
    } `json:"time_entry"`
  }

  r.Parameter.ProjectId = t.Find_projectid_from_Projects(pname)
  if r.Parameter.ProjectId < 0 {
    return "プロジェクトが見つからないぜ"
  }

  param,err := json.Marshal(r)
  payload := bytes.NewBuffer([]byte(string(param)))
  if err != nil {
    panic(err.Error)
  }

  req,err := http.NewRequest("PUT",url,payload)
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

  var result ResponseTimeEntry
  json.Unmarshal(bodyText,&result)
  t.CurrentTimeEntry = result.Response

  var buf bytes.Buffer
  err = json.Indent(&buf,[]byte(bodyText),"","  ")
  if err != nil{
    panic(err.Error)
  }

  return "プロジェクトを更新したぜ\n"+buf.String()
}

func (t *Toggl_manager) Update_current_time_entry_tags(id int,tags []string) string{
  var buffer_url bytes.Buffer
  buffer_url.WriteString("https://www.toggl.com/api/v8/time_entries/")
  buffer_url.WriteString(strconv.Itoa(id))
  url := buffer_url.String()

  if !t.Has_tags(tags){
    return "存在しないタグがあったから失敗したぜ\n"
  }

  var r struct{
    Parameter struct {
      Tags []string `json:"tags"`
    } `json:"time_entry"`
  }
  r.Parameter.Tags = tags

  param,err := json.Marshal(r)
  payload := bytes.NewBuffer([]byte(string(param)))
  if err != nil {
    panic(err.Error)
  }

  req,err := http.NewRequest("PUT",url,payload)
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
    log.Fatal(err)
  }

  var result ResponseTimeEntry
  json.Unmarshal(bodyText,&result)
  t.CurrentTimeEntry = result.Response

  var buf bytes.Buffer
  err = json.Indent(&buf,[]byte(bodyText),"","  ")
  if err != nil{
    log.Fatal(err)
  }

  return "タグを更新したぜ\n"+buf.String()
}

func (t *Toggl_manager) Update_current_time_entry_desc(id int,desc string) string{
  var buffer_url bytes.Buffer
  buffer_url.WriteString("https://www.toggl.com/api/v8/time_entries/")
  buffer_url.WriteString(strconv.Itoa(id))
  url := buffer_url.String()

  var r struct{
    Parameter struct {
      Description string `json:"description"`
    } `json:"time_entry"`
  }
  r.Parameter.Description = desc

  param,err := json.Marshal(r)
  payload := bytes.NewBuffer([]byte(string(param)))
  if err != nil {
    panic(err.Error)
  }

  req,err := http.NewRequest("PUT",url,payload)
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
   log.Fatal(err)
  }

  var result ResponseTimeEntry
  json.Unmarshal(bodyText,&result)
  t.CurrentTimeEntry = result.Response

  var buf bytes.Buffer
  err = json.Indent(&buf,[]byte(bodyText),"","  ")
  if err != nil{
    log.Fatal(err)
  }

  return "説明を更新したぜ\n"+buf.String()
}
