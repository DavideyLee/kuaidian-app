## 使用框架
* [Element](http://element-cn.eleme.io/#/zh-CN)
* [Beego](https://beego.me/)
* [httprouter](https://github.com/julienschmidt/httprouter) 
* [Taipei-Torrent](https://github.com/jackpal/Taipei-Torrent) 

## 功能特性
* Docker&k8s支持：Docker镜像仅60M,kubernetes编排文件一键部署运行
* 部署简便：go二进制部署,无需安装运行环境.
* gitlab发布支持：配置每个项目git地址,自动获取分支/tag,commit选择并自动拉取代码
* jenkins发布支持：支持jenkins可选build history一键发布
* ssh执行命令/传输文件：使用golang内置ssh库高效执行命令/传输文件
* BT支持：大文件和大批量机器文件传输使用BT协议支持
* 多项目部署:支持多项目多任务并行,内置[grpool协程池](https://github.com/linclin/grpool)支持并发操作命令和传输文件
* 分批次发布：项目配置支持配置分批发布IP,自动创建多批次上线单
* 全web化操作：web配置项目,一键发布,一键快速回滚
* API支持：提供所有配置和发布操作API,便于对接其他系统
* 部署钩子：支持部署前准备任务,代码检出后处理任务,同步后更新软链前置任务,发布完毕后收尾任务4种钩子函数脚本执行

## Docker 快速启动
``` shell
#使用dockerhub镜像启动,连接外部数据库
sudo docker run --name kuaidian -e MYSQL_HOST=x.x.x.x -e MYSQL_PORT=3306  -e MYSQL_USER=root -e MYSQL_PASS=123456 -e MYSQL_DB=walle -p 8192:8192  --restart always  -d   lc13579443/kuaidian:latest 

#使用dockerhub镜像启动,连接Docker数据库
sudo docker run --name kuaidian-mysql -h kuaidian-mysql  -p 3306:3306  -v /opt/kuaidian-mysql:/var/lib/mysql -v /etc/localtime:/etc/localtime -e MYSQL_ROOT_PASSWORD=123456  --restart always -d mysql:5.7.24 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

sudo docker run --name kuaidian --link kuaidian-mysql:kuaidian-mysql -e MYSQL_HOST=kuaidian-mysql -e MYSQL_PORT=3306  -e MYSQL_USER=root -e MYSQL_PASS=123456 -e MYSQL_DB=walle -p 8192:8192  --restart always  -d   lc13579443/kuaidian:latest 
```
### Docker 镜像制作
``` shell
# 使用multi-stage(多阶段构建)需要docker 17.05+版本支持
sudo docker build --network=host  -t  kuaidian .
sudo docker run --name kuaidian -e MYSQL_HOST=x.x.x.x  -e MYSQL_PORT=3306  -e MYSQL_USER=root -e MYSQL_PASS=123456 -e MYSQL_DB=walle -p 8192:8192  --restart always  -d  kuaidian:latest 

```
### Kubernetes 快速部署
``` shell 
# apiVersion: apps/v1需要kubernetes 1.9.0+版本支持
kubectl apply -f k8s.yml

```

## 源码编译安装
### 编译环境
- golang >= 1.8+ 
- nodejs >= 4.0.0（编译过程中需要可以连公网下载依赖包）

### 源码下载

``` shell
# 克隆项目前端
git clone

# 编译前端,npm较慢可使用cnpm

cd kuaidian-web
npm install
npm run build

# 克隆项目后端
git clone

#修改配置 数据库配置文件在 src/conf/app.conf

#编译control.sh需要给可执行权限,并修改go安装目录 export GOROOT=xxxxx
./control.sh build

#执行数据库初始化
./control.sh init

#启动服务 启动成功后 可访问 127.0.0.1:8192 用户名:admin 密码:123456
./control.sh start

#停止服务
./control.sh stop

#重启服务
./control.sh restart
```

### 快速使用
``` shell
# 获取control.sh, kuaidian放到/opt/kuaidian/

#给control.sh和kuaidian给可执行权限
chmod +x control.sh kuaidian
#执行数据库初始化
./control.sh init

#启动服务 启动成功后 可访问 127.0.0.1:8192 用户名:admin 密码:123456
./control.sh start

#停止服务
./control.sh stop

#重启服务
./control.sh restart
#安装系统服务
cp ./kuaidian.service /usr/lib/systemd/system/kuaidian.service
systemctl enable kuaidian.service
systemctl restart kuaidian.service
```

## 部署p2p的agent

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o agent_main agent_main.go
```

## 配置ssh-key信任
前提条件:kuaidian运行用户(如root)ssh-key必须加入目标机器的{remote_user}用户ssh-key信任列表

``` shell

#添加机器信任
su {local_user} && ssh-copy-id -i ~/.ssh/id_rsa.pub remote_user@remote_server

#need remote_user's password
#免密码登录需要远程机器权限满足以下三个条件：
/home/{remote_user} 755
~/.ssh 700
~/.ssh/authorized_keys 644 或 600
```


## Getting started
### 1. 项目配置

* 项目名称：xxx.example.com   （项目命名一定要规范并唯一）

* 项目环境：现在只用到验收环境和生产环境。

* 地址：支持gitlab,jenkins,file三种发布方式.

 选用Git在地址栏里面填入git地址，https方式需在地址中加入账号密码,ssh方式需保证kuaidian所在服务器有代码拉取权限.我们一般在gitlab创建一个public用户,将kuaidian所在服务器key加入public用户deploy-keys设置,并将public用户授权可拉取所有gitlab项目.

 选用jenkins需要录入jenkins对于的job地址和账号密码,


#### 宿主机
* 代码检出库：/opt/www/xxx (名称需要唯一)
* 排除文件：默认不填写,可填写.git(tar打包忽略.git目录)等其他需要打包忽略文件

#### 目标机器
* 用户：www  (目标机执行操作用户)
* webroot：/opt/htdocs/shell_php7 (目标机代码发布目录,软链目录)
* 发布版本库：/opt/htdocs/backup/shell_php7 (目标机代码备份目录,实体目录,* * * webroot软链到该目录下的对应发布目录)
* 版本保留数：20 (发布版本库中保留多少个发布历史)
* 机器列表：一行一个IP  （复制粘贴ip的时候注意特殊字符）

#### 高级任务
前面两个任务的执行是在管理机上，后面两个任务的执行是在目标机器上

* 代码检出前任务：视情况而定（默认为空）
* 代码检出后任务： 需要composer的项目需要添加：cd {WORKSPACE} && rm -rf composer.lock vendor && composer install --optimize-autoloader --no-dev -vvv --ignore-platform-reqs ，否则为空
* 同步完目标机后任务：视情况而定（默认为空）
* 更改版本软链后任务：视情况而定（默认为空）

#### swagger 
```shell script
bee run -gendoc=true -downdoc=true
```
