### 一、背景
- 之前使用的翻墙工具需要再次付费，想要免费怎么办？

### 二、解决办法

- 1.买国外的vps自行搭建服务器，算了下光买一个vps主机的价格，还不如直接续租vpn呢
- 2.使用免费的翻墙流量的软件，通常是需要签到（本文采用的办法）


### 三、安装环境 
- 1.使用树莓派作为vpn实现科学上网，树莓派烧录的是debian-pi plus++版本，**注意：树莓派的版本基本上都是arm架构**  
- 2.基于 [v2free机场](https://v2free.net/) ，文档也足够详细，非常ns，正常注册都会有1g左右的流量，然后每天签到会获得300~500M流量，对于本人已经足够用。


### 四、安装

- 使用的是linux for clash arm7版本,都是在root用户下操作的（[参考链接](https://v2free.net/doc/#/linux/clash)）

####  &emsp;    **安装clash**

- 1.下载linux arm clash版本

>wget -O clash.gz https://github.com/Dreamacro/clash/releases/download/v1.9.0/clash-linux-armv7-v1.9.0.gz  

- 2.解压到当前文件夹下  

>gzip -f clash.gz -d    

- 3.添加可执行权限  

>chmod +x clash  

- 4.移动到 /usr/local/bin 目录  

>mv clash /usr/local/bin/clash  

- 5.clash -v 命令看是否安装成功,出现以下命令说明安装成功  

>Clash v1.9.0 linux arm with go1.17.5 Sun Jan  2 03:13:32 UTC 2022

####   &emsp;  **添加配置文件**
- 1.生成config.yml文件 在~/.config/clash/config.yml

>clash

- 2.下载v2free机场的配置文件，覆盖掉config.yml

>wget -U "Mozilla/6.0" -O ~/.config/clash/config.yaml  你的Clash订阅链接网址

- 3.下载.mmdb文件，Country.mmdb为全球IP库，可以实现各个国家的IP信息解析和地理定位，没有这个文件clash是无法运行的 

>wget -O Country.mmdb https://www.sub-speeder.com/client-download/Country.mmdb

- 4.启动clash

>clash

&emsp; `ps:如果端口already bind，请在config.yml中配置新的端口即可。`

- 5.安装clash的webui界面（可选）

>wget https://github.com/haishanh/yacd/archive/gh-pages.zip  

>unzip gh-pages.zip  

>mv yacd-gh-pages/ dashboard/  

>在config.yml中添加两行 secret: "" external-ui: dashboard  

>在浏览器里面访问 http://serverip:9090/ui/ 来调试 Clash>>



#### &emsp;**添加开机自启动**
- 1.在 /lib/systemd/system/ 目录下创建 clash@.service 文件

>nano /lib/systemd/system/clash@.service  


&emsp;`写入以下内容（别修改）然后保存：`

	
		[Unit]
		Description=A rule based proxy in Go for %i.
		After=network.target
		
		[Service]
		Type=simple
		User=%i
		Restart=on-abort
		ExecStart=/usr/local/bin/clash
		
		[Install]
		WantedBy=multi-user.target
		
- 2.重新加载 systemd 模块

>systemctl daemon-reload  

- 3.启动clash服务，user 表示的是当前用户名，本文这里是root

>systemctl start clash@user  

- 4.设置开机自启动

>systemctl enable clash@user  

- 5.尝试一下

>reboot

#### &emsp;**chrome配置SwitchyOmega**
- 1.下载Proxy SwitchyOmega

>https://proxy-switchyomega.com/file/Proxy-SwitchyOmega-Chromium-2.5.15.crx

- 2.将crx文件后缀名修改为.zip格式，然后解压到当前文件夹下


- 3.打开chrome 设置-》扩展-》开发者模式，加载刚才的解压文件夹    


- 4.本地树莓派的地址是局域网地址192.168.1.107，之前配置的内网地址是10.8.0.14  


![avatar](https://blog.hexiefamily.xin/assets/switchomega.png)  

- 5.这样就实现了远程也可以科学上网的需求  



### 五、未来改进

- 追加自动脚本，每日签到
- v2free 更新，添加定时脚本更新订阅地址