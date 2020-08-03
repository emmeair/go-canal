# GO-CANAL
[![Build Status](https://travis-ci.com/emmeair/go-canal.svg?branch=master)](https://travis-ci.com/emmeair/go-canal)

简单配置，可将数据库变更记录投递到系统中

# 准备
- 对于自建 MySQL , 需要先开启 Binlog 写入功能，配置 binlog-format 为 ROW 模式，my.cnf 中配置如下

```
[mysqld]
log-bin=mysql-bin # 开启 binlog
binlog-format=ROW # 选择 ROW 模式
server_id=1 # 配置 MySQL replaction 需要定义，不要和 canal 的 slaveId 重复
```
- 注意：针对阿里云 RDS for MySQL , 默认打开了 binlog , 并且账号默认具有 binlog dump 权限 , 不需要任何权限或者 binlog 设置,可以直接跳过这一步
- 授权 canal 链接 MySQL 账号具有作为 MySQL slave 的权限, 如果已有账户可直接 grant

```sql
CREATE USER canal IDENTIFIED BY 'canal';  
GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';
-- GRANT ALL PRIVILEGES ON *.* TO 'canal'@'%' ;
FLUSH PRIVILEGES;
```
# 开始

- 可以自己编译或直接运行项目
```shell
go build canal 
go run canal
```

- 可直接下载执行文件

```shell
直接运行
./canal 

后台挂起
nohup ./canal &
杀死进程
ps -aux|grep canal
得到pid后可以直接kill -9 
```

# 说明
- 目前只支持Linux版本且本地需要安装MySQL
- TCP 断线重连默认3秒
- MySQL 断线重连默认1秒





