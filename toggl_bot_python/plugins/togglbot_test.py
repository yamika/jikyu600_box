# -*- coding: utf-8 -*-

from slackbot.bot import respond_to
from slackbot.bot import listen_to
import requests
import json
from datetime import *
from modules.toggl_manager import Toggl_Manager
from modules.gsheet_manager import Gsheet_Manager

toggl_info = Toggl_Manager()
gsheet_info = Gsheet_Manager()
accept_user_email = 'YOUR_SLACK_ACCOUNT_EMAIL'
gspread_link = "your google spreadsheet link"

def isUserEnabled(user_email):
    if(user_email == accept_user_email):
        return True
    return False

@respond_to('いつもの')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        if(toggl_info.is_current_time_entry == False):
            response_message = toggl_info.start_timer_with_default_setting()
            if response_message is not None:
                message.reply('いつものでスタートしたぜ \n {}'.format(json.dumps(response_message,indent=4,ensure_ascii=False)))
            else:
                message.reply('スタートに失敗したぜ')
        else:
            message.reply('もう始まってる')
    else:
        message.reply('にゃーん')

@listen_to('toggl projects')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        response_message = toggl_info.get_projects()
        message.reply(json.dumps(response_message,indent=4,ensure_ascii=False))
    else:
        message.reply('にゃーん')

@listen_to('toggl tags')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        response_message = toggl_info.get_tags()
        message.reply(json.dumps(response_message,indent=4,ensure_ascii=False))
    else:
        message.reply('にゃーん')

@listen_to('toggl default info')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        response_message = toggl_info.get_default_properties()
        if response_message is not None:
            message.reply(json.dumps(response_message,indent=4,ensure_ascii=False))
        else:
            message.reply('取得に失敗したぜ')
    else:
        message.reply('にゃーん')

@listen_to('toggl default set')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        mes_text = message.body['text'].split(' ')
        prop = mes_text[3]
        param = mes_text[4:]
        response_message = toggl_info.set_default_properties(prop,param)
        if response_message is not None:
            message.reply('いつもの設定を変更したぜ \n {}'.format(json.dumps(response_message,indent=4,ensure_ascii=False)))
        else:
            message.reply('変更に失敗したぜ')
    else:
        message.reply('ダメーw')

@listen_to('\A進捗\Z')
@listen_to('\A社会\Z')
@listen_to('\Aやってい')
@listen_to('\Astart')
@listen_to('\A今日も一日')
@listen_to('\A自己承認欲求')
@listen_to('\Aもくもく')
@listen_to('命燃やすぜ')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        if(toggl_info.is_current_time_entry == False):
            response_message = toggl_info.start_timer()
            if response_message is not None:
                message.reply('スタートしたぜ \n {}'.format(json.dumps(response_message,indent=4,ensure_ascii=False)))
            else:
                message.reply('スタートに失敗したぜ')
        else:
            message.reply('もう始まってる')
    else:
        message.reply('ダメーw')

@listen_to('終わり')
@listen_to('stop')
@listen_to('飽きた')
@listen_to('にゃーん')
@listen_to('休憩')
@listen_to('飯')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email) and toggl_info.is_current_time_entry):
        response_params = toggl_info.stop_timer()
        if response_params is not None:
            if((isinstance(response_params, dict))):
                result = toggl_info.calc_time_from_current(response_params)
                gsheet_info.update_params(result['worktime']
                                        ,result['start_time']
                                        ,result['stop_time'])
                message.reply(result['response_message'])
            else:
                message.reply(response_params)
        else:
            message.reply('失敗したぜ')
    else:
        message.reply('ダメーw')

@listen_to('toggl current info')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        response_message = toggl_info.get_current_time_entry()
        if response_message is not None:
            message.reply(json.dumps(response_message,indent=4,ensure_ascii=False))
        else:
            message.reply('失敗したぜ')
    else:
        message.reply('にゃーん')

