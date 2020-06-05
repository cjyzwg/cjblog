### 搭建OpenVPN （目的是用本地的mac连上公司的windows电脑）

1.wget https://git.io/vpn -O openvpn-install.sh && bash openvpn-install.sh

#### 出现的问题
- 客户端连上 但是上不了外网  
nano /etc/openvpn/server.conf  
```conf
    这部分改掉，改成下面
    ;push "redirect-gateway def1 bypass-dhcp"
    push "dhcp-option DNS 223.5.5.5"
    push "dhcp-option DNS 223.6.6.6"
    原因：
    push "redirect-gateway def1 bypass-dhcp" 配置是下发了网关  下发了网关的话，客户端访问公网就不走公网网关了，都走VPN 
```

- 客户端连上了，但是互相ping不通  
```conf
在server.conf 中加一个
client-to-client
```

- 完整的server.conf 是这样的
```conf
port 1194
proto udp
dev tun
sndbuf 0
rcvbuf 0
ca ca.crt
cert server.crt
key server.key
dh dh.pem
auth SHA512
tls-auth ta.key 0
topology subnet
server 10.8.0.0 255.255.255.0
client-to-client
ifconfig-pool-persist ipp.txt
;push "redirect-gateway def1 bypass-dhcp"
push "dhcp-option DNS 223.5.5.5"
push "dhcp-option DNS 223.6.6.6"
keepalive 10 120
cipher AES-256-CBC
user nobody
group nogroup
persist-key
persist-tun
status openvpn-status.log
verb 3
crl-verify crl.pem
```

- Windows 连上上不了外网：  
windows:
tracert 114.114.114.114 最后到114.114.114.114 说明可以
mac:
traceroute 默认是UDP协议，traceroute -I 114.114.114.114 也是可以的
tracert是成功连网的。系统防火墙一般都都是阻挡的入方向，对出方向默认都是放行的，除非手动配置出方向的，才会被阻断。如果不能连接的话，一般都是在自身网关这里出问题的，或者对方禁止访问的，也有可能是网段问题（比如冲突），自身的安全策略问题可能性比较小，除非自己手动配置了出方向安全策略或者默认出方向全部拒绝。

问题解决：
client.ovpn
```conf
client
dev tun
proto udp
sndbuf 0
rcvbuf 0
remote 47.102.84.178 1194
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
auth SHA512
cipher AES-256-CBC
;setenv opt block-outside-dns
key-direction 1
verb 3
需要将 block-outside-dns  删掉
如果设置了 block-outside-dns  ，OpenVPN 会添加 Windows 防火墙记录，拦掉除 tap 以外的所有网络接口上的 DNS 请求。
如果想全部走vpn 则使用这个。  
```

#### 同时修改之后 /etc/init.d/openvpn restart  

- 中途卡住的问题
1.openvpn-install.sh 中要去调用github的代码，阿里云实例竟然获取不到github上的数据，只能手动下载，再放到服务器中，用了文件挂载 sshfs排查问题。
2.apt-get update 一直显示 build.openvpn.net connect time out
是因为/etc/apt/sources.list.d 中有一些缓存，需要直接删除，再处理。