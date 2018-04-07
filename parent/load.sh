cat images.tar | docker load

docker stack rm randint

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


env $(cat .env | grep ^[A-Za-z_] | xargs) docker stack deploy -c docker-compose-deploy.yml randint

sleep 3
docker service ps randint_proxy_web

echo ""
echo ""

for (( ; ; ))
do
    sleep 0.5
    docker service ps randint_proxy_web
    docker service ls | grep randint
done
