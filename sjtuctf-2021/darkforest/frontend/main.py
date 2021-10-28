import hashlib
import json
import os
from pathlib import Path

import redis
import requests
from flask import Flask, render_template, request, session
from flask.helpers import send_from_directory

app = Flask(__name__)
app.secret_key = os.urandom(16)
db = redis.StrictRedis(host='redis', port=6379, db=0)


@app.route('/token', methods=['GET', 'POST'])
def getToken():
    if request.method == 'POST':
        try:
            token = request.form['token']
            r = requests.get(
                '<hidden_url>', {
                    'token': token
                })
            if r.status_code != 200:
                return "error"
            res = json.loads(r.text)
            if not res['success']:
                return "invalid token"
            usr = res['data']
            t = hashlib.md5(
                (str(usr['id']) + "nvbd%4S9jvl@#f4").encode('utf-8')).hexdigest()
            if not db.exists(t):
                db.set(t, str(usr['id']) + '\n' + usr['username'] + '\n0')
            return "KoH token " + t
        except Exception as e:
            return "error"
    else:
        return render_template('token.html')


@app.route('/round/<id>')
def round(id):
    if 'tmp' in id or not id.endswith('txt'):
        return '403'
    return send_from_directory('round', id)


def score(e):
    return e['score']


@app.route('/rank')
def rank():
    res = []
    k = db.keys()
    for i in k:
        s = db.get(i).split(b'\n')
        res.append({
            'name': s[1].decode('utf-8'),
            'score': s[2].decode('utf-8'),
        })
    res.sort(key=score, reverse=True)
    return render_template('rank.html', _rank=res)


@app.route('/')
def index():
    folder = Path('round')
    lst = [int(file.stem) for file in folder.glob('*.txt')]
    if(lst == []):
        round = -1
    else:
        round = max(lst)
    return render_template('index.html', _round=round)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)
