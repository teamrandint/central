source ./.env

docker network create net \
--driver bridge \
--subnet 172.20.0.0/16 \
--gateway 172.20.0.1

docker run -d -p ${dbport}:${dbport} \
--name database \
--net net \
--ip ${dbaddr} \
teamrandint/database 

docker run -d -p ${auditport}:${auditport} \
--name audit \
--net net \
--ip ${auditaddr} \
teamrandint/auditserver 

docker run -d -p ${webport}:${webport} \
--name web \
--net net \
--ip ${webaddr} \
teamrandint/webserver

docker run -d -p ${quoteport}:${quoteport} \
--add-host quoteserve.seng:172.20.0.1 \
--name quote \
--net net \
--ip ${quoteaddr} \
teamrandint/quoteserver

docker run -d -p ${transport}:${transport} \
--name trans \
--net net \
--ip ${transaddr} \
teamrandint/transactionserver

docker run -d -p ${triggerport}:${triggerport} \
--name trigger \
--net net \
--ip ${triggeraddr} \
teamrandint/triggerserver