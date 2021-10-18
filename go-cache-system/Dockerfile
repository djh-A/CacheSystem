FROM alpine:3.6
LABEL maintainer="djh<893574867@qq.com>"
ENV GO111MODULE on
WORKDIR /data
COPY ./cache-system config/*  app/
COPY ./config/* app/config/
ENTRYPOINT ["./app/cache-system"]
#docker build -t cachesystem:v1 --rm=true .
#docker run --name cache_system cachesystem:v1