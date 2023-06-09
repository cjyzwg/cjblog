---
title: dnsmasq配置局域网域名
date: 2022-09-05 11:43:37
categories:
  - DNSMASQ
---
### 安装dnsmasq
- apt-get install -y dnsmasq
#### /etc/dnsmasq.conf
```conf
bogus-priv
no-resolv
server=100.100.2.138
server=100.100.2.136
interface=tun0
listen-address=127.0.0.1
bind-interfaces
addn-hosts=/etc/internal_ips
```
#### /etc/internal_ips
10.8.0.1 main.int  
10.8.0.2 cj.int  
10.8.0.3 cjwindows.int  
10.8.0.4 s1.int  
10.8.0.5 s2.int   
10.8.0.6 s3.int   
10.8.0.7 s4.int    

#### dnsmasq 配置
![avatar](https://blog.hexiefamily.xin/assets/dnsmasq.jpg)