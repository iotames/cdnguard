## 简介

守护CDN存储空间，拦截网络攻击

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
