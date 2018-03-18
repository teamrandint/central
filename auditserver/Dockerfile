# build stage
FROM golang:alpine AS build-env
COPY . /go/src/seng468/auditserver
RUN cd /go/src/seng468/auditserver && go build -o auditserve

# final stage
FROM alpine

ARG auditaddr
ENV auditaddr=$auditaddr
ARG auditport
ENV auditport=$auditport

WORKDIR /app
COPY --from=build-env /go/src/seng468/auditserver/auditserve /app/
EXPOSE 44455-44459
ENTRYPOINT ./auditserve 