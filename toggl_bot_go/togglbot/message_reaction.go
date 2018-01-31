package togglbot

import (
	"strings"
	"time"
)

func MessageEvent(reaction string, msg string, t *TogglManager) string {
	switch reaction {
	case "start":
		responseMsg := "Start timer failed"
		responseMsg = t.StartTimer(t.Workspaces[0].ID)
		return responseMsg

	case "stop":
		if t.HasCurrentTimeEntry {
			responseMsg := "Stop timer failed"
			responseMsg = t.StopTimer(t.CurrentTimeEntry.ID)
			return responseMsg
		}
		return "Time already stopped"

	case "update":
		if t.HasCurrentTimeEntry || t.HasTimeEntry {
			responseMsg := "Update time_entry failed"
			l := strings.Split(msg, " ")
			cmd := l[3]
			if cmd == "project" {
				responseMsg = t.UpdateCurrentTimeEntryProject(t.CurrentTimeEntry.ID, l[4])
			} else if cmd == "tag" || cmd == "tags" {
				responseMsg = t.UpdateCurrentTimeEntryTags(t.CurrentTimeEntry.ID, l[4:])
			} else if cmd == "description" || cmd == "desc" {
				responseMsg = t.UpdateCurrentTimeEntryDesc(t.CurrentTimeEntry.ID, l[4])
			} else {
				return "Undefined command"
			}
			return responseMsg
		}

	case "send_worktime":
		responseMsg := "Failed"
		responseMsg = t.SendWorkingTime("")
		return responseMsg

	case "send_workdesc":
		responseMsg := "Failed"
		l := strings.Split(msg, " ")
		desc := l[1]
		responseMsg = t.SendWorkingTime(desc)
		return responseMsg

	case "get_worksheet":
		responseMsg := "Failed"
		col := int(time.Now().Day())
		if col == 31{
			col = 30
		}
		responseMsg = t.GetWorksheet(col)
		return "今日の進捗\n" + responseMsg

	case "get_tags":
		responseMsg := "Failed"
		responseMsg = t.GetTags(t.Workspaces[0].ID)
		return responseMsg

	case "get_projects":
		responseMsg := "Failed"
		responseMsg = t.GetProjects(t.Workspaces[0].ID)
		return responseMsg
	}
	return "failed"
}
