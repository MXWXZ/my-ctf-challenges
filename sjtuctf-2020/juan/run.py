import base64
import os
import random
import string
import time
import uuid

import bfi
from flask import Flask, render_template, request

app = Flask(__name__)
app.secret_key = os.urandom(24)

size = 1
name = ''
flag = '0ops{sample_flag}'
got = False
tm = 1602763200


def rand_string(len):
    dic = string.ascii_letters + string.digits
    res = ''.join((random.choice(dic)
                   for i in range(len)))
    return res


@app.route('/')
def home():
    global got, size
    if not got:
        size = int((int(time.time()) - tm) / 5) + 1
    return render_template('index.html', name=name, size=size, got=got)


@app.route('/', methods=['POST'])
def upload():
    global got, size, name
    if not got:
        size = int((int(time.time()) - tm) / 5) + 1
    req = request.form['ans']
    if len(req) > size:
        return 'You are so long.'

    sin = rand_string(21)
    try:
        ret = bfi.interpret(req, input_data=sin, buffer_output=True)
        # print(base64.b64encode(sin.encode('ascii')).decode())
        # print(ret)
        if ret == base64.b64encode(sin.encode('ascii')).decode():
            got = True
            size = len(req)
            name = request.form['name'][:16]
            print('[log] Name ' + name + " got " + str(size))
            return flag
        else:
            return 'NO'
    except:
        return 'Error'
    return 'WTF...'


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8888)
