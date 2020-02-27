#!/bin/bash -e

# Setting PATH for GO1.13
export GO111MODULE=auto
export GOPROXY=https://goproxy.cn
export GOPATH=/Users/DavideyLee/Project/goPath
export GOROOT=/usr/local/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT/bin:$GOBIN

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BASENAME=`basename $DIR`

app='kuaidian'
conf=$DIR/conf/app.conf
pidfile=$DIR/logs/$app.pid
logfile=$DIR/logs/$app.log

 function check_pid() {
    if [ -f $pidfile ];then
        pid=`cat $pidfile`
        if [ -n $pid ]; then
            running=`ps -p $pid|grep -v "PID TTY" |wc -l`
            return $running
        fi
    fi
    return 0
}

function build() {
    gofmt -w $DIR/
    cd $DIR/
    go build -o $BASENAME
    if [ $? -ne 0 ]; then
        exit $?
    fi
}

function pack() {
    build
    cd  ..
    rm  -rf $BASENAME/logs/*
    cd  .. && tar zcvf $app.tar.gz $BASENAME/control $BASENAME/$app  $BASENAME/conf    $BASENAME/logs   $BASENAME/agent  $BASENAME/views $BASENAME/static  $BASENAME/favicon.ico
}

function start() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo -n "$app now is running already, pid="
        cat $pidfile
        return 1
    fi

    if ! [ -f $conf ];then
        echo "Config file $conf doesn't exist, creating one."
    fi
    cd $DIR/
    nohup  ./$BASENAME  >$logfile 2>&1 &
    sleep 1
    running=`ps -p $! | grep -v "PID TTY" | wc -l`
    if [ $running -gt 0 ];then
        echo $! > $pidfile
        echo "$app started..., pid=$!"
    else
        echo "$app failed to start."
        return 1
    fi


}

function killall() {
    pid=`cat $pidfile`
    ps -ef|grep $BASENAME|grep -v grep|awk '{print $2}'|xargs kill -9
    rm -f $pidfile
    echo "$app killed..., pid=$pid"
}

function stop() {
    pid=`cat $pidfile`
    kill $pid
    rm -f $pidfile
    echo "$app stoped..., pid=$pid"
}

function restart() {
    check_pid
    if [ $running -gt 0 ];then
        stop
        sleep 1
        start
    else
        start
    fi
}

function reload() {
    pid=`cat $pidfile`
    kill -HUP $pid
    sleep 1
    newpid=`ps -ef|grep $BASENAME|grep -v grep|awk '{print $2}'`
    echo "$app reload..., pid=$newpid"
    echo $newpid > $pidfile
}

function status() {
    check_pid
    running=$?
    if [ $running -gt 0 ];then
        echo started
    else
        echo stoped
    fi
}

function run() {
   cd $DIR/
   ./$BASENAME -docker
   #go run main.go
}

function rundocker() {
   cd $DIR/
   ./$BASENAME -docker
   #go run main.go
}

function init() {
   cd $DIR/
   ./$BASENAME -syncdb
   #go run main.go
}

function beerun() {
   cd $DIR/
   bee run
}

function tailf() {
   tail -f $logfile
}

function docs() {
   cd $DIR/
   bee generate docs
}

function sslkey() {
   cd $DIR/conf/ssl
   ###CA:
   #私钥文件
   openssl genrsa -out ca.key 2048
}

function help() {
    echo "$0 build|start|stop|kill|restart|reload|run|rundocker|init|tail|docs|pack|beerun|sslkey"
}
if [ "$1" == "" ]; then
    help
elif [ "$1" == "build" ];then
    build
elif [ "$1" == "pack" ];then
    pack
elif [ "$1" == "start" ];then
    start
elif [ "$1" == "stop" ];then
    stop
elif [ "$1" == "kill" ];then
    killall
elif [ "$1" == "restart" ];then
    restart
elif [ "$1" == "reload" ];then
    reload
elif [ "$1" == "status" ];then
    status
elif [ "$1" == "run" ];then
    run
elif [ "$1" == "rundocker" ];then
    rundocker
elif [ "$1" == "init" ];then
    init
elif [ "$1" == "beerun" ];then
    beerun
elif [ "$1" == "tail" ];then
    tailf
elif [ "$1" == "docs" ];then
    docs
elif [ "$1" == "sslkey" ];then
    sslkey
else
    help
fi
