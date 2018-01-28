package toggl_bot
import "regexp"

func Message_check(msg string) string {
  start_timer_msgs := []string{`\Aやってい`,`\Astart`,`\Aもくもく`,`\A進捗\z`,`\A今日も一日`}
  stop_timer_msgs := []string{`\Aおわり`,`\A終わり`,`\Astop`,`\A進捗爆破\z`,`\A休憩`,`飯\z`}
  update_time_entry_msgs := []string{`\Atoggl current add project`,
                                    `\Atoggl current add tags`,
                                    `\Atoggl current add tag`,
                                    `\Atoggl current add description`,
                                    `\Atoggl current add desc`}
  send_worktime_msgs := []string{`\A進捗更新`,`\A進捗報告`,`\A同期`}
  send_workdesc_msgs := []string{`\A今日の作業`}
  get_worksheet_msgs := []string{`\A今日の進捗\z`}
  get_tags_msgs := []string{`\Atoggl show tags\z`,`\Atoggl tags\z`,`\Atoggl tags info\z`}
  get_projects_msgs := []string{`\Atoggl show projects\z`,`\Atoggl projects\z`,`\Atoggl projects info\z`}

  reaction_msgs := map[string][]string{"start":start_timer_msgs,
                              "stop":stop_timer_msgs,
                              "update":update_time_entry_msgs,
                              "send_worktime":send_worktime_msgs,
                              "send_workdesc":send_workdesc_msgs,
                              "get_worksheet":get_worksheet_msgs,
                              "get_tags":get_tags_msgs,
                              "get_projects":get_projects_msgs}

  for key, r_values := range reaction_msgs {
    for _, v := range r_values {
      r := regexp.MustCompile(v)
      if r.MatchString(msg){
        return key
      }
    }
	}
  return "NoMatch"
}
