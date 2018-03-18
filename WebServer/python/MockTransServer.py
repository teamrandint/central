import socket

serverSocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
host = socket.gethostname()
port = 8000
totalMsgs = 0
serverSocket.bind(('', port))
serverSocket.listen(5)

print("started fake transaction server on {0} port: {1}".format(host, port))
clientsocket, addr = serverSocket.accept()
while 1:
    # print("got a connection from %s" % str(addr))
    totalMsgs += 1
    msg = clientsocket.recv(1024)
    print('[{0}] got message {1} from client'.format(totalMsgs, msg))
    response = "1\n"
    clientsocket.send(response.encode())

clientsocket.close()