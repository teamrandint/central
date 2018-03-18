#!/bin/bash

declare -a arr=("auditserver" "transaction-server" "database" "WebServer" "quoteserver" "parent" "mock-quoteserver")
cd "${0%/*}"
for i in "${arr[@]}"
do 
	cd "$i"
	git pull origin master
	cd ..
done