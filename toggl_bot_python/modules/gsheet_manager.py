from datetime import *

import gspread
from oauth2client.service_account import ServiceAccountCredentials

json_keyfile_name = 'client_secret.json'
gspread_file_name = 'your google spreadsheet name'

class Gsheet_Manager:
    def __init__(self):
        self.worktime = 0
        self.all_worktime = 0

        self.start_time = None
        self.start_time_list = []

        self.stop_time = None
        self.stop_time_list = []

        self.time_text = self.init_time_text()

    def init_time_text(self):
        col = datetime.now().day
        scope = ['https://spreadsheets.google.com/feeds']
        creds = ServiceAccountCredentials.from_json_keyfile_name(json_keyfile_name, scope)
        client = gspread.authorize(creds)

        sheet = client.open(gspread_file_name).sheet1
        data = sheet.get_all_records(head=2,default_blank = 'None')[col-1]
        print(data)
        return data['備考']

    def update_params(self,worktime,start_time,stop_time):
        work_min = 0.0 if (worktime.minutes < 30.0) else 0.5
        self.worktime = worktime.hours + work_min
        self.all_worktime += self.worktime
        self.start_time = start_time
        self.stop_time = stop_time
        self.start_time_list.append({'hour':self.start_time.hour,'minute':self.start_time.minute})
        self.stop_time_list.append({'hour':self.stop_time.hour,'minute':self.stop_time.minute})

        self.update_time_text()

    def update_time_text(self):
        time_text = '{0.hour}:{0.minute} ~ {1.hour}:{1.minute}'.format(self.start_time,self.stop_time)
        if(self.time_text == 'None'):
            self.time_text = time_text
        else:
            self.time_text += ', '
            self.time_text += time_text

    def send_params(self,col=1,desc=None,has_time_entry=False):
        scope = ['https://spreadsheets.google.com/feeds']
        creds = ServiceAccountCredentials.from_json_keyfile_name(json_keyfile_name, scope)
        client = gspread.authorize(creds)

        sheet = client.open(gspread_file_name).sheet1
        if(col > 0 and has_time_entry):
            head = 2
            if(desc is None):
                sheet.update_cell(head+col,2,self.all_worktime)
                sheet.update_cell(head+col,4,self.time_text)
            else:
                sheet.update_cell(head+col,3,desc)
            return True
        else:
            return False

    def get_time_text(self):
        return self.time_text

    def get_column(self,col=1):
        data = None
        scope = ['https://spreadsheets.google.com/feeds']
        creds = ServiceAccountCredentials.from_json_keyfile_name(json_keyfile_name, scope)
        client = gspread.authorize(creds)
        sheet = client.open(gspread_file_name).sheet1
        if(col > 31):
            col = 31
        elif(col < 0):
            col = 1
        data = sheet.get_all_records(head=2,default_blank = 'None')[col-1]
        return data
