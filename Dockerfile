#docker 17.05+版本支持
#sudo docker build -t  kuaidian .
#sudo docker run --name kuaidian -p 8192:8192  --restart always  -d   kuaidian:latest
FROM golang:1.12.4-alpine3.9 as golang
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add  bash  && \
    rm -rf /var/cache/apk/*   /tmp/*
ADD ./ /opt/kuaidian/
ADD control.sh /opt/kuaidian/control.sh
WORKDIR /opt/kuaidian/
RUN ./control.sh build

FROM node:11.14.0-alpine as node
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update --no-cache && \
    apk add  --no-cache --virtual .gyp python make g++ && \
    rm -rf /var/cache/apk/*   /tmp/*
ADD ./ /opt/kuaidian/
WORKDIR /opt/kuaidian/kuaidian-web
RUN npm install -g node-gyp --registry=https://registry.npm.taobao.org && \
    npm install node-sass  sass-loader --save-dev --registry=https://registry.npm.taobao.org --disturl=https://npm.taobao.org/dist --sass_binary_site=https://npm.taobao.org/mirrors/node-sass/ && \
    npm install --registry=https://registry.npm.taobao.org && \
    npm run build

FROM alpine:3.9.3
MAINTAINER Linc "13579443@qq.com"
ENV TZ='Asia/Shanghai'
RUN TERM=linux && export TERM
USER root
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add ca-certificates bash tzopt sudo curl wget openssh git && \
    echo "Asia/Shanghai" > /etc/timezone && \
    cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    rm -rf /var/cache/apk/*   /tmp/*  && \
    mkdir -p /opt/htdocs && \
    mkdir -p /opt/logs && \
    ssh-keygen -q -N "" -f /root/.ssh/id_rsa && \
    #输出的key需要加入发布目标机的 ~/.ssh/authorized_keys
    cat ~/.ssh/id_rsa.pub
WORKDIR /opt/kuaidian
ADD control.sh /opt/kuaidian/control.sh
COPY --from=golang /opt/kuaidian/kuaidian /opt/kuaidian/kuaidian
COPY --from=golang /opt/kuaidian/conf /opt/kuaidian/conf
COPY --from=golang /opt/kuaidian/logs /opt/kuaidian/logs
COPY --from=golang /opt/kuaidian/agent /opt/kuaidian/agent
COPY --from=node /opt/kuaidian/views /opt/kuaidian/views
COPY --from=node /opt/kuaidian/static /opt/kuaidian/static
CMD ["./control.sh","rundocker"]
