# GO-CANAL - PHP demo

go-canal 的 php 简单调用

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
- 修改 demo 文件中监听ip和端口（可跳过）
```php
    // ip 如果有服务器资源也可换成其他ip 
    $address = "127.0.0.1";
    // 端口 随意更换
    $port = 9501;
```
- 修改配置文件
```json
{
  "schema": [ 
    "test_tt"
  ],
  "mysqlInfo": {
    "addr": "ip:3306",
    "user": "canal",
    "password": "canal"
  },
  "server": {
    "network": "tcp",
    "addr": "ip:9501"
  }
}
```

```shell
    运行 demo.php 文件 创建一个 socket 服务端
    php demo.php
```
- 确保配置文件中推送地址和 demo 文件中监听的地址一致，修改配置文件后应重启go-canal

####在数据库中随意修改或者插入一条数据，挂起的demo程序将输出一条类似数据
```
Request : {"Action":"update","ColumnData":{"admin_name":"123","created_at":null,"id":1,"password":"111","updated_at":null,"username":"111"},"SchemaName":"caopan","TableName":"admin"}
```

# 说明
- 本例子仅供参考使用，具体逻辑开发者可自行改写
- 可使用 WorkerMan Swoole 等 socket 集成包开发


