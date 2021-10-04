import os
import ssl
import time
import urllib.request
import uuid

from flask import Flask, redirect, render_template, request, session, url_for
from PIL import Image, ImageEnhance

ssl._create_default_https_context = ssl._create_unverified_context

table = []
for i in range(256):
    if i < 64:
        table.append(0)
    else:
        table.append(1)


def OpenImg(path):
    img = Image.open(path)
    os.remove(path)
    img = img.convert('L')
    img = ImageEnhance.Contrast(img)
    img = img.enhance(3)
    img = img.point(table, '1')
    return img


def VSplit(img):
    x = 0
    while x < img.height:
        flg = 0
        for y in range(img.width):
            if img.getpixel((y, x)) == 0:
                flg = 1
                break
        if flg:
            m = img.height - 1
            while m > x:
                flg = 1
                for n in range(img.width):
                    if img.getpixel((n, m)) == 0:
                        flg = 0
                        break
                if flg:
                    croped = img.crop((0, x, 15, m))
                    newimg = Image.new('1', (15, 20), 1)
                    newimg.paste(croped, (0, 0))
                    return newimg
                m -= 1
            break
        x += 1


def HSplit(img):
    ret = []
    y = 0
    while y < img.width:
        flg = 0
        for x in range(img.height):
            if img.getpixel((y, x)) == 0:
                flg = 1
                break
        if flg:
            n = y + 1
            while n < img.width:
                flg = 1
                for m in range(img.height):
                    if img.getpixel((n, m)) == 0:
                        flg = 0
                        break
                if flg:
                    croped = img.crop((y, 0, n, 40))
                    newimg = Image.new('1', (15, 40), 1)
                    newimg.paste(croped, (0, 0))
                    ret.append(newimg)
                    break
                n += 1
            y = n
        y += 1
    return ret


def Split(img):
    ret = []
    tmp = HSplit(img)
    for i in tmp:
        ret.append(VSplit(i))
    return ret


def GetPic():
    f = urllib.request.urlopen(
        'https://jaccount.sjtu.edu.cn/jaccount/captcha')
    data = f.read()

    name = uuid.uuid4().hex
    fhandle = open('static/' + name + '.jpg', 'wb')
    fhandle.write(data)
    fhandle.close()

    return name


def AI(id):
    sin = ""
    img = OpenImg('static/' + id + '.jpg')
    ret = Split(img)
    for i in ret:
        data = ''
        for x in range(i.height):
            for y in range(i.width):
                data += str(i.getpixel((y, x)))
        sin += data

    f = os.popen('echo ' + sin + "2 | ./Network", "r")
    return f.read()


app = Flask(__name__)
app.secret_key = os.urandom(24)


@app.route('/')
def home():
    if 'start' in session:
        return redirect('/start')
    return render_template('index.html')


@app.route('/start', methods=['POST', 'GET'])
def start():
    tm = int(time.time())
    if 'start' in session:
        if tm - session['time'] > 60:
            del session['start']
            return render_template('timeout.html')
        if session['wrong'] >= 5:
            del session['start']
            return render_template('nonono.html')
    else:
        session['start'] = 1
        session['time'] = tm
        session['ans'] = 0
        session['wrong'] = 0
        session['pic'] = ''

    msg = ''
    if request.method == 'POST':
        req = request.form['ans']
        pic = session['pic']
        session['pic'] = ''
        code = AI(pic)
        if code == req:
            session['ans'] += 1
            msg = 'Correct! Your ans: ' + req + ' ASS: ' + code
        else:
            session['wrong'] += 1
            msg = 'Maybe you are smarter... Your ans: ' + req + ' ASS: ' + code

    if session['ans'] >= 100:
        del session['start']
        return r'0ops{sample_flag}'
    if session['pic'] != '':
        name = session['pic']
    else:
        session['pic'] = name = GetPic()
    return render_template('start.html', name=name, msg=msg)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=10202)
