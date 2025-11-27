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

## 配置文件

在项目的 `main` 目录，会生成 `.env` 配置文件，每个配置项均带注释说明。

```conf
# 七牛云的AccessKey，SecretKey
QINIU_ACCESS_KEY="xxxxxxxxxxx"
QINIU_SECRET_KEY="xxxxxxxxxxx"

# 可用的空间名列表。逗号分隔。固定好顺序不要变，有需要可往后添加。因为数据表存储的bucket_id和顺序有关。
BUCKET_NAME_LIST="bucket123,bucket456,bucket789"
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
# 同步Bucket空间文件
./main --syncbucketfiles --bucketname=bucket123 --cdnname=qiniu

# cdnname为CDN服务商名称。可省略。默认为：qiniu
# bucketname为空间名称。必填。不可省略。
./main --syncbucketfiles --bucketname=bucket123

# 更换Bucket空间
./main --syncbucketfiles --bucketname=bucket789

# 放到后台运行
nohup ./main --syncbucketfiles > syncfiles.log 2>&1 &
```

注：如果 `--bucketname` 指定的值不在 `.env` 文件的 `BUCKET_NAME_LIST` 中，会提示错误。

### 文件迁移

1. 配置 `.env`

```conf
# 文件迁移前的空间名
MIGRATE_FROM_BUCKET = "origin-bucket"
# 文件迁移后的空间名
MIGRATE_TO_BUCKET = "target-bucket"
```

2. 更新待迁移文件列表: `qiniu_cdnauth_file_migrate_list`。 必填字段 `file_key`。`status = 0` 表示待迁移。

3. 执行迁移任务，使用 `复制` 方式迁移: `./main --filemigrate` 根据待迁移文件列表，调用API复制文件到新空间。数据表 `qiniu_cdnauth_file_migrate_list` 更新 `status = 1` 表示已复制。

4. 删除迁移成功的文件: `./main --filedelete` 根据 `qiniu_cdnauth_file_migrate_list` 表 `status = 1`，删除 `复制` 成功的文件。数据表 `qiniu_cdnauth_file_migrate_list` 更新 `status = 3` 表示原文件已删除。

```bash
# 显示待迁移文件列表
./main --showmigratefiles

# 以复制文件的方式，执行文件迁移
./main --filemigrate
# 迁移的时候，给目标路径，添加前缀目录。不能以/符号开头或结尾。
./main --filemigrate --addpredir=email

# 显示待删除文件列表
./main --showdeletefiles

# 删除迁移成功的文件
./main --filedelete
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
