## 简介

守护CDN存储空间，拦截网络攻击

## 编译

在 `main` 目录下，即项目源代码入口文件 `main.go` 所在目录，编译可执行文件。

```bash
go mod tidy
cd main

# Linux编译成main文件
go build -o main main.go

# Windows编译成main.exe
go build -o main.exe main.go
```

## 常用命令

进入 `main` 文件夹，然后执行 `./main`(Linux) 或 `main.exe`(Windows) 命令:

### 启动守护进程

```bash
# 启动监听端口为1212的守护进程
# 修改main目录下的.env文件以修改系统配置
./main
```

### 数据表清理

```bash
# 数据表清理
# SQL文件：./main/sql/prune.sql
./main --prune
```

### 添加黑名单IP

```bash
# 获取当天1点后，5点前，网络请求的referer为空，且请求最多的IP
# 该时间段内，请求数超过1600，把IP加入黑名单
# 请求数限制REQUEST_LIMIT=1600。可通过 main/.env 文件修改。
./main --addblackips
```

### 同步Bucket空间的文件

```bash
# cdnname参数默认为qiniu
# bucketname参数默认为wildto，可省略不填。
./main --syncbucketfiles --bucketname=wildto --cdnname=qiniu

./main --syncbucketfiles
```

七牛云的AccessKey，SecretKey 在 `.env` 文件中配置。

```
QINIU_ACCESS_KEY="xxxxxxxxxxx"
QINIU_SECRET_KEY="xxxxxxxxxxx"
```

同步其他空间的文件：

```bash
# 同步Bucket空间文件
# bucketname为空间名称。不添加参数，默认为：wildto。可用值为: wildto, wildto-private, buerdiy, buerdiy-staging, santic-pan, sagriatech-private, santic, newwildto
./main --syncbucketfiles --bucketname=wildto-private
./main --syncbucketfiles --bucketname=buerdiy

# 放到后台运行
nohup ./main --syncbucketfiles > syncfiles.log 2>&1 &
```

## Systemd系统服务

### 系统服务配置

```bash
vim /etc/systemd/system/cdnguard.service
```

```
[Unit]
Description=Qiniu CDN Guard
After=network.target

[Service]
WorkingDirectory=/home/yourname/cdnguard/main
ExecStart=/home/yourname/cdnguard/main/cdnguard
User=yourname
Restart=on-failure
RestartSec=15
TimeoutStopSec=30

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable cdnguard
systemctl start cdnguard
```

### 系统日志配置

```
# 查看服务的日志
journalctl -u cdnguard.service

# 查看日志已占用的空间
journalctl --disk-usage

# 设置日志最大占用空间: 500M, 2G
journalctl --vacuum-size=500M

# 设置日志最大保存时间: 10d, 1years
journalctl --vacuum-time=30d
```

日志配置：

`vim /etc/systemd/journald.conf`

```conf
[Journal]
SystemMaxUse=1G
```

重启系统日志服务，使得配置立即生效：
```bash
systemctl restart systemd-journald
```
