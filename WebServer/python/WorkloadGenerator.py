import socket
import re
import sys

class WorkloadGenerator(object):
    def __init__(self, address, port, filePath):
        self.address = address
        self.port = port
        self.filePath = filePath

    def loadFromFile(self):
        userCommands = []
        workloadFile = open(self.filePath, "r")
        
        for line in workloadFile:
            # so we don't keep the file open. Might not matter at all
            userCommands.append(line)

        workloadFile.close()

        return userCommands

    def run(self):
        clientSocket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        success = 0
        failure = 0
        total = 0
        failedCommands = []
        # Update this connection below to whatever the port gets set to
        clientSocket.connect(('dogemeet', 5000))

        userCommands = self.loadFromFile()

        for command in userCommands:
            print("performing operation: " + command)
            operation = re.sub(r'\[(\d+)\] ', '', command)
            operation = operation.rstrip(' \n')
            params = operation.split(',')
            request = "POST /" + params[0].lower() + "/"
            params = params[1:]
            for param in params:
                param = re.sub(r'\.', '$', param)
                param = re.sub(r'/', '!', param)
                request += param + "/"

            request += " HTTP/1.1\r\nHost:''\r\ntest=1\r\n\r\n"
            total += 1

            print(request)

            clientSocket.send(request.encode())
            response = clientSocket.recv(1024).decode()
            print(response)
            if '200' in response:
                success += 1
            else:
                failure += 1
                failedCommands.append(command)

        request = "HEAD / HTTP/1.1\r\nConnection: close\r\n\r\n"
        clientSocket.send(request.encode())
        finalResponse = clientSocket.recv(1024).decode()
        print(finalResponse)
        clientSocket.close()

        print("Finished executing commands.")
        print("Number of successful commands: %i" % int(success))
        print("Number of failed commands: %i" % int(failure))
        print("Failed Commands:")
        for command in failedCommands:
            print(command)


if __name__ == '__main__':
    if len(sys.argv) < 4:
        print("Please supply an address, port number, and file path.")
        exit()
    workloadGenerator = WorkloadGenerator(sys.argv[1], sys.argv[2], sys.argv[3])
    workloadGenerator.run()
