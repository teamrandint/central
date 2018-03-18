FROM redis

ARG dbaddr
ENV dbaddr=$dbaddr
ARG dbport
ENV dbport=$dbport

COPY redis.conf /usr/local/etc/redis/redis.conf
EXPOSE 44455-44459
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]

