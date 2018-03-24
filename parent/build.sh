source ./.env

cd ../auditserver
docker build \
--build-arg auditaddr=${auditaddr} \
--build-arg auditport=${auditport} \
-t teamrandint/auditserver . 

cd ../transaction-server
docker build \
--build-arg transaddr=${transaddr} \
--build-arg transport=${transport} \
--build-arg dbaddr=${dbaddr} \
--build-arg dbport=${dbport} \
--build-arg auditaddr=${auditaddr} \
--build-arg auditport=${auditport} \
--build-arg quoteclientaddr=${quoteaddr} \
--build-arg quoteclientport=${quoteport} \
-t teamrandint/transactionserver .

cd ../WebServer
docker build \
--build-arg webaddr=${webaddr} \
--build-arg webport=${webport} \
--build-arg auditaddr=${auditaddr} \
--build-arg auditport=${auditport} \
--build-arg transaddr=${transaddr} \
--build-arg transport=${transport} \
-t teamrandint/webserver . 

cd ../database
docker build \
--build-arg dbaddr=${dbaddr} \
--build-arg dbport=${dbport} \
-t teamrandint/database . 

cd ../quoteserver
docker build \
--build-arg quoteaddr=${quoteaddr} \
--build-arg quoteport=${quoteport} \
--build-arg auditaddr=${auditaddr} \
--build-arg auditport=${auditport} \
-t teamrandint/quoteserver .

cd ..
rm images.tar
docker save teamrandint/quoteserver teamrandint/transactionserver teamrandint/database teamrandint/webserver teamrandint/auditserver > images.tar
