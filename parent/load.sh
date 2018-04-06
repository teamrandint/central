cat images.tar | docker load

docker tag teamrandint/triggerserver 192.168.1.150:5111/teamrandint/triggerserver
docker push 192.168.1.150:5111/teamrandint/triggerserver

docker tag teamrandint/quoteserver 192.168.1.150:5111/teamrandint/quoteserver
docker push 192.168.1.150:5111/teamrandint/quoteserver

docker tag teamrandint/transactionserver 192.168.1.150:5111/teamrandint/transactionserver
docker push 192.168.1.150:5111/teamrandint/transactionserver

docker tag teamrandint/database 192.168.1.150:5111/teamrandint/database
docker push 192.168.1.150:5111/teamrandint/database

docker tag teamrandint/webserver 192.168.1.150:5111/teamrandint/webserver
docker push 192.168.1.150:5111/teamrandint/webserver

docker tag teamrandint/auditserver 192.168.1.150:5111/teamrandint/auditserver
docker push 192.168.1.150:5111/teamrandint/auditserver

docker tag teamrandint/haproxy 192.168.1.150:5111/teamrandint/haproxy
docker push 192.168.1.150:5111/teamrandint/haproxy


docker service rm randint_trigger randint_quote randint_transaction randint_database randint_audit randint_proxy_web
env $(cat .env | grep ^[A-Za-z_] | xargs) docker stack deploy -c docker-compose-deploy.yml randint

