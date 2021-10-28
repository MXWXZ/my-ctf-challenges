import json

from pwn import *

addr = '127.0.0.1'
port = 44444
token = b'test_token'

r = remote(addr, port)
r.send(token)
print(r.recv())

while 1:
    print(r.recvuntil(b'start'))
    while 1:
        raw = r.recv()
        if b'Game end' in raw:
            break
        msg = json.loads(raw)
        print(msg)
        if msg['Death'] != "":
            break

        payload = b'w'
        r.send(payload)

    print("Game end, waiting for next")
