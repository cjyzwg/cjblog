---
categories:
  - SSH
---
#### 访问内网局域网网站
目的是想用在家的mac访问到公司的电话机所在的网址，首先自行搭建vpn，这个网上教程还蛮多。   
本人用的是mac，公司的电脑大多数是win10  
win10 连接vpn方法如下图：  
![avatar](https://blog.hexiefamily.xin/assets/ssh3.png)  
#### ssh 连接win10：   
1. win10 开启openssh 服务端：http://jingyan.baidu.com/article/455a995057a191a1662778a3.html 百度教程你值得拥有  
2. ssh -D 4567 -N root@192.168.12.1 输入密码，即可停留   
3. 打开chrome神奇插件SwitchyOmega,配置一个这样的代理  
![avatar](https://blog.hexiefamily.xin/assets/ssh1.png)
4. 即可打开win10 所在的局域网的网站，如下图所示
![avatar](https://blog.hexiefamily.xin/assets/ssh2.png)