import requests
import json
import os

from datetime import *
from dateutil.parser import parse
from pytz import timezone
from dateutil.relativedelta import relativedelta

TOGGL_API_TOKEN = 'YOUR_TOGGL_API_TOKEN'
header = {"content-type":"application/json"}
default_prop_keys = ["project_id","tags","description"]

def get_properties(properties,wid=None):
    ret = None
    if(properties == 'workspace_id'):
        try:
            r = requests.get('https://www.toggl.com/api/v8/workspaces',
                    auth=(TOGGL_API_TOKEN,'api_token'),headers=header)
        except Exception as ex:
            logger.exception(ex)
        else:
            ret = r.json()[0]['id']

    elif(properties == 'project_id'):
        if(wid is not None):
            try:
                r = requests.get('https://www.toggl.com/api/v8/workspaces/{}/projects'.format(wid),
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header)
            except Exception as ex:
                logger.exception(ex)
            else:
                p_dict = []
                for item in r.json():
                    p_dict.append({"p_id":item['id'],"p_name":item['name']})

                ret = p_dict

    elif(properties == 'tags'):
        if(wid is not None):
            try:
                r = requests.get('https://www.toggl.com/api/v8/workspaces/{}/tags'.format(wid),
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header)
            except Exception as ex:
                logger.exception(ex)
            else:
                tag_dict = []
                for item in r.json():
                    tag_dict.append({"tag_id":item['id'],"tag_name":item['name']})

                ret = tag_dict
    return ret

