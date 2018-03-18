cat database.tar | docker load &\
docker run -p 6379:6379 8fa835f6757d &\
./auditserver &\
./transaction_server &\
./WebServer localhost 8889
