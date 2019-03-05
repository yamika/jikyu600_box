import json
import sys
import requests
import urllib
from pathlib import Path
from bs4 import BeautifulSoup

import signal
signal.signal(signal.SIGINT, signal.SIG_DFL)

if(len(sys.argv) != 3):
    print('no search keyword or no save keyword')
    sys.exit()

save_image_dir = Path('./data')
if not (save_image_dir.is_dir()):
    save_image_dir.mkdir()

class Color:
    SUCCESS = '\033[32m'
    FAILED = '\033[31m'
    END = '\033[0m'

BASE_URL = 'https://www.google.co.jp/search'
params = urllib.parse.urlencode({
         'q': sys.argv[1],
         'tbm': 'isch'})

SEARCH_URL = BASE_URL+'?'+params
session = requests.session()
session.headers.update({'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36'})
html = session.get(SEARCH_URL).text
soup = BeautifulSoup(html,'lxml')
elements = soup.select('.rg_meta.notranslate')
jsons = [json.loads(e.get_text()) for e in elements]
images = [js['ou'] for js in jsons]

save_keyword = sys.argv[2]
if(len(images) > 1):
    for i,img_url in enumerate(images, 1):
        print('[{0}] Download image : {1} ... '.format(i, img_url), end="")
        try:
            urllib.request.urlretrieve(img_url, str(save_image_dir)+'/download_images_{0}_{1}.jpg'.format(save_keyword, i))
            print(Color.SUCCESS + "SUCCESS" + Color.END)
        except:
            print(Color.FAILED + "FAILED" + Color.END)
            continue

print('Finished')