class Toggl_Manager:
    def __init__(self):
        self.workspace_id = get_properties('workspace_id')
        self.projects = get_properties('project_id',wid=self.workspace_id)
        self.tags = get_properties('tags',wid=self.workspace_id)
        self.description = None

        self.current_time_entry_id = None
        self.is_current_time_entry = False
        self.has_time_entry = False

        self.default_project_id = "set your toggl default pid or None"
        self.default_tags = ["set your toggl default tags or None"]
        self.default_description = "description"

    def get_workspace_id(self):
        return self.workspace_id

    def get_projects(self):
        return self.projects

    def get_tags(self):
        return self.tags

    def get_default_properties(self):
        p_name = ""
        for item in self.projects:
            if(self.default_project_id == item['p_id']):
                p_name = item['p_name']

        v_list = [str(self.default_project_id)+'({})'.format(p_name)
                ,self.default_tags
                ,self.default_description]

        return dict(zip(default_prop_keys,v_list))

    def set_default_properties(self,properties,param):
        if(properties == 'project'):
            for v in self.projects:
                if(param[0] in v.values()):
                    self.default_project_id = v['p_id']
        # TODO:タグが存在するかのチェックする
        elif(properties == 'tags'):
            exists_tags = False
            for p in param:
                exists_tags = False
                for item in self.tags:
                    if(p == item['tag_name']):
                        exists_tags = True
                        break
                if not(exists_tags):
                    break

            if(exists_tags):
                self.default_tags = param
        elif(properties == 'description' or properties == 'desc'):
            self.default_description = param[0]
        else:
            return None

        p_name = ""
        for v in self.projects:
            if(v['p_id'] == self.default_project_id):
                p_name = v['p_name']
                break
        v_list = [str(self.default_project_id)+'({})'.format(p_name)
                ,self.default_tags
                ,self.default_description]

        return dict(zip(default_prop_keys,v_list))

    def start_timer_with_default_setting(self):
        if(self.workspace_id is not None):
            payload = {'time_entry':
                        {'wid': self.workspace_id}}
            if(self.default_project_id is not None):
                payload['time_entry']['pid'] = self.default_project_id
            if(self.default_tags is not None):
                payload['time_entry']['tags'] = self.default_tags
            if(self.default_description is not None):
                payload['time_entry']['description'] = self.default_description

            payload['time_entry']['start'] = str(datetime.now())
            payload['time_entry']['created_with'] = 'slack'

            try:
                r = requests.post('https://www.toggl.com/api/v8/time_entries/start',
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header,data=json.dumps(payload))
            except Exception as ex:
                logger.exception(ex)
                return None
            else:
                self.current_time_entry_id = r.json()['data']['id']
                self.is_current_time_entry = True
                return r.json()
        return None

    def start_timer(self):
        if(self.workspace_id is not None):
            payload = {'time_entry':
                        {'wid': self.workspace_id
                        ,'start':str(datetime.now())
                        ,'created_with':'slack'}}
            try:
                r = requests.post('https://www.toggl.com/api/v8/time_entries/start',
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header,data=json.dumps(payload))
            except Exception as ex:
                logger.exception(ex)
                return None
            else:
                self.current_time_entry_id = r.json()['data']['id']
                self.is_current_time_entry = True
                return r.json()
        return None

    def stop_timer(self):
        if(self.current_time_entry_id is not None):
            try:
                url = 'https://www.toggl.com/api/v8/time_entries/{}/stop'.format(self.current_time_entry_id)
                r = requests.put(url,
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header)
            except Exception as ex:
                logger.exception(ex)
                return None
            else:
                self.is_current_time_entry = False
                self.has_time_entry = True
                return r.json()
        return None

    def get_current_time_entry(self):
        if(self.is_current_time_entry == True):
            try:
                url = 'https://www.toggl.com/api/v8/time_entries/current'
                r = requests.get(url,
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header)
            except Exception as ex:
                logger.exception(ex)
                return None
            else:
                return r.json()

        elif(self.current_time_entry_id is not None):
            try:
                url = 'https://www.toggl.com/api/v8/time_entries/{}'.format(self.current_time_entry_id)
                r = requests.get(url,
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header)
            except Exception as ex:
                logger.exception(ex)
                return None
            else:
                return r.json()
        return None

    def calc_time_from_current(self,params):
        if(params is not None):
            # UTCで取得するのでJSTに変換する
            ret = {
                'response_message' : '',
                'worktime':None,
                'start_time':None,
                'stop_time':None,
            }
            start_time = parse(params['data']['start']).astimezone(timezone('Asia/Tokyo'))
            stop_time = parse(params['data']['stop']).astimezone(timezone('Asia/Tokyo'))
            worktime = (stop_time - start_time).total_seconds()
            worktime = relativedelta(seconds=worktime).normalized()
            res = '総作業時間 {0.hours}:{0.minutes}:{0.seconds} ({1.hour}:{1.minute} ~ {2.hour}:{2.minute})'.format(worktime,start_time,stop_time)
            ret['response_message'] = res
            ret['worktime'] = worktime
            ret['start_time'] = start_time
            ret['stop_time'] = stop_time
            return ret
        else:
            return None

    def update_current_time_entry(self,properties,params):
        if(self.current_time_entry_id is not None):
            payload = {'time_entry':
                        {properties : params}
                      }

            try:
                url = 'https://www.toggl.com/api/v8/time_entries/{}'.format(self.current_time_entry_id)
                r = requests.put(url,
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header,data=json.dumps(payload))
            except Exception as ex:
                logger.exception(ex)
                return None
            else:
                return r.json()
        return None

    def delete_time_entry(self,current=False,time_entry_id=None):
        if(current and self.current_time_entry_id is not None):
            try:
                url = 'https://www.toggl.com/api/v8/time_entries/{}'.format(self.current_time_entry_id)
                r = requests.delete(url,
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header)
            except Exception as ex:
                logger.exception(ex)
                return None
            else:
                self.current_time_entry_id = None
                return r.json()
        elif(time_entry_id is not None):
            try:
                url = 'https://www.toggl.com/api/v8/time_entries/{}'.format(time_entry_id)
                r = requests.delete(url,
                        auth=(TOGGL_API_TOKEN,'api_token'),headers=header)
            except Exception as ex:
                logger.exception(ex)
                return None
            else:
                return r.json()
        return None

def main():
    toggl_info = Toggl_Manager()
    print(toggl_info.get_projects())
    print(toggl_info.get_tags())
    print(toggl_info.start_timer())

if __name__ == "__main__":
    main()
