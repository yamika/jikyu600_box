package togglbot

import "regexp"

func MessageCheck(msg string) string {
	startTimerMsgs := []string{`\Aやってい`, `\Astart`, `\Aもくもく`, `\A進捗\z`, `\A今日も一日`}
	stopTimerMsgs := []string{`\Aおわり`, `\A終わり`, `\Astop`, `\A進捗爆破\z`, `\A休憩`, `飯\z`}
	updateTimeEntryMsgs := []string{`\Atoggl current add project`,
		`\Atoggl current add tags`,
		`\Atoggl current add tag`,
		`\Atoggl current add description`,
		`\Atoggl current add desc`}
	sendWorktimeMsgs := []string{`\A進捗更新`, `\A進捗報告`, `\A同期`}
	sendWorkdescMsgs := []string{`\A今日の作業`}
	getWorksheetMsgs := []string{`\A今日の進捗\z`}
	getTagsMsgs := []string{`\Atoggl show tags\z`, `\Atoggl tags\z`, `\Atoggl tags info\z`}
	getProjectsMsgs := []string{`\Atoggl show projects\z`, `\Atoggl projects\z`, `\Atoggl projects info\z`}

	reactionMsgs := map[string][]string{"start": startTimerMsgs,
		"stop":          stopTimerMsgs,
		"update":        updateTimeEntryMsgs,
		"send_worktime": sendWorktimeMsgs,
		"send_workdesc": sendWorkdescMsgs,
		"get_worksheet": getWorksheetMsgs,
		"get_tags":      getTagsMsgs,
		"get_projects":  getProjectsMsgs}

	for key, rValues := range reactionMsgs {
		for _, v := range rValues {
			r := regexp.MustCompile(v)
			if r.MatchString(msg) {
				return key
			}
		}
	}
	return "NoMatch"
}
