# 项目说明
NAT打洞工具的Tracker服务器

## 功能 
- 简单的请求频率限制
- 每个服务器端最多提交10个ip地址，每个ip地址最多127字节
- 对服务端的服务端id和用户id进行短时间绑定，防止冒用，绑定时间可配置
- 管理员可以查看服务端的信息

- /server/list?password=xxx 查看服务端连接情况
- /server/asso?password=xxx 查看服务端id和用户id关联

## 配置文件
配置文件在/conf/app.conf
```
adminPasswd #管理员密码

ClientGetAddrListRateDuration #客户端查询服务器地址周期
ClientGetAddrListRateLimit #客户端查询服务器地址次数限制

ClientRequestConnectRateDuration 客户端请求服务器发起主动连接周期
ClientRequestConnectRateLimit 客户端请求服务器发起主动连接限制

ServerWSConnectRateDuration 服务端发起Websocket周期
ServerWSConnectRateLimit 服务端发起Websocket次数限制

ServerPushAddrListRateDuration 服务端更新地址周期
ServerPushAddrListRateLimit 服务端更新地址次数限制

KeepServerIDUserIDAssociation = 3600
```
频率限制周期为秒，在一个周期内达到限制次数时就不能请求了，下一个周期次数清空。
除了服务端更新地址频率限制，其它均为用户的ip地址。
服务端更新地址频率限制为服务端id。
