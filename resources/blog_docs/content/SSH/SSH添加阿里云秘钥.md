---
title: SSH添加阿里云秘钥
categories:
  - SSH
---
### 阿里云服务器添加秘钥
1.登录控制台，创建秘钥对，并下载为a.pem  
2.ssh-keygen -y -f a.pem会返回ssh-rsa  
如果该命令失败，请运行chmod 400 my-key-pair.pem命令更改权限  
3.登录到服务器a.com中，echo >> ~/.ssh/authorized_keys  
4.nano /etc/ssh/sshd_config
PermitRootLogin yes
当你完成全部设置，并以密钥方式登录成功后，再禁用密码登录：
PasswordAuthentication no  
5.service sshd restart

#### 访问：
ssh -i a.pem root@a.com 访问
如果想省事的话，ssh-add -k a.pem，重启之后需要重新添加
直接ssh root@a.com 即可访问
#### 出现 Could not open a connection to your authentication agent解决办法或者Permission denied (publickey).：
>eval `ssh-agent -s`  