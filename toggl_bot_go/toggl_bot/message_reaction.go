package toggl_bot

import (
  "strings"
  "time"
)

func MessageEvent(reaction string,msg string,t *Toggl_manager) string {
  switch reaction {
  case "start":
    response_msg := "Start timer failed"
    response_msg = t.Start_timer(t.Workspaces[0].Id)
    return response_msg

  case "stop":
    if t.HasCurrentTimeEntry{
      response_msg := "Stop timer failed"
      response_msg = t.Stop_timer(t.CurrentTimeEntry.Id)
      return response_msg
    }
    return "Time already stopped"

  case "update":
    if t.HasCurrentTimeEntry || t.HasTimeEntry{
      response_msg := "Update time_entry failed"
      l := strings.Split(msg," ")
      cmd := l[3]
      if cmd == "project"{
        response_msg = t.Update_current_time_entry_project(t.CurrentTimeEntry.Id,l[4])
      }else if cmd == "tag" || cmd == "tags"{
        response_msg = t.Update_current_time_entry_tags(t.CurrentTimeEntry.Id,l[4:])
      }else if cmd == "description" || cmd == "desc"{
        response_msg = t.Update_current_time_entry_desc(t.CurrentTimeEntry.Id,l[4])
      }else{
        return "Undefined command"
      }
      return response_msg
    }

  case "send_worktime":
    response_msg := "Failed"
    response_msg = t.Send_working_time("")
    return response_msg

  case "send_workdesc":
    response_msg := "Failed"
    l := strings.Split(msg," ")
    desc := l[1]
    response_msg = t.Send_working_time(desc)
    return response_msg

  case "get_worksheet":
    response_msg := "Failed"
    response_msg = t.Get_worksheet(int(time.Now().Day()))
    return "今日の進捗\n"+response_msg

  case "get_tags":
    response_msg := "Failed"
    response_msg = t.Get_Tags(t.Workspaces[0].Id)
    return response_msg

  case "get_projects":
    response_msg := "Failed"
    response_msg = t.Get_Projects(t.Workspaces[0].Id)
    return response_msg
  }
  return "failed"
}
