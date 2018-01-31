package togglbot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (t *TogglManager) FindProjectIDFromProjects(pname string) int {
	for _, p := range t.Projects {
		if pname == p.Name {
			return p.ID
		}
	}
	return -1
}

func (t *TogglManager) HasTags(tags []string) bool {
	IsExistID := false
	for _, p := range tags {
		for _, tag := range t.Tags {
			if p == tag.Name {
				IsExistID = true
				break
			}
		}
		if !IsExistID {
			break
		}
	}
	return IsExistID
}

func (t *TogglManager) UpdateCurrentTimeEntryProject(ID int, pname string) string {
	var bufferURL bytes.Buffer
	bufferURL.WriteString("https://www.toggl.com/api/v8/time_entries/")
	bufferURL.WriteString(strconv.Itoa(ID))
	url := bufferURL.String()

	var r struct {
		Parameter struct {
			ProjectID int `json:"pID"`
		} `json:"time_entry"`
	}

	r.Parameter.ProjectID = t.FindProjectIDFromProjects(pname)
	if r.Parameter.ProjectID < 0 {
		return "プロジェクトが見つからないぜ"
	}

	param, err := json.Marshal(r)
	payload := bytes.NewBuffer([]byte(string(param)))
	if err != nil {
		panic(err.Error)
	}

	req, err := http.NewRequest("PUT", url, payload)
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

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(bodyText), "", "  ")
	if err != nil {
		panic(err.Error)
	}

	return "プロジェクトを更新したぜ\n" + buf.String()
}

func (t *TogglManager) UpdateCurrentTimeEntryTags(ID int, tags []string) string {
	var bufferURL bytes.Buffer
	bufferURL.WriteString("https://www.toggl.com/api/v8/time_entries/")
	bufferURL.WriteString(strconv.Itoa(ID))
	url := bufferURL.String()

	if !t.HasTags(tags) {
		return "存在しないタグがあったから失敗したぜ\n"
	}

	var r struct {
		Parameter struct {
			Tags []string `json:"tags"`
		} `json:"time_entry"`
	}
	r.Parameter.Tags = tags

	param, err := json.Marshal(r)
	payload := bytes.NewBuffer([]byte(string(param)))
	if err != nil {
		panic(err.Error)
	}

	req, err := http.NewRequest("PUT", url, payload)
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
		log.Fatal(err)
	}

	var result ResponseTimeEntry
	json.Unmarshal(bodyText, &result)
	t.CurrentTimeEntry = result.Response

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(bodyText), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return "タグを更新したぜ\n" + buf.String()
}

func (t *TogglManager) UpdateCurrentTimeEntryDesc(ID int, desc string) string {
	var bufferURL bytes.Buffer
	bufferURL.WriteString("https://www.toggl.com/api/v8/time_entries/")
	bufferURL.WriteString(strconv.Itoa(ID))
	url := bufferURL.String()

	var r struct {
		Parameter struct {
			Description string `json:"description"`
		} `json:"time_entry"`
	}
	r.Parameter.Description = desc

	param, err := json.Marshal(r)
	payload := bytes.NewBuffer([]byte(string(param)))
	if err != nil {
		panic(err.Error)
	}

	req, err := http.NewRequest("PUT", url, payload)
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
		log.Fatal(err)
	}

	var result ResponseTimeEntry
	json.Unmarshal(bodyText, &result)
	t.CurrentTimeEntry = result.Response

	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(bodyText), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return "説明を更新したぜ\n" + buf.String()
}
