---
categories:
  - GOLANG
---
### mac下安装sshpass
>$ brew install https://raw.githubusercontent.com/kadwanev/bigboybrew/master/Library/Formula/sshpass.rb    

>$ cd /usr/local/bin  

>$ ln -s sshpass ../Cellar/sshpass/1.05/bin/sshpass  

#### deploy.sh
```shell
#!/usr/bin/env bash
password="xxxxx"
echo 'update code'
echo 'pack'
cd $GOPATH/src/learnbeego/blog/
bee pack -be GOOS=linux
echo 'upload'
sshpass -p $password scp blog.tar.gz root@test.com:/var/www/html/blog
echo 'restart'
sshpass -p $password ssh root@test.com "cd ~ && ./restart.sh"
```
#### restart.sh
```shell
#! /bin/bash
#默认进入的是登录用户的目录
cd /var/www/html/blog
tar -xzvf blog.tar.gz
#remove conf of dev
systemctl restart blog.service
```
#### 服务器使用systemd 部署
>$ nano /etc/systemd/system/dd-bi-go.service
```php
[Unit]
Description=blog
After=blog.service
[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/var/www/html/blog
ExecStart=/var/www/html/blog/blog
Restart=always
[Install]
WantedBy=multi-user.target
```
>$ systemctl start blog.service即可  

##### 重新加载配置文件
>$ sudo systemctl daemon-reload

##### 重启相关服务
>$ sudo systemctl restart foobar