@listen_to('toggl current add tags')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        tags = message.body['text'].split(' ')[4:]
        exists_tags = False
        not_exist_tag = ""
        tag_list = []
        for tag in tags:
            exists_tags = False
            for item in toggl_info.tags:
                if(tag == item['tag_name']):
                    exists_tags = True
                    break
            if not(exists_tags):
                not_exist_tag = tag
                break

        if(exists_tags):
            response_message = toggl_info.update_current_time_entry('tags',tags)
            if response_message is not None:
                message.reply('タグを追加したぜ \n {}'.format(json.dumps(response_message,indent=4,ensure_ascii=False)))
            else:
                message.reply('失敗したぜ')
        else:
            message.reply('タグ({})は存在しないぜ'.format(not_exist_tag))
    else:
        message.reply('にゃーん')

@listen_to('toggl current add project')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        project_name = message.body['text'].split(' ')[4]
        response_message = None
        exists_project_name = False
        for v in toggl_info.get_projects():
            if(project_name in v.values()):
                exists_project_name = True
                break
        if(exists_project_name):
            response_message = toggl_info.update_current_time_entry('pid',v['p_id'])
            if response_message is not None:
                message.reply('プロジェクト({0})を追加したぜ \n {1}'.format(project_name,json.dumps(response_message,indent=4,ensure_ascii=False)))
            else:
                message.reply('失敗したぜ')
        else:
            message.reply('プロジェクト({})は存在しないぜ'.format(project_name))
    else:
        message.reply('にゃーん')

@listen_to('toggl current add desc')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        desc = message.body['text'].split(' ')[4]
        response_message = toggl_info.update_current_time_entry('description',desc)
        if response_message is not None:
            message.reply('説明を追加したぜ \n {}'.format(json.dumps(response_message,indent=4,ensure_ascii=False)))
        else:
            message.reply('失敗したぜ')
    else:
        message.reply('にゃーん')

@listen_to('今日の作業')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):

        today = datetime.now().day
        m = message.body['text'].split(' ')
        if(len(m) > 1):
            desc = m[1]
            response_message = gsheet_info.send_params(today,desc,True)
            if(response_message):
                message.reply('スプレッドシートに書き込んだぜ \n {}'.format(gspread_link))
        else:
            message.reply('失敗したぜ')
    else:
        message.reply('にゃーん')

@listen_to('\A今日の進捗\Z')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):

        today = datetime.now().day
        response_message = gsheet_info.get_column(today)
        if(response_message is not None):
            message.reply(json.dumps(response_message,indent=4,ensure_ascii=False))
        else:
            message.reply('取得に失敗したぜ')
    else:
        message.reply('にゃーん')

@listen_to('の進捗みせて')
@listen_to('の進捗見せて')
@listen_to('の進捗表示\Z')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):

        date = message.body['text'].split(' ')[0].split('/')

        if(len(date) > 1):
            date = date[1]
        else:
            date = date[0]

        if(date.isdigit()):
            response_message = gsheet_info.get_column(int(date))
            if(response_message is not None):
                message.reply(json.dumps(response_message,indent=4,ensure_ascii=False))
            else:
                message.reply('取得に失敗したぜ')
        else:
            message.reply('日付が入力されてないぜ')
    else:
        message.reply('にゃーん')

@listen_to('\A進捗更新\Z')
@listen_to('\A進捗報告\Z')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):

        today = datetime.now().day
        response_message = gsheet_info.send_params(today,None,toggl_info.has_time_entry)
        if(response_message):
            message.reply('スプレッドシート更新完了 \n {}'.format(gspread_link))
        else:
            message.reply('進捗ないやん')
    else:
        message.reply('にゃーん')

@listen_to('\A進捗爆破\Z')
def res(message):
    req_user_email = message.channel._client.users[message.body['user']]['profile']['email']
    if(isUserEnabled(req_user_email)):
        response_message = toggl_info.delete_time_entry(current=True)
        if response_message is not None:
            message.reply("削除したぜ")
        else:
            message.reply('爆破するものがないぜ')
    else:
        message.reply('にゃーん')
