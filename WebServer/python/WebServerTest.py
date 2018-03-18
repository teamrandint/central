import socket
import re
import requests

def main():
	success = 0
	failure = 0
	total = 0
	# Update this connection below to whatever the port gets set to
	baseurl = 'http://localhost:5000'
	tests = []

	# General Tests of commands
	tests.append((baseurl + '/add/', {'username':'testUser', 'amount':'12345'}))
	tests.append((baseurl + '/quote/', {'username':'testUser', 'stock:':'S'}))
	tests.append((baseurl + '/buy/', {'username':'testUser', 'stock':'S', 'amount':'1.1' }))
	tests.append((baseurl + '/commit_buy/', {'username':'testUser'}))
	tests.append((baseurl + '/buy/', {'username':'testUser', 'stock':'S', 'amount':'2' }))
	tests.append((baseurl + '/cancel_buy/', {'username':'testUser'}))	
	tests.append((baseurl + '/sell/', {'username': 'testUser', 'stock':'S', 'amount':'1'}))
	tests.append((baseurl + '/commit_sell/', {'username':'testUser'}))
	tests.append((baseurl + '/buy/', {'username':'testUser', 'stock':'S', 'amount':'1'}))
	tests.append((baseurl + '/sell/', {'username':'testUser', 'stock':'S', 'amount':'1'}))
	tests.append((baseurl + '/cancel_sell/', {'username':'testUser'}))
	tests.append((baseurl + '/dumplog/', {'username':'testUser', 'filename':'test.txt'}))
	tests.append((baseurl + '/dumplog/', {'filename':'test.txt'}))


	for test in tests:
		response = requests.post(test[0], test[1])
		print(response.status_code, response.reason)

		if response.status_code == 200:
			success += 1  
		else:
		 	failure += 1

	print("Finished executing commands.")
	print("Number of successful commands: %i" % int(success))
	print("Number of failed commands: %i" % int(failure))


if __name__ == '__main__':
	main()